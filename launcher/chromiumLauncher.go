package launcher

import (
	"context"
	"fmt"
	"github.com/NineNineFive/go-local-web-gui/fileserver"
	"github.com/NineNineFive/go-local-web-gui/utils"
	"github.com/NineNineFive/go-local-web-gui/utils/net"
	"github.com/NineNineFive/go-local-web-gui/utils/random"
	"github.com/NineNineFive/go-local-web-gui/utils/slice"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"
)

type ChromiumLauncher struct {
	Location      string
	Domain        string
	PreferredPort int
	PortMin       int
	PortMax       int
	port          int
	portAsString  string
}

var DefaultChromiumLauncher = ChromiumLauncher{
	Location:      "/var/lib/snapd/desktop/applications/chromium_chromium.desktop", // TODO: check if better location or can be customised
	Domain:        "localhost",
	PortMin:       11430,
	PreferredPort: 11451,
	PortMax:       11500,
}

func (launcher *ChromiumLauncher) launchForLinux(waitgroup *sync.WaitGroup) (bool, *sync.WaitGroup) {

	// 1) Check if chromium is installed
	_, err := os.Stat(launcher.Location)
	if err != nil {
		// TODO: add linux warning for not installed chromium
		//messageboxw.WarningYouNeedToInstallChromium()
		fmt.Println("Need to install Chromium")
		os.Exit(0)
	} else {
	}

	openBackendAllowed := launcher.findAndSetAvailablePort()

	if openBackendAllowed {
		fileserver.SetServerAddress(launcher.Domain + ":" + launcher.portAsString)
		fmt.Println("selected address with port: http://" + fileserver.GetServerAddress())
		// Start frontend
		go func() {
			//cmd := exec.Command("chromium", "--temp-profile", "--app=http://"+fileserver.GetServerAddress())
			cmd := exec.Command("chromium", "--incognito", "--temp-profile", "--app=http://"+fileserver.GetServerAddress())
			cmd.Stdout = os.Stdout
			err := cmd.Run()
			if err != nil {
				log.Fatal(err)
			}
			fileserver.Shutdown(context.Background())
		}()
		return true, waitgroup
	}
	return false, waitgroup
}

func (launcher *ChromiumLauncher) findAndSetAvailablePort() bool {
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
					fmt.Println("not many ports, please wait") // not many ports warning
					//messageboxw.WarningManyPortsNotAvailable(launcher.PortMin, launcher.PortMax)
				} else if len(ports) == 1 {
					fmt.Println("no ports") // no ports warning
					//messageboxw.WarningNoPortsAvailable()
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

func (launcher *ChromiumLauncher) GetPort() int {
	return launcher.port
}

func (launcher *ChromiumLauncher) GetPortAsString() string {
	return launcher.portAsString
}
