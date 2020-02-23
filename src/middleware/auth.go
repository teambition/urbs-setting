package middleware

import (
	"github.com/teambition/gear"
)

// Auth 验证请求者身份，如果验证失败，则返回 401 的 gear.HTTPError
func Auth(ctx *gear.Context) error {
	// TODO
	return nil
}
