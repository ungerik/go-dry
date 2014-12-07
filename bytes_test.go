package dry

import (
	"testing"
)

func Test_BytesMap(t *testing.T) {
	upper := func(b byte) byte {
		return b - ('a' - 'A')
	}
	result := BytesMap(upper, []byte("hello"))
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

func Test_BytesFilter(t *testing.T) {
	azFunc := func(b byte) bool {
		return b >= 'A' && b <= 'Z'
	}
	result := BytesFilter(azFunc, []byte{1, 2, 3, 'A', 'f', 'R', 123})
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
