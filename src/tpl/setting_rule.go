package tpl

import (
	"time"

	"github.com/teambition/urbs-setting/src/schema"
	"github.com/teambition/urbs-setting/src/service"
)

// SettingRuleBody ...
type SettingRuleBody struct {
	schema.PercentRule
	Value string `json:"value"`
}

// SettingRuleInfo ...
type SettingRuleInfo struct {
	ID         int64       `json:"-"`
	HID        string      `json:"hid"`
	SettingHID string      `json:"settingHID"`
	Kind       string      `json:"kind"`
	Rule       interface{} `json:"rule"`
	Value      string      `json:"value"`
	Release    int64       `json:"release"`
	CreatedAt  time.Time   `json:"createdAt"`
	UpdatedAt  time.Time   `json:"updatedAt"`
}

// SettingRuleInfoFrom ...
func SettingRuleInfoFrom(settingRule schema.SettingRule) SettingRuleInfo {
	return SettingRuleInfo{
		ID:         settingRule.ID,
		HID:        service.IDToHID(settingRule.ID, "setting_rule"),
		SettingHID: service.IDToHID(settingRule.SettingID, "setting"),
		Kind:       settingRule.Kind,
		Rule:       schema.ToRuleObject(settingRule.Kind, settingRule.Rule),
		Value:      settingRule.Value,
		Release:    settingRule.Release,
		CreatedAt:  settingRule.CreatedAt,
		UpdatedAt:  settingRule.UpdatedAt,
	}
}

// SettingRulesInfoFrom ...
func SettingRulesInfoFrom(settingRules []schema.SettingRule) []SettingRuleInfo {
	res := make([]SettingRuleInfo, len(settingRules))
	for i, l := range settingRules {
		res[i] = SettingRuleInfoFrom(l)
	}
	return res
}

// SettingRulesInfoRes ...
type SettingRulesInfoRes struct {
	SuccessResponseType
	Result []SettingRuleInfo `json:"result"` // 空数组也保留
}

// SettingRuleInfoRes ...
type SettingRuleInfoRes struct {
	SuccessResponseType
	Result SettingRuleInfo `json:"result"` // 空数组也保留
}
