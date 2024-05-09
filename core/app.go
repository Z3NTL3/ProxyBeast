/*
/*

     ProxyBeast GUI

The ultimate proxy checker
       by @z3ntl3

    [proxy.pix4.dev]

License: GNU
Note: Please do give us a star on Github, if you like ProxyBeast

[App core]
*/

package core

import (
	"ProxyBeast/globals"
	"context"
	"fmt"
	"os"
	"path"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type (
	// App instance
	App struct {
		ctx context.Context
	}

	// Describes some operation
	Operation = string

	// Event listeners
	EventListeners struct {
		Name string
		Exec func(optionalData ...interface{})
		Cancel func()
	}

	// Event listeners aggregation
	EventGroup []*EventListeners
)

var APP = New()

const (
	SaveFile Operation = "dialog_file_save"
	OpenFile Operation = "dialog_open_file"
)

func New() *App {
	return &App{}
}

func (a *App) GetCtx() context.Context {
	return a.ctx
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx

	// Obtain current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return
	}
	
	// Alias CWD
	globals.RootDir = cwd

	// If <cwd>/saves is not resolvable, then mkdir.
	if _, err := os.Stat(path.Join(cwd,"saves")); err != nil || os.IsNotExist(err) {
		if err = os.Mkdir(path.Join(cwd,"saves"), os.ModeDir); err != nil {
			println(err)
		}
	} 

}

// Triggers when all resources are loaded
func (a *App) DomReady(ctx context.Context) {
	a.ctx = ctx

	var events *EventGroup = &EventGroup{
		{
			Name: "dialog",
			Exec: func(optionalData ...interface{}) {
				props, ok := optionalData[0].(string)
				if !ok {
					runtime.LogError(a.ctx, ErrPropsInvalid.Error())
					return
				}

				opts := runtime.OpenDialogOptions{
					DefaultDirectory: path.Join(globals.RootDir),
					Title: "ProxyBeast - File dialog",
					Filters: []runtime.FileFilter{
						{
							DisplayName: "Select file (.txt)",
							Pattern: "*.txt",
						},
					},
				}
				if props == SaveFile {
					opts.DefaultDirectory = path.Join(opts.DefaultDirectory, "saves")
				}
				
				loc, err := a.dialog(opts)
				if err != nil {
					runtime.EventsEmit(a.ctx, "error", err)
				}
				fmt.Println(loc)
			},
		},
	}
	events.register_eventListeners(a.ctx)
	if _, err := os.Stat(path.Join(globals.RootDir, "saves")); err != nil || os.IsNotExist(err) {
		runtime.EventsEmit(a.ctx, "svdir_failure")
	}
}

func(a *App) dialog(opts runtime.OpenDialogOptions) (string, error){
	return runtime.OpenFileDialog(a.ctx, opts)
}

// Registers listeners for events
func (g *EventGroup) register_eventListeners(ctx context.Context) {
	for _, event := range *g {
		event.Cancel = runtime.EventsOn(ctx, event.Name, event.Exec)
	}
}