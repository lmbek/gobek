package launcher

import (
	"context"
	"fmt"
	"github.com/NineNineFive/go-local-web-gui/fileserver"
	"github.com/NineNineFive/go-local-web-gui/utils"
	"github.com/NineNineFive/go-local-web-gui/utils/net"
	"github.com/NineNineFive/go-local-web-gui/utils/random"
	"github.com/NineNineFive/go-local-web-gui/utils/slice"
	messageboxw "github.com/NineNineFive/go-local-web-gui/windows/win32api"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

type ChromeLauncher struct {
	Location                string
	LocationCMD             string
	FrontendInstallLocation string
	Domain                  string
	PreferredPort           int
	PortMin                 int
	PortMax                 int
	port                    int    // will be set doing runtime
	portAsString            string // will be set doing runtime
}

var DefaultChromeLauncher = ChromeLauncher{
	Location:                "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe",
	LocationCMD:             "C:\\\"Program Files\"\\Google\\Chrome\\Application\\chrome.exe",
	FrontendInstallLocation: os.Getenv("localappdata") + "\\Google\\Chrome\\InstalledApps\\" + "DefaultOrganisationName" + "\\" + "DefaultProjectName",
	Domain:                  "localhost",
	PortMin:                 11430,
	PreferredPort:           11451,
	PortMax:                 11500,
}

// launchChromeForWindows
// Check if chrome.exe is installed in program files (default location)
// If it is not installed then give a windows warning and exit
// Then check if this application is already installed in chrome localappdata
// if it is not installed continue (application will shut down, because frontend was not allowed to open, as backend should stop if frontend stops)
// Then continue - else check if frontend is open
// If frontend is allowed to open, because it is not already open
// Then start frontend
func (launcher *ChromeLauncher) launchForWindows() bool {
	// assert chrome is installed
	launcher.assertChromeIsInstalled()

	// check if application is already open
	frontendAlreadyOpen := launcher.isApplicationOpen()

	// open frontend if not already open
	if frontendAlreadyOpen == false {
		// get random port
		openFrontendAllowed := launcher.findAndSetAvailablePort()

		// if port found, frontend and backend is allowed to start
		if openFrontendAllowed {
			// set server address and print selected address with port
			fileserver.SetServerAddress(launcher.Domain + ":" + launcher.portAsString)
			fmt.Println("selected address with port: http://" + fileserver.GetServerAddress())

			// Start frontend by starting a new Chrome process
			path := "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe"
			cmd := exec.Command(path, "--app=http://"+fileserver.GetServerAddress(), "--user-data-dir="+launcher.FrontendInstallLocation)
			cmd.Start()

			// Set up a signal handler to gracefully shutdown the program, when it should shutdown
			signalHandler := make(chan os.Signal, 1)
			signal.Notify(signalHandler, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT) // TODO: when closing from task manager, it doesn't catch the signal

			// running through terminal (termination)
			go func() {
				<-signalHandler // waiting for termination
				cmd.Process.Kill()
				fileserver.Shutdown(context.Background())
			}()

			// running through process (close window)
			go func() {
				cmd.Wait() // waiting for close window
				fileserver.Shutdown(context.Background())
			}()

			// successfully launched the frontend
			return true
		}
	}

	// return false, if reached here (the frontend did not launch)
	return false
}

func (launcher *ChromeLauncher) findAndSetAvailablePort() bool {
	var portLength int
	// portMin needs to be 0 or above, and the preferredPort needs to be (portMin or above) or (portMax or below)
	if launcher.PortMin >= 0 && (launcher.PreferredPort >= launcher.PortMin || launcher.PortMax <= launcher.PreferredPort) {
		var prefPort int
		// it needs to be made into: make array that holds numbers from (example) 30995 to 31111
		portLength = launcher.PortMax - launcher.PortMin
		ports := make([]int, portLength)
		for i := 0; i < portLength; i++ {
			ports[i] = i + launcher.PortMin
			if ports[i] == launcher.PreferredPort {
				prefPort = i
			}
		}
		// set random seed
		random.SetRandomSeed(time.Now().UnixNano())
		n := 0
		for len(ports) > 0 {
			n++
			// Take random int in array and uses it as port, remove it from array after use

			randomInt := random.GetInt(0, len(ports)-1)
			if n == 1 {
				randomInt = prefPort
			}
			launcher.port = ports[randomInt]
			launcher.portAsString = utils.IntegerToString(launcher.port)
			// test port
			if net.IsPortUsed(launcher.Domain, launcher.portAsString) {
				fmt.Println(launcher.portAsString)
				if n == 5 {
					messageboxw.WarningManyPortsNotAvailable(launcher.PortMin, launcher.PortMax)
				} else if len(ports) == 1 {
					messageboxw.WarningNoPortsAvailable()
				}
				ports = slice.RemoveIndex(ports, randomInt)
				continue // use different port
			} else {
				return true // use this port
			}
		}
	} else {
		fmt.Println("PortMax should be higher than PortMin, and they should both be above 0")
		return false
	}
	return false
}

func (launcher *ChromeLauncher) assertChromeIsInstalled() {
	// check if chrome.exe is installed
	_, err := os.Stat(launcher.Location)

	// if not installed give warning
	if err != nil {
		messageboxw.WarningYouNeedToInstallChrome()
		os.Exit(0)
	}
}

func (launcher *ChromeLauncher) isApplicationInstalled() bool {
	// check if this application is installed
	_, err := os.Stat(launcher.FrontendInstallLocation)

	// if it is not installed continue - else check if frontend is opened already
	if err != nil {
		// ignore error message and warnings, return false as it is not installed
		return false
	} else {
		return true
	}
}

func (launcher *ChromeLauncher) isApplicationOpen() bool {
	var alreadyOpen bool
	isInstalled := launcher.isApplicationInstalled()

	if isInstalled {
		// check if frontend is opened, by checking if we can rename its folder (is it locked?)
		// TODO: this can be optimized, so we better can check if frontend is already open.
		// Currently it can open multiple frontends, if it is installing (because it takes 2 seconds to install)
		err := os.Rename(launcher.FrontendInstallLocation, launcher.FrontendInstallLocation) // check lock
		if err != nil {
			fmt.Println("Frontend Already open... assuming Backend is too") // it is locked
			fmt.Println("Otherwise close the open Frontend before launching")
			fmt.Println("Both needs to not be running in order to start the program")
			alreadyOpen = true
		} else { // If it could rename, then it is not locked, open frontend (as it is not already open)
			alreadyOpen = false
		}
	} else {
		alreadyOpen = false // TODO: we should probably rework this - we can wait for it to be installed (wait 1 second) and try again, or we can rework how we look if application is already working, entirely
	}

	return alreadyOpen
}
