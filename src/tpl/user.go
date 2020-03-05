package tpl

import (
	"github.com/teambition/gear"
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
		if !validIDNameReg.MatchString(uid) {
			return gear.ErrBadRequest.WithMsgf("invalid user: %s", uid)
		}
	}
	return nil
}
