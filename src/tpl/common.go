package tpl

import (
	"crypto/rand"
	"fmt"
	"regexp"
	"strings"

	"github.com/teambition/gear"
)

var validIDReg = regexp.MustCompile(`^[0-9A-Za-z._=-]{3,63}$`)
var validHIDReg = regexp.MustCompile(`^[0-9A-Za-z_=-]{24}$`)

// https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#dns-subdomain-names
var validNameReg = regexp.MustCompile(`^[0-9a-z][0-9a-z.-]{0,61}[0-9a-z]$`)

// Should be subset of DNS-1035 label
// https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#dns-label-names
var validLabelReg = regexp.MustCompile(`^[0-9a-z][0-9a-z-]{0,61}[0-9a-z]$`)

// RandUID for testing
func RandUID() string {
	return fmt.Sprintf("uid-%x", randBytes(8))
}

// RandName for testing
func RandName() string {
	return fmt.Sprintf("name-%x", randBytes(8))
}

// RandLabel for testing
func RandLabel() string {
	return fmt.Sprintf("label-%x", randBytes(8))
}

func randBytes(size int) []byte {
	b := make([]byte, size)
	if _, err := rand.Read(b); err != nil {
		panic("crypto-go: rand.Read() failed, " + err.Error())
	}
	return b
}

// ErrorResponseType 定义了标准的 API 接口错误时返回数据模型
type ErrorResponseType struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// SuccessResponseType 定义了标准的 API 接口成功时返回数据模型
type SuccessResponseType struct {
	TotalSize     int         `json:"totalSize,omitempty"`
	NextPageToken string      `json:"nextPageToken"`
	Result        interface{} `json:"result"`
}

// ResponseType ...
type ResponseType struct {
	ErrorResponseType
	SuccessResponseType
}

// BoolRes ...
type BoolRes struct {
	SuccessResponseType
	Result bool `json:"result"`
}

// UIDURL ...
type UIDURL struct {
	UID string `json:"uid" param:"uid"`
}

// Validate 实现 gear.BodyTemplate。
func (t *UIDURL) Validate() error {
	if !validIDReg.MatchString(t.UID) {
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

// UIDPaginationURL ...
type UIDPaginationURL struct {
	Pagination
	UID string `json:"uid" param:"uid"`
}

// Validate 实现 gear.BodyTemplate。
func (t *UIDPaginationURL) Validate() error {
	if !validIDReg.MatchString(t.UID) {
		return gear.ErrBadRequest.WithMsgf("invalid uid: %s", t.UID)
	}
	if err := t.Pagination.Validate(); err != nil {
		return err
	}
	return nil
}

// NameDescBody ...
type NameDescBody struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}

// Validate 实现 gear.BodyTemplate。
func (t *NameDescBody) Validate() error {
	if !validNameReg.MatchString(t.Name) {
		return gear.ErrBadRequest.WithMsgf("invalid name: %s", t.Name)
	}
	if len(t.Desc) > 1022 {
		return gear.ErrBadRequest.WithMsgf("desc too long: %d (<= 1022)", len(t.Desc))
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
		if !validIDReg.MatchString(uid) {
			return gear.ErrBadRequest.WithMsgf("invalid user: %s", uid)
		}
	}
	for _, uid := range t.Groups {
		if !validIDReg.MatchString(uid) {
			return gear.ErrBadRequest.WithMsgf("invalid group: %s", uid)
		}
	}
	if len(t.Value) > 255 {
		return gear.ErrBadRequest.WithMsgf("value too long: %d (<= 255)", len(t.Value))
	}
	return nil
}

// StringToSlice ...
func StringToSlice(s string) []string {
	if s == "" {
		return make([]string, 0)
	}
	return strings.Split(s, ",")
}

// StringSliceHas ...
func StringSliceHas(sl []string, v string) bool {
	for _, s := range sl {
		if v == s {
			return true
		}
	}
	return false
}
