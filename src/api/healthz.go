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
func (h *Healthz) Get(ctx *gear.Context) error {
	return ctx.OkJSON(map[string]interface{}{
		"sql_db": h.blls.Models.Healthz.DBStats(ctx),
	})
}
