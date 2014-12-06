package dry

import (
	"strings"
	"testing"
)

func Test_MapS(t *testing.T) {
	result := MapS(strings.TrimSpace, []string{"  a  ", " b ", "c", "  d", "e  "})
	correct := []string{"a", "b", "c", "d", "e"}
	if len(result) != len(correct) {
		t.Fail()
	}
	for i, _ := range result {
		if result[i] != correct[i] {
			t.Fail()
		}
	}
}

func Test_FilterS(t *testing.T) {
	hFunc := func(s string) bool {
		return strings.HasPrefix(s, "h")
	}
	result := FilterS(hFunc, []string{"cheese", "mouse", "hi", "there", "horse"})
	correct := []string{"hi", "horse"}
	if len(result) != len(correct) {
		t.Fail()
	}
	for i, _ := range result {
		if result[i] != correct[i] {
			t.Fail()
		}
	}
}
