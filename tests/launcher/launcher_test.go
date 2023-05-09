package launcher

import (
	"errors"
	"fmt"
	"github.com/lmbek/gobek/launcher"
	"github.com/lmbek/gobek/tests/helpers"
	"os"
	"sync"
	"testing"
	"time"
)

// TODO: in version 0.7.0 this should be splitted into each of the chrome/chromium launchers instead

// SIMPLE VERSION
func TestStart(test *testing.T) {
	//tests.PrintGotFatalError(test, "NOT IMPLEMENTED YET!")
	waitgroup := &sync.WaitGroup{}
	waitgroup.Add(1)

	launcher.InitDefault()

	go func() {
		// shutdown app after 3 seconds
		time.Sleep(time.Second * 3)
		err := launcher.Shutdown()
		if err != nil {
			fmt.Println(err)
		}

		// we have to wait a few seconds to let the fileserver show down gracefully before running the next test
		waitgroup.Done()
	}()

	result := launcher.Start("./../_frontend", launcher.DefaultChromeLauncher, launcher.DefaultChromiumLauncher).Error()
	expected := errors.New("http: Server closed").Error()
	helpers.StandardTestChecking(test, result, expected)

	waitgroup.Wait()
}

func TestStart2(test *testing.T) {
	//tests.PrintGotFatalError(test, "NOT IMPLEMENTED YET!")
	waitgroup := &sync.WaitGroup{}
	waitgroup.Add(1)

	launcher.InitDefault()

	go func() {
		// shutdown app after 3 seconds
		time.Sleep(time.Second * 3)
		err := launcher.Shutdown()
		if err != nil {
			fmt.Println(err)
		}

		waitgroup.Done()
	}()

	// since we have a http.HandleFunc with a registration on / , we need to only do it once. This will prevent multiple registrations
	result := launcher.Start("./../_frontend", launcher.DefaultChromeLauncher, launcher.DefaultChromiumLauncher).Error()
	expected := errors.New("http: Server closed").Error()
	helpers.StandardTestChecking(test, result, expected)
	waitgroup.Wait()
}

// ADVANCED VERSION
func TestStartWithCustomSettings(test *testing.T) {
	//tests.PrintGotFatalError(test, "NOT IMPLEMENTED YET!")
	waitgroup := &sync.WaitGroup{}
	waitgroup.Add(1)

	launcher.InitDefault()

	go func() {
		// shutdown app after 3 seconds
		time.Sleep(time.Second * 3)
		err := launcher.Shutdown()
		if err != nil {
			fmt.Println(err)
		}

		waitgroup.Done()
	}()

	// For windows we need a organisation name and project name
	var organisationName = "NewOrganisationName" // put in organisation name
	var projectName = "NewProjectName"           // put in project name

	var frontendPath = "./../_frontend" // this should be set to where frontend files is (frontend folder: html, css, javascript...)

	var chromeLauncher = launcher.ChromeLauncher{
		Location:                os.Getenv("programfiles") + "\\Google\\Chrome\\Application\\chrome.exe",
		FrontendInstallLocation: os.Getenv("localappdata") + "\\Google\\Chrome\\InstalledApps\\" + organisationName + "\\" + projectName,
	}

	var chromiumLauncher = launcher.ChromiumLauncher{
		Location: "/var/lib/snapd/desktop/applications/chromium_chromium.desktop", // TODO: check if better location or can be customised
		Domain:   "localhost",
	}

	result := launcher.Start(frontendPath, chromeLauncher, chromiumLauncher).Error()
	expected := errors.New("http: Server closed").Error()
	helpers.StandardTestChecking(test, result, expected)
	waitgroup.Wait()
}
