package api

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/DavidCai1993/request"
	"github.com/doug-martin/goqu/v9"
	"github.com/stretchr/testify/assert"
	"github.com/teambition/urbs-setting/src/schema"
	"github.com/teambition/urbs-setting/src/service"
	"github.com/teambition/urbs-setting/src/tpl"
)

var time2020 = time.Unix(1577836800, 0)

func createUsers(tt *TestTools, count int) (users []schema.User, err error) {
	uids := make([]string, count)
	for i := 0; i < count; i++ {
		uids[i] = tpl.RandUID()
	}

	res, err := request.Post(fmt.Sprintf("%s/v1/users:batch", tt.Host)).
		Set("Content-Type", "application/json").
		Send(tpl.UsersBody{Users: uids}).
		End()

	if err == nil {
		res.Content() // close http client
		err = tt.DB.From("urbs_user").
			Where(goqu.Ex{"uid": uids}).Order(goqu.C("id").Desc()).
			Executor().ScanStructs(&users)
	}
	return
}

func cleanupUserLabels(db *goqu.Database, uid string) error {
	_, err := db.Exec("update `urbs_user` set `labels` = '', `active_at` = 0 where `uid` = ?", uid)
	return err
}

func TestUserAPIs(t *testing.T) {
	tt, cleanup := SetUpTestTools()
	defer cleanup()

	uid1 := tpl.RandUID()
	uid2 := tpl.RandUID()

	t.Run(`"POST /v1/users:batch"`, func(t *testing.T) {
		t.Run("should work", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Post(fmt.Sprintf("%s/v1/users:batch", tt.Host)).
				Set("Content-Type", "application/json").
				Send(tpl.UsersBody{Users: []string{uid1}}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.BoolRes{}
			res.JSON(&json)
			assert.True(json.Result)
		})

		t.Run("should work with duplicate user", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Post(fmt.Sprintf("%s/v1/users:batch", tt.Host)).
				Set("Content-Type", "application/json").
				Send(tpl.UsersBody{Users: []string{uid1, uid2}}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.BoolRes{}
			res.JSON(&json)
			assert.True(json.Result)
		})

		t.Run(`should 400 if no user`, func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Post(fmt.Sprintf("%s/v1/users:batch", tt.Host)).
				Set("Content-Type", "application/json").
				Send(tpl.UsersBody{Users: []string{}}).
				End()
			assert.Nil(err)
			assert.Equal(400, res.StatusCode)

			json := tpl.ResponseType{}
			res.JSON(&json)
			assert.Equal("BadRequest", json.Error)
			assert.Nil(json.Result)
		})
	})

	t.Run(`"GET /v1/users"`, func(t *testing.T) {
		t.Run(`should work`, func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Get(fmt.Sprintf("%s/v1/users", tt.Host)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			text, err := res.Text()
			assert.Nil(err)
			assert.False(strings.Contains(text, `"id"`))

			json := tpl.UsersRes{}
			_, err = res.JSON(&json)

			assert.Nil(err)
			assert.True(len(json.Result) > 0)
			assert.True(json.TotalSize > 0)
			assert.NotEqual("", json.Result[0].UID)
			assert.True(json.Result[0].ActiveAt >= 0)
			assert.True(json.Result[0].CreatedAt.After(time2020))
		})
	})

	t.Run(`"GET /v1/users/:uid+:exists"`, func(t *testing.T) {
		t.Run(`should work`, func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Get(fmt.Sprintf("%s/v1/users/%s:exists", tt.Host, uid1)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.BoolRes{}
			res.JSON(&json)
			assert.True(json.Result)
		})

		t.Run(`should work if not exists`, func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Get(fmt.Sprintf("%s/v1/users/%s:exists", tt.Host, tpl.RandUID())).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.BoolRes{}
			res.JSON(&json)
			assert.False(json.Result)
		})
	})

	t.Run("user, group, label", func(t *testing.T) {
		group, users, err := createGroupWithUsers(tt, 4)
		assert.Nil(t, err)

		product, err := createProduct(tt)
		assert.Nil(t, err)

		label, err := createLabel(tt, product.Name)
		assert.Nil(t, err)

		label1, err := createLabel(tt, product.Name)
		assert.Nil(t, err)

		t.Run(`"GET /users/:uid/labels:cache" for invalid user`, func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Get(fmt.Sprintf("%s/users/%s/labels:cache?product=%s", tt.Host, tpl.RandUID(), product.Name)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.CacheLabelsInfoRes{}
			_, err = res.JSON(&json)
			assert.Nil(err)
			assert.Equal(0, len(json.Result))
			assert.True(json.Timestamp > 0)
		})

		t.Run(`"GET /users/:uid/labels:cache" when no label`, func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Get(fmt.Sprintf("%s/users/%s/labels:cache?product=%s", tt.Host, users[0].UID, product.Name)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.CacheLabelsInfoRes{}
			_, err = res.JSON(&json)
			assert.Nil(err)
			assert.Equal(0, len(json.Result))
			assert.True(json.Timestamp > 0)
		})

		t.Run(`"GET /users/:uid/labels:cache" when user label exists`, func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Post(fmt.Sprintf("%s/v1/products/%s/labels/%s:assign", tt.Host, product.Name, label.Name)).
				Set("Content-Type", "application/json").
				Send(tpl.UsersGroupsBody{
					Users: schema.GetUsersUID(users),
				}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)
			res.Content() // close http client

			res, err = request.Get(fmt.Sprintf("%s/users/%s/labels:cache?product=%s", tt.Host, users[1].UID, product.Name)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.CacheLabelsInfoRes{}
			_, err = res.JSON(&json)
			assert.Nil(err)
			assert.Equal(1, len(json.Result))
			assert.True(json.Timestamp > 0)
			assert.Equal(label.Name, json.Result[0].Label)
			assert.Equal(0, len(json.Result[0].Clients))
			assert.Equal(0, len(json.Result[0].Channels))
		})

		t.Run(`"PUT /users/:uid/labels:cache"`, func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Get(fmt.Sprintf("%s/users/%s/labels:cache?product=%s", tt.Host, users[1].UID, product.Name)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.CacheLabelsInfoRes{}
			_, err = res.JSON(&json)
			assert.Nil(err)
			assert.Equal(1, len(json.Result))
			assert.True(json.Timestamp > 0)
			t1 := json.Timestamp

			time.Sleep(time.Millisecond * 1100)
			res, err = request.Put(fmt.Sprintf("%s/v1/users/%s/labels:cache", tt.Host, users[1].UID)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json2 := tpl.UserRes{}
			_, err = res.JSON(&json2)
			assert.Nil(err)
			assert.Equal(users[1].UID, json2.Result.UID)
			assert.True(time.Now().UTC().Unix()-json2.Result.ActiveAt <= 1)

			res, err = request.Get(fmt.Sprintf("%s/users/%s/labels:cache?product=%s", tt.Host, users[1].UID, product.Name)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json = tpl.CacheLabelsInfoRes{}
			_, err = res.JSON(&json)
			assert.Nil(err)
			assert.True(json.Timestamp >= (t1 + 1))
		})

		t.Run(`"GET /users/:uid/labels:cache" when group label exists`, func(t *testing.T) {
			assert := assert.New(t)

			time.Sleep(time.Second)
			res, err := request.Post(fmt.Sprintf("%s/v1/products/%s/labels/%s:assign", tt.Host, product.Name, label1.Name)).
				Set("Content-Type", "application/json").
				Send(tpl.UsersGroupsBody{
					Groups: []string{group.UID},
				}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)
			res.Content() // close http client

			// users[1] labels from cache
			res, err = request.Get(fmt.Sprintf("%s/users/%s/labels:cache?product=%s", tt.Host, users[1].UID, product.Name)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.CacheLabelsInfoRes{}
			_, err = res.JSON(&json)
			assert.Nil(err)
			assert.Equal(1, len(json.Result))
			assert.True(json.Timestamp > 0)
			assert.Equal(label.Name, json.Result[0].Label)
			assert.Equal(0, len(json.Result[0].Clients))
			assert.Equal(0, len(json.Result[0].Channels))

			// users[2] get all labels
			res, err = request.Get(fmt.Sprintf("%s/users/%s/labels:cache?product=%s", tt.Host, users[2].UID, product.Name)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json = tpl.CacheLabelsInfoRes{}
			_, err = res.JSON(&json)
			assert.Nil(err)
			assert.Equal(2, len(json.Result))
			assert.True(json.Timestamp > 0)
			assert.Equal(label1.Name, json.Result[0].Label)
			assert.Equal(0, len(json.Result[0].Clients))
			assert.Equal(0, len(json.Result[0].Channels))

			assert.Equal(label.Name, json.Result[1].Label)

			time.Sleep(time.Second)
			res, err = request.Post(fmt.Sprintf("%s/v1/products/%s/labels/%s:assign", tt.Host, product.Name, label.Name)).
				Set("Content-Type", "application/json").
				Send(tpl.UsersGroupsBody{
					Groups: []string{group.UID},
				}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)
			res.Content() // close http client

			// users[2] get all labels
			res, err = request.Get(fmt.Sprintf("%s/users/%s/labels:cache?product=%s", tt.Host, users[3].UID, product.Name)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json = tpl.CacheLabelsInfoRes{}
			_, err = res.JSON(&json)
			assert.Nil(err)
			assert.Equal(2, len(json.Result))
			assert.True(json.Timestamp > 0)
			assert.Equal(label.Name, json.Result[0].Label)
			assert.Equal(0, len(json.Result[0].Clients))
			assert.Equal(0, len(json.Result[0].Channels))

			assert.Equal(label1.Name, json.Result[1].Label)
		})

		t.Run(`"GET /v1/users/:uid/labels"`, func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Get(fmt.Sprintf("%s/v1/users/%s/labels", tt.Host, users[3].UID)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			text, err := res.Text()
			assert.Nil(err)
			assert.False(strings.Contains(text, `"id"`))

			json := tpl.MyLabelsRes{}
			_, err = res.JSON(&json)

			assert.Nil(err)
			assert.Equal(1, len(json.Result))
			assert.Equal(1, json.TotalSize)
			assert.Equal("", json.NextPageToken)
			assert.Equal(service.IDToHID(label.ID, "label"), json.Result[0].HID)
			assert.Equal(product.Name, json.Result[0].Product)
			assert.Equal(label.Name, json.Result[0].Name)
			assert.True(json.Result[0].AssignedAt.After(time2020))
			assert.True(json.Result[0].Release > 0)
		})

		t.Run(`"DELETE /v1/products/:product/labels/:label/users/:uid" should work`, func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Delete(fmt.Sprintf("%s/v1/products/%s/labels/%s/users/%s", tt.Host, product.Name, label.Name, users[3].UID)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.BoolRes{}
			_, err = res.JSON(&json)
			assert.Nil(err)
			assert.True(json.Result)

			res, err = request.Delete(fmt.Sprintf("%s/v1/products/%s/labels/%s/groups/%s", tt.Host, product.Name, label.Name, group.UID)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json = tpl.BoolRes{}
			_, err = res.JSON(&json)
			assert.Nil(err)
			assert.True(json.Result)

			assert.Nil(cleanupUserLabels(tt.DB, users[3].UID))
			res, err = request.Get(fmt.Sprintf("%s/users/%s/labels:cache?product=%s", tt.Host, users[3].UID, product.Name)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json2 := tpl.CacheLabelsInfoRes{}
			_, err = res.JSON(&json2)
			assert.Nil(err)
			assert.True(json2.Timestamp > 0)
			assert.Equal(1, len(json2.Result))
			assert.Equal(label1.Name, json2.Result[0].Label)
		})
	})

	t.Run("user, group, setting", func(t *testing.T) {
		group, users, err := createGroupWithUsers(tt, 4)
		assert.Nil(t, err)

		product, err := createProduct(tt)
		assert.Nil(t, err)

		module, err := createModule(tt, product.Name)
		assert.Nil(t, err)

		setting0, err := createSetting(tt, product.Name, module.Name, "a", "b")
		assert.Nil(t, err)

		setting1, err := createSetting(tt, product.Name, module.Name, "true", "false")
		assert.Nil(t, err)

		t.Run(`"GET /v1/users/:uid/settings" for invalid user`, func(t *testing.T) {
			assert := assert.New(t)

			res, _ := request.Get(fmt.Sprintf("%s/v1/users/%s/settings", tt.Host, tpl.RandUID())).
				End()
			assert.Equal(404, res.StatusCode)
			res.Content() // close http client
		})

		t.Run(`"GET /v1/users/:uid/settings" when without settings`, func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Get(fmt.Sprintf("%s/v1/users/%s/settings", tt.Host, users[0].UID)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			text, err := res.Text()
			assert.Nil(err)
			assert.False(strings.Contains(text, `"id"`))

			json := tpl.MySettingsRes{}
			_, err = res.JSON(&json)

			assert.Nil(err)
			assert.Equal(0, len(json.Result))
			assert.Equal(0, json.TotalSize)
			assert.Equal("", json.NextPageToken)
		})

		t.Run(`"GET /v1/users/:uid/settings" should work`, func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Post(fmt.Sprintf("%s/v1/products/%s/modules/%s/settings/%s:assign", tt.Host, product.Name, module.Name, setting0.Name)).
				Set("Content-Type", "application/json").
				Send(tpl.UsersGroupsBody{
					Users: []string{users[0].UID},
					Value: "a",
				}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)
			res.Content() // close http client

			time.Sleep(time.Millisecond * 10)
			res, err = request.Get(fmt.Sprintf("%s/v1/users/%s/settings", tt.Host, users[0].UID)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			text, err := res.Text()
			assert.Nil(err)
			assert.False(strings.Contains(text, `"id"`))

			json := tpl.MySettingsRes{}
			_, err = res.JSON(&json)

			assert.Nil(err)
			assert.Equal(1, len(json.Result))
			assert.Equal("", json.NextPageToken)

			data := json.Result[0]
			assert.Equal(service.IDToHID(setting0.ID, "setting"), data.HID)
			assert.Equal(module.Name, data.Module)
			assert.Equal(setting0.Name, data.Name)
			assert.Equal("a", data.Value)
			assert.Equal("", data.LastValue)
			assert.True(data.AssignedAt.After(time2020))
			assert.True(data.Release > 0)
		})

		t.Run(`"GET /v1/users/:uid/settings:unionAll" for invalid user`, func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Get(fmt.Sprintf("%s/v1/users/%s/settings:unionAll?product=%s", tt.Host, tpl.RandUID(), product.Name)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			text, err := res.Text()
			assert.Nil(err)
			assert.False(strings.Contains(text, `"id"`))

			json := tpl.MySettingsRes{}
			_, err = res.JSON(&json)

			assert.Nil(err)
			assert.Equal(0, len(json.Result))
			assert.Equal("", json.NextPageToken)
		})

		t.Run(`"GET /v1/users/:uid/settings:unionAll" should work`, func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Get(fmt.Sprintf("%s/v1/users/%s/settings:unionAll?product=%s", tt.Host, users[0].UID, product.Name)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			text, err := res.Text()
			assert.Nil(err)
			assert.False(strings.Contains(text, `"id"`))

			json := tpl.MySettingsRes{}
			_, err = res.JSON(&json)

			assert.Nil(err)
			assert.Equal(1, len(json.Result))
			assert.Equal("", json.NextPageToken)

			data := json.Result[0]
			assert.Equal(service.IDToHID(setting0.ID, "setting"), data.HID)
			assert.Equal(product.Name, data.Product)
			assert.Equal(module.Name, data.Module)
			assert.Equal(setting0.Name, data.Name)
			assert.Equal("a", data.Value)
			assert.Equal("", data.LastValue)
			assert.True(data.AssignedAt.After(time2020))
			assert.True(data.Release > 0)

			time.Sleep(time.Millisecond * 10)
			res, err = request.Post(fmt.Sprintf("%s/v1/products/%s/modules/%s/settings/%s:assign", tt.Host, product.Name, module.Name, setting1.Name)).
				Set("Content-Type", "application/json").
				Send(tpl.UsersGroupsBody{
					Groups: []string{group.UID},
					Value:  "true",
				}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)
			res.Content() // close http client

			res, err = request.Get(fmt.Sprintf("%s/v1/users/%s/settings", tt.Host, users[0].UID)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			text, err = res.Text()
			assert.Nil(err)
			assert.False(strings.Contains(text, `"id"`))

			json = tpl.MySettingsRes{}
			_, err = res.JSON(&json)

			assert.Nil(err)
			assert.Equal(1, len(json.Result))
			assert.Equal("", json.NextPageToken)

			res, err = request.Get(fmt.Sprintf("%s/v1/users/%s/settings:unionAll?product=%s", tt.Host, users[0].UID, product.Name)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			text, err = res.Text()
			assert.Nil(err)
			assert.False(strings.Contains(text, `"id"`))

			json = tpl.MySettingsRes{}
			_, err = res.JSON(&json)

			assert.Nil(err)
			assert.Equal(2, len(json.Result))
			assert.Equal("", json.NextPageToken)

			data = json.Result[0]
			assert.Equal(service.IDToHID(setting1.ID, "setting"), data.HID)
			assert.Equal(module.Name, data.Module)
			assert.Equal(setting1.Name, data.Name)
			assert.Equal("true", data.Value)
			assert.Equal("", data.LastValue)
			assert.True(data.AssignedAt.After(time2020))
			assert.True(data.Release > 0)

			assert.Equal(service.IDToHID(setting0.ID, "setting"), json.Result[1].HID)
		})

		t.Run(`"GET /v1/users/:uid/settings:unionAll" with channel or client query should work`, func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Put(fmt.Sprintf("%s/v1/products/%s/modules/%s/settings/%s",
				tt.Host, product.Name, module.Name, setting1.Name)).
				Set("Content-Type", "application/json").
				Send(tpl.SettingUpdateBody{Channels: &[]string{"stable", "beta"}, Clients: &[]string{"ios"}}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)
			res.Content() // close http client

			res, err = request.Get(fmt.Sprintf("%s/v1/users/%s/settings:unionAll?product=%s&channel=beta", tt.Host, users[0].UID, product.Name)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.MySettingsRes{}
			_, err = res.JSON(&json)

			assert.Nil(err)
			assert.Equal(1, len(json.Result))

			res, err = request.Get(fmt.Sprintf("%s/v1/users/%s/settings:unionAll?product=%s&channel=beta&client=ios", tt.Host, users[0].UID, product.Name)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json = tpl.MySettingsRes{}
			_, err = res.JSON(&json)

			assert.Nil(err)
			assert.Equal(2, len(json.Result))

			res, err = request.Get(fmt.Sprintf("%s/v1/users/%s/settings:unionAll?product=%s&channel=beta&client=android", tt.Host, users[0].UID, product.Name)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json = tpl.MySettingsRes{}
			_, err = res.JSON(&json)

			assert.Nil(err)
			assert.Equal(1, len(json.Result))

			res, err = request.Get(fmt.Sprintf("%s/v1/users/%s/settings:unionAll?product=%s&channel=dev", tt.Host, users[0].UID, product.Name)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json = tpl.MySettingsRes{}
			_, err = res.JSON(&json)

			assert.Nil(err)
			assert.Equal(1, len(json.Result))
		})

		t.Run(`"PUT /v1/products/:product/modules/:module/settings/:setting/users/:uid+:rollback" should work`, func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Put(fmt.Sprintf("%s/v1/products/%s/modules/%s/settings/%s/users/%s:rollback", tt.Host, product.Name, module.Name, setting0.Name, users[0].UID)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)
			res.Content() // close http client

			res, err = request.Get(fmt.Sprintf("%s/v1/users/%s/settings", tt.Host, users[0].UID)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			text, err := res.Text()
			assert.Nil(err)
			assert.False(strings.Contains(text, `"id"`))

			json := tpl.MySettingsRes{}
			_, err = res.JSON(&json)

			assert.Nil(err)
			assert.Equal(1, len(json.Result))
			assert.Equal("", json.NextPageToken)

			data := json.Result[0]
			assert.Equal(service.IDToHID(setting0.ID, "setting"), data.HID)
			assert.Equal("", data.Value)
			assert.Equal("", data.LastValue)
		})

		t.Run(`"DELETE /v1/products/:product/modules/:module/settings/:setting/users/:uid" should work`, func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Delete(fmt.Sprintf("%s/v1/products/%s/modules/%s/settings/%s/users/%s", tt.Host, product.Name, module.Name, setting0.Name, users[0].UID)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)
			res.Content() // close http client

			res, err = request.Get(fmt.Sprintf("%s/v1/users/%s/settings", tt.Host, users[0].UID)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			text, err := res.Text()
			assert.Nil(err)
			assert.False(strings.Contains(text, `"id"`))

			json := tpl.MySettingsRes{}
			_, err = res.JSON(&json)

			assert.Nil(err)
			assert.Equal(0, len(json.Result))
			assert.Equal("", json.NextPageToken)
		})
	})
}
