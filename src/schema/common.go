package schema

import (
	"errors"

	"github.com/teambition/urbs-setting/src/conf"
	"github.com/teambition/urbs-setting/src/util"
)

func init() {
	HIDer = &hIDer{
		label:   util.NewHID([]byte("label" + conf.Config.HIDKey)),
		setting: util.NewHID([]byte("setting" + conf.Config.HIDKey)),
	}
}

// HIDer 全局 HID 转换器，目前仅支持 schema.Label,  schema.setting 的 ID 转换。
var HIDer *hIDer

type hIDer struct {
	label   *util.HID
	setting *util.HID
}

// HID 从对象的 ID 转换为 HID，如果对象无效或者 ID （int64 > 0）不合法，则返回空字符串。
func (h *hIDer) HID(obj interface{}) string {
	switch v := obj.(type) {
	case *Label:
		return h.label.ToHex(v.ID)
	case *Setting:
		return h.setting.ToHex(v.ID)
	default:
		return ""
	}
}

// PutHID 把 HID 设置到对象上，如果 HID 字符串不合法或者对象不合法，则返回错误。
func (h *hIDer) PutHID(obj interface{}, hid string) error {
	var id int64 = -1
	switch v := obj.(type) {
	case *Label:
		if id = h.label.ToInt64(hid); id > 0 {
			v.ID = id
			return nil
		}

	case *Setting:
		if id = h.setting.ToInt64(hid); id > 0 {
			v.ID = id
			return nil
		}
	}

	if id == 0 {
		return errors.New("invalid hid")
	}
	return errors.New("unrecognized struct")
}
