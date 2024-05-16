/*

     ProxyBeast GUI

The ultimate proxy checker
       by @z3ntl3

    [proxy.pix4.dev]

License: GNU
Note: If you've liked ProxyBeast, please consider starring our Github repository.

[App core]
*/

package core

import (
	"context"
	"fmt"
	"os"
	"path"
	"time"

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
		Name   string
		Exec   func(optionalData ...interface{})
		Cancel func()
	}

	// Event listeners aggregation
	EventGroup []*EventListeners
)

var (
	// main app instance, global
	APP = New()
)

const (
	SaveFile  Operation = "dialog_save_file"
	InputFile Operation = "dialog_input_file"
)

func New() *App {
	return &App{}
}

func (a *App) GetCtx() context.Context {
	return a.ctx
}

func (a *App) Startup(ctx context.Context) {
	var err error 
	defer func(err_ *error) {
		// fatal
		if *err_ != nil {
			runtime.EventsEmit(a.ctx, Fire_FatalError, (*err_).Error())
		}
	}(&err)
	
	a.ctx = ctx
	runtime.WindowCenter(ctx)

	// Obtain current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return
	}

	// Alias CWD
	RootDir = cwd

	// If <cwd>/saves is not resolvable, then try to mkdir.
	if _, err := os.Stat(path.Join(cwd, "saves")); err != nil || os.IsNotExist(err) {
		if err = os.Mkdir(path.Join(cwd, "saves"), os.ModeDir); err != nil {
			return // fatal
		}
	}

	AppSettings.Init()

	MX.fd_pool = make(chan FD_Pool, 20)
	MX.worker_pool = make(chan Workers, AppSettings.Pool.Workers.Size)

	MX.Register(context.WithCancel(context.Background()))
}

// Triggers when all resources are loaded
func (a *App) DomReady(ctx context.Context) {
	a.ctx = ctx

	if _, err := os.Stat(
		path.Join(RootDir, "saves"),
	); err != nil || os.IsNotExist(err) || RootDir == "" {
		runtime.EventsEmit(a.ctx, Fire_ErrSvdirEvent)
		return
	}

	var events *EventGroup = &EventGroup{
		{
			Name: OnCancelScan,
			Exec: a.cancel_scan,
		},
		{
			Name: OnSettingsModifyTimeout,
			Exec: a.modify_default_timeout,
		},
		{
			Name: OnDialog,
			Exec: a.dialog_exec,
		}, {
			Name: OnStartScan,
			Exec: func(data ...interface{}) {
				defer func() {
					err := recover()
					if err != nil {
						runtime.EventsEmit(a.ctx, Fire_ErrEvent, err.(error).Error())
					}
				}()
				MX.StartScan(a.ctx, data[0].(string))
			},
		},
	}
	events.register_eventListeners(a.ctx)
}

func (a *App) dialog(opts runtime.OpenDialogOptions) (string, error) {
	return runtime.OpenFileDialog(a.ctx, opts)
}

// Registers listeners for events
func (g *EventGroup) register_eventListeners(ctx context.Context) {
	for _, event := range *g {
		event.Cancel = runtime.EventsOn(ctx, event.Name, event.Exec)
	}
}

func (a *App) dialog_exec(optionalData ...interface{}) {
	var err error
	defer func(err_ *error) {
		if *err_ != nil {
			fmt.Println((*err_).Error())
			runtime.EventsEmit(a.ctx, Fire_ErrEvent, (*err_).Error())
		}
	}(&err)

	props, ok := optionalData[0].(string)
	if !ok {
		err = ErrPropsInvalid
		return
	}

	opts := runtime.OpenDialogOptions{
		DefaultDirectory: path.Join(RootDir),
		Title:            "ProxyBeast - File dialog",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Select file (.txt)",
				Pattern:     "*.txt",
			},
		},
	}
	if props == SaveFile {
		opts.DefaultDirectory = path.Join(opts.DefaultDirectory, "saves")
	}

	loc, err := a.dialog(opts)
	if err != nil {
		return
	}

	f, err := OpenFileRDO(loc)
	if err != nil {
		return
	}

	FD[props].Close()
	FD[props] = f

	runtime.EventsEmit(a.ctx, props, path.Base(loc))
}

func (a *App) cancel_scan(...interface{}) {
	MX.Cancel()
}

func (a *App) modify_default_timeout(data ...interface{}) {
	timeout, ok := data[0].(string)
	if !ok {
		runtime.EventsEmit(APP.ctx, Fire_ErrEvent, ErrTimeoutString)
		return
	}

	time, err := time.ParseDuration(timeout)
	if err != nil {
		runtime.EventsEmit(APP.ctx, Fire_ErrEvent, err.Error())
		return
	}

	DefaultTimeout = time
}
