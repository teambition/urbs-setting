package model

import (
	"database/sql"
	"fmt"

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

func deleteUserAndGroupLabels(db *sql.DB, labelIDs []int64) {
	var err error
	if len(labelIDs) > 0 {
		if _, err = db.Exec("delete from `user_label` where `label_id` in (?)", labelIDs); err == nil {
			_, err = db.Exec("delete from `group_label` where `label_id` in (?)", labelIDs)
		}
	}
	if err != nil {
		logging.Err(fmt.Sprintf("deleteUserAndGroupLabels with label_id(%v) error: %v", labelIDs, err))
	}
}

func deleteUserAndGroupSettings(db *sql.DB, settingIDs []int64) {
	var err error
	if len(settingIDs) > 0 {
		if _, err = db.Exec("delete from `user_setting` where `setting_id` in (?)", settingIDs); err == nil {
			_, err = db.Exec("delete from `user_setting` where `setting_id` in (?)", settingIDs)
		}
	}
	if err != nil {
		logging.Err(fmt.Sprintf("deleteUserAndGroupSettings with setting_id(%v) error: %v", settingIDs, err))
	}
}
