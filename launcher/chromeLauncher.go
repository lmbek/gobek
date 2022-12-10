package launcher

import (
	"context"
	"fileserver"
	"fmt"
	"os"
	"sync"
	"time"
	"utils"
	"utils/net"
	"utils/random"
	"utils/slice"
	goCmd "windows/go-cmd"
	messageboxw "windows/win32api"
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

func (launcher *ChromeLauncher) launchChromeForWindows(waitgroup *sync.WaitGroup) (bool, *sync.WaitGroup) {
	var openFrontendAllowed bool

	_, err := os.Stat(launcher.Location)
	if err != nil {
		messageboxw.WarningYouNeedToInstallChrome()
		os.Exit(0)
	} else {
	}

	// 1) Check if app is installed
	_, err = os.Stat(launcher.FrontendInstallLocation)

	if err != nil { // it is not installed
		openFrontendAllowed = true
	} else { // it is installed
		err = os.Rename(launcher.FrontendInstallLocation, launcher.FrontendInstallLocation) // check lock
		if err != nil {
			// TODO: We need a GO version of the C# WindowsProcessor project
			// TODO: Improvements to be made with this: good open and close of app, and maybe better boot time
			fmt.Println("Frontend Already open... assuming Backend is too") // it is locked
			fmt.Println("Otherwise close the open Frontend before launching")
			fmt.Println("Both needs to not be running in order to start the program")
			openFrontendAllowed = false
		} else { // If it could rename, then it is not locked, open frontend
			openFrontendAllowed = true
		}
	}

	if openFrontendAllowed {
		// get random port
		openBackendAllowed := launcher.findAndSetAvailablePort()
		if openBackendAllowed {
			fileserver.SetServerAddress(launcher.Domain + ":" + launcher.portAsString)
			fmt.Println("selected address with port: http://" + fileserver.GetServerAddress())
			// Start frontend
			go func() {
				commandPrompt := goCmd.NewCmd(
					"cmd.exe",
					"/C",
					"start "+launcher.LocationCMD+" --app=http://"+fileserver.GetServerAddress()+" --user-data-dir="+launcher.FrontendInstallLocation)
				statusChan := commandPrompt.Start()
				<-statusChan // blocking
				fileserver.Shutdown(context.Background())
				waitgroup.Wait()
			}()
			return true, waitgroup
		}
		return false, waitgroup
	}
	return false, waitgroup
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
