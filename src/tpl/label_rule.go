package tpl

import (
	"time"

	"github.com/teambition/urbs-setting/src/schema"
	"github.com/teambition/urbs-setting/src/service"
)

// LabelRuleBody ...
type LabelRuleBody struct {
	schema.PercentRule
}

// LabelRuleInfo ...
type LabelRuleInfo struct {
	ID        int64       `json:"-"`
	HID       string      `json:"hid"`
	LabelHID  string      `json:"labelHID"`
	Kind      string      `json:"kind"`
	Rule      interface{} `json:"rule"`
	Release   int64       `json:"release"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

// LabelRuleInfoFrom ...
func LabelRuleInfoFrom(labelRule schema.LabelRule) LabelRuleInfo {
	return LabelRuleInfo{
		ID:        labelRule.ID,
		HID:       service.IDToHID(labelRule.ID, "label_rule"),
		LabelHID:  service.IDToHID(labelRule.LabelID, "label"),
		Kind:      labelRule.Kind,
		Rule:      schema.ToRuleObject(labelRule.Kind, labelRule.Rule),
		Release:   labelRule.Release,
		CreatedAt: labelRule.CreatedAt,
		UpdatedAt: labelRule.UpdatedAt,
	}
}

// LabelRulesInfoFrom ...
func LabelRulesInfoFrom(labelRules []schema.LabelRule) []LabelRuleInfo {
	res := make([]LabelRuleInfo, len(labelRules))
	for i, l := range labelRules {
		res[i] = LabelRuleInfoFrom(l)
	}
	return res
}

// LabelRulesInfoRes ...
type LabelRulesInfoRes struct {
	SuccessResponseType
	Result []LabelRuleInfo `json:"result"` // 空数组也保留
}

// LabelRuleInfoRes ...
type LabelRuleInfoRes struct {
	SuccessResponseType
	Result LabelRuleInfo `json:"result"` // 空数组也保留
}
