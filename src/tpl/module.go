package tpl

import (
	"github.com/teambition/urbs-setting/src/schema"
)

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
