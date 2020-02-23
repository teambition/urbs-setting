package bll

import (
	"github.com/teambition/urbs-setting/src/model"
	"github.com/teambition/urbs-setting/src/util"
)

func init() {
	util.DigProvide(NewBlls)
}

// Blls ...
type Blls struct {
	User    *User
	Group   *Group
	Product *Product
	Label   *Label
	Module  *Module
	Setting *Setting
	Models  *model.Models
}

// NewBlls ...
func NewBlls(models *model.Models) *Blls {
	return &Blls{
		User:    &User{ms: models},
		Group:   &Group{ms: models},
		Product: &Product{ms: models},
		Label:   &Label{ms: models},
		Module:  &Module{ms: models},
		Setting: &Setting{ms: models},
		Models:  models,
	}
}
