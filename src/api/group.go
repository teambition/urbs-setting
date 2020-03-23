package api

import (
	"github.com/teambition/gear"
	"github.com/teambition/urbs-setting/src/bll"
	"github.com/teambition/urbs-setting/src/schema"
	"github.com/teambition/urbs-setting/src/service"
	"github.com/teambition/urbs-setting/src/tpl"
)

// Group ..
type Group struct {
	blls *bll.Blls
}

// List ..
func (a *Group) List(ctx *gear.Context) error {
	req := tpl.GroupsURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}

	res, err := a.blls.Group.List(ctx, req.Kind, req.Pagination)
	if err != nil {
		return err
	}

	return ctx.OkJSON(res)
}

// ListLables ..
func (a *Group) ListLables(ctx *gear.Context) error {
	req := tpl.UIDPaginationURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}

	res, err := a.blls.Group.ListLables(ctx, req.UID, req.Pagination)
	if err != nil {
		return err
	}

	return ctx.OkJSON(res)
}

// ListMembers ..
func (a *Group) ListMembers(ctx *gear.Context) error {
	req := tpl.UIDPaginationURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}

	res, err := a.blls.Group.ListMembers(ctx, req.UID, req.Pagination)
	if err != nil {
		return err
	}

	return ctx.OkJSON(res)
}

// ListSettings ..
func (a *Group) ListSettings(ctx *gear.Context) error {
	req := tpl.UIDProductURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}

	res, err := a.blls.Group.ListSettings(ctx, req.UID, req.Product, req.Pagination)
	if err != nil {
		return err
	}

	return ctx.OkJSON(res)
}

// CheckExists ..
func (a *Group) CheckExists(ctx *gear.Context) error {
	req := tpl.UIDURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}

	res := tpl.BoolRes{Result: a.blls.Group.CheckExists(ctx, req.UID)}
	return ctx.OkJSON(res)
}

// BatchAdd ..
func (a *Group) BatchAdd(ctx *gear.Context) error {
	req := tpl.GroupsBody{}
	if err := ctx.ParseBody(&req); err != nil {
		return err
	}

	if err := a.blls.Group.BatchAdd(ctx, req.Groups); err != nil {
		return err
	}

	return ctx.OkJSON(tpl.BoolRes{Result: true})
}

// Update ..
func (a *Group) Update(ctx *gear.Context) error {
	req := tpl.UIDURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}

	body := tpl.GroupUpdateBody{}
	if err := ctx.ParseBody(&body); err != nil {
		return err
	}

	res, err := a.blls.Group.Update(ctx, req.UID, body)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}

// Delete ..
func (a *Group) Delete(ctx *gear.Context) error {
	req := tpl.UIDURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}
	if err := a.blls.Group.Delete(ctx, req.UID); err != nil {
		return err
	}

	return ctx.OkJSON(tpl.BoolRes{Result: true})
}

// BatchAddMembers ..
func (a *Group) BatchAddMembers(ctx *gear.Context) error {
	req := tpl.UIDURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}

	body := tpl.UsersBody{}
	if err := ctx.ParseBody(&body); err != nil {
		return err
	}

	if err := a.blls.Group.BatchAddMembers(ctx, req.UID, body.Users); err != nil {
		return err
	}

	return ctx.OkJSON(tpl.BoolRes{Result: true})
}

// RemoveMembers ..
func (a *Group) RemoveMembers(ctx *gear.Context) error {
	req := tpl.GroupMembersURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}
	if err := a.blls.Group.RemoveMembers(ctx, req.UID, req.User, req.SyncLt); err != nil {
		return err
	}

	return ctx.OkJSON(tpl.BoolRes{Result: true})
}

// RemoveLable ..
func (a *Group) RemoveLable(ctx *gear.Context) error {
	req := tpl.UIDHIDURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}
	label := &schema.Label{}
	label.ID = service.HIDToID(req.HID, "label")
	if label.ID == 0 {
		return gear.ErrBadRequest.WithMsgf("invalid hid: %s", req.HID)
	}
	if err := a.blls.Group.RemoveLable(ctx, req.UID, label.ID); err != nil {
		return err
	}
	return ctx.OkJSON(tpl.BoolRes{Result: true})
}

// RollbackSetting 回退当前设置值到上一个值
// 更新值请用 POST /products/:product/modules/:module/settings/:setting+:assign 接口
func (a *Group) RollbackSetting(ctx *gear.Context) error {
	req := tpl.UIDHIDURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}
	setting := &schema.Setting{}
	setting.ID = service.HIDToID(req.HID, "setting")
	if setting.ID == 0 {
		return gear.ErrBadRequest.WithMsgf("invalid setting hid: %s", req.HID)
	}
	if err := a.blls.Group.RollbackSetting(ctx, req.UID, setting.ID); err != nil {
		return err
	}
	return ctx.OkJSON(tpl.BoolRes{Result: true})
}

// RemoveSetting ..
func (a *Group) RemoveSetting(ctx *gear.Context) error {
	req := tpl.UIDHIDURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}
	setting := &schema.Setting{}
	setting.ID = service.HIDToID(req.HID, "setting")
	if setting.ID == 0 {
		return gear.ErrBadRequest.WithMsgf("invalid hid: %s", req.HID)
	}
	if err := a.blls.Group.RemoveSetting(ctx, req.UID, setting.ID); err != nil {
		return err
	}
	return ctx.OkJSON(tpl.BoolRes{Result: true})
}
