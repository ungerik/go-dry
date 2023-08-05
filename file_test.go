package dry

import (
	"strings"
	"testing"
)

func Test_FileGetString(t *testing.T) {
	_, err := FileGetString("invalid_file")
	if err == nil {
		t.Fail()
	}

	str, err := FileGetString("LICENSE")
	if err != nil {
		t.Error(err)
	}
	if !strings.HasPrefix(str, "The MIT License (MIT)") {
		t.Fail()
	}

	str, err = FileGetString("https://raw.githubusercontent.com/ungerik/go-dry/master/LICENSE")
	if err != nil {
		t.Error(err)
	}
	if !strings.HasPrefix(str, "The MIT License (MIT)") {
		t.Fail()
	}
}

func Test_FileIsDir(t *testing.T) {
	if FileIsDir("testfile.txt") {
		t.Fail()
	}
	if !FileIsDir(".") {
		t.Fail()
	}
}
