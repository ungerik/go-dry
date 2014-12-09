package dry

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"
)

type CountingReader struct {
	Reader    io.Reader
	BytesRead int
}

func (r *CountingReader) Read(p []byte) (n int, err error) {
	n, err = r.Reader.Read(p)
	r.BytesRead += n
	return n, err
}

type CountingWriter struct {
	Writer       io.Writer
	BytesWritten int
}

func (r *CountingWriter) Write(p []byte) (n int, err error) {
	n, err = r.Writer.Write(p)
	r.BytesWritten += n
	return n, err
}

type CountingReadWriter struct {
	ReadWriter   io.ReadWriter
	BytesRead    int
	BytesWritten int
}

func (rw *CountingReadWriter) Read(p []byte) (n int, err error) {
	n, err = rw.ReadWriter.Read(p)
	rw.BytesRead += n
	return n, err
}

func (rw *CountingReadWriter) Write(p []byte) (n int, err error) {
	n, err = rw.ReadWriter.Write(p)
	rw.BytesWritten += n
	return n, err
}

// ReadBinary wraps binary.Read with a CountingReader and returns
// the acutal bytes read by it.
func ReadBinary(r io.Reader, order binary.ByteOrder, data interface{}) (n int, err error) {
	countingReader := CountingReader{Reader: r}
	err = binary.Read(&countingReader, order, data)
	return countingReader.BytesRead, err
}

// WriteFull calls writer.Write until all of data is written,
// or an is error returned.
func WriteFull(data []byte, writer io.Writer) (n int, err error) {
	dataSize := len(data)
	for n = 0; n < dataSize; {
		m, err := writer.Write(data[n:])
		n += m
		if err != nil {
			return n, err
		}
	}
	return dataSize, nil
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
