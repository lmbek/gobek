package fileserver

import (
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"github.com/lmbek/gobek/fileserver/gzipResponse"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var FrontendPath = "./frontend" // should be set doing runtime by main.go
var Server = http.Server{
	Addr:              "", // should be set doing runtime by main.go with fileserver.SetServerAddress
	Handler:           nil,
	TLSConfig:         nil,
	ReadTimeout:       5 * time.Second,
	ReadHeaderTimeout: 20 * time.Second,
	WriteTimeout:      10 * time.Second,
	IdleTimeout:       0,
	MaxHeaderBytes:    0,
	TLSNextProto:      nil,
	ConnState:         nil,
	ErrorLog:          nil,
	BaseContext:       nil,
	ConnContext:       nil,
}
var ServerGraceShutdownTime = 5 * time.Second

func SetServerAddress(address string) {
	fmt.Println("address set to: " + address)
	Server.Addr = address
}

func GetServerAddress() string {
	return Server.Addr
}

func ServeFileServer(response http.ResponseWriter, request *http.Request) {

	fileSystem := FileSystem{http.Dir(FrontendPath), 2}
	fileServerSystem := http.FileServer(fileSystem)

	response = setHeaders(response, request)

	if !strings.Contains(request.Header.Get("Accept-Encoding"), "gzip") {
		fileServerSystem.ServeHTTP(response, request)
	} else {
		response.Header().Set("Content-Encoding", "gzip")
		gzipping := gzip.NewWriter(response)
		defer gzipping.Close()
		fileServerSystem.ServeHTTP(gzipResponse.Writer{Writer: gzipping, ResponseWriter: response}, request)
	}
}

func setHeaders(response http.ResponseWriter, request *http.Request) http.ResponseWriter {
	// Headers can be set here
	// Add Cache Cache-Control: max-age=31536000, immutable

	// response.Header().Add("Cache-Control", "max-age=31536000, immutable")

	// Check if the requested file has a ".css" extension
	if strings.HasSuffix(request.URL.Path, ".css") {
		response.Header().Set("Content-Type", "text/css")
	}

	response.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
	response.Header().Set("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
	return response
}

func Start() error {
	err := Server.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}

func Shutdown(serverContext context.Context) error {
	serverContext, cancel := context.WithTimeout(serverContext, ServerGraceShutdownTime)
	defer cancel()
	err := Server.Shutdown(serverContext)
	if err != nil {
		log.Println("Failed to gracefully shutdown the server:", err)
		return err
	}
	log.Println("Server has been shut down gracefully")
	return nil
}

func GracefulStart() error {
	err := Start()

	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return errors.New("ErrServerClosed: " + err.Error())
	} else if err != nil {
		return err
	} else {
		_, closeChannel := CreateChannel()
		defer closeChannel()
		log.Println("Application stopped gracefully")
		return Shutdown(context.Background())
	}
}

func CreateChannel() (chan os.Signal, func()) {
	stopChannel := make(chan os.Signal, 1)
	signal.Notify(stopChannel, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	return stopChannel, func() {
		close(stopChannel)
	}
}
