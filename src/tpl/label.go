package tpl

import (
	"github.com/teambition/gear"
	"github.com/teambition/urbs-setting/src/schema"
)

// LabelsURL ...
type LabelsURL struct {
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
	return nil
}

// LabelBody ...
type LabelBody struct {
	Name string `json:"Name"`
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
	HID      string `json:"hid"`
	Product  string `json:"product"`
	Name     string `json:"name"`
	Desc     string `json:"desc"`
	Channels string `json:"channels"`
	Clients  string `json:"clients"`
}

// LabelsInfoRes ...
type LabelsInfoRes struct {
	ResponseType
	Result []LabelInfo `json:"result"` // 空数组也保留
}

// CacheLabelsInfoRes ...
type CacheLabelsInfoRes struct {
	ResponseType
	Timestamp int64                   `json:"timestamp"` // labels 数组生成时间
	Result    []schema.UserCacheLabel `json:"result"`    // 空数组也保留
}

// LabelRes ...
type LabelRes struct {
	ResponseType
	Result schema.Label `json:"result"` // 空数组也保留
}

// LabelsRes ...
type LabelsRes struct {
	ResponseType
	Result []schema.Label `json:"result"` // 空数组也保留
}
