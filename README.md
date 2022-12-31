# go-local-web-gui (Local Go Chrome Framework)

go-local-web (GOLW) is a simple framework made for developing localhosted software that can reuse chrome/chromium or embed chromium (in future releases). Both available in deployment for the applications.

This framework uses Chrome (Windows) or Chromium (Linux) as frontend by opening them with cmd/terminal and hosting a localhost webserver, while opening chrome/chromium with --app and --user-directory arguments. The frontend can be changed by the user in runtime, while the backend needs to be compiled/build. The API can be decoupled in future versions, so every part of the application is changeable - Sustainable development. Frontends is easy to change. Alternatives to this is embedding a chromium or webview framework into the project, which will require more space. I chose to depend on Chrome/Chromium, as they are market leaders and html/css/javascript technology frontrunners.

More Description Coming Later <br>
This project is still under development and soon in alpha

## Requirements to developers
Go 1.19+
Chrome (Windows) or Chromium (Linux)

## Requirements for users
Chrome (Windows) or Chromium (Linux)

## How to use (with go get)
first run the following in CMD (with go installed)
<code>go get github.com/NineNineFive/go-local-web-gui/</code>
Example: how to add framework to main.go
<pre>
package main

import (
	"api"
	"fileserver"
	"launcher"
	"net/http"
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

func initHTTPHandlers() {
	http.HandleFunc("/", fileserver.ServeFileServer)
	http.HandleFunc("/api/", api.ServeAPIUseGZip)
}

func launchApp() {
	switch runtime.GOOS {
	case "windows":
		initHTTPHandlers()
		launcher.StartFrontendAndBackendWindows(frontendPath, chromeLauncher)
		return
	case "darwin": // "mac"
		panic("Darwin Not Supported Yet")
		return
	case "linux": // "linux"
		initHTTPHandlers()
		launcher.StartFrontendAndBackendLinux(frontendPath, chromiumLauncher)
		return
	default: // "freebsd", "openbsd", "netbsd"
		initHTTPHandlers()
		launcher.StartFrontendAndBackendLinux(frontendPath, chromiumLauncher)
		return
	}
}
</pre>

## How to run
<code>go run main.go</code>

## How to build
<code>go build -ldflags -H=windowsgui -o NewProjectName.exe</code>