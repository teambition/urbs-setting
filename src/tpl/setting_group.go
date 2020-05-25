package tpl

import "time"

// SettingGroupInfo ...
type SettingGroupInfo struct {
	ID         int64     `json:"-" db:"id"`
	SettingHID string    `json:"settingHID"`
	AssignedAt time.Time `json:"assignedAt" db:"assigned_at"`
	Release    int64     `json:"release" db:"rls"`
	Group      string    `json:"group" db:"uid"`
	Kind       string    `json:"kind" db:"kind"`
	Desc       string    `json:"desc" db:"description"`
	Status     int64     `json:"status" db:"status"`
	Value      string    `json:"value" db:"value"`
	LastValue  string    `json:"lastValue" db:"last_value"`
}

// SettingGroupsInfoRes ...
type SettingGroupsInfoRes struct {
	SuccessResponseType
	Result []SettingGroupInfo `json:"result"`
}
