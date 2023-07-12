package gobek

import (
	"context"
	"fmt"
	"github.com/lmbek/gobek/fileserver"
	"net"
	"os/exec"
	"strconv"
)

type ChromiumLauncher struct {
	Location string
}

var DefaultChromiumLauncher = ChromiumLauncher{
	Location: "/var/lib/snapd/desktop/applications/chromium_chromium.desktop", // TODO: check if better location or can be customised
}

func (launcher *ChromiumLauncher) LaunchForLinux() bool {
	// Listen on a random available port on localhost
	listen, err := net.Listen("tcp", fileserver.GetServerAddress())
	if err != nil {
		fmt.Println(err)
	}
	addr := listen.Addr().(*net.TCPAddr)
	port := strconv.Itoa(addr.Port)
	err = listen.Close()
	if err != nil {
		fmt.Println("Could not close listening on Addr, error: ", err)
	}
	fileserver.SetServerAddress("localhost:" + port) // set random available port
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
		err = cmd.Wait()
		if err != nil {
			fmt.Println(err)
		}
		// shutting down file server, graceful shutdown probably not needed, as api can still finish, probably
		err = fileserver.Shutdown(context.Background())
		if err != nil {
			fmt.Println("warning, could not shut server down: " + err.Error())
		}
	}()

	return true
}
