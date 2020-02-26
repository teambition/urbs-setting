package api

import (
	"github.com/teambition/gear"
	"github.com/teambition/urbs-setting/src/bll"
	"github.com/teambition/urbs-setting/src/tpl"
)

// Setting ..
type Setting struct {
	blls *bll.Blls
}

// List ..
func (a *Setting) List(ctx *gear.Context) error {
	req := tpl.ProductModuleURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}
	res, err := a.blls.Setting.List(ctx, req.Product, req.Module)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}

// Create ..
func (a *Setting) Create(ctx *gear.Context) error {
	req := tpl.ProductModuleURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}

	body := tpl.NameDescBody{}
	if err := ctx.ParseBody(&body); err != nil {
		return err
	}

	res, err := a.blls.Setting.Create(ctx, req.Product, req.Module, body.Name, body.Desc)
	if err != nil {
		return err
	}

	return ctx.OkJSON(res)
}

// Update ..
func (a *Setting) Update(ctx *gear.Context) error {
	// TODO
	return nil
}

// Offline ..
func (a *Setting) Offline(ctx *gear.Context) error {
	req := tpl.ProductModuleSettingURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}
	res, err := a.blls.Setting.Offline(ctx, req.Product, req.Module, req.Setting)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}

// Online ..
func (a *Setting) Online(ctx *gear.Context) error {
	// TODO
	return nil
}

// Assign ..
func (a *Setting) Assign(ctx *gear.Context) error {
	req := tpl.ProductModuleSettingURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}

	body := tpl.UsersGroupsBody{}
	if err := ctx.ParseBody(&body); err != nil {
		return err
	}

	if err := a.blls.Setting.Assign(ctx, req.Product, req.Module, req.Setting, body.Value, body.Users, body.Groups); err != nil {
		return err
	}
	return ctx.OkJSON(tpl.BoolRes{Result: true})
}

// Delete ..
func (a *Setting) Delete(ctx *gear.Context) error {
	req := tpl.ProductModuleSettingURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}
	res, err := a.blls.Setting.Delete(ctx, req.Product, req.Module, req.Setting)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}
