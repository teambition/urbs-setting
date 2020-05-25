package tpl

import "time"

// LabelUserInfo ...
type LabelUserInfo struct {
	ID         int64     `json:"-" db:"id"`
	LabelHID   string    `json:"labelHID"`
	AssignedAt time.Time `json:"assignedAt" db:"assigned_at"`
	Release    int64     `json:"release" db:"rls"`
	User       string    `json:"user" db:"uid"`
}

// LabelUsersInfoRes ...
type LabelUsersInfoRes struct {
	SuccessResponseType
	Result []LabelUserInfo `json:"result"`
}
