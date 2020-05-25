package tpl

import "time"

// LabelGroupInfo ...
type LabelGroupInfo struct {
	ID         int64     `json:"-" db:"id"`
	LabelHID   string    `json:"labelHID"`
	AssignedAt time.Time `json:"assignedAt" db:"assigned_at"`
	Release    int64     `json:"release" db:"rls"`
	Group      string    `json:"group" db:"uid"`
	Kind       string    `json:"kind" db:"kind"`
	Desc       string    `json:"desc" db:"description"`
	Status     int64     `json:"status" db:"status"`
}

// LabelGroupsInfoRes ...
type LabelGroupsInfoRes struct {
	SuccessResponseType
	Result []LabelGroupInfo `json:"result"`
}
