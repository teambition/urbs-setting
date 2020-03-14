package service

import (
	"github.com/teambition/urbs-setting/src/conf"
	"github.com/teambition/urbs-setting/src/util"
)

func init() {
	hIDer[""] = util.NewHID([]byte(conf.Config.HIDKey))
	hIDer["label"] = util.NewHID([]byte("label" + conf.Config.HIDKey))
	hIDer["setting"] = util.NewHID([]byte("setting" + conf.Config.HIDKey))
}

// HIDer 全局 HID 转换器，目前仅支持 schema.Label,  schema.setting 的 ID 转换。
var hIDer = map[string]*util.HID{}

// IDToHID 把 int64 的 ID 转换为 string HID，如果对象无效或者 ID （int64 > 0）不合法，则返回空字符串。
func IDToHID(id int64, kind ...string) string {
	k := ""
	if len(kind) > 0 {
		k = kind[0]
	}
	if h, ok := hIDer[k]; ok {
		return h.ToHex(id)
	}
	return ""
}

// HIDToID 把 string 的 HID 转换为 int64 ID，如果 HID 字符串不合法或者对象不合法，则返回 0。
func HIDToID(hid string, kind ...string) int64 {
	k := ""
	if len(kind) > 0 {
		k = kind[0]
	}
	if h, ok := hIDer[k]; ok {
		if id := h.ToInt64(hid); id > 0 {
			return id
		}
	}
	return 0
}
