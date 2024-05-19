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
	"encoding/json"
	"net/url"
	"os"
	"path"
	"slices"
	"strings"
	"time"
)

type Settings struct {
	Settings ApplicationSettings
}

type ApplicationSettings struct {
	Store struct {
		Timeout string 			`json:"timeout"`
		Judge   *url.URL      	`json:"judge"`
		AllTime struct {
			Scans   uint64 `json:"scans"`
			Proxies uint64 `json:"proxies"`
		} `json:"all_time"`

		Pool struct {
			Workers struct {
				Size uint32 `json:"size"`
			} `json:"workers"`
		} `json:"pool"`
	} `json:"store"`
}

var AppSettings *ApplicationSettings = &ApplicationSettings{}
var pattern = []string{"ms", "s"}

func (app *ApplicationSettings) Init() (err error) {
	return app.Get()
}

func (app *ApplicationSettings) Get() (err error){
	cwd, err := os.Getwd()
	if err != nil {return}

	var raw []byte

	location := path.Join(cwd, DefaultConfigFile)
	_, err = os.Stat(location)
	
	// [core.DefaultConfigFile] does not exist
	if os.IsNotExist(err) {
		file, err := os.Create(location)
		if err != nil {
			return err
		}
		defer file.Close()
		
		if err = app.Defaults(); err != nil {
			return err
		}

		raw, err = json.MarshalIndent(app, "", "  ")
		if err != nil {
			return err
		}

		if _, err = file.Write(raw); err != nil {
			return err
		}

		raw = nil
	} else if err != nil {
		// unknown fatal error
		return 
	} else {
		// exists
		raw, err = os.ReadFile(location)
		if err != nil {
			return 
		}

		if err = json.Unmarshal(raw, app); err != nil {
			return 
		}
	}

	return app.Patch()
}

func (app *ApplicationSettings) Defaults() (err error) {
	uri, err := url.Parse(DefaultJudge)
	if err != nil {return}
	{
		app.Store.Judge = uri
		app.Store.Pool.Workers.Size = DefaultPoolSize
		app.Store.Timeout = "8s"
	}

	return
}

func (app *ApplicationSettings) Patch() error {
	var err error 
	if app.Store.Judge == nil {
		uri, err := url.Parse(DefaultJudge)
		if err != nil {
			return err
		}

		app.Store.Judge = uri
	}

	if app.Store.Pool.Workers.Size == 0 {
		app.Store.Pool.Workers.Size = DefaultPoolSize
	}

	if !slices.ContainsFunc(pattern, func(patt string) bool {
		return strings.Contains(app.Store.Timeout, patt)
	}) {
		app.Store.Timeout = "8s"
	}

	DefaultPoolSize = app.Store.Pool.Workers.Size
	DefaultJudge = app.Store.Judge.String()
	DefaultTimeout, err = time.ParseDuration(app.Store.Timeout)

	if err != nil {return err}
	return nil
}

func (app *ApplicationSettings) SetPoolSize(size uint32) {
	app.Store.Pool.Workers.Size = size
}

func (app *ApplicationSettings) SetTimeout(timeout string) error {
	ok := false
	for _, patt := range pattern {
		if strings.Contains(timeout,patt) {
			ok = true
			break
		}
	}

	if !ok { return ErrTimeoutString }

	app.Store.Timeout = timeout 
	return nil
}