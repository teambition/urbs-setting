package api

import (
	"log"

	"github.com/teambition/gear"
	"github.com/teambition/gear/middleware/requestid"

	"github.com/teambition/urbs-setting/src/conf"
	"github.com/teambition/urbs-setting/src/logging"
	"github.com/teambition/urbs-setting/src/util"
)

// AppName 服务名
var AppName = "urbs-setting"

// AppVersion 服务版本
var AppVersion = "unknown"

// BuildTime 镜像生成时间
var BuildTime = "unknown"

// GitSHA1 镜像对应 git commit id
var GitSHA1 = "unknown"

// GetVersion ...
func GetVersion() map[string]string {
	return map[string]string{
		"name":      AppName,
		"version":   AppVersion,
		"buildTime": BuildTime,
		"gitSHA1":   GitSHA1,
	}
}

// NewApp ...
func NewApp() *gear.App {
	app := gear.New()

	app.Set(gear.SetTrustedProxy, true)
	app.Set(gear.SetBodyParser, gear.DefaultBodyParser(2<<22)) // 8MB
	// ignore TLS handshake error
	app.Set(gear.SetLogger, log.New(gear.DefaultFilterWriter(), "", 0))

	// used for health check, so ingore logger
	app.Use(func(ctx *gear.Context) error {
		if ctx.Path == "/" || ctx.Path == "/version" {
			return ctx.OkJSON(GetVersion())
		}

		return nil
	})

	app.Use(requestid.New())

	logging.SetLevel(conf.Config.Logger.Level)
	logging.Logger.SetJSONLog()
	app.UseHandler(logging.Logger)

	err := util.DigInvoke(func(routers []*gear.Router) error {
		for _, router := range routers {
			app.UseHandler(router)
		}
		return nil
	})

	if err != nil {
		logging.Panic(err)
	}

	return app
}
