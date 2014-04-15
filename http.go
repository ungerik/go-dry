package dry

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
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
	handler := &HTTPCompressHandler{Handler: handlerFunc}
	return func(response http.ResponseWriter, request *http.Request) {
		handler.ServeHTTP(response, request)
	}
}

// HTTPCompressHandler wraps a http.Handler so that the response gets
// gzip or deflate compressed if the Accept-Encoding header of the request allows it.
type HTTPCompressHandler struct {
	Handler http.Handler
}

func (h *HTTPCompressHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	accept := request.Header.Get("Accept-Encoding")
	if strings.Contains(accept, "gzip") {
		response.Header().Set("Content-Encoding", "gzip")
		writer := Gzip.GetWriter(response)
		defer Gzip.ReturnWriter(writer)
		response = wrappedResponseWriter{Writer: writer, ResponseWriter: response}
	} else if strings.Contains(accept, "deflate") {
		response.Header().Set("Content-Encoding", "deflate")
		writer := Deflate.GetWriter(response)
		defer Deflate.ReturnWriter(writer)
		response = wrappedResponseWriter{Writer: writer, ResponseWriter: response}
	}
	h.Handler.ServeHTTP(response, request)
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
// If rootElement is not empty, then an additional root element with this name will be wrapped around the content.
func HTTPRespondMarshalXML(response interface{}, rootElement string, responseWriter http.ResponseWriter, request *http.Request) (err error) {
	handlerFunc := HTTPCompressHandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		var data []byte
		if data, err = xml.Marshal(response); err == nil {
			responseWriter.Header().Set("Content-Type", "application/xml")
			if rootElement == "" {
				_, err = fmt.Fprintf(responseWriter, "%s%s", xml.Header, data)
			} else {
				_, err = fmt.Fprintf(responseWriter, "%s<%s>%s</%s>", xml.Header, rootElement, data, rootElement)
			}
		}
	})
	handlerFunc(responseWriter, request)
	return err
}

// HTTPRespondMarshalIndentXML marshals response as XML to responseWriter, sets Content-Type to application/xml
// and compresses the response if Content-Encoding from the request allows it.
// The XML will be marshalled indented according to xml.MarshalIndent.
// If rootElement is not empty, then an additional root element with this name will be wrapped around the content.
func HTTPRespondMarshalIndentXML(response interface{}, rootElement string, prefix, indent string, responseWriter http.ResponseWriter, request *http.Request) (err error) {
	handlerFunc := HTTPCompressHandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		var data []byte
		contentPrefix := prefix
		if rootElement != "" {
			contentPrefix += indent
		}
		if data, err = xml.MarshalIndent(response, contentPrefix, indent); err == nil {
			responseWriter.Header().Set("Content-Type", "application/xml")
			if rootElement == "" {
				_, err = fmt.Fprintf(responseWriter, "%s%s\n%s", prefix, xml.Header, data)
			} else {
				_, err = fmt.Fprintf(responseWriter, "%s%s%s<%s>\n%s\n%s</%s>", prefix, xml.Header, prefix, rootElement, data, prefix, rootElement)
			}
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
