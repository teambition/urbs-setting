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
func (q *QueryLabel) Validate() error {
	if q.UID == "" {
		return gear.ErrBadRequest.WithMsg("invalid uid")
	}
	if q.Product == "" {
		return gear.ErrBadRequest.WithMsg("product name")
	}
	return nil
}

// LabelsResponse ...
type LabelsResponse struct {
	ResponseType
	Result []schema.UserLabelInfo `json:"result"`
}
