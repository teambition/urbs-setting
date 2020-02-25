package tpl

import (
	"github.com/teambition/gear"
	"github.com/teambition/urbs-setting/src/schema"
)

// ProductURL ...
type ProductURL struct {
	Product string `json:"product" param:"product"`
}

// Validate 实现 gear.BodyTemplate。
func (t *ProductURL) Validate() error {
	if !validIDNameReg.MatchString(t.Product) {
		return gear.ErrBadRequest.WithMsgf("invalid product name: %s", t.Product)
	}
	return nil
}

// ProductLabelURL ...
type ProductLabelURL struct {
	ProductURL
	Label string `json:"label" param:"label"`
}

// Validate 实现 gear.BodyTemplate。
func (t *ProductLabelURL) Validate() error {
	if err := t.ProductURL.Validate(); err != nil {
		return err
	}
	if !validIDNameReg.MatchString(t.Label) {
		return gear.ErrBadRequest.WithMsgf("invalid label name: %s", t.Label)
	}
	return nil
}

// ProductModuleURL ...
type ProductModuleURL struct {
	ProductURL
	Module string `json:"module" param:"module"`
}

// Validate 实现 gear.BodyTemplate。
func (t *ProductModuleURL) Validate() error {
	if err := t.ProductURL.Validate(); err != nil {
		return err
	}
	if !validIDNameReg.MatchString(t.Module) {
		return gear.ErrBadRequest.WithMsgf("invalid module name: %s", t.Module)
	}
	return nil
}

// ProductModuleSettingURL ...
type ProductModuleSettingURL struct {
	ProductModuleURL
	Setting string `json:"setting" param:"setting"`
}

// Validate 实现 gear.BodyTemplate。
func (t *ProductModuleSettingURL) Validate() error {
	if err := t.ProductModuleURL.Validate(); err != nil {
		return err
	}
	if !validIDNameReg.MatchString(t.Setting) {
		return gear.ErrBadRequest.WithMsgf("invalid setting name: %s", t.Setting)
	}
	return nil
}

// ProductsRes ...
type ProductsRes struct {
	ResponseType
	Result []schema.Product `json:"result"`
}

// ProductRes ...
type ProductRes struct {
	ResponseType
	Result schema.Product `json:"result"`
}
