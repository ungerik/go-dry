package dry

import (
	"bytes"
	"fmt"
)

// PanicIfErr panics with a stack trace if any
// of the passed values is a non nil error
func PanicIfErr(values ...interface{}) {
	for _, v := range values {
		if err, ok := v.(error); ok {
			if err != nil {
				panic(fmt.Errorf("Panicking because of error: %s\nAt:\n%s\n", err, StackTrace(2)))
			}
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

// ErrorList holds a slice of errors
type ErrorList []error

// Error calls fmt.Println for of every error in the list
// and returns the concernated text.
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

func (list ErrorList) First() error {
	if len(list) == 0 {
		return nil
	}
	return list[0]
}

func (list ErrorList) Last() error {
	if len(list) == 0 {
		return nil
	}
	return list[len(list)-1]
}
