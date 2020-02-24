package tpl

import (
	"regexp"

	"github.com/teambition/gear"
)

var validIDNameReg = regexp.MustCompile(`^[0-9A-Za-z._-]{3,63}$`)

// ResponseType 定义了标准的 List 接口返回数据模型
type ResponseType struct {
	Error         string      `json:"error,omitempty"`
	Message       string      `json:"message,omitempty"`
	NextPageToken string      `json:"nextPageToken,omitempty"`
	TotalSize     uint64      `json:"totalSize,omitempty"`
	Result        interface{} `json:"result,omitempty"`
}

// UIDHIDReq ...
type UIDHIDReq struct {
	UID string `json:"uid" param:"uid"`
	HID string `json:"hid" param:"hid"`
}

// Validate 实现 gear.BodyTemplate。
func (t *UIDHIDReq) Validate() error {
	if !validIDNameReg.MatchString(t.UID) {
		return gear.ErrBadRequest.WithMsgf("invalid uid: %s", t.UID)
	}
	return nil
}
