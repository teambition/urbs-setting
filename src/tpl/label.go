package tpl

import (
	"time"

	"github.com/teambition/gear"
	"github.com/teambition/urbs-setting/src/schema"
	"github.com/teambition/urbs-setting/src/service"
)

// LabelsURL ...
type LabelsURL struct {
	Pagination
	UID     string `json:"uid" param:"uid"`
	Product string `json:"product" query:"product"`
}

// Validate 实现 gear.BodyTemplate。
func (t *LabelsURL) Validate() error {
	if !validIDNameReg.MatchString(t.UID) {
		return gear.ErrBadRequest.WithMsgf("invalid user: %s", t.UID)
	}
	if t.Product != "" && !validIDNameReg.MatchString(t.Product) {
		return gear.ErrBadRequest.WithMsgf("invalid product name: %s", t.Product)
	}
	if err := t.Pagination.Validate(); err != nil {
		return err
	}
	return nil
}

// LabelBody ...
type LabelBody struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}

// Validate 实现 gear.BodyTemplate。
func (t *LabelBody) Validate() error {
	if !validLabelReg.MatchString(t.Name) {
		return gear.ErrBadRequest.WithMsgf("invalid label: %s", t.Name)
	}
	if len(t.Desc) > 1022 {
		return gear.ErrBadRequest.WithMsgf("desc too long: %d (<= 1022)", len(t.Desc))
	}
	return nil
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
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	OfflineAt *time.Time `json:"offline_at"`
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
		CreatedAt: label.CreatedAt,
		UpdatedAt: label.UpdatedAt,
		OfflineAt: label.OfflineAt,
	}
}

// LabelInfosFrom create a slice of LabelInfo from a slice of schema.Label
func LabelInfosFrom(labels []schema.Label, product string) []LabelInfo {
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

// CacheLabelsInfoRes ...
type CacheLabelsInfoRes struct {
	SuccessResponseType
	Timestamp int64                   `json:"timestamp"` // labels 数组生成时间
	Result    []schema.UserCacheLabel `json:"result"`    // 空数组也保留
}
