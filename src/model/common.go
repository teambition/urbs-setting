package model

import (
	"github.com/jinzhu/gorm"
	"github.com/teambition/urbs-setting/src/logging"
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

func deleteUserAndGroupLabels(db *gorm.DB, labelIDs []int64) {
	var err error
	if len(labelIDs) > 0 {
		if err = db.Exec("delete from `user_label` where `label_id` in ( ? )", labelIDs).Error; err == nil {
			err = db.Exec("delete from `group_label` where `label_id` in ( ? )", labelIDs).Error
		}
	}
	if err != nil {
		logging.Errf("deleteUserAndGroupLabels with label_id(%v) error: %v", labelIDs, err)
	}
}

func deleteUserAndGroupSettings(db *gorm.DB, settingIDs []int64) {
	var err error
	if len(settingIDs) > 0 {
		if err = db.Exec("delete from `user_setting` where `setting_id` in ( ? )", settingIDs).Error; err == nil {
			err = db.Exec("delete from `group_setting` where `setting_id` in ( ? )", settingIDs).Error
		}
	}
	if err != nil {
		logging.Errf("deleteUserAndGroupSettings with setting_id(%v) error: %v", settingIDs, err)
	}
}
