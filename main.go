/*
/*
     ProxyBeast GUI

The ultimate proxy checker
       by @z3ntl3

    [proxy.pix4.dev]

License: GNU
Note: Please do give us a star on Github, if you like ProxyBeast

[App index]
*/

package main

import (
	"embed"
	"log"

	"ProxyBeast/core"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS


func main() {
	// App configuration
	if err := wails.Run(&options.App{
		Title:  "ProxyBeast", // Window title
		Width:  1024, // Window width
		Height: 768, // Window height
		AssetServer: &assetserver.Options{
			Assets: assets, // Webview2 assets
		},
		Bind: []any{ // Porting vanilla Go to JS
			core.APP, // core app instance
		},
		DisableResize: true, // Do not allow shrink/expanding window
		OnDomReady: core.APP.DomReady,
		OnStartup: core.APP.Startup, 
		Windows: &windows.Options{
			Theme: windows.Dark,
		},
		Mac: &mac.Options{
			Appearance: mac.NSAppearanceNameDarkAqua,
		},	
	}); err != nil {
		log.Fatalf("[Error]: %s", err)
	}
}
