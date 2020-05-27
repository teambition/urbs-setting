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

func createSetting(tt *TestTools, productName, moduleName string, values ...string) (setting schema.Setting, err error) {
	name := tpl.RandName()
	res, err := request.Post(fmt.Sprintf("%s/v1/products/%s/modules/%s/settings", tt.Host, productName, moduleName)).
		Set("Content-Type", "application/json").
		Send(tpl.NameDescBody{Name: name, Desc: name}).
		End()

	if err == nil && len(values) > 0 {
		res.Content() // close http client
		res, err = request.Put(fmt.Sprintf("%s/v1/products/%s/modules/%s/settings/%s", tt.Host, productName, moduleName, name)).
			Set("Content-Type", "application/json").
			Send(tpl.SettingUpdateBody{Values: &values}).
			End()
	}

	var product schema.Product
	var module schema.Module
	if err == nil {
		res.Content() // close http client
		_, err = tt.DB.ScanStruct(&product, "select * from `urbs_product` where `name` = ? limit 1", productName)
	}

	if err == nil {
		_, err = tt.DB.ScanStruct(&module, "select * from `urbs_module` where `product_id` = ? and `name` = ? limit 1", product.ID, moduleName)
	}

	if err == nil {
		_, err = tt.DB.ScanStruct(&setting, "select * from `urbs_setting` where `module_id` = ? and `name` = ? limit 1", module.ID, name)
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
			assert.True(strings.Contains(text, `"offlineAt":null`))
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
			res.Content() // close http client
		})

		t.Run(`should return 400`, func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Post(fmt.Sprintf("%s/v1/products/%s/modules/%s/settings", tt.Host, product.Name, module.Name)).
				Set("Content-Type", "application/json").
				Send(tpl.NameDescBody{Name: "aB", Desc: "test"}).
				End()
			assert.Nil(err)
			assert.Equal(400, res.StatusCode)
			res.Content() // close http client
		})
	})

	t.Run(`"GET /v1/products/:product/settings"`, func(t *testing.T) {
		t.Run("should work", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Get(fmt.Sprintf("%s/v1/products/%s/settings", tt.Host, product.Name)).
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
			assert.Equal(product.Name, data.Product)
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
			assert.True(strings.Contains(text, `"offlineAt":null`))
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
			assert.NotNil(json2.Result)
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
			res.Content() // close http client
		})
	})

	t.Run(`POST "/v1/products/:product/modules/:module/settings/:setting+:assign"`, func(t *testing.T) {
		module, err := createModule(tt, product.Name)
		assert.Nil(t, err)

		setting, err := createSetting(tt, product.Name, module.Name, "a", "b")
		assert.Nil(t, err)

		users, err := createUsers(tt, 3)
		assert.Nil(t, err)

		group, err := createGroup(tt)
		assert.Nil(t, err)

		t.Run("should work", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Post(fmt.Sprintf("%s/v1/products/%s/modules/%s/settings/%s:assign", tt.Host, product.Name, module.Name, setting.Name)).
				Set("Content-Type", "application/json").
				Send(tpl.UsersGroupsBody{
					Users:  schema.GetUsersUID(users[0:2]),
					Groups: []string{group.UID},
					Value:  "a",
				}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.SettingReleaseInfoRes{}
			res.JSON(&json)
			assert.Equal(int64(1), json.Result.Release)
			assert.Equal("a", json.Result.Value)
			assert.Equal(group.UID, json.Result.Groups[0])

			var count int64
			_, err = tt.DB.ScanVal(&count, "select count(*) from `user_setting` where `setting_id` = ?", setting.ID)
			assert.Nil(err)
			assert.Equal(int64(2), count)

			_, err = tt.DB.ScanVal(&count, "select count(*) from `group_setting` where `setting_id` = ?", setting.ID)
			assert.Nil(err)
			assert.Equal(int64(1), count)

			res, err = request.Get(fmt.Sprintf("%s/v1/users/%s/settings:unionAll?product=%s", tt.Host, users[0].UID, product.Name)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json2 := tpl.MySettingsRes{}
			_, err = res.JSON(&json2)

			assert.Equal(1, len(json2.Result))

			data := json2.Result[0]
			assert.Equal("a", data.Value)
			assert.Equal("", data.LastValue)
		})

		t.Run("should work with duplicate data", func(t *testing.T) {
			assert := assert.New(t)

			uids := []string{users[0].UID, users[2].UID}
			res, err := request.Post(fmt.Sprintf("%s/v1/products/%s/modules/%s/settings/%s:assign", tt.Host, product.Name, module.Name, setting.Name)).
				Set("Content-Type", "application/json").
				Send(tpl.UsersGroupsBody{
					Users: uids,
					Value: "b",
				}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.SettingReleaseInfoRes{}
			res.JSON(&json)
			result := json.Result
			assert.Equal(int64(2), result.Release)
			assert.Equal("b", result.Value)
			assert.Equal(0, len(result.Groups))
			assert.True(tpl.StringSliceHas(result.Users, users[0].UID))
			assert.True(tpl.StringSliceHas(result.Users, users[2].UID))

			var count int64
			_, err = tt.DB.ScanVal(&count, "select count(*) from `user_setting` where `setting_id` = ?", setting.ID)
			assert.Nil(err)
			assert.Equal(int64(3), count)

			_, err = tt.DB.ScanVal(&count, "select count(*) from `group_setting` where `setting_id` = ?", setting.ID)
			assert.Nil(err)
			assert.Equal(int64(1), count)

			res, err = request.Get(fmt.Sprintf("%s/v1/users/%s/settings:unionAll?product=%s", tt.Host, users[0].UID, product.Name)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json2 := tpl.MySettingsRes{}
			_, err = res.JSON(&json2)

			assert.Equal(1, len(json2.Result))

			data := json2.Result[0]
			assert.Equal("b", data.Value)
			assert.Equal("a", data.LastValue)
		})
	})

	t.Run(`GET "/v1/products/:product/modules/:module/settings/:setting/users"`, func(t *testing.T) {
		module, err := createModule(tt, product.Name)
		assert.Nil(t, err)

		setting, err := createSetting(tt, product.Name, module.Name, "x", "y")
		assert.Nil(t, err)

		users, err := createUsers(tt, 1)
		assert.Nil(t, err)

		t.Run("should work", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Post(fmt.Sprintf("%s/v1/products/%s/modules/%s/settings/%s:assign", tt.Host, product.Name, module.Name, setting.Name)).
				Set("Content-Type", "application/json").
				Send(tpl.UsersGroupsBody{
					Users: schema.GetUsersUID(users),
					Value: "x",
				}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.SettingReleaseInfoRes{}
			res.JSON(&json)
			assert.Equal(int64(1), json.Result.Release)
			assert.Equal("x", json.Result.Value)
			assert.Equal(users[0].UID, json.Result.Users[0])

			res, err = request.Get(fmt.Sprintf("%s/v1/products/%s/modules/%s/settings/%s/users", tt.Host, product.Name, module.Name, setting.Name)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json2 := tpl.SettingUsersInfoRes{}
			res.JSON(&json2)
			assert.Equal(1, json2.TotalSize)
			assert.Equal(1, len(json2.Result))
			assert.Equal(setting.ID, service.HIDToID(json2.Result[0].SettingHID, "setting"))
			assert.Equal(users[0].UID, json2.Result[0].User)
			assert.Equal(int64(1), json2.Result[0].Release)
			assert.Equal("x", json2.Result[0].Value)
		})
	})

	t.Run(`GET "/v1/products/:product/modules/:module/settings/:setting/groups"`, func(t *testing.T) {
		module, err := createModule(tt, product.Name)
		assert.Nil(t, err)

		setting, err := createSetting(tt, product.Name, module.Name, "x", "y")
		assert.Nil(t, err)

		group, err := createGroup(tt)
		assert.Nil(t, err)

		t.Run("should work", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Post(fmt.Sprintf("%s/v1/products/%s/modules/%s/settings/%s:assign", tt.Host, product.Name, module.Name, setting.Name)).
				Set("Content-Type", "application/json").
				Send(tpl.UsersGroupsBody{
					Groups: []string{group.UID},
					Value:  "x",
				}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.SettingReleaseInfoRes{}
			res.JSON(&json)
			assert.Equal(int64(1), json.Result.Release)
			assert.Equal("x", json.Result.Value)
			assert.Equal(group.UID, json.Result.Groups[0])

			res, err = request.Get(fmt.Sprintf("%s/v1/products/%s/modules/%s/settings/%s/groups", tt.Host, product.Name, module.Name, setting.Name)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json2 := tpl.SettingGroupsInfoRes{}
			res.JSON(&json2)
			assert.Equal(1, json2.TotalSize)
			assert.Equal(1, len(json2.Result))
			assert.Equal(setting.ID, service.HIDToID(json2.Result[0].SettingHID, "setting"))
			assert.Equal(group.UID, json2.Result[0].Group)
			assert.Equal(group.Kind, json2.Result[0].Kind)
			assert.Equal(int64(1), json2.Result[0].Release)
			assert.Equal("x", json2.Result[0].Value)
		})
	})

	t.Run(`POST "/v1/products/:product/modules/:module/settings/:setting+:recall"`, func(t *testing.T) {
		module, err := createModule(tt, product.Name)
		assert.Nil(t, err)

		setting, err := createSetting(tt, product.Name, module.Name, "x", "y")
		assert.Nil(t, err)

		group, err := createGroup(tt)
		assert.Nil(t, err)

		t.Run("should work", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Post(fmt.Sprintf("%s/v1/products/%s/modules/%s/settings/%s:assign", tt.Host, product.Name, module.Name, setting.Name)).
				Set("Content-Type", "application/json").
				Send(tpl.UsersGroupsBody{
					Groups: []string{group.UID},
					Value:  "x",
				}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.SettingReleaseInfoRes{}
			res.JSON(&json)
			assert.Equal(int64(1), json.Result.Release)
			assert.Equal("x", json.Result.Value)
			assert.Equal(group.UID, json.Result.Groups[0])

			var count int64
			_, err = tt.DB.ScanVal(&count, "select count(*) from `group_setting` where `setting_id` = ?", setting.ID)
			assert.Nil(err)
			assert.Equal(int64(1), count)

			res, err = request.Post(fmt.Sprintf("%s/v1/products/%s/modules/%s/settings/%s:recall", tt.Host, product.Name, module.Name, setting.Name)).
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

			_, err = tt.DB.ScanVal(&count, "select count(*) from `group_setting` where `setting_id` = ?", setting.ID)
			assert.Nil(err)
			assert.Equal(int64(0), count)
		})
	})

	t.Run(`"PUT /v1/products/:product/modules/:module/settings/:setting+:offline"`, func(t *testing.T) {
		product, err := createProduct(tt)
		assert.Nil(t, err)

		module, err := createModule(tt, product.Name)
		assert.Nil(t, err)

		setting, err := createSetting(tt, product.Name, module.Name, "x", "y")
		assert.Nil(t, err)

		users, err := createUsers(tt, 3)
		assert.Nil(t, err)

		group, err := createGroup(tt)
		assert.Nil(t, err)

		t.Run("should work", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Post(fmt.Sprintf("%s/v1/products/%s/modules/%s/settings/%s:assign", tt.Host, product.Name, module.Name, setting.Name)).
				Set("Content-Type", "application/json").
				Send(tpl.UsersGroupsBody{
					Users:  schema.GetUsersUID(users),
					Groups: []string{group.UID},
					Value:  "x",
				}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.SettingReleaseInfoRes{}
			res.JSON(&json)
			result := json.Result
			assert.Equal(int64(1), result.Release)
			assert.Equal("x", result.Value)

			var count int64
			_, err = tt.DB.ScanVal(&count, "select count(*) from `user_setting` where `setting_id` = ?", setting.ID)
			assert.Nil(err)
			assert.Equal(int64(3), count)

			_, err = tt.DB.ScanVal(&count, "select count(*) from `group_setting` where `setting_id` = ?", setting.ID)
			assert.Nil(err)
			assert.Equal(int64(1), count)

			res, err = request.Put(fmt.Sprintf("%s/v1/products/%s/modules/%s/settings/%s:offline", tt.Host, product.Name, module.Name, setting.Name)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json2 := tpl.BoolRes{}
			res.JSON(&json2)
			assert.True(json2.Result)

			time.Sleep(time.Millisecond * 100)
			_, err = tt.DB.ScanVal(&count, "select count(*) from `user_setting` where `setting_id` = ?", setting.ID)
			assert.Nil(err)
			assert.Equal(int64(0), count)

			_, err = tt.DB.ScanVal(&count, "select count(*) from `group_setting` where `setting_id` = ?", setting.ID)
			assert.Nil(err)
			assert.Equal(int64(0), count)

			assert.Nil(setting.OfflineAt)
			s := setting
			_, err = tt.DB.ScanStruct(&s, "select * from `urbs_setting` where `id` = ? limit 1", s.ID)
			assert.NotNil(s.OfflineAt)
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
	})

	t.Run(`setting rules`, func(t *testing.T) {
		product, err := createProduct(tt)
		assert.Nil(t, err)

		module, err := createModule(tt, product.Name)
		assert.Nil(t, err)

		setting, err := createSetting(tt, product.Name, module.Name, "x", "y")
		assert.Nil(t, err)

		users, err := createUsers(tt, 1)
		assert.Nil(t, err)
		user := users[0]

		var rule tpl.SettingRuleInfo

		t.Run(`"POST /v1/products/:product/modules/:module/settings/:setting/rules" should work`, func(t *testing.T) {
			assert := assert.New(t)
			res, err := request.Post(fmt.Sprintf("%s/v1/products/%s/modules/%s/settings/%s/rules", tt.Host, product.Name, module.Name, setting.Name)).
				Set("Content-Type", "application/json").
				Send(map[string]interface{}{
					"kind":  "userPercent",
					"value": "y",
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

			json := tpl.SettingRuleInfoRes{}
			res.JSON(&json)
			data := json.Result
			assert.True(service.HIDToID(data.HID, "setting_rule") > int64(0))
			assert.Equal(setting.ID, service.HIDToID(data.SettingHID, "setting"))
			assert.Equal("userPercent", data.Kind)
			assert.True(data.CreatedAt.UTC().Unix() > int64(0))
			assert.True(data.UpdatedAt.UTC().Unix() > int64(0))
			assert.Equal(int64(1), data.Release)
			assert.Equal("y", data.Value)

			rule = data
		})

		t.Run(`"POST /v1/products/:product/modules/:module/settings/:setting/rules" should return 409`, func(t *testing.T) {
			assert := assert.New(t)
			res, err := request.Post(fmt.Sprintf("%s/v1/products/%s/modules/%s/settings/%s/rules", tt.Host, product.Name, module.Name, setting.Name)).
				Set("Content-Type", "application/json").
				Send(map[string]interface{}{
					"kind":  "userPercent",
					"value": "y",
					"rule": map[string]interface{}{
						"value": 100,
					},
				}).
				End()
			assert.Nil(err)
			assert.Equal(409, res.StatusCode)
			res.Content() // close http client
		})

		t.Run(`"GET /v1/users/:uid/settings:unionAll" should apply rules`, func(t *testing.T) {
			assert := assert.New(t)
			res, err := request.Get(fmt.Sprintf("%s/v1/users/%s/settings:unionAll?product=%s", tt.Host, user.UID, product.Name)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.MySettingsRes{}
			_, err = res.JSON(&json)

			assert.Nil(err)
			assert.Equal(0, len(json.Result))

			time.Sleep(time.Millisecond * 100)
			res, err = request.Get(fmt.Sprintf("%s/v1/users/%s/settings:unionAll?product=%s", tt.Host, user.UID, product.Name)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			text, err := res.Text()
			assert.Nil(err)
			assert.False(strings.Contains(text, `"id"`))

			json = tpl.MySettingsRes{}
			_, err = res.JSON(&json)

			assert.Nil(err)
			assert.Equal(1, len(json.Result))
			assert.Equal("", json.NextPageToken)

			data := json.Result[0]
			assert.Equal(service.IDToHID(setting.ID, "setting"), data.HID)
			assert.Equal(module.Name, data.Module)
			assert.Equal(setting.Name, data.Name)
			assert.Equal("y", data.Value)
			assert.True(data.AssignedAt.After(time2020))
		})

		t.Run(`"GET /v1/products/:product/modules/:module/settings/:setting/rules" should work`, func(t *testing.T) {
			assert := assert.New(t)
			res, err := request.Get(fmt.Sprintf("%s/v1/products/%s/modules/%s/settings/%s/rules", tt.Host, product.Name, module.Name, setting.Name)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			text, err := res.Text()
			assert.Nil(err)
			assert.False(strings.Contains(text, `"id"`))

			json := tpl.SettingRulesInfoRes{}
			_, err = res.JSON(&json)

			assert.Nil(err)
			assert.Equal(1, json.TotalSize)
			assert.Equal(1, len(json.Result))
			assert.Equal("", json.NextPageToken)

			data := json.Result[0]
			assert.True(service.HIDToID(data.HID, "setting_rule") > int64(0))
			assert.Equal(setting.ID, service.HIDToID(data.SettingHID, "setting"))
			assert.Equal("userPercent", data.Kind)
			assert.True(data.CreatedAt.UTC().Unix() > int64(0))
			assert.True(data.UpdatedAt.UTC().Unix() > int64(0))
			assert.Equal(int64(1), data.Release)
			assert.Equal("y", data.Value)
		})

		t.Run(`"PUT /v1/products/:product/modules/:module/settings/:setting/rules/:hid" should work`, func(t *testing.T) {
			assert := assert.New(t)
			res, err := request.Put(fmt.Sprintf("%s/v1/products/%s/modules/%s/settings/%s/rules/%s", tt.Host, product.Name, module.Name, setting.Name, rule.HID)).
				Set("Content-Type", "application/json").
				Send(map[string]interface{}{
					"kind":  "userPercent",
					"value": "x",
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

			json := tpl.SettingRuleInfoRes{}
			_, err = res.JSON(&json)

			assert.Nil(err)

			data := json.Result
			assert.Equal(rule.HID, data.HID)
			assert.Equal(rule.SettingHID, data.SettingHID)
			assert.Equal("userPercent", data.Kind)
			assert.Equal(int64(2), data.Release)
			assert.Equal("x", data.Value)
		})

		t.Run(`"DELETE /v1/products/:product/modules/:module/settings/:setting/rules/:hid" should work`, func(t *testing.T) {
			assert := assert.New(t)
			res, err := request.Delete(fmt.Sprintf("%s/v1/products/%s/modules/%s/settings/%s/rules/%s", tt.Host, product.Name, module.Name, setting.Name, rule.HID)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.BoolRes{}
			_, err = res.JSON(&json)
			assert.Nil(err)
			assert.True(json.Result)

			res, err = request.Get(fmt.Sprintf("%s/v1/products/%s/modules/%s/settings/%s/rules", tt.Host, product.Name, module.Name, setting.Name)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json2 := tpl.SettingRulesInfoRes{}
			_, err = res.JSON(&json2)

			assert.Nil(err)
			assert.Equal(0, len(json2.Result))

			res, err = request.Delete(fmt.Sprintf("%s/v1/products/%s/modules/%s/settings/%s/rules/%s", tt.Host, product.Name, module.Name, setting.Name, rule.HID)).
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
