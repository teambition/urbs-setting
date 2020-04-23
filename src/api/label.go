package api

import (
	"github.com/teambition/gear"
	"github.com/teambition/urbs-setting/src/bll"
	"github.com/teambition/urbs-setting/src/service"
	"github.com/teambition/urbs-setting/src/tpl"
)

// Label ..
type Label struct {
	blls *bll.Blls
}

// List ..
func (a *Label) List(ctx *gear.Context) error {
	req := tpl.ProductPaginationURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}
	res, err := a.blls.Label.List(ctx, req.Product, req.Pagination)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}

// Create ..
func (a *Label) Create(ctx *gear.Context) error {
	req := tpl.ProductURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}

	body := tpl.LabelBody{}
	if err := ctx.ParseBody(&body); err != nil {
		return err
	}

	res, err := a.blls.Label.Create(ctx, req.Product, body.Name, body.Desc)
	if err != nil {
		return err
	}

	return ctx.OkJSON(res)
}

// Update ..
func (a *Label) Update(ctx *gear.Context) error {
	req := tpl.ProductLabelURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}

	body := tpl.LabelUpdateBody{}
	if err := ctx.ParseBody(&body); err != nil {
		return err
	}

	res, err := a.blls.Label.Update(ctx, req.Product, req.Label, body)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}

// Offline ..
func (a *Label) Offline(ctx *gear.Context) error {
	req := tpl.ProductLabelURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}
	res, err := a.blls.Label.Offline(ctx, req.Product, req.Label)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}

// Assign ..
func (a *Label) Assign(ctx *gear.Context) error {
	req := tpl.ProductLabelURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}

	body := tpl.UsersGroupsBody{}
	if err := ctx.ParseBody(&body); err != nil {
		return err
	}

	res, err := a.blls.Label.Assign(ctx, req.Product, req.Label, body.Users, body.Groups)
	if err != nil {
		return err
	}
	return ctx.OkJSON(tpl.LabelReleaseInfoRes{Result: *res})
}

// Delete ..
func (a *Label) Delete(ctx *gear.Context) error {
	req := tpl.ProductLabelURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}
	res, err := a.blls.Label.Delete(ctx, req.Product, req.Label)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}

// Recall ..
func (a *Label) Recall(ctx *gear.Context) error {
	req := tpl.ProductLabelURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}

	body := tpl.RecallBody{}
	if err := ctx.ParseBody(&body); err != nil {
		return err
	}

	res, err := a.blls.Label.Recall(ctx, req.Product, req.Label, body.Release)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}

// CreateRule ..
func (a *Label) CreateRule(ctx *gear.Context) error {
	req := tpl.ProductLabelURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}

	body := tpl.LabelRuleBody{}
	if err := ctx.ParseBody(&body); err != nil {
		return err
	}

	res, err := a.blls.Label.CreateRule(ctx, req.Product, req.Label, body)
	if err != nil {
		return err
	}

	return ctx.OkJSON(res)
}

// ListRules ..
func (a *Label) ListRules(ctx *gear.Context) error {
	req := tpl.ProductLabelURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}
	res, err := a.blls.Label.ListRules(ctx, req.Product, req.Label)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}

// UpdateRule ..
func (a *Label) UpdateRule(ctx *gear.Context) error {
	req := tpl.HIDRuleHIDURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}

	labelID := service.HIDToID(req.HID, "label")
	if labelID <= 0 {
		return gear.ErrBadRequest.WithMsgf("invalid label hid: %s", req.HID)
	}

	ruleID := service.HIDToID(req.RuleHID, "label_rule")
	if ruleID <= 0 {
		return gear.ErrBadRequest.WithMsgf("invalid label_rule hid: %s", req.RuleHID)
	}

	body := tpl.LabelRuleBody{}
	if err := ctx.ParseBody(&body); err != nil {
		return err
	}

	res, err := a.blls.Label.UpdateRule(ctx, labelID, ruleID, body)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}

// DeleteRule ..
func (a *Label) DeleteRule(ctx *gear.Context) error {
	req := tpl.HIDRuleHIDURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}

	labelID := service.HIDToID(req.HID, "label")
	if labelID <= 0 {
		return gear.ErrBadRequest.WithMsgf("invalid label hid: %s", req.HID)
	}

	ruleID := service.HIDToID(req.RuleHID, "label_rule")
	if ruleID <= 0 {
		return gear.ErrBadRequest.WithMsgf("invalid label_rule hid: %s", req.RuleHID)
	}

	res, err := a.blls.Label.DeleteRule(ctx, labelID, ruleID)
	if err != nil {
		return err
	}
	return ctx.OkJSON(res)
}
