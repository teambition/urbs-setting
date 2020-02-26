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

// ListLablesInCache ..
func (a *User) ListLablesInCache(ctx *gear.Context) error {
	req := tpl.LabelsURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}

	if req.Product == "" {
		return gear.ErrBadRequest.WithMsg("product required")
	}

	res, err := a.blls.User.ListLablesInCache(ctx, req.UID, req.Product, req.Client, req.Channel)
	if err != nil {
		return err
	}

	return ctx.OkJSON(res)
}

// ListSettings ..
func (a *User) ListSettings(ctx *gear.Context) error {
	// TODO
	return nil
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
	if err := service.HIDer.PutHID(label, req.HID); err != nil {
		return err
	}
	if err := a.blls.User.RemoveLable(ctx, req.UID, label.ID); err != nil {
		return err
	}
	return ctx.OkJSON(tpl.BoolRes{Result: true})
}

// UpdateSetting ..
func (a *User) UpdateSetting(ctx *gear.Context) error {
	// TODO
	return nil
}

// RemoveSetting ..
func (a *User) RemoveSetting(ctx *gear.Context) error {
	req := tpl.UIDHIDURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}
	setting := &schema.Setting{}
	if err := service.HIDer.PutHID(setting, req.HID); err != nil {
		return err
	}
	if err := a.blls.User.RemoveSetting(ctx, req.UID, setting.ID); err != nil {
		return err
	}
	return ctx.OkJSON(tpl.BoolRes{Result: true})
}
