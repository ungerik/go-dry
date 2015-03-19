package dry

import (
	"bytes"
	"fmt"
)

// PanicIfErr panics with a stack trace if any
// of the passed values is a non nil error
func PanicIfErr(values ...interface{}) {
	for _, v := range values {
		if err, _ := v.(error); err != nil {
			panic(fmt.Errorf("Panicking because of error: %s\nAt:\n%s\n", err, StackTrace(2)))
		}
	}
}

// Nop is a dummy function that can be called in source files where
// other debug functions are constantly added and removed.
// That way import "github.com/ungerik/go-quick" won't cause an error when
// no other debug function is currently used.
// Arbitrary objects can be passed as arguments to avoid "declared and not used"
// error messages when commenting code out and in.
// The result is a nil interface{} dummy value.
func Nop(dummiesIn ...interface{}) (dummyOut interface{}) {
	return nil
}

// Error returns r as error, converting it when necessary
func AsError(r interface{}) error {
	if r == nil {
		return nil
	}
	if err, _ := r.(error); err != nil {
		return err
	} else {
		return fmt.Errorf("%v", r)
	}
}

// Returns the first non nil error, or nil
func FirstError(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

// Returns the last non nil error, or nil
func LastError(errs ...error) error {
	for i := len(errs) - 1; i >= 0; i-- {
		err := errs[i]
		if err != nil {
			return err
		}
	}
	return nil
}

// AsErrorList checks if err is already an ErrorList
// and returns it if this is the case.
// Else an ErrorList with err as element is created.
// Useful if a function potentially returns an ErrorList as error
// and you want to avoid creating nested ErrorLists.
func AsErrorList(err error) ErrorList {
	if list, ok := err.(ErrorList); ok {
		return list
	}
	return ErrorList{err}
}

/*
ErrorList holds a slice of errors.

Usage example:

	func maybeError() (int, error) {
		return
	}

	func main() {
		e := NewErrorList(maybeError())
		e.Collect(maybeError())
		e.Collect(maybeError())

		if e.Err() != nil {
			fmt.Println("Some calls of maybeError() returned errors:", e)
		} else {
			fmt.Println("No call of maybeError() returned an error")
		}
	}
*/
type ErrorList []error

// NewErrorList returns an ErrorList where Collect has been called for args.
// The returned list will be nil if there was no non nil error in args.
// Note that alle methods of ErrorList can be called with a nil ErrorList.
func NewErrorList(args ...interface{}) (list ErrorList) {
	list.Collect(args...)
	return list
}

// Error calls fmt.Println for of every error in the list
// and returns the concernated text.
// Can be called for a nil ErrorList.
func (list ErrorList) Error() string {
	if len(list) == 0 {
		return "Empty ErrorList"
	}
	var buf bytes.Buffer
	for _, err := range list {
		fmt.Fprintln(&buf, err)
	}
	return buf.String()
}

// Err returns the list if it is not empty,
// or nil if it is empty.
// Can be called for a nil ErrorList.
func (list ErrorList) Err() error {
	if len(list) == 0 {
		return nil
	}
	return list
}

// First returns the first error in the list or nil.
// Can be called for a nil ErrorList.
func (list ErrorList) First() error {
	if len(list) == 0 {
		return nil
	}
	return list[0]
}

// Last returns the last error in the list or nil.
// Can be called for a nil ErrorList.
func (list ErrorList) Last() error {
	if len(list) == 0 {
		return nil
	}
	return list[len(list)-1]
}

// Collect adds any non nil errors in args to the list.
func (list *ErrorList) Collect(args ...interface{}) {
	for _, a := range args {
		if err, _ := a.(error); err != nil {
			*list = append(*list, err)
		}
	}
}
