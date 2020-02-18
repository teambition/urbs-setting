package api

import (
	"github.com/teambition/gear"
	"go.uber.org/dig"

	"github.com/teambition/urbs-setting/src/util"
)

func init() {
	util.DigProvide(NewRouter)
}

// digAPIs ..
type digAPIs struct {
	dig.In

	Healthz *Healthz
}

// NewRouter ...
func NewRouter(apis digAPIs) *gear.Router {
	router := gear.NewRouter()
	router.Get("/healthz", apis.Healthz.Get)

	return router
}
