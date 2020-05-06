package api

import (
	"log"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/teambition/gear"

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

	app.Set(gear.SetParseError, func(err error) gear.HTTPError {
		msg := err.Error()

		if gorm.IsRecordNotFoundError(err) {
			return gear.ErrNotFound.WithMsg(msg)
		}
		if strings.Contains(msg, "Error 1062: Duplicate") {
			return gear.ErrConflict.WithMsg(msg)
		}

		return gear.ParseError(err)
	})

	// used for health check, so ingore logger
	app.Use(func(ctx *gear.Context) error {
		if ctx.Path == "/" || ctx.Path == "/version" {
			return ctx.OkJSON(GetVersion())
		}

		return nil
	})

	if app.Env() != "test" {
		app.UseHandler(logging.AccessLogger)
	}

	err := util.DigInvoke(func(routers []*gear.Router) error {
		for _, router := range routers {
			app.UseHandler(router)
		}
		return nil
	})

	if err != nil {
		logging.Panicf("DigInvoke error: %v", err)
	}

	return app
}
