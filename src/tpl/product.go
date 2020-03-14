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

// ProductPaginationURL ...
type ProductPaginationURL struct {
	Pagination
	Product string `json:"product" param:"product"`
}

// Validate 实现 gear.BodyTemplate。
func (t *ProductPaginationURL) Validate() error {
	if !validIDNameReg.MatchString(t.Product) {
		return gear.ErrBadRequest.WithMsgf("invalid product name: %s", t.Product)
	}
	if err := t.Pagination.Validate(); err != nil {
		return err
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
	if !validLabelReg.MatchString(t.Label) {
		return gear.ErrBadRequest.WithMsgf("invalid label: %s", t.Label)
	}
	return nil
}

// ProductModuleURL ...
type ProductModuleURL struct {
	Pagination
	ProductURL
	Module string `json:"module" param:"module"`
}

// Validate 实现 gear.BodyTemplate。
func (t *ProductModuleURL) Validate() error {
	if !validIDNameReg.MatchString(t.Module) {
		return gear.ErrBadRequest.WithMsgf("invalid module name: %s", t.Module)
	}
	if err := t.ProductURL.Validate(); err != nil {
		return err
	}
	if err := t.Pagination.Validate(); err != nil {
		return err
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
	SuccessResponseType
	Result []schema.Product `json:"result"`
}

// ProductRes ...
type ProductRes struct {
	SuccessResponseType
	Result schema.Product `json:"result"`
}
