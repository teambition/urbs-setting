package middleware

import (
	"time"

	otgo "github.com/open-trust/ot-go-lib"
	"github.com/teambition/gear"
	auth "github.com/teambition/gear-auth"
	authjwt "github.com/teambition/gear-auth/jwt"
	"github.com/teambition/urbs-setting/src/conf"
	"github.com/teambition/urbs-setting/src/logging"
)

func init() {
	// otgo.Debugging = logging.Logger // 开启 otgo debug 日志

	otConf := conf.Config.OpenTrust
	keys := conf.Config.AuthKeys
	if len(keys) > 0 {
		Auther = auth.New(authjwt.StrToKeys(keys...)...)
		Auther.JWT().SetExpiresIn(time.Minute * 10)
	}
	if err := otConf.OTID.Validate(); err == nil {
		otVerifier, err = otgo.NewVerifier(conf.Config.GlobalCtx, otConf.OTID, false,
			otConf.DomainPublicKeys...)
		if err != nil {
			logging.Panicf("Parse Open Trust config failed: %s", err)
		}
	}

	if otVerifier == nil && Auther == nil {
		logging.Warningf("`auth_keys` is empty, Auth middleware will not be executed.")
	}
}

var otVerifier *otgo.Verifier

// Auther 是基于 JWT 的身份验证，当 config.auth_keys 配置了才会启用
var Auther *auth.Auth

// Auth 验证请求者身份，如果验证失败，则返回 401 的 gear.HTTPError
func Auth(ctx *gear.Context) error {
	if otVerifier != nil {
		token := otgo.ExtractTokenFromHeader(ctx.Req.Header)
		if token == "" {
			return gear.ErrUnauthorized.WithMsg("invalid authorization token")
		}

		vid, err := otVerifier.ParseOTVID(token)
		if err != nil {
			if Auther != nil { // 兼容老的 jwt 验证
				return oldAuth(ctx)
			}
			return gear.ErrUnauthorized.WithMsg("authorization token verification failed")
		}

		logging.AccessLogger.SetTo(ctx, "otSub", vid.ID.String())
		return nil
	}
	return oldAuth(ctx)
}

func oldAuth(ctx *gear.Context) error {
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
