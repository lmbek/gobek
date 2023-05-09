# gobek (Local GUI / Go Chrome Framework)

gobek is a simple framework made for developing localhosted software that can reuse chrome/chromium or embed chromium (in future releases). Both available in deployment for the applications.

This framework uses Chrome (Windows) or Chromium (Linux) as frontend by opening them with cmd/terminal and hosting a localhost webserver, while opening chrome/chromium with --app and --user-directory arguments. The frontend can be changed by the user in runtime, while the backend needs to be compiled/build. The API can be decoupled in future versions, so every part of the application is changeable - Sustainable development. Frontends is easy to change. Alternatives to this is embedding a chromium or webview framework into the project, which will require more space. I chose to depend on Chrome/Chromium, as they are market leaders and html/css/javascript technology frontrunners.

Feel free to use this piece of software, I will be happy to assist you

I am currently working on this project, it will be updated and maintained. I consider it production ready. 

This project is used by Beksoft ApS for projects such as:
* BekCMS
* PingPong Game made in Three.js
* Several local webbased software projects

Write to me at lars@beksoft.dk if you want to have your project listed

## Contributors
Lars Morten Bek (https://github.com/lmbek)

## Requirements to developers
Go 1.20+
Chrome (Windows) or Chromium (Linux)
macOS/DARWIN NOT SUPPORTED YET

## Requirements for users
Chrome (Windows) or Chromium (Linux)
macOS/DARWIN NOT SUPPORTED YET

## How to use (download example project)
The best way to start using the project is to download the example project at:
https://github.com/lmbek/gobek-example

This example project uses this package and combines it with a local api
Then the Go api is being developed and customized by you together with the frontend (JavaScript, HTML, CSS)

## How to use (with go get)
first run the following in CMD (with go installed)
<code>go get github.com/lmbek/gobek</code>
Example: how to add framework to main.go
<pre>
package main

import (
	"fmt"
	"github.com/lmbek/gobek/launcher"
	"os"
)

// For windows we need a organisation name and project name
var organisationName = "NewOrganisationName" // put in organisation name
var projectName = "NewProjectName"           // put in project name

var frontendPath = "./frontend" // this should be set to where frontend files is (frontend folder: html, css, javascript...)

var chromeLauncher = launcher.ChromeLauncher{
	Location:                os.Getenv("programfiles") + "\\Google\\Chrome\\Application\\chrome.exe",
	FrontendInstallLocation: os.Getenv("localappdata") + "\\Google\\Chrome\\InstalledApps\\" + "DefaultOrganisationName" + "\\" + "DefaultProjectName",
}

var chromiumLauncher = launcher.DefaultChromiumLauncher // default chrome or chromium launcher settings can be used like this

/*
	// Otherwise they can also be customized like this

	var chromiumLauncher = launcher.ChromiumLauncher{
		Location:      "/var/lib/snapd/desktop/applications/chromium_chromium.desktop", // TODO: check if better location or can be customised
		Domain:        "localhost",
	}
*/

func main() {
    // add your own API
	//api.Init()
	/*
		var once sync.Once
		once.Do(func() {
			http.HandleFunc("/", fileserver.ServeFileServer)
			http.HandleFunc("/api/", api.ServeAPIUseGZip)
		})
		err := launcher.Start(frontendPath, chromeLauncher, chromiumLauncher) // serves "/" as fileserver.ServeFileServer. If you want to manage "/", then use launcher.StartCustom() instead
		if err != nil {
			fmt.Println(err)
		}
	*/
	err := launcher.StartDefault(frontendPath, chromeLauncher, chromiumLauncher)
	if err != nil {
		fmt.Println(err)
	}
}

</pre>

## How to test
<code>go test ./tests/...</code>

## How to run
make your own main.go by following https://github.com/lmbek/gobek-example

## How to apply manifest and logo to executible
Use something like goversioninfo: https://github.com/josephspurrier/goversioninfo

## How to build
<code>go build -ldflags -H=windowsgui -o NewProjectName.exe</code>

## How to make setup file and update functionality
Coming later

## Further plans
Huge changes coming to version 0.7.0 that will change the whole project structure, this will break all versions before 0.7.0 when upgrading