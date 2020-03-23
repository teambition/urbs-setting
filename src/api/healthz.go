package api

import (
	"github.com/teambition/gear"
	"github.com/teambition/urbs-setting/src/bll"
)

// Healthz ..
type Healthz struct {
	blls *bll.Blls
}

// Get ..
func (a *Healthz) Get(ctx *gear.Context) error {
	stats := a.blls.Models.Healthz.DBStats(ctx)
	return ctx.OkJSON(map[string]interface{}{
		"db_connect": stats.OpenConnections > 0,
	})
}
