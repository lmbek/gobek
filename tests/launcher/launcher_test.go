package launcher_test

import (
	"github.com/NineNineFive/go-local-web-gui/launcher"
	"testing"
)

// SIMPLE VERSION
func TestStart(test *testing.T) {
	var err error
	err = launcher.Start("./frontend", launcher.DefaultChromeLauncher, launcher.DefaultChromiumLauncher)
	if err != nil {
		test.Errorf("Expected no error, but got: %v", err)
	}
}

/*
// ADVANCED VERSION
func TestStartWithCustomSettings(test *testing.T) {
	// For windows we need a organisation name and project name
	var organisationName = "NewOrganisationName" // put in organisation name
	var projectName = "NewProjectName"           // put in project name

	var frontendPath = "./frontend" // this should be set to where frontend files is (frontend folder: html, css, javascript...)

	var chromeLauncher = launcher.ChromeLauncher{
		Location:                "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe",
		LocationCMD:             "C:\\\"Program Files\"\\Google\\Chrome\\Application\\chrome.exe",
		FrontendInstallLocation: os.Getenv("localappdata") + "\\Google\\Chrome\\InstalledApps\\" + organisationName + "\\" + projectName,
		Domain:                  "localhost",
		PortMin:                 11430,
		PreferredPort:           11451,
		PortMax:                 11500,
	}

	var chromiumLauncher = launcher.ChromiumLauncher{
		Location:      "/var/lib/snapd/desktop/applications/chromium_chromium.desktop", // TODO: check if better location or can be customised
		Domain:        "localhost",
		PortMin:       11430,
		PreferredPort: 11451,
		PortMax:       11500,
	}

	var err error

	err = launcher.Start(frontendPath, chromeLauncher, chromiumLauncher)
	if err != nil {
		test.Errorf("Expected no error, but got: %v", err)
	}
}
*/
