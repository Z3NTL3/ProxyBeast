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
	"os"
	path "path/filepath"
	"sync/atomic"

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

func (a *App) OnExit(ctx context.Context) {
	AppSettings.Sync() // try to update appstore 
}

func (a *App) Startup(ctx context.Context) {
	var err error 
	defer func() {
		// fatal
		if err != nil {
			runtime.EventsEmit(a.ctx, Fire_FatalError, err.Error())
		}
	}()
	
	a.ctx = ctx
	runtime.WindowCenter(ctx)

	cwd, err := os.Getwd()
	if err != nil {
		return
	}

	// Alias CWD
	RootDir = cwd

	// If <cwd>/saves is not resolvable, then try to mkdir.
	if _, err = os.Stat(path.Join(cwd, "saves")); err != nil || os.IsNotExist(err) {
		if err = os.Mkdir(path.Join(cwd, "saves"), os.ModeDir); err != nil {
			return // fatal
		}
	}

	AppSettings.Init()

	MX.fd_pool = make(chan FD_Pool, 20)
	MX.worker_pool = make(chan Workers, AppSettings.Store.Pool.Workers.Size)

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
				MX.StartScan(a.ctx, data[0].(string)) // may panic
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
	defer func() {
		if err != nil {
			runtime.EventsEmit(a.ctx, Fire_ErrEvent, err.Error())
		}
	}()

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
	runtime.EventsEmit(a.ctx, Fire_MsgEvent, "Cancel called, killing goroutines please wait")
	MX.Cancel()
}

func (a *App) modify_default_timeout(data ...interface{}) {
	var err error
	
	defer func(){
		if err != nil {
			runtime.EventsEmit(APP.ctx, Fire_ErrEvent, err.Error())
		}
	}()

	timeout, ok := data[0].(string)
	if !ok {
		err = ErrTimeoutString
		return
	}

	if err = AppSettings.SetTimeout(timeout); err != nil {
		return
	}

	if err = AppSettings.Patch(); err != nil {
		return
	}
}

func (a *App) Refresh() string {
	get := AppSettings.Get()
	if get != nil { return get.Error()}
	return ""
}

func (a *App) GetJudge() string {
	return AppSettings.Store.Judge.String()
}

func (a *App) GetPoolSize() uint32 {
	return AppSettings.Store.Pool.Workers.Size
}

func (a *App) GetTimeoutString() string {
	return AppSettings.Store.Timeout
}

func (a *App) GetAllTimeChecks() uint64 {
	return atomic.LoadUint64(&AppSettings.Store.AllTime.Proxies)
}

func (a *App) GetAllTimeScans() uint64 {
	return AppSettings.Store.AllTime.Scans
}