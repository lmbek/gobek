package gzipResponse

import (
	"bytes"
	"compress/gzip"
	"github.com/lmbek/gobek/fileserver/gzipResponse"
	"io/ioutil"
	"net/http/httptest"
	"testing"
)

func TestWriter(t *testing.T) {
	// Create a mock http.ResponseWriter
	responseRecorder := httptest.NewRecorder()

	// SETTING header to GZip
	responseRecorder.Header().Set("Content-Encoding", "gzip")

	// Create a new gzipResponse.Writer, wrapping the mock response writer
	gzipWriter := gzipResponse.Writer{gzip.NewWriter(responseRecorder), responseRecorder}

	// Write some content to the gzip writer
	content := []byte("test content")
	bytesWritten, err := gzipWriter.Write(content)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Flush and close the gzip writer
	err = gzipWriter.Close()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify that the content was compressed and written correctly
	if bytesWritten == 0 {
		t.Errorf("no bytes were written")
	}

	// Decode the compressed content
	gzipReader, err := gzip.NewReader(bytes.NewReader(responseRecorder.Body.Bytes()))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	defer gzipReader.Close()

	// Read the decompressed content
	decompressedContent, err := ioutil.ReadAll(gzipReader)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify that the decompressed content matches the original content
	if !bytes.Equal(decompressedContent, content) {
		t.Errorf("unexpected content: expected %q, got %q", content, decompressedContent)
	}

	// Verify that the response headers were set correctly
	if responseRecorder.Header().Get("Content-Encoding") != "gzip" {
		t.Errorf("unexpected content encoding: expected %q, got %q", "gzip", responseRecorder.Header().Get("Content-Encoding"))
	}
}
