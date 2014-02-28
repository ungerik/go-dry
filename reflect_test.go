package dry

import (
	// "strings"
	// "fmt"
	"testing"
)

func Test_ReflectSort(t *testing.T) {
	ints := []int{3, 5, 0, 2, 1, 4}
	ReflectSort(ints, func(a, b int) bool {
		return a < b
	})
	for i := range ints {
		if i != ints[i] {
			t.Fail()
		}
	}

	strings := []string{"aaa", "bbb", "abb", "aab"}
	ReflectSort(strings, func(a, b string) bool {
		return a < b
	})
	if strings[0] != "aaa" {
		t.Fail()
	}
	if strings[1] != "aab" {
		t.Fail()
	}
	if strings[2] != "abb" {
		t.Fail()
	}
	if strings[3] != "bbb" {
		t.Fail()
	}
}
