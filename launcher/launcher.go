package launcher

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

func InitDefault() {
	// since we have a http.HandleFunc with a registration on / , we need to only do it once. This will prevent multiple registrations
	once.Do(func() {
		http.HandleFunc("/", fileserver.ServeFileServer)
	})
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
		return errors.New("Darwin Not Supported Yet")
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

	}

	return nil
}

func StartDefault(frontendPath string, chromeLauncher ChromeLauncher, chromiumLauncher ChromiumLauncher) error {
	InitDefault()
	return Start(frontendPath, chromeLauncher, chromiumLauncher)
}

func Start(frontendPath string, chromeLauncher ChromeLauncher, chromiumLauncher ChromiumLauncher) error {
	switch runtime.GOOS {
	case "windows":
		return StartOnWindows(frontendPath, chromeLauncher)
	case "darwin": // "mac"
		return errors.New("Darwin Not Supported Yet")
	case "linux": // "linux"
		return StartOnLinux(frontendPath, chromiumLauncher)
	default: // "freebsd", "openbsd", "netbsd"
		return StartOnLinux(frontendPath, chromiumLauncher)
	}
}

// StartOnWindows
// start frontend (chrome) and backend (http localhost) on Windows
func StartOnWindows(frontendPath string, chromeLauncher ChromeLauncher) error {
	fmt.Println("Attempting to start on: " + runtime.GOOS + ", " + runtime.GOARCH)
	// Start Frontend
	// TODO: look into reworking this into looking for if same process is already open, and make it possible to allow for multiple clients of same in the chromeLauncher struct (bool)
	launched := chromeLauncher.launchForWindows()
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
	launched := chromiumLauncher.launchForLinux()
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
	fmt.Println("Server is now running!")
	fileserver.FrontendPath = frontendPath
	return fileserver.GracefulStart()
}
