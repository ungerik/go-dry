package dry

import (
	"testing"
)

func Test_MapB(t *testing.T) {
	upper := func(b byte) byte {
		return b - ('a' - 'A')
	}
	result := MapB(upper, []byte("hello"))
	correct := []byte("HELLO")
	if len(result) != len(correct) {
		t.Fail()
	}
	for i, _ := range result {
		if result[i] != correct[i] {
			t.Fail()
		}
	}
}

func Test_FilterB(t *testing.T) {
	azFunc := func(b byte) bool {
		return b >= 'A' && b <= 'Z'
	}
	result := FilterB(azFunc, []byte{1, 2, 3, 'A', 'f', 'R', 123})
	correct := []byte{'A', 'R'}
	if len(result) != len(correct) {
		t.Fail()
	}
	for i, _ := range result {
		if result[i] != correct[i] {
			t.Fail()
		}
	}
}
