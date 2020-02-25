package tpl

import (
	"time"

	"github.com/teambition/gear"
)

// GroupsBody ...
type GroupsBody struct {
	Groups []GroupBody `json:"groups"`
}

// GroupBody ...
type GroupBody struct {
	UID  string `json:"uid"`
	Desc string `json:"desc"`
}

// Validate 实现 gear.BodyTemplate。
func (t *GroupsBody) Validate() error {
	if len(t.Groups) == 0 {
		return gear.ErrBadRequest.WithMsg("groups emtpy")
	}
	for _, g := range t.Groups {
		if !validIDNameReg.MatchString(g.UID) {
			return gear.ErrBadRequest.WithMsgf("invalid group uid: %s", g.UID)
		}
		if len(g.Desc) > 1023 {
			return gear.ErrBadRequest.WithMsgf("desc too long: %d", len(g.Desc))
		}
	}
	return nil
}

// RemoveGroupMembersURL ...
type RemoveGroupMembersURL struct {
	UID    string `json:"uid" param:"uid"`
	User   string `json:"user" query:"user"`       // 根据用户 uid 删除一个成员
	SyncLt int64  `json:"sync_lt" query:"sync_lt"` // 或根据 sync_lt 删除同步时间小于指定值的所有成员
}

// Validate 实现 gear.BodyTemplate。
func (t *RemoveGroupMembersURL) Validate() error {
	if !validIDNameReg.MatchString(t.UID) {
		return gear.ErrBadRequest.WithMsgf("invalid group uid: %s", t.UID)
	}

	if t.User != "" {
		if !validIDNameReg.MatchString(t.User) {
			return gear.ErrBadRequest.WithMsgf("invalid user uid: %s", t.User)
		}
	} else if t.SyncLt != 0 {
		if t.SyncLt < 0 || t.SyncLt > (time.Now().UTC().Unix()+3600) {
			// 较大的 SyncLt 可以删除整个群组成员！+3600 是防止把毫秒当秒用
			return gear.ErrBadRequest.WithMsgf("invalid sync_lt: %d", t.SyncLt)
		}
	} else {
		return gear.ErrBadRequest.WithMsg("user or sync_lt required")
	}
	return nil
}
