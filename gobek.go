package gobek

import (
	"errors"
	"fmt"
	"github.com/lmbek/gobek/fileserver"
	"net/http"
	"os/exec"
	"runtime"
	"sync"
)

var once sync.Once
var cmd *exec.Cmd

// StartDefault - This is the function that should be used when starting the program, if no other custom behavior is needed
// This is meant to be what should be used by users of this open source system
func StartDefault(frontendPath string, chromeLauncher ChromeLauncher, chromiumLauncher ChromiumLauncher) error {
	InitDefault()
	return Start(frontendPath, chromeLauncher, chromiumLauncher)
}

func InitDefault() {
	// since we have a http.HandleFunc with a registration on / , we need to only do it once. This will prevent multiple registrations
	once.Do(func() {
		http.HandleFunc("/", fileserver.ServeFileServer)
	})
}

func Start(frontendPath string, chromeLauncher ChromeLauncher, chromiumLauncher ChromiumLauncher) error {
	switch runtime.GOOS {
	case "windows":
		return StartOnWindows(frontendPath, chromeLauncher)
	case "darwin": // "mac"
		return errors.New("darwin not supported yet")
	case "linux": // "linux"
		return StartOnLinux(frontendPath, chromiumLauncher)
	default: // "freebsd", "openbsd", "netbsd"
		return errors.New("other os than windows and linux not supported yet")
	}
}

func Shutdown() error {
	switch runtime.GOOS {
	case "windows":
		if cmd == nil || cmd.Process == nil {
			return errors.New("chrome process is not running")
		}
		err := cmd.Process.Kill()
		if err != nil {
			fmt.Println("warning, could not wait for cmd.Process.Kill(): " + err.Error())
		}
	case "darwin": // "mac"
		return errors.New("darwin not supported yet")
	case "linux": // "linux"
		if cmd == nil || cmd.Process == nil {
			return errors.New("chromium process is not running")
		}
		// Cannot kill the chromium process bc it wont close before program termination, it is linked to our process
		err := cmd.Process.Kill()
		if err != nil {
			fmt.Println("warning, could not wait for cmd.Process.Kill(): " + err.Error())
		}
	default: // "freebsd", "openbsd", "netbsd"
		return errors.New("other os than windows and linux not supported yet")
	}

	return nil
}

// StartOnWindows
// start frontend (chrome) and backend (http localhost) on Windows
func StartOnWindows(frontendPath string, chromeLauncher ChromeLauncher) error {
	fmt.Println("Attempting to start on: " + runtime.GOOS + ", " + runtime.GOARCH)
	// Start Frontend
	// TODO: look into reworking this into looking for if same process is already open, and make it possible to allow for multiple clients of same in the chromeLauncher struct (bool)
	// right now we are checking if we can rename chrome.exe to the same name, which is not an optimal solution
	launched := chromeLauncher.LaunchForWindows()
	// Start Backend (if frontend is allowed to launch - not opened already)
	if launched {
		err := startServer(frontendPath)
		if err != nil {
			return err
		}
	}
	return nil
}

// StartOnLinux
// start frontend (chrome) and backend (http localhost) on Windows
func StartOnLinux(frontendPath string, chromiumLauncher ChromiumLauncher) error {
	fmt.Println("Attempting to start on: " + runtime.GOOS + ", " + runtime.GOARCH)
	// Start Frontend
	// TODO: look into reworking this into looking for if same process is already open, and make it possible to allow for multiple clients of same in the chromeLauncher struct (bool)
	launched := chromiumLauncher.LaunchForLinux()
	if launched {
		// Start Backend
		err := startServer(frontendPath)
		if err != nil {
			return err
		}
	}
	return nil
}

// StartServer
// Starts the backend (http localhost)
func startServer(frontendPath string) error {
	fmt.Println("http: Server starting (localhost)")
	fileserver.FrontendPath = frontendPath
	return fileserver.GracefulStart()
}
