package dry

import (
	// "strings"
	// "fmt"
	"errors"
	"testing"
)

func Test_Error(t *testing.T) {
	err := Error("TestError")
	if err == nil || err.Error() != "TestError" {
		t.Fail()
	}

	err = Error(errors.New("TestError"))
	if err == nil || err.Error() != "TestError" {
		t.Fail()
	}
}
