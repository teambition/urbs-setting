package tpl

import (
	"github.com/teambition/gear"
	"github.com/teambition/urbs-setting/src/schema"
)

// LabelsURL ...
type LabelsURL struct {
	UID     string `json:"uid" param:"uid"`
	Product string `json:"product" query:"product"`
	Client  string `json:"client" query:"client"`
	Channel string `json:"channel" query:"channel"`
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
	Result []schema.UserCacheLabel `json:"result"` // 空数组也保留
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
