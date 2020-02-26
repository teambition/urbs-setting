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

func createLabel(appHost, productName string) (*schema.Label, error) {
	name := tpl.RandName()
	res, err := request.Post(fmt.Sprintf("%s/v1/products/%s/labels", appHost, productName)).
		Set("Content-Type", "application/json").
		Send(tpl.NameDescBody{Name: name, Desc: name}).
		End()

	if err != nil {
		return nil, err
	}

	json := tpl.LabelRes{}
	res.JSON(&json)
	return &json.Result, nil
}

func TestLabelAPIs(t *testing.T) {
	tt, cleanup := SetUpTestTools()
	defer cleanup()

	product, err := createProduct(tt.Host)
	assert.Nil(t, err)

	n1 := tpl.RandName()

	t.Run(`"POST /products/:product/labels"`, func(t *testing.T) {
		t.Run("should work", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Post(fmt.Sprintf("%s/v1/products/%s/labels", tt.Host, product.Name)).
				Set("Content-Type", "application/json").
				Send(tpl.NameDescBody{Name: n1, Desc: "test"}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			text, err := res.Text()
			assert.Nil(err)
			assert.True(strings.Contains(text, `"offline_at":null`))

			json := tpl.LabelRes{}
			res.JSON(&json)
			assert.NotNil(json.Result)
			assert.Equal(n1, json.Result.Name)
			assert.Equal("test", json.Result.Desc)
			assert.Equal("", json.Result.Channels)
			assert.Equal("", json.Result.Clients)
			assert.True(json.Result.CreatedAt.UTC().Unix() > int64(0))
			assert.True(json.Result.UpdatedAt.UTC().Unix() > int64(0))
			assert.Nil(json.Result.OfflineAt)
			assert.Equal(int64(0), json.Result.Status)
		})

		t.Run(`should return 409`, func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Post(fmt.Sprintf("%s/v1/products/%s/labels", tt.Host, product.Name)).
				Set("Content-Type", "application/json").
				Send(tpl.NameDescBody{Name: n1, Desc: "test"}).
				End()
			assert.Nil(err)
			assert.Equal(409, res.StatusCode)
		})

		t.Run(`should return 400`, func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Post(fmt.Sprintf("%s/v1/products/%s/labels", tt.Host, product.Name)).
				Set("Content-Type", "application/json").
				Send(tpl.NameDescBody{Name: "ab", Desc: "test"}).
				End()
			assert.Nil(err)
			assert.Equal(400, res.StatusCode)
		})
	})

	t.Run(`"GET /products/:product/labels"`, func(t *testing.T) {
		t.Run("should work", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Get(fmt.Sprintf("%s/v1/products/%s/labels", tt.Host, product.Name)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			text, err := res.Text()
			assert.Nil(err)
			assert.True(strings.Contains(text, n1))

			json := tpl.LabelsRes{}
			res.JSON(&json)
			assert.NotNil(json.Result)
			assert.True(len(json.Result) > 0)
		})
	})

	t.Run(`"DELETE /v1/products/:product/labels/:label"`, func(t *testing.T) {
		t.Run("should conflict before offline", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Delete(fmt.Sprintf("%s/v1/products/%s/labels/%s", tt.Host, product.Name, n1)).
				End()
			assert.Nil(err)
			assert.Equal(409, res.StatusCode)
		})

		t.Run("should offline", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Put(fmt.Sprintf("%s/v1/products/%s/labels/%s:offline", tt.Host, product.Name, n1)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.BoolRes{}
			res.JSON(&json)
			assert.True(json.Result)
		})

		t.Run("should work", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Delete(fmt.Sprintf("%s/v1/products/%s/labels/%s", tt.Host, product.Name, n1)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.BoolRes{}
			res.JSON(&json)
			assert.True(json.Result)
		})

		t.Run(`should work idempotent`, func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Delete(fmt.Sprintf("%s/v1/products/%s/labels/%s", tt.Host, product.Name, n1)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.BoolRes{}
			res.JSON(&json)
			assert.False(json.Result)
		})
	})

	t.Run(`"PUT /products/:product/labels/:label+:offline"`, func(t *testing.T) {
		label, err := createLabel(tt.Host, product.Name)
		assert.Nil(t, err)

		t.Run("should work", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Put(fmt.Sprintf("%s/v1/products/%s/labels/%s:offline", tt.Host, product.Name, label.Name)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.BoolRes{}
			res.JSON(&json)
			assert.True(json.Result)

			l := *label
			assert.Nil(tt.DB.First(&l).Error)
			assert.NotNil(l.OfflineAt)
		})

		t.Run("should work idempotent", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Put(fmt.Sprintf("%s/v1/products/%s/labels/%s:offline", tt.Host, product.Name, label.Name)).
				End()
			assert.Nil(err)

			assert.Equal(200, res.StatusCode)

			json := tpl.BoolRes{}
			res.JSON(&json)
			assert.False(json.Result)

			l := *label
			assert.Nil(tt.DB.First(&l).Error)
			assert.NotNil(l.OfflineAt)
		})
	})

	t.Run(`POST "/products/:product/labels/:label+:assign"`, func(t *testing.T) {
		label, err := createLabel(tt.Host, product.Name)
		assert.Nil(t, err)

		users, err := createUsers(tt.Host, 3)
		assert.Nil(t, err)

		group, err := createGroup(tt.Host)
		assert.Nil(t, err)

		t.Run("should work", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Post(fmt.Sprintf("%s/v1/products/%s/labels/%s:assign", tt.Host, product.Name, label.Name)).
				Set("Content-Type", "application/json").
				Send(tpl.UsersGroupsBody{
					Users:  users[0:2],
					Groups: []string{group},
				}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.BoolRes{}
			res.JSON(&json)
			assert.True(json.Result)

			var count int64
			assert.Nil(tt.DB.Table(`user_label`).Where("label_id = ?", label.ID).Count(&count).Error)
			assert.Equal(int64(2), count)

			assert.Nil(tt.DB.Table(`group_label`).Where("label_id = ?", label.ID).Count(&count).Error)
			assert.Equal(int64(1), count)
		})

		t.Run("should work with duplicate data", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Post(fmt.Sprintf("%s/v1/products/%s/labels/%s:assign", tt.Host, product.Name, label.Name)).
				Set("Content-Type", "application/json").
				Send(tpl.UsersGroupsBody{
					Users: []string{users[0], users[2]},
				}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.BoolRes{}
			res.JSON(&json)
			assert.True(json.Result)

			var count int64
			assert.Nil(tt.DB.Table(`user_label`).Where("label_id = ?", label.ID).Count(&count).Error)
			assert.Equal(int64(3), count)

			assert.Nil(tt.DB.Table(`group_label`).Where("label_id = ?", label.ID).Count(&count).Error)
			assert.Equal(int64(1), count)
		})
	})
}
