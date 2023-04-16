package fileserver_test

import (
	"errors"
	"github.com/NineNineFive/go-local-web-gui/fileserver"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSetAndGetServerAddress(t *testing.T) {
	// Test SetServerAddress
	fileserver.SetServerAddress("localhost:80")

	// Check if Server.Addr is set correctly
	if fileserver.Server.Addr != "localhost:80" {
		t.Errorf("unexpected server address: got %v, want %v", fileserver.Server.Addr, "localhost:80")
	}

	// Test GetServerAddress
	addr := fileserver.GetServerAddress()

	// Check if GetServerAddress returns the correct address
	if addr != "localhost:80" {
		t.Errorf("unexpected server address: got %v, want %v", addr, "localhost:80")
	}
}

func TestServeFileServer(t *testing.T) {
	// Create a mock HTTP request with Accept-Encoding header set to "gzip"
	request, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Accept-Encoding", "gzip")

	// Create a mock HTTP response
	recorder := httptest.NewRecorder()

	// Call ServeFileServer function with mock request and response
	fileserver.ServeFileServer(recorder, request)

	// Check if response status code is 200
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("unexpected status code: got %v, want %v", status, http.StatusOK)
	}

	// Check if Content-Encoding header is set to "gzip"
	if encoding := recorder.Header().Get("Content-Encoding"); encoding != "gzip" {
		t.Errorf("unexpected content encoding: got %v, want %v", encoding, "gzip")
	}
}

// TODO: this test needs to be improved to test what happens when the program terminates
/*
func TestStart(t *testing.T) {
	http.HandleFunc("/", fileserver.ServeFileServer)
	fileserver.FrontendPath = "./frontend"

	go func() {
		err := fileserver.Start()
		if err != nil {
			t.Fatalf("Start() returned an unexpected error: %v", err)
		}
	}()

	time.Sleep(1 * time.Second) // Wait for the server to start listening

	resp, err := http.Get("http://localhost")
	if err != nil {
		t.Fatalf("Failed to connect to the server: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestShutdown(t *testing.T) {
	err := fileserver.Shutdown(context.Background())
	if err != nil {
		t.Fatalf("Failed to shutdown the server gracefully: %v", err)
	}
}
*/
// TODO: this test needs to be improved to test what happens when the program terminates
func TestGracefulStart(t *testing.T) {
	t.Run("subtest depending on TestShutdown", func(t *testing.T) {
		http.HandleFunc("/", fileserver.ServeFileServer)
		fileserver.FrontendPath = "./frontend"

		go func() {
			err := fileserver.GracefulStart()
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				t.Fatalf("GracefulStart() returned an unexpected error: %v", err)
			}
		}()

		time.Sleep(1 * time.Second) // Wait for the server to start listening

		resp, err := http.Get("http://localhost")
		if err != nil {
			t.Fatalf("Failed to connect to the server: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
		}
	})
}
