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
	Healthz *Healthz
	User    *User
	Group   *Group
	Product *Product
	Label   *Label
	Module  *Module
	Setting *Setting
}

// NewModels ...
func NewModels(sql *service.SQL) *Models {
	return &Models{
		Healthz: &Healthz{DB: sql.DB},
		User:    &User{DB: sql.DB},
		Group:   &Group{DB: sql.DB},
		Product: &Product{DB: sql.DB},
		Label:   &Label{DB: sql.DB},
		Module:  &Module{DB: sql.DB},
		Setting: &Setting{DB: sql.DB},
	}
}
