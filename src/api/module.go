package api

import (
	"github.com/teambition/gear"
	"github.com/teambition/urbs-setting/src/bll"
	"github.com/teambition/urbs-setting/src/tpl"
)

// Module ..
type Module struct {
	blls *bll.Blls
}

// List ..
func (a *Module) List(ctx *gear.Context) error {
	req := tpl.ProductPaginationURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}
	res, err := a.blls.Module.List(ctx, req.Product, req.Pagination)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}

// Create ..
func (a *Module) Create(ctx *gear.Context) error {
	req := tpl.ProductURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}

	body := tpl.NameDescBody{}
	if err := ctx.ParseBody(&body); err != nil {
		return err
	}

	res, err := a.blls.Module.Create(ctx, req.Product, body.Name, body.Desc)
	if err != nil {
		return err
	}

	return ctx.OkJSON(res)
}

// Update ..
func (a *Module) Update(ctx *gear.Context) error {
	req := tpl.ProductModuleURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}

	body := tpl.ModuleUpdateBody{}
	if err := ctx.ParseBody(&body); err != nil {
		return err
	}

	res, err := a.blls.Module.Update(ctx, req.Product, req.Module, body)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}

// Offline ..
func (a *Module) Offline(ctx *gear.Context) error {
	req := tpl.ProductModuleURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}
	res, err := a.blls.Module.Offline(ctx, req.Product, req.Module)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}
