package tpl

import (
	"time"

	"github.com/teambition/urbs-setting/src/schema"
	"github.com/teambition/urbs-setting/src/service"
)

// SettingInfo ...
type SettingInfo struct {
	ID        int64      `json:"-"`
	HID       string     `json:"hid"`
	Product   string     `json:"product"`
	Module    string     `json:"module"`
	Name      string     `json:"name"`
	Desc      string     `json:"desc"`
	Channels  []string   `json:"channels"`
	Clients   []string   `json:"clients"`
	Values    []string   `json:"values"`
	Status    int64      `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	OfflineAt *time.Time `json:"offline_at"`
}

// SettingInfoFrom create a SettingInfo from schema.Setting
func SettingInfoFrom(setting schema.Setting, product, module string) SettingInfo {
	return SettingInfo{
		ID:        setting.ID,
		HID:       service.IDToHID(setting.ID, "setting"),
		Product:   product,
		Module:    module,
		Name:      setting.Name,
		Desc:      setting.Desc,
		Channels:  StringToSlice(setting.Channels),
		Clients:   StringToSlice(setting.Clients),
		Values:    StringToSlice(setting.Values),
		Status:    setting.Status,
		CreatedAt: setting.CreatedAt,
		UpdatedAt: setting.UpdatedAt,
		OfflineAt: setting.OfflineAt,
	}
}

// SettingsInfoFrom create a slice of SettingInfo from a slice of schema.Setting
func SettingsInfoFrom(settings []schema.Setting, product, module string) []SettingInfo {
	res := make([]SettingInfo, len(settings))
	for i, l := range settings {
		res[i] = SettingInfoFrom(l, product, module)
	}
	return res
}

// SettingsInfoRes ...
type SettingsInfoRes struct {
	SuccessResponseType
	Result []SettingInfo `json:"result"` // 空数组也保留
}

// SettingInfoRes ...
type SettingInfoRes struct {
	SuccessResponseType
	Result SettingInfo `json:"result"` // 空数组也保留
}

// MySetting ...
type MySetting struct {
	ID        int64     `json:"-"`
	HID       string    `json:"hid"`
	Module    string    `json:"module"`
	Name      string    `json:"name"`
	Value     string    `json:"value"`
	LastValue string    `json:"last_value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// MySettingsRes ...
type MySettingsRes struct {
	SuccessResponseType
	Result []MySetting `json:"result"` // 空数组也保留
}
