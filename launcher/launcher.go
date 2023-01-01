package launcher

import (
	"fmt"
	"github.com/NineNineFive/go-local-web-gui/fileserver"
	"runtime"
	"sync"
)

// StartOnWindows
// start frontend (chrome) and backend (http localhost) on Windows
func StartOnWindows(frontendPath string, chromeLauncher ChromeLauncher) {
	fmt.Println("Attempting to start on: " + runtime.GOOS + ", " + runtime.GOARCH)
	// Start Frontend
	launched := chromeLauncher.launchForWindows()

	// Start Backend (if frontend is allowed to launch - not opened already)
	if launched {
		err := StartServer(frontendPath)
		if err != nil {
			fmt.Println(err)
		}
	}
}

// StartOnLinux
// start frontend (chrome) and backend (http localhost) on Windows
func StartOnLinux(frontendPath string, chromiumLauncher ChromiumLauncher) {
	fmt.Println("Attempting to start on: " + runtime.GOOS + ", " + runtime.GOARCH)
	var waitgroup *sync.WaitGroup
	waitgroup = &sync.WaitGroup{}
	waitgroup.Add(1)
	// Start Frontend
	launched, waitgroup := chromiumLauncher.launchForLinux(waitgroup)
	if launched {
		// Start Backend
		err := StartServer(frontendPath)
		if err != nil {
			fmt.Println(err)
		}
		waitgroup.Done()
	} else {
		waitgroup.Done()
	}
}

// StartServer
// Starts the backend (http localhost)
func StartServer(frontendPath string) error {
	fmt.Println("Server is now running!")
	fileserver.FrontendPath = frontendPath
	return fileserver.GracefulStart()
}
