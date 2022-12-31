package launcher

import (
	"fileserver"
	"fmt"
	"runtime"
	"sync"
)

func StartFrontendAndBackendWindows(frontendPath string, launcher ChromeLauncher) {
	fmt.Println("Attempting to start on: " + runtime.GOOS + ", " + runtime.GOARCH)
	// Start Frontend
	launched := launcher.launchChromeForWindows()

	// Start Backend (if frontend is allowed to launch - not opened already)
	if launched {
		err := StartServer(frontendPath)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func StartFrontendAndBackendLinux(frontendPath string, launcher ChromiumLauncher) {
	fmt.Println("Attempting to start on: " + runtime.GOOS + ", " + runtime.GOARCH)
	var waitgroup *sync.WaitGroup
	waitgroup = &sync.WaitGroup{}
	waitgroup.Add(1)
	// Start Frontend
	launched, waitgroup := launcher.launchChromiumForLinux(waitgroup)
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

func StartServer(frontendPath string) error {
	fmt.Println("how is this even running?")
	fileserver.FrontendPath = frontendPath
	return fileserver.GracefulStart()
}
