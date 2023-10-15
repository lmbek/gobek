# Note: Archived
This project is archived as i will not be updating it anymore - I will use webview to do the same as i did with Chrome, this will add stability as i can bundle a special version. However this version will be based on Edge, which might not have the same features as chrome. The new repository for the webview will be called BekView
<br>Go to: 
<a href="https://github.com/lmbek/bekview">
https://github.com/lmbek/bekview
</a>
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
first run the following in CMD (with go installed) <br><br>
The Go package:

    go get github.com/lmbek/gobek
<br><br>
Example: how to add framework to main.go<br>
	
 	//go:generate ./bin/windows/go-packager/GoPackager.exe
	// go generate
	// go build -ldflags -H=windowsgui -o NewProjectName.exe
	
	// ** Generating and Building **
	// The above is the generate commands to add icon and manifest to binary,
	// and the build command
	// NOTE: GoPackager.exe is a modified version of goversioninfo (https://github.com/josephspurrier/goversioninfo).
	// You can compile your own GoPackager.exe at ./bin/windows/go-packager by go building the main.go file
	
	// ** Reason Of Framework Explained **
	// This application is meant to find a port and launch a frontend in chrome (windows) or chromium (linux).
	// Meanwhile, it will open a http localhost backend with an api.
	// The API package at ./backend/api/ can be replaced to fit another application,
	// if the framework is to be used for other applications.
	// It is important to note, that the system will require a frontend directory and a data directory,
	// when the application is tested and released, these will also need to be managed (other applications can be put in)
	// This framework is very open-source, since when released to users, the users can modify the frontend files.
	// However, the backend can be kept as binary, but sent to users if they later want to modify it,
	// for example users must be able to modify the application backend, and replace the binary with the modified version.
	
	package main
	
	import (
		"fmt"
		"github.com/lmbek/gobek"
		"os"
	)
	
	// For windows, we need an organisation name and project name
	var organisationName = "NewOrganisationName" // put in organisation name
	var projectName = "NewProjectName"           // put in project name
	var chromeLauncher = gobek.ChromeLauncher{
		Location:                os.Getenv("programfiles") + "\\Google\\Chrome\\Application\\chrome.exe",
		FrontendInstallLocation: os.Getenv("localappdata") + "\\Google\\Chrome\\InstalledApps\\" + organisationName + "\\" + projectName,
	}
	var chromiumLauncher = gobek.ChromiumLauncher{
		Location: "/var/lib/snapd/desktop/applications/chromium_chromium.desktop",
	}
	var frontendPath = "./frontend" // this should be set to where frontend files is (frontend folder: html, css, javascript...)
	func main() {
		err := gobek.StartDefault(frontendPath, chromeLauncher, chromiumLauncher)
		if err != nil {
			fmt.Println(err)
		}
	}
Please do note however that using the main.go from the gobek-example project is recommended
## How to test
	go test ./tests/...

## How to run (after you made your own main.go)
    go run main.go

## How to apply manifest and logo to executible
Use something like goversioninfo: https://github.com/josephspurrier/goversioninfo
(Be aware of resource.syso and how it interacts with things like gcc, it is a learning process in itself)

## How to build
	go build -ldflags -H=windowsgui -o NewProjectName.exe

## For each project only one instance
The project is made, so you can only have one instance of the same organisation name and project name open, so if you have multiple project you are developing, please change their project names.

## For advanced users (Databases)
I have made an example project for more advanced users, where i demonstrated use of sqlite together with gobek <br>
The project can be found at: https://github.com/lmbek/gobek-sqlite-example