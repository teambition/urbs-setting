package tpl

import (
	"github.com/teambition/gear"
	"github.com/teambition/urbs-setting/src/schema"
)

// QueryLabel ...
type QueryLabel struct {
	UID     string `json:"uid" param:"uid"`
	Product string `json:"product" query:"product"`
	Client  string `json:"client" query:"client"`
	Channel string `json:"channel" query:"channel"`
}

// Validate 实现 gear.BodyTemplate。
func (t *QueryLabel) Validate() error {
	if !validIDNameReg.MatchString(t.UID) {
		return gear.ErrBadRequest.WithMsgf("invalid uid: %s", t.UID)
	}
	if t.Product != "" && !validIDNameReg.MatchString(t.Product) {
		return gear.ErrBadRequest.WithMsgf("invalid product name: %s", t.Product)
	}
	return nil
}

// LabelInfo ...
type LabelInfo struct {
	HID     string `json:"hid"`
	Product string `json:"product"`
	Name    string `json:"name"`
	Desc    string `json:"desc"`
}

// LabelsInfoRes ...
type LabelsInfoRes struct {
	ResponseType
	Result []LabelInfo `json:"result"`
}

// CacheLabelsRes ...
type CacheLabelsRes struct {
	ResponseType
	Result []schema.UserCacheLabel `json:"result"`
}
