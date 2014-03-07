package dry

import (
	"fmt"
	"reflect"
	"sort"
	"unicode"
)

// ReflectTypeOfError returns the built-in error type
func ReflectTypeOfError() reflect.Type {
	return reflect.TypeOf((*error)(nil)).Elem()
}

/*
ExportedStructFields returns a map from exported struct field names to values,
inlining anonymous sub-structs so that their field names are available
at the base level.
Example:
	type A struct {
		X int
	}
	type B Struct {
		A
		Y int
	}
	// Yields X and Y instead of A and Y:
	InlineAnonymousStructFields(reflect.ValueOf(B{}))
*/
func ReflectExportedStructFields(v reflect.Value) map[string]reflect.Value {
	t := v.Type()
	if t.Kind() != reflect.Struct {
		panic(fmt.Errorf("Expected a struct, got %s", t))
	}
	result := make(map[string]reflect.Value)
	reflectExportedStructFields(v, t, result)
	return result
}

func reflectExportedStructFields(v reflect.Value, t reflect.Type, result map[string]reflect.Value) {
	for i := 0; i < t.NumField(); i++ {
		structField := t.Field(i)
		if ReflectStructFieldIsExported(structField) {
			if structField.Anonymous && structField.Type.Kind() == reflect.Struct {
				reflectExportedStructFields(v.Field(i), structField.Type, result)
			} else {
				result[structField.Name] = v.Field(i)
			}
		}
	}
}

func ReflectNameIsExported(name string) bool {
	return name != "" && unicode.IsUpper(rune(name[0]))
}

func ReflectStructFieldIsExported(structField reflect.StructField) bool {
	return structField.PkgPath == ""
}

// ReflectSort will sort slice according to compareFunc using reflection.
// slice can be a slice of any element type including interface{}.
// compareFunc must have two arguments that are assignable from
// the slice element type or pointers to such a type.
// The result of compareFunc must be a bool indicating
// if the first argument is less than the second.
// If the element type of slice is interface{}, then the type
// of the compareFunc arguments can be any type and dynamic
// casting from the interface value or its address will be attempted.
func ReflectSort(slice, compareFunc interface{}) {
	sortable, err := newReflectSortable(slice, compareFunc)
	if err != nil {
		panic(err)
	}
	sort.Sort(sortable)
}

func newReflectSortable(slice, compareFunc interface{}) (*reflectSortable, error) {
	t := reflect.TypeOf(compareFunc)
	if t.Kind() != reflect.Func {
		return nil, fmt.Errorf("compareFunc must be a function, got %T", compareFunc)
	}
	if t.NumIn() != 2 {
		return nil, fmt.Errorf("compareFunc must take two arguments, got %d", t.NumIn())
	}
	if t.In(0) != t.In(1) {
		return nil, fmt.Errorf("compareFunc's arguments must be identical, got %s and %s", t.In(0), t.In(1))
	}
	if t.NumOut() != 1 {
		return nil, fmt.Errorf("compareFunc must have one result, got %d", t.NumOut())
	}
	if t.Out(0).Kind() != reflect.Bool {
		return nil, fmt.Errorf("compareFunc result must be bool, got %s", t.Out(0))
	}

	argType := t.In(0)
	ptrArgs := argType.Kind() == reflect.Ptr
	if ptrArgs {
		argType = argType.Elem()
	}

	sliceV := reflect.ValueOf(slice)
	if sliceV.Kind() != reflect.Slice {
		return nil, fmt.Errorf("Need slice got %T", slice)
	}
	elemT := sliceV.Type().Elem()
	if elemT != argType && elemT.Kind() != reflect.Interface {
		return nil, fmt.Errorf("Slice element type must be interface{} or %s, got %s", argType, elemT)
	}

	return &reflectSortable{
		Slice:       sliceV,
		CompareFunc: reflect.ValueOf(compareFunc),
		ArgType:     argType,
		PtrArgs:     ptrArgs,
	}, nil
}

type reflectSortable struct {
	Slice       reflect.Value
	CompareFunc reflect.Value
	ArgType     reflect.Type
	PtrArgs     bool
}

func (self *reflectSortable) Len() int {
	return self.Slice.Len()
}

func (self *reflectSortable) Less(i, j int) bool {
	arg0 := self.Slice.Index(i)
	arg1 := self.Slice.Index(j)
	if self.Slice.Type().Elem().Kind() == reflect.Interface {
		arg0 = arg0.Elem()
		arg1 = arg1.Elem()
	}
	if (arg0.Kind() == reflect.Ptr) != self.PtrArgs {
		if self.PtrArgs {
			// Expects PtrArgs for reflectSortable, but Slice is value type
			arg0 = arg0.Addr()
		} else {
			// Expects value type for reflectSortable, but Slice is PtrArgs
			arg0 = arg0.Elem()
		}
	}
	if (arg1.Kind() == reflect.Ptr) != self.PtrArgs {
		if self.PtrArgs {
			// Expects PtrArgs for reflectSortable, but Slice is value type
			arg1 = arg1.Addr()
		} else {
			// Expects value type for reflectSortable, but Slice is PtrArgs
			arg1 = arg1.Elem()
		}
	}
	return self.CompareFunc.Call([]reflect.Value{arg0, arg1})[0].Bool()
}

func (self *reflectSortable) Swap(i, j int) {
	temp := self.Slice.Index(i).Interface()
	self.Slice.Index(i).Set(self.Slice.Index(j))
	self.Slice.Index(j).Set(reflect.ValueOf(temp))
}

// InterfaceSlice converts a slice of any type into a slice of interface{}.
func InterfaceSlice(slice interface{}) []interface{} {
	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {
		panic(fmt.Errorf("InterfaceSlice: not a slice but %T", slice))
	}
	result := make([]interface{}, v.Len())
	for i := range result {
		result[i] = v.Index(i).Interface()
	}
	return result
}