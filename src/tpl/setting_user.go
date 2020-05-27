package tpl

import "time"

// SettingUserInfo ...
type SettingUserInfo struct {
	ID         int64     `json:"-" db:"id"`
	SettingHID string    `json:"settingHID"`
	AssignedAt time.Time `json:"assignedAt" db:"assigned_at"`
	Release    int64     `json:"release" db:"rls"`
	User       string    `json:"user" db:"uid"`
	Value      string    `json:"value" db:"value"`
	LastValue  string    `json:"lastValue" db:"last_value"`
}

// SettingUsersInfoRes ...
type SettingUsersInfoRes struct {
	SuccessResponseType
	Result []SettingUserInfo `json:"result"`
}
