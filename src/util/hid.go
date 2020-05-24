package util

// util 模块不要引入其它内部模块
import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"math"
)

const maxInt64u = uint64(math.MaxInt64)

// HID 基于 HMAC 算法，将内部 int64 的 ID 与 base64 URL 字符串进行相互转换。API 接口不支持 int64 ID 参数。
type HID struct {
	key []byte
}

// NewHID 根据给定的秘钥生成 HID 实例。
func NewHID(key []byte) *HID {
	return &HID{key: key}
}

// ToHex 将内部 ID（大于0的 int64）转换成24位 base64 URL 字符串。
// 如果输入值 <= 0，则返回空字符串。
func (h *HID) ToHex(i int64) string {
	if i <= 0 {
		return ""
	}
	data := make([]byte, 8)
	binary.LittleEndian.PutUint64(data, uint64(i))

	hs := hmac.New(sha1.New, h.key)
	hs.Write(data)
	sum := hs.Sum(nil)
	hs.Reset()

	copy(sum[:8], data)
	return base64.URLEncoding.EncodeToString(sum[0:18])
}

// ToInt64 将合法的24位 base64 URL 字符串转换成内部 ID（大于0的 int64）。
// 如果输入值不合法，则返回0。
func (h *HID) ToInt64(s string) int64 {
	if s == "" {
		return 0
	}
	data, err := base64.URLEncoding.DecodeString(s)
	if len(data) != 18 || err != nil {
		return 0
	}
	x := binary.LittleEndian.Uint64(data[:8])
	if x > maxInt64u {
		return 0
	}
	hs := hmac.New(sha1.New, h.key)
	hs.Write(data[:8])
	sum := hs.Sum(nil)
	hs.Reset()
	if !bytes.Equal(data[8:18], sum[8:18]) {
		return 0
	}
	return int64(x)
}
