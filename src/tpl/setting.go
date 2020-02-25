package tpl

import (
	"github.com/teambition/urbs-setting/src/schema"
)

// SettingRes ...
type SettingRes struct {
	ResponseType
	Result schema.Setting `json:"result"` // 空数组也保留
}

// SettingsRes ...
type SettingsRes struct {
	ResponseType
	Result []schema.Setting `json:"result"` // 空数组也保留
}
