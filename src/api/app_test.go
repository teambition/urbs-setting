package api

import (
	"fmt"
	"os"
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

func TestMain(m *testing.M) {
	tt, cleanup := SetUpTestTools()
	tt.DB.Exec("TRUNCATE TABLE urbs_user;")
	tt.DB.Exec("TRUNCATE TABLE urbs_group;")
	tt.DB.Exec("TRUNCATE TABLE urbs_product;")
	tt.DB.Exec("TRUNCATE TABLE urbs_label;")
	tt.DB.Exec("TRUNCATE TABLE urbs_module;")
	tt.DB.Exec("TRUNCATE TABLE urbs_setting;")
	tt.DB.Exec("TRUNCATE TABLE user_group;")
	tt.DB.Exec("TRUNCATE TABLE user_label;")
	tt.DB.Exec("TRUNCATE TABLE user_setting;")
	tt.DB.Exec("TRUNCATE TABLE group_label;")
	tt.DB.Exec("TRUNCATE TABLE group_setting;")
	cleanup()
	os.Exit(m.Run())
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
