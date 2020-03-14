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

func createProduct(tt *TestTools) (product schema.Product, err error) {
	name := tpl.RandName()
	_, err = request.Post(fmt.Sprintf("%s/v1/products", tt.Host)).
		Set("Content-Type", "application/json").
		Send(tpl.NameDescBody{Name: name, Desc: name}).
		End()

	if err == nil {
		err = tt.DB.Where("name = ?", name).First(&product).Error
	}
	return
}

func TestProductAPIs(t *testing.T) {
	tt, cleanup := SetUpTestTools()
	defer cleanup()

	n1 := tpl.RandName()

	t.Run(`"POST /v1/products"`, func(t *testing.T) {
		t.Run("should work", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Post(fmt.Sprintf("%s/v1/products", tt.Host)).
				Set("Content-Type", "application/json").
				Send(tpl.NameDescBody{Name: n1, Desc: "test"}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			text, err := res.Text()
			assert.Nil(err)
			assert.True(strings.Contains(text, `"offline_at":null`))
			assert.True(strings.Contains(text, `"deleted_at":null`))

			json := tpl.ProductRes{}
			res.JSON(&json)
			assert.NotNil(json.Result)
			assert.Equal(n1, json.Result.Name)
			assert.Equal("test", json.Result.Desc)
			assert.True(json.Result.CreatedAt.UTC().Unix() > int64(0))
			assert.True(json.Result.UpdatedAt.UTC().Unix() > int64(0))
			assert.Nil(json.Result.OfflineAt)
			assert.Nil(json.Result.DeletedAt)
			assert.Equal(int64(0), json.Result.Status)
		})

		t.Run(`should return 409`, func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Post(fmt.Sprintf("%s/v1/products", tt.Host)).
				Set("Content-Type", "application/json").
				Send(tpl.NameDescBody{Name: n1, Desc: "test"}).
				End()
			assert.Nil(err)
			assert.Equal(409, res.StatusCode)
		})

		t.Run(`should return 400`, func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Post(fmt.Sprintf("%s/v1/products", tt.Host)).
				Set("Content-Type", "application/json").
				Send(tpl.NameDescBody{Name: "ab", Desc: "test"}).
				End()
			assert.Nil(err)
			assert.Equal(400, res.StatusCode)
		})
	})

	t.Run(`"GET /v1/products"`, func(t *testing.T) {
		t.Run("should work", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Get(fmt.Sprintf("%s/v1/products", tt.Host)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			text, err := res.Text()
			assert.Nil(err)
			assert.True(strings.Contains(text, n1))

			json := tpl.ProductsRes{}
			res.JSON(&json)
			assert.NotNil(json.Result)
			assert.True(len(json.Result) > 0)
		})
	})

	t.Run(`"DELETE /v1/products/:product"`, func(t *testing.T) {
		t.Run("should conflict before offline", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Delete(fmt.Sprintf("%s/v1/products/%s", tt.Host, n1)).
				End()
			assert.Nil(err)
			assert.Equal(409, res.StatusCode)
		})

		t.Run("should offline", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Put(fmt.Sprintf("%s/v1/products/%s:offline", tt.Host, n1)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.BoolRes{}
			res.JSON(&json)
			assert.True(json.Result)
		})

		t.Run("should work", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Delete(fmt.Sprintf("%s/v1/products/%s", tt.Host, n1)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.BoolRes{}
			res.JSON(&json)
			assert.True(json.Result)
		})

		t.Run(`should idempotent`, func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Delete(fmt.Sprintf("%s/v1/products/%s", tt.Host, n1)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.BoolRes{}
			res.JSON(&json)
			assert.False(json.Result)
		})
	})

	t.Run(`"PUT /products/:product+:offline"`, func(t *testing.T) {
		product, err := createProduct(tt)
		assert.Nil(t, err)

		label, err := createLabel(tt, product.Name)
		assert.Nil(t, err)

		module, err := createModule(tt, product.Name)
		assert.Nil(t, err)

		setting, err := createSetting(tt, product.Name, module.Name)
		assert.Nil(t, err)

		t.Run("should work", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Put(fmt.Sprintf("%s/v1/products/%s:offline", tt.Host, product.Name)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.BoolRes{}
			res.JSON(&json)
			assert.True(json.Result)
		})

		t.Run("should work idempotent", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Put(fmt.Sprintf("%s/v1/products/%s:offline", tt.Host, product.Name)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.BoolRes{}
			res.JSON(&json)
			assert.False(json.Result)
		})

		t.Run("product's resource should offline", func(t *testing.T) {
			assert := assert.New(t)

			assert.Nil(label.OfflineAt)
			l := label
			assert.Nil(tt.DB.First(&l).Error)
			assert.NotNil(l.OfflineAt)

			assert.Nil(module.OfflineAt)
			m := module
			assert.Nil(tt.DB.First(&m).Error)
			assert.NotNil(m.OfflineAt)

			assert.Nil(setting.OfflineAt)
			s := setting
			assert.Nil(tt.DB.First(&s).Error)
			assert.NotNil(s.OfflineAt)
			assert.True(true)
		})

		t.Run("should not effect other data", func(t *testing.T) {
			assert := assert.New(t)

			product1, err := createProduct(tt)
			assert.Nil(err)

			product2, err := createProduct(tt)
			assert.Nil(err)

			label1, err := createLabel(tt, product1.Name)
			assert.Nil(err)

			label2, err := createLabel(tt, product2.Name)
			assert.Nil(err)

			users, err := createUsers(tt, 10)
			assert.Nil(err)

			group, err := createGroup(tt)
			assert.Nil(err)

			res, err := request.Post(fmt.Sprintf("%s/v1/products/%s/labels/%s:assign", tt.Host, product1.Name, label1.Name)).
				Set("Content-Type", "application/json").
				Send(tpl.UsersGroupsBody{
					Users:  schema.GetUsersUID(users),
					Groups: []string{group.UID},
				}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			var count int64
			assert.Nil(tt.DB.Table(`user_label`).Where("label_id = ?", label1.ID).Count(&count).Error)
			assert.Equal(int64(10), count)

			assert.Nil(tt.DB.Table(`group_label`).Where("label_id = ?", label1.ID).Count(&count).Error)
			assert.Equal(int64(1), count)

			res, err = request.Post(fmt.Sprintf("%s/v1/products/%s/labels/%s:assign", tt.Host, product2.Name, label2.Name)).
				Set("Content-Type", "application/json").
				Send(tpl.UsersGroupsBody{
					Users:  schema.GetUsersUID(users),
					Groups: []string{group.UID},
				}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			assert.Nil(tt.DB.Table(`user_label`).Where("label_id = ?", label2.ID).Count(&count).Error)
			assert.Equal(int64(10), count)

			assert.Nil(tt.DB.Table(`group_label`).Where("label_id = ?", label2.ID).Count(&count).Error)
			assert.Equal(int64(1), count)

			res, err = request.Put(fmt.Sprintf("%s/v1/products/%s:offline", tt.Host, product1.Name)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			time.Sleep(time.Second * 2)

			assert.Nil(tt.DB.First(&product1).Error)
			assert.NotNil(product1.OfflineAt)

			assert.Nil(tt.DB.First(&label1).Error)
			assert.NotNil(label1.OfflineAt)

			assert.Nil(tt.DB.Table(`user_label`).Where("label_id = ?", label1.ID).Count(&count).Error)
			assert.Equal(int64(0), count)

			assert.Nil(tt.DB.Table(`group_label`).Where("label_id = ?", label1.ID).Count(&count).Error)
			assert.Equal(int64(0), count)

			assert.Nil(tt.DB.First(&product2).Error)
			assert.Nil(product2.OfflineAt)

			assert.Nil(tt.DB.First(&label2).Error)
			assert.Nil(label2.OfflineAt)

			assert.Nil(tt.DB.Table(`user_label`).Where("label_id = ?", label2.ID).Count(&count).Error)
			assert.Equal(int64(10), count)

			assert.Nil(tt.DB.Table(`group_label`).Where("label_id = ?", label2.ID).Count(&count).Error)
			assert.Equal(int64(1), count)
		})
	})
}
