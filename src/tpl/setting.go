package tpl

import (
	"strings"
	"time"

	"github.com/teambition/gear"
	"github.com/teambition/urbs-setting/src/conf"
	"github.com/teambition/urbs-setting/src/schema"
	"github.com/teambition/urbs-setting/src/service"
)

// SettingUpdateBody ...
type SettingUpdateBody struct {
	Desc     *string   `json:"desc"`
	Channels *[]string `json:"channels"`
	Clients  *[]string `json:"clients"`
	Values   *[]string `json:"values"`
}

// Validate 实现 gear.BodyTemplate。
func (t *SettingUpdateBody) Validate() error {
	if t.Desc == nil && t.Channels == nil && t.Clients == nil && t.Values == nil {
		return gear.ErrBadRequest.WithMsgf("desc or channels or clients or values required")
	}
	if t.Desc != nil && len(*t.Desc) > 1022 {
		return gear.ErrBadRequest.WithMsgf("desc too long: %d", len(*t.Desc))
	}
	if t.Channels != nil {
		if len(*t.Channels) > 5 {
			return gear.ErrBadRequest.WithMsgf("too many channels: %d", len(*t.Channels))
		}
		if !SortStringsAndCheck(*t.Channels) {
			return gear.ErrBadRequest.WithMsgf("invalid channels: %v", *t.Channels)
		}
		for _, channel := range *t.Channels {
			if !StringSliceHas(conf.Config.Channels, channel) {
				return gear.ErrBadRequest.WithMsgf("invalid channel: %s", channel)
			}
		}
	}
	if t.Clients != nil {
		if len(*t.Clients) > 10 {
			return gear.ErrBadRequest.WithMsgf("too many clients: %d", len(*t.Clients))
		}
		if !SortStringsAndCheck(*t.Clients) {
			return gear.ErrBadRequest.WithMsgf("invalid clients: %v", *t.Clients)
		}
		for _, client := range *t.Clients {
			if !StringSliceHas(conf.Config.Clients, client) {
				return gear.ErrBadRequest.WithMsgf("invalid client: %s", client)
			}
		}
	}
	if t.Values != nil {
		if len(*t.Values) > 10 {
			return gear.ErrBadRequest.WithMsgf("too many values: %d", len(*t.Clients))
		}
		if !SortStringsAndCheck(*t.Values) {
			return gear.ErrBadRequest.WithMsgf("invalid values: %v", *t.Values)
		}
		for _, value := range *t.Values {
			if !validValueReg1.MatchString(value) && !validValueReg2.MatchString(value) {
				return gear.ErrBadRequest.WithMsgf("invalid value: %s", value)
			}
		}
	}
	return nil
}

// ToMap ...
func (t *SettingUpdateBody) ToMap() map[string]interface{} {
	changed := make(map[string]interface{})
	if t.Desc != nil {
		changed["description"] = *t.Desc
	}
	if t.Channels != nil {
		changed["channels"] = strings.Join(*t.Channels, ",")
	}
	if t.Clients != nil {
		changed["clients"] = strings.Join(*t.Clients, ",")
	}
	if t.Values != nil {
		changed["vals"] = strings.Join(*t.Values, ",")
	}
	return changed
}

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
