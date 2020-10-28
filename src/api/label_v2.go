package api

import (
	"github.com/teambition/gear"
	"github.com/teambition/urbs-setting/src/tpl"
)

// AssignV2 ..
func (a *Label) AssignV2(ctx *gear.Context) error {
	req := tpl.ProductLabelURL{}
	if err := ctx.ParseURL(&req); err != nil {
		return err
	}

	body := tpl.UsersGroupsBodyV2{}
	if err := ctx.ParseBody(&body); err != nil {
		return err
	}

	res, err := a.blls.Label.Assign(ctx, req.Product, req.Label, body.Users, body.Groups)
	if err != nil {
		return err
	}
	return ctx.OkJSON(tpl.LabelReleaseInfoRes{Result: *res})
}
