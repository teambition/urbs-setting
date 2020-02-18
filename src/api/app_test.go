package api

import (
	"testing"

	"github.com/DavidCai1993/request"
	"github.com/stretchr/testify/assert"
)

func TestApp(t *testing.T) {
	t.Run("app should work", func(t *testing.T) {
		assert := assert.New(t)

		app := NewApp()
		srv := app.Start()

		// var res = new(map[string]string)

		res, err := request.Get("http://" + srv.Addr().String()).End()
		json := map[string]string{}
		res.JSON(&json)

		assert.Nil(err)
		assert.Equal(200, res.StatusCode)
		assert.Equal("urbs-setting", json["name"])
		assert.NotEqual("", json["version"])
		assert.NotEqual("", json["gitSHA1"])
		assert.NotEqual("", json["buildTime"])
		assert.Nil(app.Close())
	})
}
