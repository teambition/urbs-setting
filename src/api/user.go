package api

import (
	"github.com/teambition/gear"
	"github.com/teambition/urbs-setting/src/bll"
	"github.com/teambition/urbs-setting/src/schema"
	"github.com/teambition/urbs-setting/src/service"
	"github.com/teambition/urbs-setting/src/tpl"
)

// User ..
type User struct {
	blls *bll.Blls
}

// ListCachedLables 返回执行 user 在 product 下所有 labels，按照 label 指派时间反序
func (a *User) ListCachedLables(ctx *gear.Context) error {
	req := tpl.UIDProductURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}

	res, err := a.blls.User.ListCachedLables(ctx, req.UID, req.Product)
	if err != nil {
		return err
	}

	return ctx.OkJSON(res)
}

// RefreshCachedLables 强制更新 user 的 labels 缓存
func (a *User) RefreshCachedLables(ctx *gear.Context) error {
	req := tpl.UIDPaginationURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}

	err := a.blls.User.RefreshCachedLables(ctx, req.UID)
	if err != nil {
		return err
	}
	res := tpl.BoolRes{Result: true}
	return ctx.OkJSON(res)
}

// ListLables 返回 user 的 labels，按照 label 指派时间正序，支持分页
func (a *User) ListLables(ctx *gear.Context) error {
	req := tpl.UIDPaginationURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}

	res, err := a.blls.User.ListLables(ctx, req.UID, req.Pagination)
	if err != nil {
		return err
	}

	return ctx.OkJSON(res)
}

// ListSettings 返回 user 的 settings，按照 setting 设置时间正序，支持分页
func (a *User) ListSettings(ctx *gear.Context) error {
	req := tpl.UIDProductURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}

	res, err := a.blls.User.ListSettings(ctx, req.UID, req.Product, req.Pagination)
	if err != nil {
		return err
	}

	return ctx.OkJSON(res)
}

// ListSettingsUnionAll 返回 user 的 settings，按照 setting 设置时间反序，支持分页
// 包含了 user 从属的 group 的 settings
func (a *User) ListSettingsUnionAll(ctx *gear.Context) error {
	req := tpl.UIDProductURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}

	res, err := a.blls.User.ListSettingsUnionAll(ctx, req.UID, req.Product, req.Pagination)
	if err != nil {
		return err
	}

	return ctx.OkJSON(res)
}

// CheckExists ..
func (a *User) CheckExists(ctx *gear.Context) error {
	req := tpl.UIDURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}

	res := tpl.BoolRes{Result: a.blls.User.CheckExists(ctx, req.UID)}
	return ctx.OkJSON(res)
}

// BatchAdd ..
func (a *User) BatchAdd(ctx *gear.Context) error {
	req := tpl.UsersBody{}
	if err := ctx.ParseBody(&req); err != nil {
		return err
	}

	if err := a.blls.User.BatchAdd(ctx, req.Users); err != nil {
		return err
	}

	return ctx.OkJSON(tpl.BoolRes{Result: true})
}

// RemoveLable ..
func (a *User) RemoveLable(ctx *gear.Context) error {
	req := tpl.UIDHIDURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}
	label := &schema.Label{}
	label.ID = service.HIDToID(req.HID, "label")
	if label.ID == 0 {
		return gear.ErrBadRequest.WithMsgf("invalid hid: %s", req.HID)
	}
	if err := a.blls.User.RemoveLable(ctx, req.UID, label.ID); err != nil {
		return err
	}
	return ctx.OkJSON(tpl.BoolRes{Result: true})
}

// RollbackSetting 回退当前设置值到上一个值
// 更新值请用 POST /products/:product/modules/:module/settings/:setting+:assign 接口
func (a *User) RollbackSetting(ctx *gear.Context) error {
	req := tpl.UIDHIDURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}
	setting := &schema.Setting{}
	setting.ID = service.HIDToID(req.HID, "setting")
	if setting.ID == 0 {
		return gear.ErrBadRequest.WithMsgf("invalid setting hid: %s", req.HID)
	}
	if err := a.blls.User.RollbackSetting(ctx, req.UID, setting.ID); err != nil {
		return err
	}
	return ctx.OkJSON(tpl.BoolRes{Result: true})
}

// RemoveSetting ..
func (a *User) RemoveSetting(ctx *gear.Context) error {
	req := tpl.UIDHIDURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}
	setting := &schema.Setting{}
	setting.ID = service.HIDToID(req.HID, "setting")
	if setting.ID == 0 {
		return gear.ErrBadRequest.WithMsgf("invalid hid: %s", req.HID)
	}
	if err := a.blls.User.RemoveSetting(ctx, req.UID, setting.ID); err != nil {
		return err
	}
	return ctx.OkJSON(tpl.BoolRes{Result: true})
}
