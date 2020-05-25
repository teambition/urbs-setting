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
	res, err := request.Post(fmt.Sprintf("%s/v1/products", tt.Host)).
		Set("Content-Type", "application/json").
		Send(tpl.NameDescBody{Name: name, Desc: name}).
		End()

	if err == nil {
		res.Content() // close http client
		_, err = tt.DB.ScanStruct(&product, "select * from `urbs_product` where `name` = ? limit 1", name)
	}
	return
}

func TestProductAPIs(t *testing.T) {
	tt, cleanup := SetUpTestTools()
	defer cleanup()

	t.Run(`"POST /v1/products"`, func(t *testing.T) {
		n1 := tpl.RandName()

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
			assert.True(strings.Contains(text, `"offlineAt":null`))
			assert.True(strings.Contains(text, `"deletedAt":null`))
			assert.False(strings.Contains(text, `"id"`))

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
			res.Content() // close http client
		})

		t.Run(`should return 400`, func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Post(fmt.Sprintf("%s/v1/products", tt.Host)).
				Set("Content-Type", "application/json").
				Send(tpl.NameDescBody{Name: "a", Desc: "test"}).
				End()
			assert.Nil(err)
			assert.Equal(400, res.StatusCode)
			res.Content() // close http client
		})
	})

	t.Run(`"GET /v1/products"`, func(t *testing.T) {
		product, err := createProduct(tt)
		assert.Nil(t, err)

		t.Run("should work", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Get(fmt.Sprintf("%s/v1/products?pageSize=1000", tt.Host)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			text, err := res.Text()
			assert.Nil(err)
			assert.True(strings.Contains(text, product.Name))
			assert.False(strings.Contains(text, `"id"`))

			json := tpl.ProductsRes{}
			res.JSON(&json)
			assert.NotNil(json.Result)
			assert.True(len(json.Result) > 0)
		})
	})

	t.Run(`"DELETE /v1/products/:product"`, func(t *testing.T) {
		product, err := createProduct(tt)
		assert.Nil(t, err)

		t.Run("should conflict before offline", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Delete(fmt.Sprintf("%s/v1/products/%s", tt.Host, product.Name)).
				End()
			assert.Nil(err)
			assert.Equal(409, res.StatusCode)
			res.Content() // close http client
		})

		t.Run("should offline", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Put(fmt.Sprintf("%s/v1/products/%s:offline", tt.Host, product.Name)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.BoolRes{}
			res.JSON(&json)
			assert.True(json.Result)
		})

		t.Run("should work", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Delete(fmt.Sprintf("%s/v1/products/%s", tt.Host, product.Name)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.BoolRes{}
			res.JSON(&json)
			assert.True(json.Result)
		})

		t.Run(`should idempotent`, func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Delete(fmt.Sprintf("%s/v1/products/%s", tt.Host, product.Name)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.BoolRes{}
			res.JSON(&json)
			assert.False(json.Result)
		})
	})

	t.Run(`"PUT /v1/products/:product"`, func(t *testing.T) {
		product, err := createProduct(tt)
		assert.Nil(t, err)

		t.Run("should work", func(t *testing.T) {
			assert := assert.New(t)

			desc := "abc"
			res, err := request.Put(fmt.Sprintf("%s/v1/products/%s", tt.Host, product.Name)).
				Set("Content-Type", "application/json").
				Send(tpl.ProductUpdateBody{
					Desc: &desc,
				}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			text, err := res.Text()
			assert.Nil(err)
			assert.True(strings.Contains(text, `"offlineAt":null`))
			assert.True(strings.Contains(text, `"deletedAt":null`))
			assert.False(strings.Contains(text, `"id"`))

			json := tpl.ProductRes{}
			res.JSON(&json)
			assert.NotNil(json.Result)
			assert.Equal(product.Name, json.Result.Name)
			assert.Equal(desc, json.Result.Desc)
			assert.True(json.Result.UpdatedAt.After(json.Result.CreatedAt))
			assert.Nil(json.Result.OfflineAt)
			assert.Nil(json.Result.DeletedAt)

			// should work idempotent
			time.Sleep(time.Millisecond * 100)
			res, err = request.Put(fmt.Sprintf("%s/v1/products/%s", tt.Host, product.Name)).
				Set("Content-Type", "application/json").
				Send(tpl.ProductUpdateBody{
					Desc: &desc,
				}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json2 := tpl.ProductRes{}
			res.JSON(&json2)
			assert.NotNil(json.Result)
			assert.True(json2.Result.UpdatedAt.Equal(json.Result.UpdatedAt))
		})

		t.Run("should 400", func(t *testing.T) {
			assert := assert.New(t)

			res, _ := request.Put(fmt.Sprintf("%s/v1/products/%s", tt.Host, product.Name)).
				Set("Content-Type", "application/json").
				Send(tpl.ProductUpdateBody{
					Desc: nil,
				}).
				End()
			assert.Equal(400, res.StatusCode)
			res.Content() // close http client
		})
	})

	t.Run(`"GET /v1/products/:product/statistics"`, func(t *testing.T) {
		product, err := createProduct(tt)
		assert.Nil(t, err)

		_, err = createLabel(tt, product.Name)
		assert.Nil(t, err)

		module, err := createModule(tt, product.Name)
		assert.Nil(t, err)

		_, err = createSetting(tt, product.Name, module.Name)
		assert.Nil(t, err)

		t.Run("should work", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Get(fmt.Sprintf("%s/v1/products/%s/statistics", tt.Host, product.Name)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.ProductStatisticsRes{}
			res.JSON(&json)
			assert.True(json.Result.Labels > 0)
			assert.True(json.Result.Modules > 0)
			assert.True(json.Result.Settings > 0)
			assert.True(json.Result.Release >= 0)
			assert.True(json.Result.Status >= 0)
		})
	})

	t.Run(`"PUT /v1/products/:product+:offline"`, func(t *testing.T) {
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
			_, err = tt.DB.ScanStruct(&l, "select * from `urbs_label` where `id` = ? limit 1", label.ID)
			assert.Nil(err)
			assert.NotNil(l.OfflineAt)

			assert.Nil(module.OfflineAt)
			m := module
			_, err = tt.DB.ScanStruct(&m, "select * from `urbs_module` where `id` = ? limit 1", module.ID)
			assert.Nil(err)
			assert.NotNil(m.OfflineAt)

			assert.Nil(setting.OfflineAt)
			s := setting
			_, err = tt.DB.ScanStruct(&s, "select * from `urbs_setting` where `id` = ? limit 1", setting.ID)
			assert.Nil(err)
			assert.NotNil(s.OfflineAt)
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
			res.Content() // close http client

			var count int64
			_, err = tt.DB.ScanVal(&count, "select count(*) from `user_label` where `label_id` = ?", label1.ID)
			assert.Nil(err)
			assert.Equal(int64(10), count)

			_, err = tt.DB.ScanVal(&count, "select count(*) from `group_label` where `label_id` = ?", label1.ID)
			assert.Nil(err)
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
			res.Content() // close http client

			_, err = tt.DB.ScanVal(&count, "select count(*) from `user_label` where `label_id` = ?", label2.ID)
			assert.Nil(err)
			assert.Equal(int64(10), count)

			_, err = tt.DB.ScanVal(&count, "select count(*) from `group_label` where `label_id` = ?", label2.ID)
			assert.Nil(err)
			assert.Equal(int64(1), count)

			res, err = request.Put(fmt.Sprintf("%s/v1/products/%s:offline", tt.Host, product1.Name)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)
			res.Content() // close http client

			time.Sleep(time.Second * 2)

			_, err = tt.DB.ScanStruct(&product1, "select * from `urbs_product` where `id` = ? limit 1", product1.ID)
			assert.Nil(err)
			assert.NotNil(product1.OfflineAt)

			_, err = tt.DB.ScanStruct(&label1, "select * from `urbs_label` where `id` = ? limit 1", label1.ID)
			assert.Nil(err)
			assert.NotNil(label1.OfflineAt)

			_, err = tt.DB.ScanVal(&count, "select count(*) from `user_label` where `label_id` = ?", label1.ID)
			assert.Nil(err)
			assert.Equal(int64(0), count)

			_, err = tt.DB.ScanVal(&count, "select count(*) from `group_label` where `label_id` = ?", label1.ID)
			assert.Nil(err)
			assert.Equal(int64(0), count)

			_, err = tt.DB.ScanStruct(&product2, "select * from `urbs_product` where `id` = ? limit 1", product2.ID)
			assert.Nil(err)
			assert.Nil(product2.OfflineAt)

			_, err = tt.DB.ScanStruct(&label2, "select * from `urbs_label` where `id` = ? limit 1", label2.ID)
			assert.Nil(err)
			assert.Nil(label2.OfflineAt)

			_, err = tt.DB.ScanVal(&count, "select count(*) from `user_label` where `label_id` = ?", label2.ID)
			assert.Nil(err)
			assert.Equal(int64(10), count)

			_, err = tt.DB.ScanVal(&count, "select count(*) from `group_label` where `label_id` = ?", label2.ID)
			assert.Nil(err)
			assert.Equal(int64(1), count)
		})
	})
}
