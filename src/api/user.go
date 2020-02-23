package api

import (
	"github.com/teambition/gear"
	"github.com/teambition/urbs-setting/src/bll"
	"github.com/teambition/urbs-setting/src/tpl"
)

// User ..
type User struct {
	blls *bll.Blls
}

// GetLables ..
func (a *User) GetLables(ctx *gear.Context) error {
	req := tpl.QueryLabel{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}

	res, err := a.blls.User.GetLables(ctx, req.UID, req.Product, req.Client, req.Channel)
	if err != nil {
		return err
	}

	return ctx.OkJSON(res)
}

// GetSettings ..
func (a *User) GetSettings(ctx *gear.Context) error {
	// TODO
	return nil
}

// CheckExists ..
func (a *User) CheckExists(ctx *gear.Context) error {
	uid := ctx.Param("uid")
	res := tpl.ResponseType{Result: false}
	if uid != "" {
		res.Result = a.blls.User.CheckExists(ctx, uid)
	}
	return ctx.OkJSON(res)
}

// BatchAdd ..
func (a *User) BatchAdd(ctx *gear.Context) error {
	return nil
}

// RemoveLable ..
func (a *User) RemoveLable(ctx *gear.Context) error {
	return nil
}

// UpdateSetting ..
func (a *User) UpdateSetting(ctx *gear.Context) error {
	return nil
}

// RemoveSetting ..
func (a *User) RemoveSetting(ctx *gear.Context) error {
	return nil
}
