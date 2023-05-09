package fileserver

import (
	"context"
	"errors"
	"fmt"
	"github.com/lmbek/gobek/fileserver"
	"github.com/lmbek/gobek/helpers"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

var once sync.Once

func TestGetServerAddress(test *testing.T) {
	// result and expected
	result := fileserver.GetServerAddress()
	expected := "localhost:0" // expect default value (localhost:0)

	// check if result is the same as expected
	helpers.StandardTestChecking(test, result, expected)
}

func TestSetServerAddress(test *testing.T) {
	original := fileserver.GetServerAddress()

	// Test SetServerAddress
	fileserver.SetServerAddress("example.com:12345")

	// result and expected
	result := fileserver.GetServerAddress()
	expected := "example.com:12345"

	// set variable back to original (so chained tests does not fail)
	fileserver.SetServerAddress(original)

	// check if result is the same as expected
	helpers.StandardTestChecking(test, result, expected)

}

func TestServeFileServer(test *testing.T) {
	// Create a mock HTTP request with Accept-Encoding header set to "gzip"
	request, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		test.Fatal(err)
	}
	request.Header.Set("Accept-Encoding", "gzip")

	// Create a mock HTTP response
	recorder := httptest.NewRecorder()

	// Call ServeFileServer function with mock request and response
	fileserver.FrontendPath = "./../_frontend" // setting frontend path as we have custom path for testing only (should work with any path)
	fileserver.ServeFileServer(recorder, request)

	// Check if Content-Encoding header is set to "gzip" by setting result and expected
	result := recorder.Header().Get("Content-Encoding")
	expected := "gzip"

	// check if result is the same as expected
	helpers.StandardTestChecking(test, result, expected)
}

func TestFileServer(test *testing.T) {
	test.Run("fileserver(run)", func(test *testing.T) {
		// Code is sensitive, this code should not run parallel with other tests, therefore using waitgroup
		waitgroup := &sync.WaitGroup{}
		done := make(chan struct{})

		// initialise fileserver defaults
		fileserver.FrontendPath = "./../_frontend"

		once.Do(func() {
			http.HandleFunc("/", fileserver.ServeFileServer)
		})

		go func() {
			// shutdown app after 3 seconds
			time.Sleep(time.Second * 3)
			err := fileserver.Shutdown(context.Background()) //app.Shutdown()
			if err != nil {
				fmt.Println(err)
			}
			close(done)
		}()

		waitgroup.Add(1)
		go func() {
			defer waitgroup.Done()
			<-done
		}()

		result := fileserver.GracefulStart().Error()
		expected := errors.New("http: Server closed").Error()
		helpers.StandardTestChecking(test, result, expected)

		// wait before running other tests as this one has operation critical procedures that requires a bit of isolation
		waitgroup.Wait()
	})

	test.Run("fileserver(shutdown)", func(test *testing.T) {
		// Code is sensitive, this code should not run parallel with other tests, therefore using waitgroup
		waitgroup := &sync.WaitGroup{}
		waitgroup.Add(1)
		result := error(nil)
		go func() {
			result = fileserver.Shutdown(context.Background())
			waitgroup.Done()
		}()

		expected := error(nil)
		helpers.StandardTestChecking(test, result, expected)
		// wait before running other tests as this one has operation critical procedures that requires a bit of isolation
		waitgroup.Wait()
	})

	test.Run("fileserver(run->shutdown)", func(test *testing.T) {
		// Code is sensitive, this code should not run parallel with other tests, therefore using waitgroup
		waitgroup := &sync.WaitGroup{}
		waitgroup.Add(1)
		done := make(chan struct{})

		// initialise fileserver defaults
		fileserver.FrontendPath = "./../_frontend"

		once.Do(func() {
			http.HandleFunc("/", fileserver.ServeFileServer)
		})

		go func() {
			// shutdown app after 3 seconds
			time.Sleep(time.Second * 3)
			err := fileserver.Shutdown(context.Background()) //app.Shutdown()
			if err != nil {
				fmt.Println(err)
			}
			close(done)
		}()

		go func() {
			defer waitgroup.Done()
			<-done
		}()

		result := fileserver.GracefulStart().Error()
		expected := errors.New("http: Server closed").Error()
		helpers.StandardTestChecking(test, result, expected)

		// wait before running other tests as this one has operation critical procedures that requires a bit of isolation
		waitgroup.Wait()
	})
}
