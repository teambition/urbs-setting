package tpl

import (
	"regexp"

	"github.com/teambition/gear"
)

var validIDNameReg = regexp.MustCompile(`^[0-9A-Za-z._-]{3,63}$`)
var validHIDReg = regexp.MustCompile(`^[0-9A-Za-z_=-]{24}$`)

// ResponseType 定义了标准的 List 接口返回数据模型
type ResponseType struct {
	Error         string      `json:"error,omitempty"`
	Message       string      `json:"message,omitempty"`
	NextPageToken string      `json:"nextPageToken,omitempty"`
	TotalSize     uint64      `json:"totalSize,omitempty"`
	Result        interface{} `json:"result,omitempty"`
}

// BoolRes ...
type BoolRes struct {
	ResponseType
	Result bool `json:"result"`
}

// UIDURL ...
type UIDURL struct {
	UID string `json:"uid" param:"uid"`
}

// Validate 实现 gear.BodyTemplate。
func (t *UIDURL) Validate() error {
	if !validIDNameReg.MatchString(t.UID) {
		return gear.ErrBadRequest.WithMsgf("invalid uid: %s", t.UID)
	}
	return nil
}

// UIDHIDURL ...
type UIDHIDURL struct {
	UIDURL
	HID string `json:"hid" param:"hid"`
}

// Validate 实现 gear.BodyTemplate。
func (t *UIDHIDURL) Validate() error {
	if err := t.UIDURL.Validate(); err != nil {
		return err
	}
	if !validHIDReg.MatchString(t.HID) {
		return gear.ErrBadRequest.WithMsgf("invalid hid: %s", t.HID)
	}
	return nil
}

// NameDescBody ...
type NameDescBody struct {
	Name string `json:"Name"`
	Desc string `json:"desc"`
}

// Validate 实现 gear.BodyTemplate。
func (t *NameDescBody) Validate() error {
	if !validIDNameReg.MatchString(t.Name) {
		return gear.ErrBadRequest.WithMsgf("invalid name: %s", t.Name)
	}
	if len(t.Desc) > 1023 {
		return gear.ErrBadRequest.WithMsgf("desc too long: %d (<= 1023)", len(t.Desc))
	}
	return nil
}

// UsersGroupsBody ...
type UsersGroupsBody struct {
	Users  []string `json:"users"`
	Groups []string `json:"groups"`
	Value  string   `json:"value"`
}

// Validate 实现 gear.BodyTemplate。
func (t *UsersGroupsBody) Validate() error {
	if len(t.Users) == 0 && len(t.Groups) == 0 {
		return gear.ErrBadRequest.WithMsg("users and groups are empty")
	}

	for _, uid := range t.Users {
		if !validIDNameReg.MatchString(uid) {
			return gear.ErrBadRequest.WithMsgf("invalid user: %s", uid)
		}
	}
	for _, uid := range t.Groups {
		if !validIDNameReg.MatchString(uid) {
			return gear.ErrBadRequest.WithMsgf("invalid group: %s", uid)
		}
	}
	if len(t.Value) > 255 {
		return gear.ErrBadRequest.WithMsgf("value too long: %d (<= 255)", len(t.Value))
	}
	return nil
}
