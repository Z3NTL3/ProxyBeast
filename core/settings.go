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
	"fmt"
	"net/url"
	"os"
	"time"
)

type Settings struct {
	Settings ApplicationSettings
}

type ApplicationSettings struct {
	Store struct {
		Judge *url.URL `json:"judge"`
		Timeout time.Duration `json:"timeout"`
		AllTime struct { 
			Scans uint64 `json:"scans"`
			Proxies uint64 `json:"proxies"`
		} `json:"all_time"`

		Pool struct {
			Workers struct {
				Size uint32 `json:"size"`
			} `json:"workers"`
		} `json:"pool"`
	}
	
}

var AppSettings *ApplicationSettings = &ApplicationSettings{}

func (app *ApplicationSettings) Init() (err error) {
	uri, err := url.Parse(DefaultJudge)
	if err != nil { return }
	{
		app.Store.Judge = uri
		app.Store.Pool.Workers.Size = DefaultPoolSize
		app.Store.Timeout = DefaultTimeout
	}

	app, err = app.GetStore()
	return 
}

func (app *ApplicationSettings) GetStore() (settings *ApplicationSettings, err error) {
	cwd, err := os.Getwd()
	if err != nil {return}

	fmt.Println(cwd)

	return
}