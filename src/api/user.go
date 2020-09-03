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

// List ..
func (a *User) List(ctx *gear.Context) error {
	req := tpl.Pagination{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}

	res, err := a.blls.User.List(ctx, req)
	if err != nil {
		return err
	}

	return ctx.OkJSON(res)
}

// ListCachedLabels 返回执行 user 在 product 下所有 labels，按照 label 指派时间反序
func (a *User) ListCachedLabels(ctx *gear.Context) error {
	req := tpl.UIDProductURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}

	res := a.blls.User.ListCachedLabels(ctx, req.UID, req.Product)
	return ctx.OkJSON(res)
}

// RefreshCachedLabels 强制更新 user 的 labels 缓存
func (a *User) RefreshCachedLabels(ctx *gear.Context) error {
	req := tpl.UIDAndProductURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}

	user, err := a.blls.User.RefreshCachedLabels(ctx, req.Product, req.UID)
	if err != nil {
		return err
	}
	res := tpl.UserRes{Result: *user}
	return ctx.OkJSON(res)
}

// ListLabels 返回 user 的 labels，按照 label 指派时间正序，支持分页
func (a *User) ListLabels(ctx *gear.Context) error {
	req := tpl.UIDPaginationURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}

	res, err := a.blls.User.ListLabels(ctx, req.UID, req.Pagination)
	if err != nil {
		return err
	}

	return ctx.OkJSON(res)
}

// ListSettings 返回 user 的 settings，按照 setting 设置时间正序，支持分页
func (a *User) ListSettings(ctx *gear.Context) error {
	req := tpl.MySettingsQueryURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}

	res, err := a.blls.User.ListSettings(ctx, req)
	if err != nil {
		return err
	}

	return ctx.OkJSON(res)
}

// ListSettingsUnionAll 返回 user 的 settings，按照 setting 设置时间反序，支持分页
// 包含了 user 从属的 group 的 settings
func (a *User) ListSettingsUnionAll(ctx *gear.Context) error {
	req := tpl.MySettingsQueryURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}

	if req.Product == "" {
		return gear.ErrBadRequest.WithMsgf("product required")
	}

	res, err := a.blls.User.ListSettingsUnionAll(ctx, req)
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

// ApplyRules ..
func (a *User) ApplyRules(ctx *gear.Context) error {
	req := &tpl.ProductURL{}
	if err := ctx.ParseURL(req); err != nil {
		return err
	}
	body := &tpl.ApplyRulesBody{}
	if err := ctx.ParseBody(body); err != nil {
		return err
	}
	err := a.blls.User.ApplyRules(ctx, req.Product, body)
	if err != nil {
		return err
	}
	return ctx.OkJSON(tpl.BoolRes{Result: true})
}
