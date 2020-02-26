package api

import (
	"fmt"
	"strings"
	"testing"

	"github.com/DavidCai1993/request"
	"github.com/stretchr/testify/assert"
	"github.com/teambition/urbs-setting/src/schema"
	"github.com/teambition/urbs-setting/src/tpl"
)

func createSetting(appHost, productName, moduleName string) (*schema.Setting, error) {
	name := tpl.RandName()
	res, err := request.Post(fmt.Sprintf("%s/v1/products/%s/modules/%s/settings", appHost, productName, moduleName)).
		Set("Content-Type", "application/json").
		Send(tpl.NameDescBody{Name: name, Desc: name}).
		End()

	if err != nil {
		return nil, err
	}

	json := tpl.SettingRes{}
	res.JSON(&json)
	return &json.Result, nil
}

func TestSettingAPIs(t *testing.T) {
	tt, cleanup := SetUpTestTools()
	defer cleanup()

	product, err := createProduct(tt.Host)
	assert.Nil(t, err)

	module, err := createModule(tt.Host, product.Name)
	assert.Nil(t, err)

	n1 := tpl.RandName()

	t.Run(`"POST /products/:product/modules/:module/settings"`, func(t *testing.T) {
		t.Run("should work", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Post(fmt.Sprintf("%s/v1/products/%s/modules/%s/settings", tt.Host, product.Name, module.Name)).
				Set("Content-Type", "application/json").
				Send(tpl.NameDescBody{Name: n1, Desc: "test"}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			text, err := res.Text()
			assert.Nil(err)
			assert.True(strings.Contains(text, `"offline_at":null`))

			json := tpl.SettingRes{}
			res.JSON(&json)
			assert.NotNil(json.Result)
			assert.Equal(n1, json.Result.Name)
			assert.Equal("test", json.Result.Desc)
			assert.Equal("", json.Result.Channels)
			assert.Equal("", json.Result.Clients)
			assert.Equal("", json.Result.Values)
			assert.True(json.Result.CreatedAt.UTC().Unix() > int64(0))
			assert.True(json.Result.UpdatedAt.UTC().Unix() > int64(0))
			assert.Nil(json.Result.OfflineAt)
			assert.Equal(int64(0), json.Result.Status)
		})

		t.Run(`should return 409`, func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Post(fmt.Sprintf("%s/v1/products/%s/modules/%s/settings", tt.Host, product.Name, module.Name)).
				Set("Content-Type", "application/json").
				Send(tpl.NameDescBody{Name: n1, Desc: "test"}).
				End()
			assert.Nil(err)
			assert.Equal(409, res.StatusCode)
		})

		t.Run(`should return 400`, func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Post(fmt.Sprintf("%s/v1/products/%s/modules/%s/settings", tt.Host, product.Name, module.Name)).
				Set("Content-Type", "application/json").
				Send(tpl.NameDescBody{Name: "ab", Desc: "test"}).
				End()
			assert.Nil(err)
			assert.Equal(400, res.StatusCode)
		})
	})

	t.Run(`"GET /products/:product/modules/:module/settings"`, func(t *testing.T) {
		t.Run("should work", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Get(fmt.Sprintf("%s/v1/products/%s/modules/%s/settings", tt.Host, product.Name, module.Name)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			text, err := res.Text()
			assert.Nil(err)
			assert.True(strings.Contains(text, n1))

			json := tpl.ModulesRes{}
			res.JSON(&json)
			assert.NotNil(json.Result)
			assert.True(len(json.Result) > 0)
		})
	})
}
