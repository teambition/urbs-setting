package tpl

import "time"

// LabelUserInfo ...
type LabelUserInfo struct {
	ID         int64     `json:"-"`
	LabelHID   string    `json:"labelHID"`
	AssignedAt time.Time `json:"assignedAt"`
	Release    int64     `json:"release"`
	User       string    `json:"user"`
}

// LabelUsersInfoRes ...
type LabelUsersInfoRes struct {
	SuccessResponseType
	Result []LabelUserInfo `json:"result"`
}
