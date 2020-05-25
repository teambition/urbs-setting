package tpl

import (
	"time"

	"github.com/teambition/gear"
	"github.com/teambition/urbs-setting/src/schema"
)

// GroupsBody ...
type GroupsBody struct {
	Groups []GroupBody `json:"groups"`
}

// GroupBody ...
type GroupBody struct {
	UID  string `json:"uid"`
	Kind string `json:"kind"`
	Desc string `json:"desc"`
}

// Validate 实现 gear.BodyTemplate。
func (t *GroupsBody) Validate() error {
	if len(t.Groups) == 0 {
		return gear.ErrBadRequest.WithMsg("groups emtpy")
	}
	for _, g := range t.Groups {
		if !validIDReg.MatchString(g.UID) {
			return gear.ErrBadRequest.WithMsgf("invalid group uid: %s", g.UID)
		}
		if !validLabelReg.MatchString(g.Kind) {
			return gear.ErrBadRequest.WithMsgf("invalid group kind: %s", g.Kind)
		}
		if len(g.Desc) > 1022 {
			return gear.ErrBadRequest.WithMsgf("desc too long: %d", len(g.Desc))
		}
	}
	return nil
}

// GroupUpdateBody ...
type GroupUpdateBody struct {
	Desc   *string `json:"desc"`
	SyncAt *int64  `json:"syncAt"`
}

// Validate 实现 gear.BodyTemplate。
func (t *GroupUpdateBody) Validate() error {
	if t.Desc == nil && t.SyncAt == nil {
		return gear.ErrBadRequest.WithMsgf("desc or kind or sync_at required")
	}
	if t.Desc != nil && len(*t.Desc) > 1022 {
		return gear.ErrBadRequest.WithMsgf("desc too long: %d", len(*t.Desc))
	}
	if t.SyncAt != nil {
		now := time.Now().Unix()
		if *t.SyncAt < (now-3600) || *t.SyncAt > (now+3600) {
			// SyncAt 应该在当前时刻前后范围内
			return gear.ErrBadRequest.WithMsgf("invalid sync_at: %d", *t.SyncAt)
		}
	}
	return nil
}

// ToMap ...
func (t *GroupUpdateBody) ToMap() map[string]interface{} {
	changed := make(map[string]interface{})
	if t.Desc != nil {
		changed["description"] = *t.Desc
	}
	if t.SyncAt != nil {
		changed["sync_at"] = *t.SyncAt
	}
	return changed
}

// GroupsURL ...
type GroupsURL struct {
	Pagination
	Kind string `json:"kind" query:"kind"`
}

// GroupMembersURL ...
type GroupMembersURL struct {
	UID    string `json:"uid" param:"uid"`
	User   string `json:"user" query:"user"`     // 根据用户 uid 删除一个成员
	SyncLt int64  `json:"syncLt" query:"syncLt"` // 或根据 sync_lt 删除同步时间小于指定值的所有成员
}

// Validate 实现 gear.BodyTemplate。
func (t *GroupMembersURL) Validate() error {
	if !validIDReg.MatchString(t.UID) {
		return gear.ErrBadRequest.WithMsgf("invalid group uid: %s", t.UID)
	}

	if t.User != "" {
		if !validIDReg.MatchString(t.User) {
			return gear.ErrBadRequest.WithMsgf("invalid user uid: %s", t.User)
		}
	} else if t.SyncLt != 0 {
		if t.SyncLt < 0 || t.SyncLt > (time.Now().UTC().Unix()+3600) {
			// 较大的 SyncLt 可以删除整个群组成员！+3600 是防止把毫秒当秒用
			return gear.ErrBadRequest.WithMsgf("invalid syncLt: %d", t.SyncLt)
		}
	} else {
		return gear.ErrBadRequest.WithMsg("user or syncLt required")
	}
	return nil
}

// GroupsRes ...
type GroupsRes struct {
	SuccessResponseType
	Result []schema.Group `json:"result"`
}

// GroupMember ...
type GroupMember struct {
	ID        int64     `json:"-" db:"id"`
	User      string    `json:"user" db:"uid"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	SyncAt    int64     `json:"syncAt" db:"sync_at"` // 归属关系同步时间戳，1970 以来的秒数，应该与 group.sync_at 相等
}

// GroupMembersRes ...
type GroupMembersRes struct {
	SuccessResponseType
	Result []GroupMember `json:"result"`
}

// GroupRes ...
type GroupRes struct {
	SuccessResponseType
	Result schema.Group `json:"result"`
}
