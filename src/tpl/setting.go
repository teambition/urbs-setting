package tpl

import (
	"strings"
	"time"

	"github.com/teambition/gear"
	"github.com/teambition/urbs-setting/src/conf"
	"github.com/teambition/urbs-setting/src/schema"
	"github.com/teambition/urbs-setting/src/service"
)

// SettingBody ...
type SettingBody struct {
	Name     string    `json:"name"`
	Desc     string    `json:"desc"`
	Channels *[]string `json:"channels"`
	Clients  *[]string `json:"clients"`
	Values   *[]string `json:"values"`
}

// Validate 实现 gear.BodyTemplate。
func (t *SettingBody) Validate() error {
	if !validNameReg.MatchString(t.Name) {
		return gear.ErrBadRequest.WithMsgf("invalid name: %s", t.Name)
	}
	if len(t.Desc) > 1022 {
		return gear.ErrBadRequest.WithMsgf("desc too long: %d (<= 1022)", len(t.Desc))
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
			if !validValueReg.MatchString(value) {
				return gear.ErrBadRequest.WithMsgf("invalid value: %s", value)
			}
		}
	}
	return nil
}

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
			if !validValueReg.MatchString(value) {
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
	Release   int64      `json:"release"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	OfflineAt *time.Time `json:"offlineAt"`
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
		Release:   setting.Release,
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
	ID         int64     `json:"-" db:"id"`
	HID        string    `json:"hid"`
	Product    string    `json:"product" db:"product"`
	Module     string    `json:"module" db:"module"`
	Name       string    `json:"name" db:"name"`
	Desc       string    `json:"desc" db:"description"`
	Value      string    `json:"value" db:"value"`
	LastValue  string    `json:"lastValue" db:"last_value"`
	Release    int64     `json:"release" db:"rls"`
	AssignedAt time.Time `json:"assignedAt" db:"assigned_at"`
}

// MySettingsRes ...
type MySettingsRes struct {
	SuccessResponseType
	Result []MySetting `json:"result"` // 空数组也保留
}

// MySettingsQueryURL ...
type MySettingsQueryURL struct {
	Pagination
	UID     string `json:"uid" param:"uid"`
	Product string `json:"product" query:"product"`
	Module  string `json:"module" query:"module"`
	Setting string `json:"setting" query:"setting"`
	Channel string `json:"channel" query:"channel"`
	Client  string `json:"client" query:"client"`
}

// Validate 实现 gear.BodyTemplate。
func (t *MySettingsQueryURL) Validate() error {
	if !validIDReg.MatchString(t.UID) {
		return gear.ErrBadRequest.WithMsgf("invalid user: %s", t.UID)
	}
	if t.Product != "" && !validNameReg.MatchString(t.Product) {
		return gear.ErrBadRequest.WithMsgf("invalid product name: %s", t.Product)
	}
	if t.Module != "" && !validNameReg.MatchString(t.Module) {
		return gear.ErrBadRequest.WithMsgf("invalid module name: %s", t.Module)
	}
	if t.Module != "" && t.Product == "" {
		return gear.ErrBadRequest.WithMsgf("product required for module: %s", t.Module)
	}
	if t.Setting != "" && !validNameReg.MatchString(t.Setting) {
		return gear.ErrBadRequest.WithMsgf("invalid setting name: %s", t.Setting)
	}
	if t.Setting != "" && t.Module == "" {
		return gear.ErrBadRequest.WithMsgf("module required for setting: %s", t.Setting)
	}

	if err := t.Pagination.Validate(); err != nil {
		return err
	}

	if t.Channel != "" && !StringSliceHas(conf.Config.Channels, t.Channel) {
		return gear.ErrBadRequest.WithMsgf("invalid channel: %s", t.Channel)
	}

	if t.Client != "" && !StringSliceHas(conf.Config.Clients, t.Client) {
		return gear.ErrBadRequest.WithMsgf("invalid client: %s", t.Client)
	}
	return nil
}

// SettingReleaseInfo ...
type SettingReleaseInfo struct {
	Release int64    `json:"release"`
	Users   []string `json:"users"`
	Groups  []string `json:"groups"`
	Value   string   `json:"value"`
}

// SettingReleaseInfoRes ...
type SettingReleaseInfoRes struct {
	SuccessResponseType
	Result SettingReleaseInfo `json:"result"` // 空数组也保留
}
