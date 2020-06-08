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
		_, err = tt.DB.ScanStruct(&product, "select * from `urbs_product` where `name` = ? limit 1", productName)
	}

	if err == nil {
		_, err = tt.DB.ScanStruct(&label, "select * from `urbs_label` where `product_id` = ? and `name` = ? limit 1", product.ID, name)
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
			assert.True(strings.Contains(text, `"offlineAt":null`))
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
			assert.True(strings.Contains(text, `"offlineAt":null`))
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
			_, err = tt.DB.ScanVal(&count, "select count(*) from `user_label` where `label_id` = ?", label.ID)
			assert.Nil(err)
			assert.Equal(int64(2), count)

			_, err = tt.DB.ScanVal(&count, "select count(*) from `group_label` where `label_id` = ?", label.ID)
			assert.Nil(err)
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
			_, err = tt.DB.ScanVal(&count, "select count(*) from `user_label` where `label_id` = ?", label.ID)
			assert.Nil(err)
			assert.Equal(int64(3), count)

			_, err = tt.DB.ScanVal(&count, "select count(*) from `group_label` where `label_id` = ?", label.ID)
			assert.Nil(err)
			assert.Equal(int64(1), count)
		})
	})

	t.Run(`GET "/v1/products/:product/labels/:label/users"`, func(t *testing.T) {
		label, err := createLabel(tt, product.Name)
		assert.Nil(t, err)

		users, err := createUsers(tt, 1)
		assert.Nil(t, err)

		t.Run("should work", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Post(fmt.Sprintf("%s/v1/products/%s/labels/%s:assign", tt.Host, product.Name, label.Name)).
				Set("Content-Type", "application/json").
				Send(tpl.UsersGroupsBody{
					Users: schema.GetUsersUID(users),
				}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.LabelReleaseInfoRes{}
			res.JSON(&json)
			assert.Equal(int64(1), json.Result.Release)
			assert.Equal(users[0].UID, json.Result.Users[0])

			res, err = request.Get(fmt.Sprintf("%s/v1/products/%s/labels/%s/users", tt.Host, product.Name, label.Name)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json2 := tpl.LabelUsersInfoRes{}
			res.JSON(&json2)
			assert.Equal(1, json2.TotalSize)
			assert.Equal(1, len(json2.Result))
			assert.Equal(label.ID, service.HIDToID(json2.Result[0].LabelHID, "label"))
			assert.Equal(users[0].UID, json2.Result[0].User)
			assert.Equal(int64(1), json2.Result[0].Release)
		})
	})

	t.Run(`GET "/v1/products/:product/labels/:label/groups"`, func(t *testing.T) {
		label, err := createLabel(tt, product.Name)
		assert.Nil(t, err)

		group, err := createGroup(tt)
		assert.Nil(t, err)

		t.Run("should work", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Post(fmt.Sprintf("%s/v1/products/%s/labels/%s:assign", tt.Host, product.Name, label.Name)).
				Set("Content-Type", "application/json").
				Send(tpl.UsersGroupsBody{
					Groups: []string{group.UID},
				}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.SettingReleaseInfoRes{}
			res.JSON(&json)
			assert.Equal(int64(1), json.Result.Release)
			assert.Equal(group.UID, json.Result.Groups[0])

			res, err = request.Get(fmt.Sprintf("%s/v1/products/%s/labels/%s/groups", tt.Host, product.Name, label.Name)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json2 := tpl.LabelGroupsInfoRes{}
			res.JSON(&json2)
			assert.Equal(1, json2.TotalSize)
			assert.Equal(1, len(json2.Result))
			assert.Equal(label.ID, service.HIDToID(json2.Result[0].LabelHID, "label"))
			assert.Equal(group.UID, json2.Result[0].Group)
			assert.Equal(group.Kind, json2.Result[0].Kind)
			assert.Equal(int64(1), json2.Result[0].Release)
		})
	})

	t.Run(`POST "/v1/products/:product/labels/:label+:recall"`, func(t *testing.T) {
		label, err := createLabel(tt, product.Name)
		assert.Nil(t, err)

		group, err := createGroup(tt)
		assert.Nil(t, err)

		t.Run("should work", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Post(fmt.Sprintf("%s/v1/products/%s/labels/%s:assign", tt.Host, product.Name, label.Name)).
				Set("Content-Type", "application/json").
				Send(tpl.UsersGroupsBody{
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
			_, err = tt.DB.ScanVal(&count, "select count(*) from `group_label` where `label_id` = ?", label.ID)
			assert.Nil(err)
			assert.Equal(int64(1), count)

			res, err = request.Post(fmt.Sprintf("%s/v1/products/%s/labels/%s:recall", tt.Host, product.Name, label.Name)).
				Set("Content-Type", "application/json").
				Send(tpl.RecallBody{
					Release: 1,
				}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json2 := tpl.BoolRes{}
			res.JSON(&json2)
			assert.True(json2.Result)

			_, err = tt.DB.ScanVal(&count, "select count(*) from `group_label` where `label_id` = ?", label.ID)
			assert.Nil(err)
			assert.Equal(int64(0), count)
		})
	})

	t.Run(`"PUT /v1/products/:product/labels/:label+:cleanup"`, func(t *testing.T) {
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

			res, err = request.Post(fmt.Sprintf("%s/v1/products/%s/labels/%s/rules", tt.Host, product.Name, label.Name)).
				Set("Content-Type", "application/json").
				Send(map[string]interface{}{
					"kind": "userPercent",
					"rule": map[string]interface{}{
						"value": 100,
					},
				}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			var count int64
			_, err = tt.DB.ScanVal(&count, "select count(*) from `user_label` where `label_id` = ?", label.ID)
			assert.Nil(err)
			assert.Equal(int64(3), count)

			_, err = tt.DB.ScanVal(&count, "select count(*) from `group_label` where `label_id` = ?", label.ID)
			assert.Nil(err)
			assert.Equal(int64(1), count)

			_, err = tt.DB.ScanVal(&count, "select count(*) from `label_rule` where `label_id` = ?", label.ID)
			assert.Nil(err)
			assert.Equal(int64(1), count)

			res, err = request.Delete(fmt.Sprintf("%s/v1/products/%s/labels/%s:cleanup", tt.Host, product.Name, label.Name)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.BoolRes{}
			res.JSON(&json)
			assert.True(json.Result)

			l := label
			_, err = tt.DB.ScanStruct(&l, "select * from `urbs_label` where `id` = ? limit 1", label.ID)
			assert.Nil(err)
			assert.Equal(int64(0), l.Status)

			_, err = tt.DB.ScanVal(&count, "select count(*) from `user_label` where `label_id` = ?", label.ID)
			assert.Nil(err)
			assert.Equal(int64(0), count)

			_, err = tt.DB.ScanVal(&count, "select count(*) from `group_label` where `label_id` = ?", label.ID)
			assert.Nil(err)
			assert.Equal(int64(0), count)

			_, err = tt.DB.ScanVal(&count, "select count(*) from `label_rule` where `label_id` = ?", label.ID)
			assert.Nil(err)
			assert.Equal(int64(0), count)
		})

		t.Run("should work idempotent", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Delete(fmt.Sprintf("%s/v1/products/%s/labels/%s:cleanup", tt.Host, product.Name, label.Name)).
				End()
			assert.Nil(err)

			assert.Equal(200, res.StatusCode)

			json := tpl.BoolRes{}
			res.JSON(&json)
			assert.True(json.Result)

			l := label
			_, err = tt.DB.ScanStruct(&l, "select * from `urbs_label` where `id` = ? limit 1", label.ID)
			assert.Nil(err)
			assert.Equal(int64(0), l.Status)
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

			res, err = request.Post(fmt.Sprintf("%s/v1/products/%s/labels/%s/rules", tt.Host, product.Name, label.Name)).
				Set("Content-Type", "application/json").
				Send(map[string]interface{}{
					"kind": "userPercent",
					"rule": map[string]interface{}{
						"value": 100,
					},
				}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			var count int64
			_, err = tt.DB.ScanVal(&count, "select count(*) from `user_label` where `label_id` = ?", label.ID)
			assert.Nil(err)
			assert.Equal(int64(3), count)

			_, err = tt.DB.ScanVal(&count, "select count(*) from `group_label` where `label_id` = ?", label.ID)
			assert.Nil(err)
			assert.Equal(int64(1), count)

			_, err = tt.DB.ScanVal(&count, "select count(*) from `label_rule` where `label_id` = ?", label.ID)
			assert.Nil(err)
			assert.Equal(int64(1), count)

			res, err = request.Put(fmt.Sprintf("%s/v1/products/%s/labels/%s:offline", tt.Host, product.Name, label.Name)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.BoolRes{}
			res.JSON(&json)
			assert.True(json.Result)

			l := label
			_, err = tt.DB.ScanStruct(&l, "select * from `urbs_label` where `id` = ? limit 1", label.ID)
			assert.Nil(err)
			assert.NotNil(l.OfflineAt)

			time.Sleep(time.Millisecond * 100)
			_, err = tt.DB.ScanVal(&count, "select count(*) from `user_label` where `label_id` = ?", label.ID)
			assert.Nil(err)
			assert.Equal(int64(0), count)

			_, err = tt.DB.ScanVal(&count, "select count(*) from `group_label` where `label_id` = ?", label.ID)
			assert.Nil(err)
			assert.Equal(int64(0), count)

			_, err = tt.DB.ScanVal(&count, "select count(*) from `label_rule` where `label_id` = ?", label.ID)
			assert.Nil(err)
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
			_, err = tt.DB.ScanStruct(&l, "select * from `urbs_label` where `id` = ? limit 1", label.ID)
			assert.Nil(err)
			assert.NotNil(l.OfflineAt)
		})
	})

	t.Run(`label rules`, func(t *testing.T) {
		product, err := createProduct(tt)
		assert.Nil(t, err)

		label, err := createLabel(tt, product.Name)
		assert.Nil(t, err)

		users, err := createUsers(tt, 1)
		assert.Nil(t, err)
		user := users[0]

		var rule tpl.LabelRuleInfo

		t.Run(`"POST /v1/products/:product/labels/:label/rules" should work`, func(t *testing.T) {
			assert := assert.New(t)
			res, err := request.Post(fmt.Sprintf("%s/v1/products/%s/labels/%s/rules", tt.Host, product.Name, label.Name)).
				Set("Content-Type", "application/json").
				Send(map[string]interface{}{
					"kind": "userPercent",
					"rule": map[string]interface{}{
						"value": 100,
					},
				}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			text, err := res.Text()
			assert.Nil(err)
			assert.True(strings.Contains(text, `"rule":{"value":100}`))
			assert.False(strings.Contains(text, `"id"`))

			json := tpl.LabelRuleInfoRes{}
			res.JSON(&json)
			data := json.Result
			assert.True(service.HIDToID(data.HID, "label_rule") > int64(0))
			assert.Equal(label.ID, service.HIDToID(data.LabelHID, "label"))
			assert.Equal("userPercent", data.Kind)
			assert.True(data.CreatedAt.UTC().Unix() > int64(0))
			assert.True(data.UpdatedAt.UTC().Unix() > int64(0))
			assert.Equal(int64(1), data.Release)

			rule = data
		})

		t.Run(`"POST /v1/products/:product/labels/:label/rules" should return 409`, func(t *testing.T) {
			assert := assert.New(t)
			res, err := request.Post(fmt.Sprintf("%s/v1/products/%s/labels/%s/rules", tt.Host, product.Name, label.Name)).
				Set("Content-Type", "application/json").
				Send(map[string]interface{}{
					"kind": "userPercent",
					"rule": map[string]interface{}{
						"value": 100,
					},
				}).
				End()
			assert.Nil(err)
			assert.Equal(409, res.StatusCode)
			res.Content() // close http client
		})

		t.Run(`"GET /users/:uid/labels:cache" should apply rules`, func(t *testing.T) {
			assert := assert.New(t)
			res, err := request.Get(fmt.Sprintf("%s/users/%s/labels:cache?product=%s", tt.Host, user.UID, product.Name)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			text, err := res.Text()
			assert.Nil(err)
			assert.False(strings.Contains(text, `"id"`))

			json := tpl.CacheLabelsInfoRes{}
			_, err = res.JSON(&json)

			assert.Nil(err)
			assert.Equal(1, len(json.Result))

			data := json.Result[0]
			assert.Equal(label.Name, data.Label)
		})

		t.Run(`"GET /v1/products/:product/labels/:label/rules" should work`, func(t *testing.T) {
			assert := assert.New(t)
			res, err := request.Get(fmt.Sprintf("%s/v1/products/%s/labels/%s/rules", tt.Host, product.Name, label.Name)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			text, err := res.Text()
			assert.Nil(err)
			assert.False(strings.Contains(text, `"id"`))

			json := tpl.LabelRulesInfoRes{}
			_, err = res.JSON(&json)

			assert.Nil(err)
			assert.Equal(1, json.TotalSize)
			assert.Equal(1, len(json.Result))
			assert.Equal("", json.NextPageToken)

			data := json.Result[0]
			assert.True(service.HIDToID(data.HID, "label_rule") > int64(0))
			assert.Equal(label.ID, service.HIDToID(data.LabelHID, "label"))
			assert.Equal("userPercent", data.Kind)
			assert.True(data.CreatedAt.UTC().Unix() > int64(0))
			assert.True(data.UpdatedAt.UTC().Unix() > int64(0))
			assert.Equal(int64(1), data.Release)
		})

		t.Run(`"PUT /v1/products/:product/labels/:label/rules/:hid" should work`, func(t *testing.T) {
			assert := assert.New(t)
			res, err := request.Put(fmt.Sprintf("%s/v1/products/%s/labels/%s/rules/%s", tt.Host, product.Name, label.Name, rule.HID)).
				Set("Content-Type", "application/json").
				Send(map[string]interface{}{
					"kind": "userPercent",
					"rule": map[string]interface{}{
						"value": 0,
					},
				}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			text, err := res.Text()
			assert.Nil(err)
			assert.True(strings.Contains(text, `"rule":{"value":0}`))
			assert.False(strings.Contains(text, `"id"`))

			json := tpl.LabelRuleInfoRes{}
			_, err = res.JSON(&json)

			assert.Nil(err)

			data := json.Result
			assert.Equal(rule.HID, data.HID)
			assert.Equal(rule.LabelHID, data.LabelHID)
			assert.Equal("userPercent", data.Kind)
			assert.Equal(int64(2), data.Release)
		})

		t.Run(`"DELETE /v1/products/:product/labels/:label/rules/:hid" should work`, func(t *testing.T) {
			assert := assert.New(t)
			res, err := request.Delete(fmt.Sprintf("%s/v1/products/%s/labels/%s/rules/%s", tt.Host, product.Name, label.Name, rule.HID)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.BoolRes{}
			_, err = res.JSON(&json)
			assert.Nil(err)
			assert.True(json.Result)

			res, err = request.Get(fmt.Sprintf("%s/v1/products/%s/labels/%s/rules", tt.Host, product.Name, label.Name)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json2 := tpl.LabelRulesInfoRes{}
			_, err = res.JSON(&json2)

			assert.Nil(err)
			assert.Equal(0, len(json2.Result))

			res, err = request.Delete(fmt.Sprintf("%s/v1/products/%s/labels/%s/rules/%s", tt.Host, product.Name, label.Name, rule.HID)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json = tpl.BoolRes{}
			_, err = res.JSON(&json)
			assert.Nil(err)
			assert.False(json.Result)
		})
	})
}
