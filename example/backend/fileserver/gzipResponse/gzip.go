package gzipResponse

import (
	"io"
	"net/http"
)

type Writer struct {
	io.Writer
	http.ResponseWriter
}

func (gzipResponse Writer) Write(b []byte) (int, error) {
	return gzipResponse.Writer.Write(b)
}
