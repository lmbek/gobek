module main

require (
	github.com/NineNineFive/go-local-web-gui/fileserver v0.0.0-20221231234257-b8c44930244f
	github.com/NineNineFive/go-local-web-gui/launcher v0.0.0-20221231234257-b8c44930244f
)

replace github.com/NineNineFive/go-local-web-gui => ./../..

go 1.19
