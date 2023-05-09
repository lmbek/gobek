package launcher

import (
	"context"
	"fileserver"
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"
)

type ChromeLauncher struct {
	Location                string
	FrontendInstallLocation string
}

var DefaultChromeLauncher = ChromeLauncher{
	Location:                os.Getenv("programfiles") + "\\Google\\Chrome\\Application\\chrome.exe",
	FrontendInstallLocation: os.Getenv("localappdata") + "\\Google\\Chrome\\InstalledApps\\" + "DefaultOrganisationName" + "\\" + "DefaultProjectName",
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
	// check if application is already open
	frontendAlreadyOpen := launcher.isApplicationOpen()

	// open frontend if not already open
	if frontendAlreadyOpen == false {
		// Listen on a random available port on localhost
		listen, err := net.Listen("tcp", fileserver.GetServerAddress())
		if err != nil {
			fmt.Println(err)
		}
		addr := listen.Addr().(*net.TCPAddr)
		port := strconv.Itoa(addr.Port)
		listen.Close()

		// set server address
		fileserver.SetServerAddress("localhost:" + port)

		// Print the port that was found
		fmt.Println("selected address with port: http://" + fileserver.GetServerAddress())

		// Start frontend by starting a new Chrome process
		path := launcher.Location

		cmd = exec.Command(path, "--app=http://"+fileserver.GetServerAddress(), "--user-data-dir="+launcher.FrontendInstallLocation)
		err = cmd.Start()
		if err != nil {
			println("Warning: Chrome could not start, is it installed?")
		}

		// Set up a signal handler to shutdown the program, when it should shutdown
		signalHandler := make(chan os.Signal, 1)
		signal.Notify(signalHandler, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT) // TODO: when closing from task manager, it doesn't catch the signal

		// TODO: Add context with timeout handler (or find out which context to use)
		// and then find out if it can stop Task Manager from Exiting the program too early
		// - we need to kill cmd process if it happens

		// TODO: here is some of the code (from https://github.com/halilylm/go-redis/blob/main/main.go)
		// ctx, cancel := context.WithTimeout(context.Background(),10 * time.Second)
		//

		// running through terminal (termination)
		go func() {
			<-signalHandler // waiting for termination
			err = cmd.Process.Kill()
			if err != nil {
				fmt.Println("warning, could not wait for cmd.Process.Kill(): " + err.Error())
			}
			err = fileserver.Shutdown(context.Background())
			if err != nil {
				fmt.Println("warning, could not shut server down: " + err.Error())
			}
		}()

		// running through process (close window)
		go func() {
			err = cmd.Wait() // waiting for close window

			if err != nil {
				fmt.Println("warning, could not wait for cmd.Wait (probably this app shut itself down internally with launcher.Shutdown()): " + err.Error())
			}

			// shutting down file server, graceful shutdown probably not needed, as api can still finish, probably
			err = fileserver.Shutdown(context.Background())
			if err != nil {
				fmt.Println("warning, could not shut server down: " + err.Error())
			}
		}()

		// successfully launched the frontend
		return true
	}
	// return false, if reached here (the frontend did not launch)
	return false
}

/*
func (launcher *ChromeLauncher) assertChromeIsInstalled() {
	// check if chrome.exe is installed
	_, err := os.Stat(launcher.Location)

	// if not installed give warning
	if err != nil {
		//messageboxw.WarningYouNeedToInstallChrome()
		os.Exit(0)
	}
}
*/

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
