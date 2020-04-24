package tpl

import "time"

// SettingUserInfo ...
type SettingUserInfo struct {
	ID         int64     `json:"-"`
	SettingHID string    `json:"settingHID"`
	AssignedAt time.Time `json:"assignedAt"`
	Release    int64     `json:"release"`
	User       string    `json:"user"`
	Value      string    `json:"value"`
	LastValue  string    `json:"lastValue"`
}

// SettingUsersInfoRes ...
type SettingUsersInfoRes struct {
	SuccessResponseType
	Result []SettingUserInfo `json:"result"`
}
