package fileserver

import (
	"github.com/NineNineFive/go-local-web-gui/fileserver"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOpen(t *testing.T) {
	fileSystem := fileserver.FileSystem{FileSystem: http.Dir("./frontend/test"), ReadDirBatchSize: 2}
	//fileServerSystem := http.FileServer(fileSystem)

	// Create a test server using the FileSystem
	server := httptest.NewServer(http.FileServer(fileSystem))
	defer server.Close()

	// Make a request for the root URL
	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("Error making request to test server: %v", err)
	}
	defer resp.Body.Close()

	// Check that the response doesn't contain the file listing
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code %d, but got %d", http.StatusNotFound, resp.StatusCode)
	}
}
