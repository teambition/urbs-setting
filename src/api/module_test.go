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

func createModule(tt *TestTools, productName string) (module schema.Module, err error) {
	name := tpl.RandName()
	_, err = request.Post(fmt.Sprintf("%s/v1/products/%s/modules", tt.Host, productName)).
		Set("Content-Type", "application/json").
		Send(tpl.NameDescBody{Name: name, Desc: name}).
		End()

	var product schema.Product
	if err == nil {
		err = tt.DB.Where("name = ?", productName).First(&product).Error
	}

	if err == nil {
		err = tt.DB.Where("product_id = ? and name = ?", product.ID, name).First(&module).Error
	}
	return
}

func TestModuleAPIs(t *testing.T) {
	tt, cleanup := SetUpTestTools()
	defer cleanup()

	product, err := createProduct(tt)
	assert.Nil(t, err)

	n1 := tpl.RandName()

	t.Run(`"POST /v1/products/:product/modules"`, func(t *testing.T) {
		t.Run("should work", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Post(fmt.Sprintf("%s/v1/products/%s/modules", tt.Host, product.Name)).
				Set("Content-Type", "application/json").
				Send(tpl.NameDescBody{Name: n1, Desc: "test"}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			text, err := res.Text()
			assert.Nil(err)
			assert.True(strings.Contains(text, `"offline_at":null`))
			assert.False(strings.Contains(text, `"id"`))

			json := tpl.ModuleRes{}
			res.JSON(&json)
			assert.NotNil(json.Result)
			assert.Equal(n1, json.Result.Name)
			assert.Equal("test", json.Result.Desc)
			assert.True(json.Result.CreatedAt.UTC().Unix() > int64(0))
			assert.True(json.Result.UpdatedAt.UTC().Unix() > int64(0))
			assert.Nil(json.Result.OfflineAt)
			assert.Equal(int64(0), json.Result.Status)
		})

		t.Run(`should return 409`, func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Post(fmt.Sprintf("%s/v1/products/%s/modules", tt.Host, product.Name)).
				Set("Content-Type", "application/json").
				Send(tpl.NameDescBody{Name: n1, Desc: "test"}).
				End()
			assert.Nil(err)
			assert.Equal(409, res.StatusCode)
		})

		t.Run(`should return 400`, func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Post(fmt.Sprintf("%s/v1/products/%s/modules", tt.Host, product.Name)).
				Set("Content-Type", "application/json").
				Send(tpl.NameDescBody{Name: ".ab", Desc: "test"}).
				End()
			assert.Nil(err)
			assert.Equal(400, res.StatusCode)
		})
	})

	t.Run(`"GET /v1/products/:product/modules"`, func(t *testing.T) {
		t.Run("should work", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Get(fmt.Sprintf("%s/v1/products/%s/modules", tt.Host, product.Name)).
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

	t.Run(`"PUT /v1/products/:product/modules/:module"`, func(t *testing.T) {
		product, err := createProduct(tt)
		assert.Nil(t, err)

		module, err := createModule(tt, product.Name)
		assert.Nil(t, err)

		t.Run("should work", func(t *testing.T) {
			assert := assert.New(t)

			desc := "abc"
			res, err := request.Put(fmt.Sprintf("%s/v1/products/%s/modules/%s", tt.Host, product.Name, module.Name)).
				Set("Content-Type", "application/json").
				Send(tpl.ModuleUpdateBody{
					Desc: &desc,
				}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			text, err := res.Text()
			assert.Nil(err)
			assert.True(strings.Contains(text, `"offline_at":null`))
			assert.False(strings.Contains(text, `"id"`))

			json := tpl.ModuleRes{}
			res.JSON(&json)
			assert.NotNil(json.Result)
			assert.Equal(module.Name, json.Result.Name)
			assert.Equal(desc, json.Result.Desc)
			assert.True(json.Result.UpdatedAt.After(json.Result.CreatedAt))
			assert.Nil(json.Result.OfflineAt)

			// should work idempotent
			time.Sleep(time.Millisecond * 100)
			res, err = request.Put(fmt.Sprintf("%s/v1/products/%s/modules/%s", tt.Host, product.Name, module.Name)).
				Set("Content-Type", "application/json").
				Send(tpl.ModuleUpdateBody{
					Desc: &desc,
				}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json2 := tpl.ModuleRes{}
			res.JSON(&json2)
			assert.NotNil(json.Result)
			assert.True(json2.Result.UpdatedAt.Equal(json.Result.UpdatedAt))
		})

		t.Run("should 400", func(t *testing.T) {
			assert := assert.New(t)

			res, _ := request.Put(fmt.Sprintf("%s/v1/products/%s/modules/%s", tt.Host, product.Name, module.Name)).
				Set("Content-Type", "application/json").
				Send(tpl.ModuleUpdateBody{
					Desc: nil,
				}).
				End()
			assert.Equal(400, res.StatusCode)
		})
	})

	t.Run(`"PUT /v1/products/:product/modules/:module+:offline"`, func(t *testing.T) {
		product, err := createProduct(tt)
		assert.Nil(t, err)

		module, err := createModule(tt, product.Name)
		assert.Nil(t, err)

		setting, err := createSetting(tt, product.Name, module.Name)
		assert.Nil(t, err)

		t.Run("should work", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Put(fmt.Sprintf("%s/v1/products/%s/modules/%s:offline", tt.Host, product.Name, module.Name)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.BoolRes{}
			res.JSON(&json)
			assert.True(json.Result)
		})

		t.Run("should work idempotent", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Put(fmt.Sprintf("%s/v1/products/%s/modules/%s:offline", tt.Host, product.Name, module.Name)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.BoolRes{}
			res.JSON(&json)
			assert.False(json.Result)
		})

		t.Run("module's resource should offline", func(t *testing.T) {
			assert := assert.New(t)

			assert.Nil(module.OfflineAt)
			m := module
			assert.Nil(tt.DB.First(&m).Error)
			assert.NotNil(m.OfflineAt)

			assert.Nil(setting.OfflineAt)
			s := setting
			assert.Nil(tt.DB.First(&s).Error)
			assert.NotNil(s.OfflineAt)
		})
	})
}
