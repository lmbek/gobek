package launcher

import (
	"errors"
	"fmt"
	"github.com/NineNineFive/go-local-web-gui/fileserver"
	"net/http"
	"runtime"
	"sync"
)

func Start(frontendPath string, chromeLauncher ChromeLauncher, chromiumLauncher ChromiumLauncher) error {
	http.HandleFunc("/", fileserver.ServeFileServer)
	return StartCustom(frontendPath, chromeLauncher, chromiumLauncher)
}

func StartCustom(frontendPath string, chromeLauncher ChromeLauncher, chromiumLauncher ChromiumLauncher) error {
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
	launched := chromeLauncher.launchForWindows()

	// Start Backend (if frontend is allowed to launch - not opened already)
	if launched {
		err := startServer(frontendPath)
		if err != nil {

			if err.Error() == "http: Server closed" {
				// if http server is closed, we assume it is not an error that caused it
				// we write to the console, that the http server is closed
				fmt.Println("HTTP Server is now closed")
				return nil
			} else {
				return err
			}
		}
	}

	return nil
}

// StartOnLinux
// start frontend (chrome) and backend (http localhost) on Windows
func StartOnLinux(frontendPath string, chromiumLauncher ChromiumLauncher) error {
	fmt.Println("Attempting to start on: " + runtime.GOOS + ", " + runtime.GOARCH)
	var waitgroup *sync.WaitGroup
	waitgroup = &sync.WaitGroup{}
	waitgroup.Add(1)
	// Start Frontend
	launched, waitgroup := chromiumLauncher.launchForLinux(waitgroup)
	if launched {
		// Start Backend
		err := startServer(frontendPath)
		if err != nil {
			fmt.Println("Error...")
			fmt.Println(err)
			waitgroup.Done()
			return err
		}
		waitgroup.Done()
	} else {
		waitgroup.Done()
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
