package api

import (
	"github.com/teambition/gear"
	"github.com/teambition/urbs-setting/src/bll"
	"github.com/teambition/urbs-setting/src/service"
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
	res, err := a.blls.Setting.List(ctx, req.Product, req.Module, req.Pagination)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}

// ListByProduct ..
func (a *Setting) ListByProduct(ctx *gear.Context) error {
	req := tpl.ProductPaginationURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}
	res, err := a.blls.Setting.List(ctx, req.Product, "", req.Pagination)
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

	body := tpl.SettingBody{}
	if err := ctx.ParseBody(&body); err != nil {
		return err
	}

	res, err := a.blls.Setting.Create(ctx, req.Product, req.Module, &body)
	if err != nil {
		return err
	}

	return ctx.OkJSON(res)
}

// Get ..
func (a *Setting) Get(ctx *gear.Context) error {
	req := tpl.ProductModuleSettingURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}
	res, err := a.blls.Setting.Get(ctx, req.Product, req.Module, req.Setting)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}

// Update ..
func (a *Setting) Update(ctx *gear.Context) error {
	req := tpl.ProductModuleSettingURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}

	body := tpl.SettingUpdateBody{}
	if err := ctx.ParseBody(&body); err != nil {
		return err
	}

	res, err := a.blls.Setting.Update(ctx, req.Product, req.Module, req.Setting, body)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
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
	res, err := a.blls.Setting.Assign(ctx, req.Product, req.Module, req.Setting, body.Value, body.Users, body.Groups)
	if err != nil {
		return err
	}
	return ctx.OkJSON(tpl.SettingReleaseInfoRes{Result: *res})
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

// Recall ..
func (a *Setting) Recall(ctx *gear.Context) error {
	req := tpl.ProductModuleSettingURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}

	body := tpl.RecallBody{}
	if err := ctx.ParseBody(&body); err != nil {
		return err
	}

	res, err := a.blls.Setting.Recall(ctx, req.Product, req.Module, req.Setting, body.Release)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}

// CreateRule ..
func (a *Setting) CreateRule(ctx *gear.Context) error {
	req := tpl.ProductModuleSettingURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}

	body := tpl.SettingRuleBody{}
	if err := ctx.ParseBody(&body); err != nil {
		return err
	}

	res, err := a.blls.Setting.CreateRule(ctx, req.Product, req.Module, req.Setting, body)
	if err != nil {
		return err
	}

	return ctx.OkJSON(res)
}

// ListRules ..
func (a *Setting) ListRules(ctx *gear.Context) error {
	req := tpl.ProductModuleSettingURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}
	res, err := a.blls.Setting.ListRules(ctx, req.Product, req.Module, req.Setting)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}

// UpdateRule ..
func (a *Setting) UpdateRule(ctx *gear.Context) error {
	req := tpl.ProductModuleSettingHIDURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}

	ruleID := service.HIDToID(req.HID, "setting_rule")
	if ruleID <= 0 {
		return gear.ErrBadRequest.WithMsgf("invalid setting_rule hid: %s", req.HID)
	}

	body := tpl.SettingRuleBody{}
	if err := ctx.ParseBody(&body); err != nil {
		return err
	}

	res, err := a.blls.Setting.UpdateRule(ctx, req.Product, req.Module, req.Setting, ruleID, body)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}

// DeleteRule ..
func (a *Setting) DeleteRule(ctx *gear.Context) error {
	req := tpl.ProductModuleSettingHIDURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}

	ruleID := service.HIDToID(req.HID, "setting_rule")
	if ruleID <= 0 {
		return gear.ErrBadRequest.WithMsgf("invalid setting_rule hid: %s", req.HID)
	}

	res, err := a.blls.Setting.DeleteRule(ctx, req.Product, req.Module, req.Setting, ruleID)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}

// ListUsers ..
func (a *Setting) ListUsers(ctx *gear.Context) error {
	req := tpl.ProductModuleSettingURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}
	res, err := a.blls.Setting.ListUsers(ctx, req.Product, req.Module, req.Setting, req.Pagination)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}

// RollbackUserSetting ..
func (a *Setting) RollbackUserSetting(ctx *gear.Context) error {
	req := tpl.ProductModuleSettingUIDURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}

	res, err := a.blls.Setting.RollbackUserSetting(ctx, req.Product, req.Module, req.Setting, req.UID)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}

// DeleteUser ..
func (a *Setting) DeleteUser(ctx *gear.Context) error {
	req := tpl.ProductModuleSettingUIDURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}

	res, err := a.blls.Setting.DeleteUser(ctx, req.Product, req.Module, req.Setting, req.UID)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}

// ListGroups ..
func (a *Setting) ListGroups(ctx *gear.Context) error {
	req := tpl.ProductModuleSettingURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}
	res, err := a.blls.Setting.ListGroups(ctx, req.Product, req.Module, req.Setting, req.Pagination)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}

// RollbackGroupSetting ..
func (a *Setting) RollbackGroupSetting(ctx *gear.Context) error {
	req := tpl.ProductModuleSettingUIDURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}

	res, err := a.blls.Setting.RollbackGroupSetting(ctx, req.Product, req.Module, req.Setting, req.UID)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}

// DeleteGroup ..
func (a *Setting) DeleteGroup(ctx *gear.Context) error {
	req := tpl.ProductModuleSettingUIDURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}

	res, err := a.blls.Setting.DeleteGroup(ctx, req.Product, req.Module, req.Setting, req.UID)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}
