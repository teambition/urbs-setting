package api

import (
	"github.com/teambition/gear"
	"github.com/teambition/urbs-setting/src/util"
)

func init() {
	util.DigProvide(NewHealthz)
}

// Healthz ..
type Healthz struct{}

// NewHealthz return an Healthz instance.
func NewHealthz() *Healthz {
	return &Healthz{}
}

// Get ..
func (h *Healthz) Get(ctx *gear.Context) error {
	// TODO
	return ctx.OkJSON(GetVersion())
}
