package tpl

import (
	"github.com/teambition/gear"
	"github.com/teambition/urbs-setting/src/schema"
)

// ModuleUpdateBody ...
type ModuleUpdateBody struct {
	Desc *string `json:"desc"`
}

// Validate 实现 gear.BodyTemplate。
func (t *ModuleUpdateBody) Validate() error {
	if t.Desc == nil {
		return gear.ErrBadRequest.WithMsgf("desc required")
	}

	if len(*t.Desc) > 1022 {
		return gear.ErrBadRequest.WithMsgf("desc too long: %d", len(*t.Desc))
	}
	return nil
}

// ToMap ...
func (t *ModuleUpdateBody) ToMap() map[string]interface{} {
	changed := make(map[string]interface{})
	if t.Desc != nil {
		changed["description"] = *t.Desc
	}
	return changed
}

// ModuleRes ...
type ModuleRes struct {
	SuccessResponseType
	Result schema.Module `json:"result"` // 空数组也保留
}

// ModulesRes ...
type ModulesRes struct {
	SuccessResponseType
	Result []schema.Module `json:"result"` // 空数组也保留
}
