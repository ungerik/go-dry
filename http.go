package dry

import (
	"compress/flate"
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type wrappedResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (wrapped wrappedResponseWriter) Write(data []byte) (int, error) {
	return wrapped.Writer.Write(data)
}

// HTTPCompressHandlerFunc wraps a http.HandlerFunc so that the response gets
// gzip or deflate compressed if the Accept-Encoding header of the request allows it.
func HTTPCompressHandlerFunc(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		accept := request.Header.Get("Accept-Encoding")
		if strings.Contains(accept, "gzip") {
			response.Header().Set("Content-Encoding", "gzip")
			writer := gzip.NewWriter(response)
			defer writer.Close()
			response = wrappedResponseWriter{Writer: writer, ResponseWriter: response}
		} else if strings.Contains(accept, "deflate") {
			response.Header().Set("Content-Encoding", "deflate")
			writer, _ := flate.NewWriter(response, flate.BestCompression)
			defer writer.Close()
			response = wrappedResponseWriter{Writer: writer, ResponseWriter: response}
		}
		handlerFunc(response, request)
	}
}

func HTTPDelete(url string) (statusCode int, statusText string, err error) {
	request, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return 0, "", err
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return 0, "", err
	}
	return response.StatusCode, response.Status, nil
}
