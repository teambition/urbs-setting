package tpl

import (
	"github.com/teambition/urbs-setting/src/schema"
)

// ModuleRes ...
type ModuleRes struct {
	ResponseType
	Result schema.Module `json:"result"` // 空数组也保留
}

// ModulesRes ...
type ModulesRes struct {
	ResponseType
	Result []schema.Module `json:"result"` // 空数组也保留
}
