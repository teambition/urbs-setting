package model

import (
	"github.com/teambition/urbs-setting/src/service"
	"github.com/teambition/urbs-setting/src/util"
)

func init() {
	util.DigProvide(NewModels)
}

// Models ...
type Models struct {
	Healthz *healthz
	User    *user
	Group   *group
	Product *product
	Label   *label
	Module  *module
	Setting *setting
}

// NewModels ...
func NewModels(sql *service.SQL) *Models {
	return &Models{
		Healthz: &healthz{DB: sql.DB},
		User:    &user{DB: sql.DB},
		Group:   &group{DB: sql.DB},
		Product: &product{DB: sql.DB},
		Label:   &label{DB: sql.DB},
		Module:  &module{DB: sql.DB},
		Setting: &setting{DB: sql.DB},
	}
}
