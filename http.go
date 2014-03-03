package dry

import (
	"compress/flate"
	"compress/gzip"
	"encoding/json"
	"encoding/xml"
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

// HTTPRespondMarshalJSON marshals response as JSON to responseWriter, sets Content-Type to application/json
// and compresses the response if Content-Encoding from the request allows it.
func HTTPRespondMarshalJSON(response interface{}, responseWriter http.ResponseWriter, request *http.Request) (err error) {
	handlerFunc := HTTPCompressHandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		var data []byte
		if data, err = json.Marshal(response); err == nil {
			responseWriter.Header().Set("Content-Type", "application/json")
			_, err = responseWriter.Write(data)
		}
	})
	handlerFunc(responseWriter, request)
	return err
}

// HTTPRespondMarshalIndentJSON marshals response as JSON to responseWriter, sets Content-Type to application/json
// and compresses the response if Content-Encoding from the request allows it.
// The JSON will be marshalled indented according to json.MarshalIndent
func HTTPRespondMarshalIndentJSON(response interface{}, prefix, indent string, responseWriter http.ResponseWriter, request *http.Request) (err error) {
	handlerFunc := HTTPCompressHandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		var data []byte
		if data, err = json.MarshalIndent(response, prefix, indent); err == nil {
			responseWriter.Header().Set("Content-Type", "application/json")
			_, err = responseWriter.Write(data)
		}
	})
	handlerFunc(responseWriter, request)
	return err
}

// HTTPRespondMarshalXML marshals response as XML to responseWriter, sets Content-Type to application/xml
// and compresses the response if Content-Encoding from the request allows it.
func HTTPRespondMarshalXML(response interface{}, responseWriter http.ResponseWriter, request *http.Request) (err error) {
	handlerFunc := HTTPCompressHandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		var data []byte
		if data, err = xml.Marshal(response); err == nil {
			responseWriter.Header().Set("Content-Type", "application/xml")
			_, err = responseWriter.Write(data)
		}
	})
	handlerFunc(responseWriter, request)
	return err
}

// HTTPRespondMarshalIndentXML marshals response as XML to responseWriter, sets Content-Type to application/xml
// and compresses the response if Content-Encoding from the request allows it.
// The XML will be marshalled indented according to xml.MarshalIndent
func HTTPRespondMarshalIndentXML(response interface{}, prefix, indent string, responseWriter http.ResponseWriter, request *http.Request) (err error) {
	handlerFunc := HTTPCompressHandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		var data []byte
		if data, err = xml.MarshalIndent(response, prefix, indent); err == nil {
			responseWriter.Header().Set("Content-Type", "application/xml")
			_, err = responseWriter.Write(data)
		}
	})
	handlerFunc(responseWriter, request)
	return err
}

// HTTPRespondText sets Content-Type to text/plain
// and compresses the response if Content-Encoding from the request allows it.
func HTTPRespondText(response string, responseWriter http.ResponseWriter, request *http.Request) (err error) {
	handlerFunc := HTTPCompressHandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Header().Set("Content-Type", "text/plain")
		_, err = responseWriter.Write([]byte(response))
	})
	handlerFunc(responseWriter, request)
	return err
}
