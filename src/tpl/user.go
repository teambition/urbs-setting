package tpl

import (
	"github.com/teambition/gear"
)

// BatchAddUsers ...
type BatchAddUsers struct {
	Users []AddUser `json:"users"`
}

// AddUser ...
type AddUser struct {
	UID string `json:"uid"`
}

// Validate 实现 gear.BodyTemplate。
func (t *BatchAddUsers) Validate() error {
	if len(t.Users) == 0 {
		return gear.ErrBadRequest.WithMsg("users emtpy")
	}
	for _, u := range t.Users {
		if !validIDNameReg.MatchString(u.UID) {
			return gear.ErrBadRequest.WithMsgf("invalid uid: %s", u.UID)
		}
	}
	return nil
}
