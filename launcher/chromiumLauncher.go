package launcher

import (
	"context"
	"fileserver"
	"fmt"
	"net"
	"os/exec"
	"strconv"
)

// TODO: the linux version will get same port and have some issues. In version 0.7.0 this should be fixed, we will add a struct to make the launchers unique

type ChromiumLauncher struct {
	Location string
	Domain   string
}

var DefaultChromiumLauncher = ChromiumLauncher{
	Location: "/var/lib/snapd/desktop/applications/chromium_chromium.desktop", // TODO: check if better location or can be customised
	Domain:   "localhost",
}

func (launcher *ChromiumLauncher) launchForLinux() bool {
	fmt.Println("i am running!")
	// Listen on a random available port on localhost
	listen, err := net.Listen("tcp", fileserver.GetServerAddress())
	if err != nil {
		fmt.Println(err)
	}
	addr := listen.Addr().(*net.TCPAddr)
	port := strconv.Itoa(addr.Port)
	listen.Close()
	fileserver.SetServerAddress(launcher.Domain + ":" + port) // set random available port
	fmt.Println("selected address with port: http://" + fileserver.GetServerAddress())
	// Start frontend

	//cmd := exec.Command("chromium", "--temp-profile", "--app=http://"+fileserver.GetServerAddress())
	cmd = exec.Command("chromium", "--incognito", "--temp-profile", "--app=http://"+fileserver.GetServerAddress())
	//cmd.Stdout = os.Stdout

	err = cmd.Start()
	if err != nil {
		fmt.Println(err)
		return false
	}

	go func() {
		cmd.Wait()
		// shutting down file server, graceful shutdown probably not needed, as api can still finish, probably
		err = fileserver.Shutdown(context.Background())
		if err != nil {
			fmt.Println("warning, could not shut server down: " + err.Error())
		}
	}()

	return true
}

/*
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
*/
