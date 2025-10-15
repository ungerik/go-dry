package dry

import (
	"reflect"
	"strings"
	"testing"
)

func Test_StringMap(t *testing.T) {
	result := StringMap(strings.TrimSpace, []string{"  a  ", " b ", "c", "  d", "e  "})
	correct := []string{"a", "b", "c", "d", "e"}
	if len(result) != len(correct) {
		t.Fail()
	}
	for i := range result {
		if result[i] != correct[i] {
			t.Fail()
		}
	}
}

func Test_StringFilter(t *testing.T) {
	hFunc := func(s string) bool {
		return strings.HasPrefix(s, "h")
	}
	result := StringFilter(hFunc, []string{"cheese", "mouse", "hi", "there", "horse"})
	correct := []string{"hi", "horse"}
	if len(result) != len(correct) {
		t.Fail()
	}
	for i := range result {
		if result[i] != correct[i] {
			t.Fail()
		}
	}
}

func Test_StringFindBetween(t *testing.T) {
	s := "Hello <em>World</em>!"

	between, remainder, found := StringFindBetween(s, "<em>", "</em>")
	if between != "World" {
		t.Fail()
	}
	if remainder != "!" {
		t.Fail()
	}
	if !found {
		t.Fail()
	}

	between, remainder, found = StringFindBetween(s, "l", "l")
	if between != "" {
		t.Fail()
	}
	if remainder != "o <em>World</em>!" {
		t.Fail()
	}
	if !found {
		t.Fail()
	}

	between, remainder, found = StringFindBetween(s, "<i>", "</i>")
	if between != "" {
		t.Fail()
	}
	if remainder != "Hello <em>World</em>!" {
		t.Fail()
	}
	if found {
		t.Fail()
	}

}

func Test_StringStripHTMLTags(t *testing.T) {
	withHTML := "<div>Hello > World <br/> <im src='xxx'/>"
	skippedHTML := "Hello > World  "

	if StringStripHTMLTags(withHTML) != skippedHTML {
		t.Fail()
	}
}

func Test_StringReplaceHTMLTags(t *testing.T) {
	withHTML := "<div>Hello > World <br/> <im src='xxx'/>"
	replacedHTML := "xxHello > World xx xx"

	if StringReplaceHTMLTags(withHTML, "xx") != replacedHTML {
		t.Fail()
	}
}

func Test_TwoSlicesSubtraction(t *testing.T) {
	A := []string{"apple", "orange", "banana", "peach", "plum"}
	B := []string{"melon", "banana", "guava", "plum"}
	wanted := []string{"apple", "orange", "peach"}
	result := TwoSlicesSubtraction(A, B)
	if !reflect.DeepEqual(wanted, result) {
		t.Errorf("wanted: %v, but got: %v", wanted, result)
	}
}

func Test_StringAddURLParam(t *testing.T) {
	tests := []struct {
		url      string
		name     string
		value    string
		expected string
	}{
		{"http://example.com", "foo", "bar", "http://example.com?foo=bar"},
		{"http://example.com?existing=param", "foo", "bar", "http://example.com?existing=param&foo=bar"},
		{"http://example.com/path", "key", "value", "http://example.com/path?key=value"},
		{"http://example.com/path?a=1", "b", "2", "http://example.com/path?a=1&b=2"},
	}

	for _, tt := range tests {
		result := StringAddURLParam(tt.url, tt.name, tt.value)
		if result != tt.expected {
			t.Errorf("StringAddURLParam(%q, %q, %q) = %q, want %q",
				tt.url, tt.name, tt.value, result, tt.expected)
		}
	}
}

func Test_StringSplitOnce(t *testing.T) {
	tests := []struct {
		s        string
		sep      string
		wantPre  string
		wantPost string
	}{
		{"hello:world", ":", "hello", "world"},
		{"one:two:three", ":", "one", "two:three"},
		{"no-separator", ":", "no-separator", ""},
		{"", ":", "", ""},
		{"start:", ":", "start", ""},
		{":end", ":", "", "end"},
	}

	for _, tt := range tests {
		pre, post := StringSplitOnce(tt.s, tt.sep)
		if pre != tt.wantPre || post != tt.wantPost {
			t.Errorf("StringSplitOnce(%q, %q) = (%q, %q), want (%q, %q)",
				tt.s, tt.sep, pre, post, tt.wantPre, tt.wantPost)
		}
	}
}
