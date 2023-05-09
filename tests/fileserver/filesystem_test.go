package fileserver

import (
	"fmt"
	"github.com/lmbek/gobek/fileserver"
	"github.com/lmbek/gobek/tests/helpers"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOpen(test *testing.T) {
	// Code is sensitive, this code should not run parallel with other tests, therefore using waitgroup
	//waitgroup := &sync.WaitGroup{}
	//done := make(chan struct{})

	test.Run("no index no files", noIndexNoFiles)
	test.Run("with index", withIndex)
	test.Run("with index with slash", withIndexWithSlash)
	test.Run("no index no files", withIndexNoFiles)

	// anonymous function test (unique from the others)
	test.Run("with index no files - with default go fileserver", func(t *testing.T) {
		// testing fileserver.FileSystem
		fileSystem := http.FileServer(http.Dir("./../_frontend/some-folder-with-files-no-index/"))

		// start a http server and test file system for correct responses
		server := httptest.NewServer(fileSystem)
		defer server.Close()

		// Make a request for the root URL
		response, err := http.Get(server.URL)
		if err != nil {
			fmt.Printf("\t\tError making request to test server: %v \n", err)
		}
		defer response.Body.Close()

		// result and expected
		result := response.StatusCode
		expected := http.StatusOK

		// check if result is the same as expected
		helpers.StandardTestChecking(test, result, expected)
	})
}

// Test to ensure status 404: there is no file listing / indexing returned when calling domain/some-folder, when no index.html file is found
func noIndexNoFiles(test *testing.T) {
	// testing fileserver.FileSystem
	fileSystem := fileserver.FileSystem{FileSystem: http.Dir("./../_frontend/some-empty-folder"), ReadDirBatchSize: 2}

	// start a http server and test file system for correct responses
	response := newServerHttpGet(fileSystem)

	// result and expected
	result := response.StatusCode
	expected := http.StatusNotFound

	// check if result is the same as expected
	helpers.StandardTestChecking(test, result, expected)
}

// Test to ensure we get status 200, if folder has index.html
func withIndex(test *testing.T) {
	// testing fileserver.FileSystem
	fileSystem := fileserver.FileSystem{FileSystem: http.Dir("./../_frontend/some-folder-with-index"), ReadDirBatchSize: 2}
	//fileServerSystem := http.FileServer(fileSystem)

	// start a http server and test file system for correct responses
	response := newServerHttpGet(fileSystem)

	// result and expected
	result := response.StatusCode
	expected := http.StatusOK

	// check if result is the same as expected
	helpers.StandardTestChecking(test, result, expected)
}

// Test to ensure we get status 200, if folder has index.html (when testing with /)
func withIndexWithSlash(test *testing.T) {
	// testing fileserver.FileSystem
	fileSystem := fileserver.FileSystem{FileSystem: http.Dir("./../_frontend/some-folder-with-index/"), ReadDirBatchSize: 2}
	//fileServerSystem := http.FileServer(fileSystem)

	// start a http server and test file system for correct responses
	response := newServerHttpGet(fileSystem)

	// result and expected
	result := response.StatusCode
	expected := http.StatusOK

	// check if result is the same as expected
	if result != expected {
		test.Errorf("\t\tExpected %d, but got %d", expected, result)
	} else {
		fmt.Printf("\t\tGot expected: %v \n", expected)
	}
}

// Test to ensure we get status 404, if folder has no index.html but still files
func withIndexNoFiles(test *testing.T) {
	// testing fileserver.FileSystem
	fileSystem := fileserver.FileSystem{FileSystem: http.Dir("./../_frontend/some-folder-with-files-no-index"), ReadDirBatchSize: 2}
	//fileServerSystem := http.FileServer(fileSystem)

	// start a http server and test file system for correct responses
	response := newServerHttpGet(fileSystem)

	// result and expected
	result := response.StatusCode
	expected := http.StatusNotFound

	// check if result is the same as expected
	if result != expected {
		test.Errorf("\t\tExpected %d, but got %d", expected, result)
	} else {
		fmt.Printf("\t\tGot expected: %v \n", expected)
	}
}

// Helper function
// repeatable func for less duplicate code
// newServerHttpGet - creates a test server using the file system and make a request for the URL

func newServerHttpGet(fileSystem fileserver.FileSystem) *http.Response {
	fileServer := http.FileServer(fileSystem)
	// start a http server and test file system for correct responses
	server := httptest.NewServer(fileServer)
	defer server.Close()

	// Make a request for the root URL
	response, err := http.Get(server.URL)
	if err != nil {
		fmt.Printf("\t\tError making request to test server: %v \n", err)
	}
	defer response.Body.Close()

	return response
}
