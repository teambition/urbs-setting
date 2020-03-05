package bll

import (
	"context"

	"github.com/teambition/gear"
	"github.com/teambition/urbs-setting/src/model"
	"github.com/teambition/urbs-setting/src/tpl"
)

// Group ...
type Group struct {
	ms *model.Models
}

// List 返回群组列表，TODO：支持分页
func (b *Group) List(ctx context.Context) (*tpl.GroupsRes, error) {
	groups, err := b.ms.Group.Find(ctx)
	if err != nil {
		return nil, err
	}
	res := &tpl.GroupsRes{Result: groups}
	return res, nil
}

// ListLables ...
func (b *Group) ListLables(ctx context.Context, uid, product string) (*tpl.LabelsInfoRes, error) {
	group, err := b.ms.Group.FindByUID(ctx, uid, "id")
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, gear.ErrNotFound.WithMsgf("group %s not found", uid)
	}

	labels, err := b.ms.Group.FindLables(ctx, group.ID, product)
	if err != nil {
		return nil, err
	}
	res := &tpl.LabelsInfoRes{Result: labels}
	return res, nil
}

// ListMembers ...
func (b *Group) ListMembers(ctx context.Context, uid string) (*tpl.GroupMembersRes, error) {
	group, err := b.ms.Group.FindByUID(ctx, uid, "id")
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, gear.ErrNotFound.WithMsgf("group %s not found", uid)
	}

	members, err := b.ms.Group.FindMembers(ctx, group.ID)
	if err != nil {
		return nil, err
	}
	res := &tpl.GroupMembersRes{Result: members}
	return res, nil
}

// CheckExists ...
func (b *Group) CheckExists(ctx context.Context, uid string) bool {
	group, _ := b.ms.Group.FindByUID(ctx, uid, "id")
	return group != nil
}

// BatchAdd ...
func (b *Group) BatchAdd(ctx context.Context, groups []tpl.GroupBody) error {
	return b.ms.Group.BatchAdd(ctx, groups)
}

// BatchAddMembers 批量给群组添加成员，如果用户未加入系统，则会自动加入
func (b *Group) BatchAddMembers(ctx context.Context, uid string, users []string) error {
	group, err := b.ms.Group.FindByUID(ctx, uid, "id, `sync_at`")
	if err != nil {
		return err
	}
	if group == nil {
		return gear.ErrNotFound.WithMsgf("group %s not found", uid)
	}

	if err := b.ms.User.BatchAdd(ctx, users); err != nil {
		return err
	}

	return b.ms.Group.BatchAddMembers(ctx, group, users)
}

// RemoveMembers ...
func (b *Group) RemoveMembers(ctx context.Context, uid, userUID string, syncLt int64) error {
	group, _ := b.ms.Group.FindByUID(ctx, uid, "id")
	if group == nil {
		return gear.ErrNotFound.WithMsgf("Group not found: %s", uid)
	}

	var userID int64
	if user, _ := b.ms.User.FindByUID(ctx, userUID, "id"); user != nil {
		userID = user.ID
	}
	return b.ms.Group.RemoveMembers(ctx, group.ID, userID, syncLt)
}

// RemoveLable ...
func (b *Group) RemoveLable(ctx context.Context, uid string, lableID int64) error {
	group, _ := b.ms.Group.FindByUID(ctx, uid, "id")
	if group == nil {
		return gear.ErrNotFound.WithMsgf("Group not found: %s", uid)
	}
	return b.ms.Group.RemoveLable(ctx, group.ID, lableID)
}

// RemoveSetting ...
func (b *Group) RemoveSetting(ctx context.Context, uid string, settingID int64) error {
	group, _ := b.ms.Group.FindByUID(ctx, uid, "id")
	if group == nil {
		return gear.ErrNotFound.WithMsgf("Group not found: %s", uid)
	}

	return b.ms.Group.RemoveSetting(ctx, group.ID, settingID)
}
