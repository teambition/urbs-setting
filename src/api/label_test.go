package api

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/DavidCai1993/request"
	"github.com/stretchr/testify/assert"
	"github.com/teambition/urbs-setting/src/schema"
	"github.com/teambition/urbs-setting/src/service"
	"github.com/teambition/urbs-setting/src/tpl"
)

func createLabel(tt *TestTools, productName string) (label schema.Label, err error) {
	name := tpl.RandLabel()
	res, err := request.Post(fmt.Sprintf("%s/v1/products/%s/labels", tt.Host, productName)).
		Set("Content-Type", "application/json").
		Send(tpl.LabelBody{Name: name, Desc: name}).
		End()

	var product schema.Product
	if err == nil {
		res.Content() // close http client
		err = tt.DB.Where("name = ?", productName).First(&product).Error
	}

	if err == nil {
		err = tt.DB.Where("product_id = ? and name = ?", product.ID, name).First(&label).Error
	}
	return
}

func TestLabelAPIs(t *testing.T) {
	tt, cleanup := SetUpTestTools()
	defer cleanup()

	product, err := createProduct(tt)
	assert.Nil(t, err)

	n1 := tpl.RandLabel()

	t.Run(`"POST /v1/products/:product/labels"`, func(t *testing.T) {
		t.Run("should work", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Post(fmt.Sprintf("%s/v1/products/%s/labels", tt.Host, product.Name)).
				Set("Content-Type", "application/json").
				Send(tpl.LabelBody{Name: n1, Desc: "test"}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			text, err := res.Text()
			assert.Nil(err)
			assert.True(strings.Contains(text, `"offline_at":null`))
			assert.True(strings.Contains(text, `"hid"`))
			assert.False(strings.Contains(text, `"id"`))

			json := tpl.LabelInfoRes{}
			res.JSON(&json)
			assert.NotNil(json.Result)
			data := json.Result
			assert.NotEqual("", data.HID)
			assert.Equal(n1, data.Name)
			assert.Equal("test", data.Desc)
			assert.Equal(product.Name, data.Product)
			assert.Equal([]string{}, data.Channels)
			assert.Equal([]string{}, data.Clients)
			assert.True(data.CreatedAt.UTC().Unix() > int64(0))
			assert.True(data.UpdatedAt.UTC().Unix() > int64(0))
			assert.Nil(data.OfflineAt)
			assert.Equal(int64(0), data.Status)
		})

		t.Run(`should return 409`, func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Post(fmt.Sprintf("%s/v1/products/%s/labels", tt.Host, product.Name)).
				Set("Content-Type", "application/json").
				Send(tpl.LabelBody{Name: n1, Desc: "test"}).
				End()
			assert.Nil(err)
			assert.Equal(409, res.StatusCode)
			res.Content() // close http client
		})

		t.Run(`should return 400`, func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Post(fmt.Sprintf("%s/v1/products/%s/labels", tt.Host, product.Name)).
				Set("Content-Type", "application/json").
				Send(tpl.LabelBody{Name: "_abc", Desc: "test"}).
				End()
			assert.Nil(err)
			assert.Equal(400, res.StatusCode)
			res.Content() // close http client
		})
	})

	t.Run(`"GET /v1/products/:product/labels"`, func(t *testing.T) {
		t.Run("should work", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Get(fmt.Sprintf("%s/v1/products/%s/labels", tt.Host, product.Name)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			text, err := res.Text()
			assert.Nil(err)
			assert.True(strings.Contains(text, n1))
			assert.False(strings.Contains(text, `"id"`))

			json := tpl.LabelsInfoRes{}
			res.JSON(&json)
			assert.NotNil(json.Result)
			assert.True(json.TotalSize > 0)
			data := json.Result[0]
			assert.NotEqual("", data.HID)
			assert.NotEqual("", data.Name)
			assert.NotEqual("", data.Desc)
			assert.NotEqual("", data.Product)
			assert.Equal([]string{}, data.Channels)
			assert.Equal([]string{}, data.Clients)
			assert.True(data.CreatedAt.UTC().Unix() > int64(0))
			assert.True(data.UpdatedAt.UTC().Unix() > int64(0))
			assert.Nil(data.OfflineAt)
			assert.Equal(int64(0), data.Status)
		})
	})

	t.Run(`"PUT /v1/products/:product/labels/:label"`, func(t *testing.T) {
		product, err := createProduct(tt)
		assert.Nil(t, err)

		label, err := createLabel(tt, product.Name)
		assert.Nil(t, err)

		t.Run("should work", func(t *testing.T) {
			assert := assert.New(t)

			desc := "abc"
			res, err := request.Put(fmt.Sprintf("%s/v1/products/%s/labels/%s", tt.Host, product.Name, label.Name)).
				Set("Content-Type", "application/json").
				Send(tpl.LabelUpdateBody{
					Desc: &desc,
				}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			text, err := res.Text()
			assert.Nil(err)
			assert.True(strings.Contains(text, `"offline_at":null`))
			assert.False(strings.Contains(text, `"id"`))

			json := tpl.LabelInfoRes{}
			res.JSON(&json)
			assert.NotNil(json.Result)
			data := json.Result
			assert.Equal(service.IDToHID(label.ID, "label"), data.HID)
			assert.Equal(label.Name, data.Name)
			assert.Equal(desc, data.Desc)
			assert.Equal(product.Name, data.Product)
			assert.Equal([]string{}, data.Channels)
			assert.Equal([]string{}, data.Clients)
			assert.True(data.CreatedAt.UTC().Unix() > int64(0))
			assert.True(data.UpdatedAt.UTC().Unix() > int64(0))
			assert.Nil(data.OfflineAt)
			assert.Equal(int64(0), data.Status)

			// should work idempotent
			time.Sleep(time.Millisecond * 100)
			res, err = request.Put(fmt.Sprintf("%s/v1/products/%s/labels/%s", tt.Host, product.Name, label.Name)).
				Set("Content-Type", "application/json").
				Send(tpl.LabelUpdateBody{
					Desc: &desc,
				}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json2 := tpl.LabelInfoRes{}
			res.JSON(&json2)
			assert.NotNil(json.Result)
			assert.True(json2.Result.UpdatedAt.Equal(json.Result.UpdatedAt))
		})

		t.Run("should work with Channels", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Put(fmt.Sprintf("%s/v1/products/%s/labels/%s", tt.Host, product.Name, label.Name)).
				Set("Content-Type", "application/json").
				Send(tpl.LabelUpdateBody{
					Channels: &[]string{"stable"},
				}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.SettingInfoRes{}
			res.JSON(&json)
			assert.NotNil(json.Result)
			assert.Equal(label.Name, json.Result.Name)
			assert.Equal([]string{"stable"}, json.Result.Channels)
		})

		t.Run("should work with Clients", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Put(fmt.Sprintf("%s/v1/products/%s/labels/%s", tt.Host, product.Name, label.Name)).
				Set("Content-Type", "application/json").
				Send(tpl.LabelUpdateBody{
					Channels: &[]string{"stable", "beta"},
					Clients:  &[]string{"ios"},
				}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.SettingInfoRes{}
			res.JSON(&json)
			assert.NotNil(json.Result)
			assert.Equal(label.Name, json.Result.Name)
			assert.Equal([]string{"beta", "stable"}, json.Result.Channels)
			assert.Equal([]string{"ios"}, json.Result.Clients)
		})

		t.Run("should 400", func(t *testing.T) {
			assert := assert.New(t)

			res, _ := request.Put(fmt.Sprintf("%s/v1/products/%s/labels/%s", tt.Host, product.Name, label.Name)).
				Set("Content-Type", "application/json").
				Send(tpl.LabelUpdateBody{
					Desc: nil,
				}).
				End()
			assert.Equal(400, res.StatusCode)
			res.Content() // close http client
		})
	})

	t.Run(`"DELETE /v1/products/:product/labels/:label"`, func(t *testing.T) {
		t.Run("should conflict before offline", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Delete(fmt.Sprintf("%s/v1/products/%s/labels/%s", tt.Host, product.Name, n1)).
				End()
			assert.Nil(err)
			assert.Equal(409, res.StatusCode)
			res.Content() // close http client
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

	t.Run(`POST "/v1/products/:product/labels/:label+:assign"`, func(t *testing.T) {
		label, err := createLabel(tt, product.Name)
		assert.Nil(t, err)

		users, err := createUsers(tt, 3)
		assert.Nil(t, err)

		group, err := createGroup(tt)
		assert.Nil(t, err)

		t.Run("should work", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Post(fmt.Sprintf("%s/v1/products/%s/labels/%s:assign", tt.Host, product.Name, label.Name)).
				Set("Content-Type", "application/json").
				Send(tpl.UsersGroupsBody{
					Users:  schema.GetUsersUID(users[0:2]),
					Groups: []string{group.UID},
				}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.LabelReleaseInfoRes{}
			res.JSON(&json)
			assert.Equal(int64(1), json.Result.Release)
			assert.Equal(group.UID, json.Result.Groups[0])

			var count int64
			assert.Nil(tt.DB.Table(`user_label`).Where("label_id = ?", label.ID).Count(&count).Error)
			assert.Equal(int64(2), count)

			assert.Nil(tt.DB.Table(`group_label`).Where("label_id = ?", label.ID).Count(&count).Error)
			assert.Equal(int64(1), count)
		})

		t.Run("should work with duplicate data", func(t *testing.T) {
			assert := assert.New(t)

			uids := []string{users[0].UID, users[2].UID}
			res, err := request.Post(fmt.Sprintf("%s/v1/products/%s/labels/%s:assign", tt.Host, product.Name, label.Name)).
				Set("Content-Type", "application/json").
				Send(tpl.UsersGroupsBody{
					Users: uids,
				}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.LabelReleaseInfoRes{}
			res.JSON(&json)
			assert.Equal(int64(2), json.Result.Release)
			assert.Equal(0, len(json.Result.Groups))
			assert.True(tpl.StringSliceHas(json.Result.Users, users[0].UID))
			assert.True(tpl.StringSliceHas(json.Result.Users, users[2].UID))

			var count int64
			assert.Nil(tt.DB.Table(`user_label`).Where("label_id = ?", label.ID).Count(&count).Error)
			assert.Equal(int64(3), count)

			assert.Nil(tt.DB.Table(`group_label`).Where("label_id = ?", label.ID).Count(&count).Error)
			assert.Equal(int64(1), count)
		})
	})

	t.Run(`"PUT /v1/products/:product/labels/:label+:offline"`, func(t *testing.T) {
		label, err := createLabel(tt, product.Name)
		assert.Nil(t, err)

		users, err := createUsers(tt, 3)
		assert.Nil(t, err)

		group, err := createGroup(tt)
		assert.Nil(t, err)

		t.Run("should work", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Post(fmt.Sprintf("%s/v1/products/%s/labels/%s:assign", tt.Host, product.Name, label.Name)).
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
			assert.Nil(tt.DB.Table(`user_label`).Where("label_id = ?", label.ID).Count(&count).Error)
			assert.Equal(int64(3), count)

			assert.Nil(tt.DB.Table(`group_label`).Where("label_id = ?", label.ID).Count(&count).Error)
			assert.Equal(int64(1), count)

			res, err = request.Put(fmt.Sprintf("%s/v1/products/%s/labels/%s:offline", tt.Host, product.Name, label.Name)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.BoolRes{}
			res.JSON(&json)
			assert.True(json.Result)

			l := label
			assert.Nil(tt.DB.First(&l).Error)
			assert.NotNil(l.OfflineAt)

			time.Sleep(time.Millisecond * 100)
			assert.Nil(tt.DB.Table(`user_label`).Where("label_id = ?", label.ID).Count(&count).Error)
			assert.Equal(int64(0), count)

			assert.Nil(tt.DB.Table(`group_label`).Where("label_id = ?", label.ID).Count(&count).Error)
			assert.Equal(int64(0), count)
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

			l := label
			assert.Nil(tt.DB.First(&l).Error)
			assert.NotNil(l.OfflineAt)
		})
	})
}
