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

// GetLables ...
func (b *Group) GetLables(ctx context.Context, uid, product string) (*tpl.LabelsInfoRes, error) {
	group, err := b.ms.Group.FindByUID(ctx, uid, "id")
	if err != nil {
		return nil, err
	}

	res := &tpl.LabelsInfoRes{Result: []tpl.LabelInfo{}}
	if group == nil {
		return res, nil // group 不存在，返回空
	}

	labels, err := b.ms.Group.GetLables(ctx, group.ID, product)
	if err != nil {
		return nil, err
	}
	res.Result = labels
	return res, nil
}

// CheckExists ...
func (b *Group) CheckExists(ctx context.Context, uid string) bool {
	group, _ := b.ms.Group.FindByUID(ctx, uid, "id")
	return group != nil
}

// BatchAdd ...
func (b *Group) BatchAdd(ctx context.Context, groups []tpl.AddGroup) error {
	return b.ms.Group.BatchAdd(ctx, groups)
}

// BatchAddMembers ...
func (b *Group) BatchAddMembers(ctx context.Context, uid string, users []tpl.AddUser) error {
	uids := make([]string, len(users))
	for i, user := range users {
		uids[i] = user.UID
	}

	group, err := b.ms.Group.FindByUID(ctx, uid, "id, `sync_at`")
	if err != nil {
		return err
	}
	if group == nil {
		return gear.ErrNotFound.WithMsgf("group %s not found", uid)
	}

	return b.ms.Group.BatchAddMembers(ctx, group, uids)
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
