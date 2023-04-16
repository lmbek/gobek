package gzipResponse

import (
	"compress/gzip"
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

func (gzipWriter Writer) Close() error {
	// Flush the gzip writer
	err := gzipWriter.Writer.(*gzip.Writer).Flush()
	if err != nil {
		return err
	}

	// Close the gzip writer
	return gzipWriter.Writer.(*gzip.Writer).Close()
}
