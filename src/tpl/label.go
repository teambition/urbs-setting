package tpl

import (
	"strings"
	"time"

	"github.com/teambition/gear"
	"github.com/teambition/urbs-setting/src/conf"
	"github.com/teambition/urbs-setting/src/schema"
	"github.com/teambition/urbs-setting/src/service"
)

// LabelBody ...
type LabelBody struct {
	Name     string    `json:"name"`
	Desc     string    `json:"desc"`
	Channels *[]string `json:"channels"`
	Clients  *[]string `json:"clients"`
}

// Validate 实现 gear.BodyTemplate。
func (t *LabelBody) Validate() error {
	if !validLabelReg.MatchString(t.Name) {
		return gear.ErrBadRequest.WithMsgf("invalid label: %s", t.Name)
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
	return nil
}

// LabelUpdateBody ...
type LabelUpdateBody struct {
	Desc     *string   `json:"desc"`
	Channels *[]string `json:"channels"`
	Clients  *[]string `json:"clients"`
}

// Validate 实现 gear.BodyTemplate。
func (t *LabelUpdateBody) Validate() error {
	if t.Desc == nil && t.Channels == nil && t.Clients == nil {
		return gear.ErrBadRequest.WithMsgf("desc or channels or clients required")
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
	return nil
}

// ToMap ...
func (t *LabelUpdateBody) ToMap() map[string]interface{} {
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
	return changed
}

// LabelInfo ...
type LabelInfo struct {
	ID        int64      `json:"-"`
	HID       string     `json:"hid"`
	Product   string     `json:"product"`
	Name      string     `json:"name"`
	Desc      string     `json:"desc"`
	Channels  []string   `json:"channels"`
	Clients   []string   `json:"clients"`
	Status    int64      `json:"status"`
	Release   int64      `json:"release"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	OfflineAt *time.Time `json:"offlineAt"`
}

// LabelInfoFrom create a LabelInfo from schema.Label
func LabelInfoFrom(label schema.Label, product string) LabelInfo {
	return LabelInfo{
		ID:        label.ID,
		HID:       service.IDToHID(label.ID, "label"),
		Product:   product,
		Name:      label.Name,
		Desc:      label.Desc,
		Channels:  StringToSlice(label.Channels),
		Clients:   StringToSlice(label.Clients),
		Status:    label.Status,
		Release:   label.Release,
		CreatedAt: label.CreatedAt,
		UpdatedAt: label.UpdatedAt,
		OfflineAt: label.OfflineAt,
	}
}

// LabelsInfoFrom create a slice of LabelInfo from a slice of schema.Label
func LabelsInfoFrom(labels []schema.Label, product string) []LabelInfo {
	res := make([]LabelInfo, len(labels))
	for i, l := range labels {
		res[i] = LabelInfoFrom(l, product)
	}
	return res
}

// LabelsInfoRes ...
type LabelsInfoRes struct {
	SuccessResponseType
	Result []LabelInfo `json:"result"` // 空数组也保留
}

// LabelInfoRes ...
type LabelInfoRes struct {
	SuccessResponseType
	Result LabelInfo `json:"result"` // 空数组也保留
}

// MyLabel ...
type MyLabel struct {
	ID         int64     `json:"-" db:"id"`
	HID        string    `json:"hid"`
	Product    string    `json:"product" db:"product"`
	Name       string    `json:"name" db:"name"`
	Desc       string    `json:"desc" db:"description"`
	Release    int64     `json:"release" db:"rls"`
	AssignedAt time.Time `json:"assignedAt" db:"assigned_at"`
}

// MyLabelsRes ...
type MyLabelsRes struct {
	SuccessResponseType
	Result []MyLabel `json:"result"` // 空数组也保留
}

// CacheLabelsInfoRes ...
type CacheLabelsInfoRes struct {
	SuccessResponseType
	Timestamp int64                   `json:"timestamp"` // labels 数组生成时间
	Result    []schema.UserCacheLabel `json:"result"`    // 空数组也保留
}

// LabelReleaseInfo ...
type LabelReleaseInfo struct {
	Release int64    `json:"release"`
	Users   []string `json:"users"`
	Groups  []string `json:"groups"`
}

// LabelReleaseInfoRes ...
type LabelReleaseInfoRes struct {
	SuccessResponseType
	Result LabelReleaseInfo `json:"result"` // 空数组也保留
}
