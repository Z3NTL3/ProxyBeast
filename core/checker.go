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
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//todo
func (c *Controller) StartChecking(ctx context.Context, proto string) {
	var err error
	defer func(err_ *error) {
		if *err_ != nil {
			runtime.EventsEmit(APP.ctx, Fire_ErrEvent, (*err_).Error())
		}
	}(&err)

	selection := Proxy(proto)

	if selection != "" && !(&selection).IsValid(Scheme) {
		runtime.EventsEmit(APP.ctx, Fire_ProtoUnknown)
		return
	}

	scheme := Scheme
	if proto != "" { scheme = NoScheme }

	err = FD.Validate(scheme)
	if err != nil {return}
}