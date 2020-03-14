package api

import (
	"github.com/teambition/gear"
	"github.com/teambition/urbs-setting/src/bll"
	"github.com/teambition/urbs-setting/src/tpl"
)

// Product ..
type Product struct {
	blls *bll.Blls
}

// List ..
func (a *Product) List(ctx *gear.Context) error {
	req := tpl.Pagination{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}

	res, err := a.blls.Product.List(ctx, req)
	if err != nil {
		return err
	}

	return ctx.OkJSON(res)
}

// Create ..
func (a *Product) Create(ctx *gear.Context) error {
	body := tpl.NameDescBody{}
	if err := ctx.ParseBody(&body); err != nil {
		return err
	}

	res, err := a.blls.Product.Create(ctx, body.Name, body.Desc)
	if err != nil {
		return err
	}

	return ctx.OkJSON(res)
}

// Update ..
func (a *Product) Update(ctx *gear.Context) error {
	// TODO
	return nil
}

// Offline ..
func (a *Product) Offline(ctx *gear.Context) error {
	req := tpl.ProductURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}
	res, err := a.blls.Product.Offline(ctx, req.Product)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}

// Online ..
func (a *Product) Online(ctx *gear.Context) error {
	// TODO
	return nil
}

// Delete ..
func (a *Product) Delete(ctx *gear.Context) error {
	req := tpl.ProductURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}
	res, err := a.blls.Product.Delete(ctx, req.Product)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}
