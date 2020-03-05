package middleware

import (
	"time"

	"github.com/teambition/gear"
	auth "github.com/teambition/gear-auth"
	authjwt "github.com/teambition/gear-auth/jwt"
	"github.com/teambition/urbs-setting/src/conf"
	"github.com/teambition/urbs-setting/src/logging"
)

func init() {
	keys := conf.Config.AuthKeys
	if len(keys) > 0 {
		Auther = auth.New(authjwt.StrToKeys(keys...)...)
		Auther.JWT().SetExpiresIn(time.Minute * 10)
	} else {
		logging.Warningf("`auth_keys` is empty, Auth middleware will not be executed.")
	}
}

// Auther 是基于 JWT 的身份验证，当 config.auth_keys 配置了才会启用
var Auther *auth.Auth

// Auth 验证请求者身份，如果验证失败，则返回 401 的 gear.HTTPError
func Auth(ctx *gear.Context) error {
	if Auther != nil {
		claims, err := Auther.FromCtx(ctx)
		if err != nil {
			return err
		}
		if sub, ok := claims.Subject(); ok {
			logging.AccessLogger.SetTo(ctx, "jwt_sub", sub)
		}
		if jti, ok := claims.JWTID(); ok {
			logging.AccessLogger.SetTo(ctx, "jwt_id", jti)
		}
	}
	return nil
}
