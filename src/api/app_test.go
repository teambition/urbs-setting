package api

import (
	"fmt"
	"testing"

	"github.com/DavidCai1993/request"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/teambition/gear"
	"github.com/teambition/urbs-setting/src/service"
	"github.com/teambition/urbs-setting/src/util"
)

type TestTools struct {
	DB   *gorm.DB
	App  *gear.App
	Host string
}

func SetUpTestTools() (tt *TestTools, cleanup func()) {
	tt = &TestTools{}
	tt.App = NewApp()
	srv := tt.App.Start()
	tt.Host = "http://" + srv.Addr().String()

	err := util.DigInvoke(func(sql *service.SQL) error {
		tt.DB = sql.DB
		return nil
	})
	if err != nil {
		panic(err)
	}

	return tt, func() {
		srv.Close()
	}
}

func TestApp(t *testing.T) {
	tt, cleanup := SetUpTestTools()
	defer cleanup()

	t.Run(`app should work`, func(t *testing.T) {
		assert := assert.New(t)

		res, err := request.Get(tt.Host).End()
		json := map[string]string{}
		res.JSON(&json)

		assert.Nil(err)
		assert.Equal(200, res.StatusCode)
		assert.Equal("urbs-setting", json["name"])
		assert.NotEqual("", json["version"])
		assert.NotEqual("", json["gitSHA1"])
		assert.NotEqual("", json["buildTime"])
	})

	t.Run(`"GET /version" should work`, func(t *testing.T) {
		assert := assert.New(t)

		res, err := request.Get(fmt.Sprintf("%s/version", tt.Host)).End()
		json := map[string]string{}
		res.JSON(&json)

		assert.Nil(err)
		assert.Equal(200, res.StatusCode)
		assert.Equal("urbs-setting", json["name"])
		assert.NotEqual("", json["version"])
		assert.NotEqual("", json["gitSHA1"])
		assert.NotEqual("", json["buildTime"])
	})
}
