package tpl

import (
	"github.com/teambition/gear"
	"github.com/teambition/urbs-setting/src/schema"
)

// UsersBody ...
type UsersBody struct {
	Users []string `json:"users"`
}

// Validate 实现 gear.BodyTemplate。
func (t *UsersBody) Validate() error {
	if len(t.Users) == 0 {
		return gear.ErrBadRequest.WithMsg("users emtpy")
	}
	for _, uid := range t.Users {
		if !validIDReg.MatchString(uid) {
			return gear.ErrBadRequest.WithMsgf("invalid user: %s", uid)
		}
	}
	return nil
}

// UsersRes ...
type UsersRes struct {
	SuccessResponseType
	Result []schema.User `json:"result"`
}

// UserRes ...
type UserRes struct {
	SuccessResponseType
	Result schema.User `json:"result"`
}

// ApplyRulesBody ...
type ApplyRulesBody struct {
	UsersBody
	Kind string `json:"kind"`
}

// Validate 实现 gear.BodyTemplate。
func (t *ApplyRulesBody) Validate() error {
	if err := t.UsersBody.Validate(); err != nil {
		return err
	}
	if t.Kind == "" || !StringSliceHas(schema.RuleKinds, t.Kind) {
		return gear.ErrBadRequest.WithMsgf("invalid kind: %s", t.Kind)
	}
	return nil
}
