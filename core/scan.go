package core

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// todo
func (c *Controller) Scan(ctx context.Context, proto string) {
	selection := Proxy(proto)

	if selection != "none" && !(&selection).IsValid() {
		runtime.EventsEmit(APP.ctx, Fire_ProtoUnknown)
		return
	}
}