//go:generate ./generate/GoPackager.exe
// go generate
// go build -ldflags -H=windowsgui -o NewProjectName.exe

// ** Generating and Building **
// The above is the generate commands to add icon and manifest to binary,
// and the build command
// NOTE: GoPackager.exe is a modified version of goversioninfo.
// You can find goversioninfo at https://github.com/josephspurrier/goversioninfo

// ** Reason Of Framework Explained **
// This application is meant to find a port and launch a frontend in chrome (windows) or chromium (linux).
// Meanwhile, it will open a http localhost backend with an api.
// The API package at /development/backend/api/ can be replaced to fit another application,
// if the framework is to be used for other applications.
// It is important to note, that the system will require a frontend directory and a data directory,
// when the application is tested and released, these will also need to be managed (other applications can be put in)
// This framework is very open-source, since when released to users, the users can modify the frontend files.
// However the backend can be kept as binary, but sent to users if they later want to modify it,
// for example users must be able to modify the application backend, and replace the binary with the modified version.

package main

import (
	"GoLocalChromeFramework/backend/launcher"
	"os"
	"runtime"
)

var projectName = "NewProjectName"
var organisationName = "NewCompanyName"

// var address = "localhost:10995"
var frontendPath = "./frontend"

var chromeLauncher = launcher.ChromeLauncher{
	Location:                "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe",
	LocationCMD:             "C:\\\"Program Files\"\\Google\\Chrome\\Application\\chrome.exe",
	FrontendInstallLocation: os.Getenv("localappdata") + "\\Google\\Chrome\\InstalledApps\\" + organisationName + "\\" + projectName,
	Domain:                  "localhost",
	PortMin:                 19430,
	PreferredPort:           19451,
	PortMax:                 19500,
}

var chromiumLauncher = launcher.ChromiumLauncher{
	Location:      "/var/lib/snapd/desktop/applications/chromium_chromium.desktop", // TODO: check if better location or can be customised
	Domain:        "localhost",
	PortMin:       19430,
	PreferredPort: 19451,
	PortMax:       19500,
}

func main() {
	launchApp()
}

func launchApp() {
	switch runtime.GOOS {
	case "windows":
		launcher.StartFrontendAndBackendWindows(frontendPath, chromeLauncher)
		return
	case "darwin": // "mac"
		panic("Darwin Not Supported Yet")
		return
	case "linux": // "linux"
		launcher.StartFrontendAndBackendLinux(frontendPath, chromiumLauncher)
		return
	default: // "freebsd", "openbsd", "netbsd"
		launcher.StartFrontendAndBackendLinux(frontendPath, chromiumLauncher)
		return
	}
}
