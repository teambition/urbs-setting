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

func createSetting(tt *TestTools, productName, moduleName string) (setting schema.Setting, err error) {
	name := tpl.RandName()
	_, err = request.Post(fmt.Sprintf("%s/v1/products/%s/modules/%s/settings", tt.Host, productName, moduleName)).
		Set("Content-Type", "application/json").
		Send(tpl.NameDescBody{Name: name, Desc: name}).
		End()

	var product schema.Product
	var module schema.Module
	if err == nil {
		err = tt.DB.Where("name = ?", productName).First(&product).Error
	}

	if err == nil {
		err = tt.DB.Where("product_id = ? and name = ?", product.ID, moduleName).First(&module).Error
	}

	if err == nil {
		err = tt.DB.Where("module_id = ? and name = ?", module.ID, name).First(&setting).Error
	}
	return
}

func TestSettingAPIs(t *testing.T) {
	tt, cleanup := SetUpTestTools()
	defer cleanup()

	product, err := createProduct(tt)
	assert.Nil(t, err)

	module, err := createModule(tt, product.Name)
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
			assert.False(strings.Contains(text, `"id"`))

			json := tpl.SettingInfoRes{}
			res.JSON(&json)
			assert.NotNil(json.Result)
			assert.Equal(n1, json.Result.Name)
			assert.Equal("test", json.Result.Desc)
			assert.Equal([]string{}, json.Result.Channels)
			assert.Equal([]string{}, json.Result.Clients)
			assert.Equal([]string{}, json.Result.Values)
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
			assert.False(strings.Contains(text, `"id"`))

			json := tpl.ModulesRes{}
			res.JSON(&json)
			assert.NotNil(json.Result)
			assert.True(len(json.Result) > 0)
		})
	})
}
