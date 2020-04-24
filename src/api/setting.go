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
	req := tpl.HIDRuleHIDURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}

	settingID := service.HIDToID(req.HID, "setting")
	if settingID <= 0 {
		return gear.ErrBadRequest.WithMsgf("invalid setting hid: %s", req.HID)
	}

	ruleID := service.HIDToID(req.RuleHID, "setting_rule")
	if ruleID <= 0 {
		return gear.ErrBadRequest.WithMsgf("invalid setting_rule hid: %s", req.RuleHID)
	}

	body := tpl.SettingRuleBody{}
	if err := ctx.ParseBody(&body); err != nil {
		return err
	}

	res, err := a.blls.Setting.UpdateRule(ctx, settingID, ruleID, body)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}

// DeleteRule ..
func (a *Setting) DeleteRule(ctx *gear.Context) error {
	req := tpl.HIDRuleHIDURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}

	settingID := service.HIDToID(req.HID, "setting")
	if settingID <= 0 {
		return gear.ErrBadRequest.WithMsgf("invalid setting hid: %s", req.HID)
	}

	ruleID := service.HIDToID(req.RuleHID, "setting_rule")
	if ruleID <= 0 {
		return gear.ErrBadRequest.WithMsgf("invalid setting_rule hid: %s", req.RuleHID)
	}

	res, err := a.blls.Setting.DeleteRule(ctx, settingID, ruleID)
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
