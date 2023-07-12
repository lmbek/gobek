package launcher

import (
	"errors"
	"fmt"
	"github.com/lmbek/gobek"
	"github.com/lmbek/gobek/tests/helpers"
	"os"
	"sync"
	"testing"
	"time"
)

// main.go edition (to show methods that can be used in your own main.go)
func TestMainStart(test *testing.T) {
	waitgroup := &sync.WaitGroup{}
	waitgroup.Add(1)

	// shutdown for test after 3 seconds
	go func() {
		// shutdown app after 3 seconds
		time.Sleep(time.Second * 3)
		err := gobek.Shutdown()
		if err != nil {
			fmt.Println(err)
		}

		waitgroup.Done()
	}()

	// CODE STARTS HERE
	// For windows we need a organisation name and project name
	var organisationName = "NewOrganisationName" // put in organisation name
	var projectName = "NewProjectName"           // put in project name

	var frontendPath = "./../_frontend" // this should be set to where frontend files is (frontend folder: html, css, javascript...)

	var chromeLauncher = gobek.ChromeLauncher{
		Location:                os.Getenv("programfiles") + "\\Google\\Chrome\\Application\\chrome.exe",
		FrontendInstallLocation: os.Getenv("localappdata") + "\\Google\\Chrome\\InstalledApps\\" + organisationName + "\\" + projectName,
	}

	var chromiumLauncher = gobek.ChromiumLauncher{
		Location: "/var/lib/snapd/desktop/applications/chromium_chromium.desktop", // TODO: check if better location or can be customised
	}

	result := gobek.StartDefault(frontendPath, chromeLauncher, chromiumLauncher).Error()

	expected := errors.New("http: Server closed").Error()
	helpers.StandardTestChecking(test, result, expected)
	waitgroup.Wait()
}

// SIMPLE VERSION
func TestStart(test *testing.T) {
	time.Sleep(time.Second * 1) // wait 1 second because this test fails bc of concurrent operations
	//tests.PrintGotFatalError(test, "NOT IMPLEMENTED YET!")
	waitgroup := &sync.WaitGroup{}
	waitgroup.Add(1)

	gobek.InitDefault()

	go func() {
		// shutdown app after 3 seconds
		time.Sleep(time.Second * 3)
		err := gobek.Shutdown()
		if err != nil {
			fmt.Println(err)
		}

		// we have to wait a few seconds to let the fileserver show down gracefully before running the next test
		waitgroup.Done()
	}()

	result := gobek.Start("./../_frontend", gobek.DefaultChromeLauncher, gobek.DefaultChromiumLauncher).Error()
	expected := errors.New("http: Server closed").Error()
	helpers.StandardTestChecking(test, result, expected)

	waitgroup.Wait()
}

func TestStart2(test *testing.T) {
	time.Sleep(time.Second * 1) // wait 1 second because this test fails bc of concurrent operations
	//tests.PrintGotFatalError(test, "NOT IMPLEMENTED YET!")
	waitgroup := &sync.WaitGroup{}
	waitgroup.Add(1)

	gobek.InitDefault()

	go func() {
		// shutdown app after 3 seconds
		time.Sleep(time.Second * 3)
		err := gobek.Shutdown()
		if err != nil {
			fmt.Println(err)
		}

		waitgroup.Done()
	}()

	// since we have a http.HandleFunc with a registration on / , we need to only do it once. This will prevent multiple registrations
	result := gobek.Start("./../_frontend", gobek.DefaultChromeLauncher, gobek.DefaultChromiumLauncher).Error()
	expected := errors.New("http: Server closed").Error()
	helpers.StandardTestChecking(test, result, expected)
	waitgroup.Wait()
}

// ADVANCED VERSION
func TestStartWithCustomSettings(test *testing.T) {
	time.Sleep(time.Second * 1) // wait 1 second because this test fails bc of concurrent operations
	//tests.PrintGotFatalError(test, "NOT IMPLEMENTED YET!")
	waitgroup := &sync.WaitGroup{}
	waitgroup.Add(1)

	gobek.InitDefault()

	go func() {
		// shutdown app after 3 seconds
		time.Sleep(time.Second * 3)
		err := gobek.Shutdown()
		if err != nil {
			fmt.Println(err)
		}

		waitgroup.Done()
	}()

	// For windows we need a organisation name and project name
	var organisationName = "NewOrganisationName" // put in organisation name
	var projectName = "NewProjectName"           // put in project name

	var frontendPath = "./../_frontend" // this should be set to where frontend files is (frontend folder: html, css, javascript...)

	var chromeLauncher = gobek.ChromeLauncher{
		Location:                os.Getenv("programfiles") + "\\Google\\Chrome\\Application\\chrome.exe",
		FrontendInstallLocation: os.Getenv("localappdata") + "\\Google\\Chrome\\InstalledApps\\" + organisationName + "\\" + projectName,
	}

	var chromiumLauncher = gobek.ChromiumLauncher{
		Location: "/var/lib/snapd/desktop/applications/chromium_chromium.desktop", // TODO: check if better location or can be customised
	}

	result := gobek.Start(frontendPath, chromeLauncher, chromiumLauncher).Error()
	expected := errors.New("http: Server closed").Error()
	helpers.StandardTestChecking(test, result, expected)
	waitgroup.Wait()
}
