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

## How to run
<code>go run main.go</code>
