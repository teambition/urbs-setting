package api

import (
	"fmt"
	"testing"
	"time"

	"github.com/DavidCai1993/request"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/teambition/urbs-setting/src/service"
	"github.com/teambition/urbs-setting/src/tpl"
)

func createUsers(appHost string, count int) ([]string, error) {
	uids := make([]string, count)
	for i := 0; i < count; i++ {
		uids[i] = tpl.RandUID()
	}

	_, err := request.Post(fmt.Sprintf("%s/v1/users:batch", appHost)).
		Set("Content-Type", "application/json").
		Send(tpl.UsersBody{Users: uids}).
		End()

	if err != nil {
		return nil, err
	}
	return uids, nil
}

func cleanupUserLabels(db *gorm.DB, uid string) error {
	return db.Table("urbs_user").Where("uid = ?", uid).Updates(map[string]interface{}{
		"labels": "", "active_at": 1}).Error
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
		group, users, err := createGroupWithUsers(tt.Host, 4)
		assert.Nil(t, err)

		product, err := createProduct(tt.Host)
		assert.Nil(t, err)

		label, err := createLabel(tt.Host, product.Name)
		assert.Nil(t, err)

		label1, err := createLabel(tt.Host, product.Name)
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

			res, err := request.Get(fmt.Sprintf("%s/users/%s/labels:cache?product=%s", tt.Host, users[0], product.Name)).
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
					Users: users,
				}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			res, err = request.Get(fmt.Sprintf("%s/users/%s/labels:cache?product=%s", tt.Host, users[1], product.Name)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.CacheLabelsInfoRes{}
			_, err = res.JSON(&json)
			assert.Nil(err)
			assert.Equal(1, len(json.Result))
			assert.True(json.Timestamp > 0)
			assert.Equal(label.Name, json.Result[0].Label)
			assert.Equal("", json.Result[0].Clients)
			assert.Equal("", json.Result[0].Channels)
		})

		t.Run(`"GET /users/:uid/labels:cache" when group label exists`, func(t *testing.T) {
			assert := assert.New(t)

			time.Sleep(time.Second)
			res, err := request.Post(fmt.Sprintf("%s/v1/products/%s/labels/%s:assign", tt.Host, product.Name, label1.Name)).
				Set("Content-Type", "application/json").
				Send(tpl.UsersGroupsBody{
					Groups: []string{group},
				}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			// users[1] lables from cache
			res, err = request.Get(fmt.Sprintf("%s/users/%s/labels:cache?product=%s", tt.Host, users[1], product.Name)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.CacheLabelsInfoRes{}
			_, err = res.JSON(&json)
			assert.Nil(err)
			assert.Equal(1, len(json.Result))
			assert.True(json.Timestamp > 0)
			assert.Equal(label.Name, json.Result[0].Label)
			assert.Equal("", json.Result[0].Clients)
			assert.Equal("", json.Result[0].Channels)

			// users[2] get all lables
			res, err = request.Get(fmt.Sprintf("%s/users/%s/labels:cache?product=%s", tt.Host, users[2], product.Name)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json = tpl.CacheLabelsInfoRes{}
			_, err = res.JSON(&json)
			assert.Nil(err)
			assert.Equal(2, len(json.Result))
			assert.True(json.Timestamp > 0)
			assert.Equal(label1.Name, json.Result[0].Label)
			assert.Equal("", json.Result[0].Clients)
			assert.Equal("", json.Result[0].Channels)

			assert.Equal(label.Name, json.Result[1].Label)

			time.Sleep(time.Second)
			res, err = request.Post(fmt.Sprintf("%s/v1/products/%s/labels/%s:assign", tt.Host, product.Name, label.Name)).
				Set("Content-Type", "application/json").
				Send(tpl.UsersGroupsBody{
					Groups: []string{group},
				}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			// users[2] get all lables
			res, err = request.Get(fmt.Sprintf("%s/users/%s/labels:cache?product=%s", tt.Host, users[3], product.Name)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json = tpl.CacheLabelsInfoRes{}
			_, err = res.JSON(&json)
			assert.Nil(err)
			assert.Equal(2, len(json.Result))
			assert.True(json.Timestamp > 0)
			assert.Equal(label.Name, json.Result[0].Label)
			assert.Equal("", json.Result[0].Clients)
			assert.Equal("", json.Result[0].Channels)

			assert.Equal(label1.Name, json.Result[1].Label)
		})

		t.Run(`"GET /users/:uid/labels"`, func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Get(fmt.Sprintf("%s/v1/users/%s/labels", tt.Host, users[3])).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.LabelsInfoRes{}
			_, err = res.JSON(&json)
			assert.Nil(err)
			assert.Equal(1, len(json.Result))
			assert.Equal(service.HIDer.HID(label), json.Result[0].HID)
			assert.Equal(product.Name, json.Result[0].Product)
			assert.Equal(label.Name, json.Result[0].Name)
			assert.Equal("", json.Result[0].Clients)
			assert.Equal("", json.Result[0].Channels)
		})

		t.Run(`Delete label should work`, func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Delete(fmt.Sprintf("%s/v1/users/%s/labels/%s", tt.Host, users[3], service.HIDer.HID(label))).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.BoolRes{}
			_, err = res.JSON(&json)
			assert.Nil(err)
			assert.True(json.Result)

			res, err = request.Delete(fmt.Sprintf("%s/v1/groups/%s/labels/%s", tt.Host, group, service.HIDer.HID(label))).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json = tpl.BoolRes{}
			_, err = res.JSON(&json)
			assert.Nil(err)
			assert.True(json.Result)

			assert.Nil(cleanupUserLabels(tt.DB, users[3]))
			res, err = request.Get(fmt.Sprintf("%s/users/%s/labels:cache?product=%s", tt.Host, users[3], product.Name)).
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
}
