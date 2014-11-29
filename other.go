package dry

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"
)

func RandSeedWithTime() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// WriteAllBytes calls writer.Write until all of data is written,
// or an is error returned.
func WriteAllBytes(data []byte, writer io.Writer) error {
	n, err := writer.Write(data)
	if err != nil {
		return err
	}
	dataSize := len(data)
	for i := n; i < dataSize; i += n {
		n, err = writer.Write(data[i:])
		if err != nil {
			return err
		}
	}
	return nil
}

// ReadLine reads unbuffered until a newline '\n' byte and removes
// an optional carriege return '\r' at the end of the line.
// In case of an error, the string up to the error is returned.
func ReadLine(reader io.Reader) (line string, err error) {
	buffer := bytes.NewBuffer(make([]byte, 0, 4096))
	p := make([]byte, 1)
	for {
		_, err = reader.Read(p)
		if err != nil || p[0] == '\n' {
			break
		}
		buffer.Write(p)
	}
	line = strings.TrimSuffix(buffer.String(), "\r")
	return line, err
}

// WaitForStdin blocks until input is available from os.Stdin.
// The first byte from os.Stdin is returned as result.
// If there are println arguments, then fmt.Println will be
// called with those before reading from os.Stdin.
func WaitForStdin(println ...interface{}) byte {
	if len(println) > 0 {
		fmt.Println(println...)
	}
	buffer := make([]byte, 1)
	os.Stdin.Read(buffer)
	return buffer[0]
}
