package bll

import (
	"context"

	"github.com/teambition/urbs-setting/src/model"
	"github.com/teambition/urbs-setting/src/tpl"
)

// Group ...
type Group struct {
	ms *model.Models
}

// List 返回群组列表，TODO：支持分页
func (b *Group) List(ctx context.Context, kind string, pg tpl.Pagination) (*tpl.GroupsRes, error) {
	groups, total, err := b.ms.Group.Find(ctx, kind, pg)
	if err != nil {
		return nil, err
	}
	res := &tpl.GroupsRes{Result: groups}
	res.TotalSize = total
	if len(res.Result) > pg.PageSize {
		res.NextPageToken = tpl.IDToPageToken(res.Result[pg.PageSize].ID)
		res.Result = res.Result[:pg.PageSize]
	}
	return res, nil
}

// ListLabels ...
func (b *Group) ListLabels(ctx context.Context, uid string, pg tpl.Pagination) (*tpl.LabelsInfoRes, error) {
	group, err := b.ms.Group.Acquire(ctx, uid)
	if err != nil {
		return nil, err
	}

	labels, total, err := b.ms.Group.FindLabels(ctx, group.ID, pg)
	if err != nil {
		return nil, err
	}

	res := &tpl.LabelsInfoRes{Result: labels}
	res.TotalSize = total
	if len(res.Result) > pg.PageSize {
		res.NextPageToken = tpl.IDToPageToken(res.Result[pg.PageSize].ID)
		res.Result = res.Result[:pg.PageSize]
	}
	return res, nil
}

// ListMembers ...
func (b *Group) ListMembers(ctx context.Context, uid string, pg tpl.Pagination) (*tpl.GroupMembersRes, error) {
	group, err := b.ms.Group.Acquire(ctx, uid)
	if err != nil {
		return nil, err
	}

	members, total, err := b.ms.Group.FindMembers(ctx, group.ID, pg)
	if err != nil {
		return nil, err
	}

	res := &tpl.GroupMembersRes{Result: members}
	res.TotalSize = total
	if len(res.Result) > pg.PageSize {
		res.NextPageToken = tpl.IDToPageToken(res.Result[pg.PageSize].ID)
		res.Result = res.Result[:pg.PageSize]
	}
	return res, nil
}

// ListSettings ...
func (b *Group) ListSettings(ctx context.Context, uid string, pg tpl.Pagination) (*tpl.MySettingsRes, error) {
	group, err := b.ms.Group.Acquire(ctx, uid)
	if err != nil {
		return nil, err
	}

	settings, total, err := b.ms.Group.FindSettings(ctx, group.ID, pg)
	if err != nil {
		return nil, err
	}

	res := &tpl.MySettingsRes{Result: settings}
	res.TotalSize = total
	if len(res.Result) > pg.PageSize {
		res.NextPageToken = tpl.IDToPageToken(res.Result[pg.PageSize].ID)
		res.Result = res.Result[:pg.PageSize]
	}
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
	group, err := b.ms.Group.Acquire(ctx, uid)
	if err != nil {
		return err
	}

	if err = b.ms.User.BatchAdd(ctx, users); err != nil {
		return err
	}

	return b.ms.Group.BatchAddMembers(ctx, group, users)
}

// RemoveMembers ...
func (b *Group) RemoveMembers(ctx context.Context, uid, userUID string, syncLt int64) error {
	group, err := b.ms.Group.Acquire(ctx, uid)
	if err != nil {
		return err
	}

	var userID int64
	if userUID != "" {
		if user, _ := b.ms.User.FindByUID(ctx, userUID, "id"); user != nil {
			userID = user.ID
		}
	}

	return b.ms.Group.RemoveMembers(ctx, group.ID, userID, syncLt)
}

// RemoveLabel ...
func (b *Group) RemoveLabel(ctx context.Context, uid string, labelID int64) error {
	group, err := b.ms.Group.Acquire(ctx, uid)
	if err != nil {
		return err
	}
	return b.ms.Label.RemoveGroupLabel(ctx, group.ID, labelID)
}

// RollbackSetting ...
func (b *Group) RollbackSetting(ctx context.Context, uid string, settingID int64) error {
	group, err := b.ms.Group.Acquire(ctx, uid)
	if err != nil {
		return err
	}

	return b.ms.Setting.RollbackGroupSetting(ctx, group.ID, settingID)
}

// RemoveSetting ...
func (b *Group) RemoveSetting(ctx context.Context, uid string, settingID int64) error {
	group, err := b.ms.Group.Acquire(ctx, uid)
	if err != nil {
		return err
	}

	return b.ms.Setting.RemoveGroupSetting(ctx, group.ID, settingID)
}

// Update ...
func (b *Group) Update(ctx context.Context, uid string, body tpl.GroupUpdateBody) (*tpl.GroupRes, error) {
	group, err := b.ms.Group.Acquire(ctx, uid)
	if err != nil {
		return nil, err
	}
	group, err = b.ms.Group.Update(ctx, group.ID, body.ToMap())
	if err != nil {
		return nil, err
	}
	return &tpl.GroupRes{Result: *group}, nil
}

// Delete ...
func (b *Group) Delete(ctx context.Context, uid string) error {
	group, _ := b.ms.Group.FindByUID(ctx, uid, "id")
	if group == nil {
		return nil
	}
	return b.ms.Group.Delete(ctx, group.ID)
}
