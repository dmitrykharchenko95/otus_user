package response

import (
	"bytes"
	"net/http"
	"strings"
)

type ResponseWriter struct {
	http.ResponseWriter
	code int
	buf  bytes.Buffer
}

func New(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{ResponseWriter: w, code: http.StatusOK}
}

func (rw *ResponseWriter) WriteHeader(code int) {
	rw.code = code
	rw.ResponseWriter.WriteHeader(code)
}
func (rw *ResponseWriter) Write(b []byte) (int, error) {
	if rw.code == 0 {
		rw.code = http.StatusOK
	}
	size, err := rw.ResponseWriter.Write(b)
	rw.buf.Write(b)

	return size, err
}

func (rw *ResponseWriter) GetStatus() int {
	return rw.code
}

func (rw *ResponseWriter) GetBody() string {
	return strings.TrimSpace(rw.buf.String())
}
