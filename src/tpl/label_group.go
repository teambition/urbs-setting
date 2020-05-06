package tpl

import "time"

// LabelGroupInfo ...
type LabelGroupInfo struct {
	ID         int64     `json:"-"`
	LabelHID   string    `json:"labelHID"`
	AssignedAt time.Time `json:"assignedAt"`
	Release    int64     `json:"release"`
	Group      string    `json:"group"`
	Kind       string    `json:"kind"`
	Desc       string    `json:"desc"`
	Status     int64     `json:"status"`
}

// LabelGroupsInfoRes ...
type LabelGroupsInfoRes struct {
	SuccessResponseType
	Result []LabelGroupInfo `json:"result"`
}
