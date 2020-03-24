package api

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/DavidCai1993/request"
	"github.com/stretchr/testify/assert"
	"github.com/teambition/urbs-setting/src/schema"
	"github.com/teambition/urbs-setting/src/tpl"
)

func createSetting(tt *TestTools, productName, moduleName string, values ...string) (setting schema.Setting, err error) {
	name := tpl.RandName()
	_, err = request.Post(fmt.Sprintf("%s/v1/products/%s/modules/%s/settings", tt.Host, productName, moduleName)).
		Set("Content-Type", "application/json").
		Send(tpl.NameDescBody{Name: name, Desc: name}).
		End()

	if err == nil && len(values) > 0 {
		_, err = request.Put(fmt.Sprintf("%s/v1/products/%s/modules/%s/settings/%s", tt.Host, productName, moduleName, name)).
			Set("Content-Type", "application/json").
			Send(tpl.SettingUpdateBody{Values: &values}).
			End()
	}

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

	t.Run(`"POST /v1/products/:product/modules/:module/settings"`, func(t *testing.T) {
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
				Send(tpl.NameDescBody{Name: "aB", Desc: "test"}).
				End()
			assert.Nil(err)
			assert.Equal(400, res.StatusCode)
		})
	})

	t.Run(`"GET /v1/products/:product/modules/:module/settings"`, func(t *testing.T) {
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

			json := tpl.SettingsInfoRes{}
			res.JSON(&json)
			assert.NotNil(json.Result)
			assert.True(json.TotalSize > 0)
			data := json.Result[0]
			assert.NotEqual("", data.HID)
			assert.NotEqual("", data.Name)
			assert.NotEqual("", data.Product)
			assert.NotEqual("", data.Module)
			assert.Equal([]string{}, data.Channels)
			assert.Equal([]string{}, data.Clients)
			assert.Equal([]string{}, data.Values)
			assert.True(data.CreatedAt.UTC().Unix() > int64(0))
			assert.True(data.UpdatedAt.UTC().Unix() > int64(0))
			assert.Nil(data.OfflineAt)
			assert.Equal(int64(0), data.Status)
		})
	})

	t.Run(`"GET /v1/products/:product/modules/:module/settings/:setting"`, func(t *testing.T) {
		t.Run("should work", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Get(fmt.Sprintf("%s/v1/products/%s/modules/%s/settings/%s",
				tt.Host, product.Name, module.Name, n1)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			text, err := res.Text()
			assert.Nil(err)
			assert.True(strings.Contains(text, n1))
			assert.False(strings.Contains(text, `"id"`))

			json := tpl.SettingInfoRes{}
			res.JSON(&json)
			assert.NotNil(json.Result)
			data := json.Result
			assert.NotEqual("", data.HID)
			assert.Equal(n1, data.Name)
			assert.Equal(product.Name, data.Product)
			assert.Equal(module.Name, data.Module)
			assert.Equal([]string{}, data.Channels)
			assert.Equal([]string{}, data.Clients)
			assert.Equal([]string{}, data.Values)
			assert.True(data.CreatedAt.UTC().Unix() > int64(0))
			assert.True(data.UpdatedAt.UTC().Unix() > int64(0))
			assert.Nil(data.OfflineAt)
			assert.Equal(int64(0), data.Status)
		})
	})

	t.Run(`"PUT /v1/products/:product/modules/:module/settings/:setting"`, func(t *testing.T) {
		product, err := createProduct(tt)
		assert.Nil(t, err)

		module, err := createModule(tt, product.Name)
		assert.Nil(t, err)

		setting, err := createSetting(tt, product.Name, module.Name, "a", "b")
		assert.Nil(t, err)

		t.Run("should work", func(t *testing.T) {
			assert := assert.New(t)

			desc := "abc"
			res, err := request.Put(fmt.Sprintf("%s/v1/products/%s/modules/%s/settings/%s", tt.Host, product.Name, module.Name, setting.Name)).
				Set("Content-Type", "application/json").
				Send(tpl.SettingUpdateBody{
					Desc: &desc,
				}).
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
			assert.Equal(setting.Name, json.Result.Name)
			assert.Equal(desc, json.Result.Desc)
			assert.True(json.Result.UpdatedAt.After(json.Result.CreatedAt))
			assert.Nil(json.Result.OfflineAt)

			// should work idempotent
			time.Sleep(time.Millisecond * 100)
			res, err = request.Put(fmt.Sprintf("%s/v1/products/%s/modules/%s/settings/%s", tt.Host, product.Name, module.Name, setting.Name)).
				Set("Content-Type", "application/json").
				Send(tpl.SettingUpdateBody{
					Desc: &desc,
				}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json2 := tpl.SettingInfoRes{}
			res.JSON(&json2)
			assert.NotNil(json.Result)
			assert.True(json2.Result.UpdatedAt.Equal(json.Result.UpdatedAt))
		})

		t.Run("should work with Channels", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Put(fmt.Sprintf("%s/v1/products/%s/modules/%s/settings/%s", tt.Host, product.Name, module.Name, setting.Name)).
				Set("Content-Type", "application/json").
				Send(tpl.SettingUpdateBody{
					Channels: &[]string{"stable"},
				}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.SettingInfoRes{}
			res.JSON(&json)
			assert.NotNil(json.Result)
			assert.Equal(setting.Name, json.Result.Name)
			assert.Equal([]string{"stable"}, json.Result.Channels)
		})

		t.Run("should work with Clients", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Put(fmt.Sprintf("%s/v1/products/%s/modules/%s/settings/%s", tt.Host, product.Name, module.Name, setting.Name)).
				Set("Content-Type", "application/json").
				Send(tpl.SettingUpdateBody{
					Channels: &[]string{"stable", "beta"},
					Clients:  &[]string{"ios"},
				}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.SettingInfoRes{}
			res.JSON(&json)
			assert.NotNil(json.Result)
			assert.Equal(setting.Name, json.Result.Name)
			assert.Equal([]string{"beta", "stable"}, json.Result.Channels)
			assert.Equal([]string{"ios"}, json.Result.Clients)
		})

		t.Run("should work with Values", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Put(fmt.Sprintf("%s/v1/products/%s/modules/%s/settings/%s", tt.Host, product.Name, module.Name, setting.Name)).
				Set("Content-Type", "application/json").
				Send(tpl.SettingUpdateBody{
					Clients: &[]string{"ios", "web"},
					Values:  &[]string{"b", "a", "c"},
				}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.SettingInfoRes{}
			res.JSON(&json)
			assert.NotNil(json.Result)
			assert.Equal(setting.Name, json.Result.Name)
			assert.Equal([]string{"beta", "stable"}, json.Result.Channels)
			assert.Equal([]string{"ios", "web"}, json.Result.Clients)
			assert.Equal([]string{"a", "b", "c"}, json.Result.Values)
		})

		t.Run("should 400", func(t *testing.T) {
			assert := assert.New(t)

			res, _ := request.Put(fmt.Sprintf("%s/v1/products/%s/modules/%s/settings/%s", tt.Host, product.Name, module.Name, setting.Name)).
				Set("Content-Type", "application/json").
				Send(tpl.SettingUpdateBody{
					Desc: nil,
				}).
				End()
			assert.Equal(400, res.StatusCode)
		})
	})

	t.Run(`"PUT /v1/products/:product/modules/:module/settings/:setting+:offline"`, func(t *testing.T) {
		product, err := createProduct(tt)
		assert.Nil(t, err)

		module, err := createModule(tt, product.Name)
		assert.Nil(t, err)

		setting, err := createSetting(tt, product.Name, module.Name)
		assert.Nil(t, err)

		t.Run("should work", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Put(fmt.Sprintf("%s/v1/products/%s/modules/%s/settings/%s:offline", tt.Host, product.Name, module.Name, setting.Name)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.BoolRes{}
			res.JSON(&json)
			assert.True(json.Result)
		})

		t.Run("should work idempotent", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Put(fmt.Sprintf("%s/v1/products/%s/modules/%s/settings/%s:offline", tt.Host, product.Name, module.Name, setting.Name)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.BoolRes{}
			res.JSON(&json)
			assert.False(json.Result)
		})

		t.Run("module's resource should offline", func(t *testing.T) {
			assert := assert.New(t)

			assert.Nil(setting.OfflineAt)
			s := setting
			assert.Nil(tt.DB.First(&s).Error)
			assert.NotNil(s.OfflineAt)
		})
	})
}
