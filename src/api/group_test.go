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

func createGroup(tt *TestTools) (group schema.Group, err error) {
	uid := tpl.RandUID()
	res, err := request.Post(fmt.Sprintf("%s/v1/groups:batch", tt.Host)).
		Set("Content-Type", "application/json").
		Send(tpl.GroupsBody{Groups: []tpl.GroupBody{
			{UID: uid, Kind: "org", Desc: uid},
		}}).
		End()

	if err == nil {
		res.Content() // close http client
		err = tt.DB.Where("uid= ?", uid).First(&group).Error
	}
	return
}

func createGroupWithUsers(tt *TestTools, count int) (group schema.Group, users []schema.User, err error) {
	groupUID := tpl.RandUID()
	userUIDs := make([]string, count)
	for i := 0; i < count; i++ {
		userUIDs[i] = tpl.RandUID()
	}

	res, err := request.Post(fmt.Sprintf("%s/v1/groups:batch", tt.Host)).
		Set("Content-Type", "application/json").
		Send(tpl.GroupsBody{Groups: []tpl.GroupBody{
			{UID: groupUID, Kind: "org", Desc: groupUID},
		}}).
		End()

	if err == nil {
		res.Content() // close http client
		res, err = request.Post(fmt.Sprintf("%s/v1/groups/%s/members:batch", tt.Host, groupUID)).
			Set("Content-Type", "application/json").
			Send(tpl.UsersBody{Users: userUIDs}).
			End()
	}

	if err == nil {
		res.Content() // close http client
		err = tt.DB.Where("uid= ?", groupUID).First(&group).Error
	}

	if err == nil {
		err = tt.DB.Where("uid in ( ? )", userUIDs).Order("id").Find(&users).Error
	}
	return
}

func TestGroupAPIs(t *testing.T) {
	tt, cleanup := SetUpTestTools()
	defer cleanup()

	uid1 := tpl.RandUID()
	uid2 := tpl.RandUID()

	t.Run(`"POST /v1/groups:batch"`, func(t *testing.T) {
		t.Run("should work", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Post(fmt.Sprintf("%s/v1/groups:batch", tt.Host)).
				Set("Content-Type", "application/json").
				Send(tpl.GroupsBody{Groups: []tpl.GroupBody{
					{UID: uid1, Kind: "org", Desc: "test"},
				}}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.BoolRes{}
			res.JSON(&json)
			assert.True(json.Result)

			var count int64
			assert.Nil(tt.DB.Table(`urbs_group`).Where("uid = ?", uid1).Count(&count).Error)
			assert.Equal(int64(1), count)

			var group schema.Group
			err = tt.DB.Where("uid= ?", uid1).First(&group).Error
			assert.Nil(err)
			assert.Equal("org", group.Kind)
			assert.Equal(group.CreatedAt, group.UpdatedAt)
		})

		t.Run("should work with duplicate group", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Post(fmt.Sprintf("%s/v1/groups:batch", tt.Host)).
				Set("Content-Type", "application/json").
				Send(tpl.GroupsBody{Groups: []tpl.GroupBody{
					{UID: uid1, Kind: "org", Desc: "test"},
					{UID: uid2, Kind: "project", Desc: "test"},
				}}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.BoolRes{}
			res.JSON(&json)
			assert.True(json.Result)

			var count int64
			assert.Nil(tt.DB.Table(`urbs_group`).Where("uid = ?", uid1).Count(&count).Error)
			assert.Equal(int64(1), count)
			assert.Nil(tt.DB.Table(`urbs_group`).Where("uid = ?", uid2).Count(&count).Error)
			assert.Equal(int64(1), count)

			var group schema.Group
			err = tt.DB.Where("uid= ?", uid1).First(&group).Error
			assert.Nil(err)
			assert.Equal("org", group.Kind)
			assert.Equal(group.CreatedAt, group.UpdatedAt)
		})

		t.Run(`should 400 if no group`, func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Post(fmt.Sprintf("%s/v1/groups:batch", tt.Host)).
				Set("Content-Type", "application/json").
				Send(tpl.GroupsBody{Groups: []tpl.GroupBody{}}).
				End()
			assert.Nil(err)
			assert.Equal(400, res.StatusCode)

			json := tpl.ResponseType{}
			res.JSON(&json)
			assert.Equal("BadRequest", json.Error)
			assert.Nil(json.Result)
		})

		t.Run(`should 400 if no kind`, func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Post(fmt.Sprintf("%s/v1/groups:batch", tt.Host)).
				Set("Content-Type", "application/json").
				Send(tpl.GroupsBody{Groups: []tpl.GroupBody{
					{UID: uid2, Kind: "", Desc: "test"},
				}}).
				End()
			assert.Nil(err)
			assert.Equal(400, res.StatusCode)

			json := tpl.ResponseType{}
			res.JSON(&json)
			assert.Equal("BadRequest", json.Error)
			assert.Nil(json.Result)
		})
	})

	t.Run(`"GET /v1/groups"`, func(t *testing.T) {
		t.Run("should work", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Get(fmt.Sprintf("%s/v1/groups", tt.Host)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			text, err := res.Text()
			assert.Nil(err)
			assert.True(strings.Contains(text, uid1))
			assert.True(strings.Contains(text, `"kind":"org"`))
			assert.True(strings.Contains(text, `"kind":"project"`))
			assert.False(strings.Contains(text, `"id"`))

			json := tpl.GroupsRes{}
			res.JSON(&json)
			assert.True(len(json.Result) > 0)
			assert.Equal("project", json.Result[0].Kind)
		})

		t.Run("should work with kind", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Get(fmt.Sprintf("%s/v1/groups?kind=org", tt.Host)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			text, err := res.Text()
			assert.Nil(err)
			assert.True(strings.Contains(text, uid1))
			assert.True(strings.Contains(text, `"kind":"org"`))
			assert.False(strings.Contains(text, `"kind":"project"`))
			assert.False(strings.Contains(text, `"id"`))

			json := tpl.GroupsRes{}
			res.JSON(&json)
			assert.True(len(json.Result) > 0)
			assert.Equal("org", json.Result[0].Kind)
		})
	})

	t.Run(`"GET /v1/groups/:uid+:exists"`, func(t *testing.T) {
		t.Run(`should work`, func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Get(fmt.Sprintf("%s/v1/groups/%s:exists", tt.Host, uid2)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.BoolRes{}
			res.JSON(&json)
			assert.True(json.Result)
		})

		t.Run(`should work if not exists`, func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Get(fmt.Sprintf("%s/v1/groups/%s:exists", tt.Host, tpl.RandUID())).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.BoolRes{}
			res.JSON(&json)
			assert.False(json.Result)
		})
	})

	t.Run(`"PUT /v1/groups/:uid"`, func(t *testing.T) {
		group, err := createGroup(tt)
		assert.Nil(t, err)

		t.Run("should work", func(t *testing.T) {
			assert := assert.New(t)

			desc := "abc"
			res, err := request.Put(fmt.Sprintf("%s/v1/groups/%s", tt.Host, group.UID)).
				Set("Content-Type", "application/json").
				Send(tpl.GroupUpdateBody{
					Desc: &desc,
				}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			text, err := res.Text()
			assert.Nil(err)
			assert.False(strings.Contains(text, `"id"`))

			json := tpl.GroupRes{}
			res.JSON(&json)
			assert.NotNil(json.Result)
			assert.Equal(group.UID, json.Result.UID)
			assert.Equal(desc, json.Result.Desc)
			assert.True(json.Result.UpdatedAt.After(json.Result.CreatedAt))

			// should work idempotent
			time.Sleep(time.Millisecond * 100)
			res, err = request.Put(fmt.Sprintf("%s/v1/groups/%s", tt.Host, group.UID)).
				Set("Content-Type", "application/json").
				Send(tpl.GroupUpdateBody{
					Desc: &desc,
				}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json2 := tpl.GroupRes{}
			res.JSON(&json2)
			assert.NotNil(json2.Result)
			assert.True(json2.Result.UpdatedAt.Equal(json.Result.UpdatedAt))
		})

		t.Run(`should work with syncAt`, func(t *testing.T) {
			assert := assert.New(t)

			syncAt := group.SyncAt + 1
			res, err := request.Put(fmt.Sprintf("%s/v1/groups/%s", tt.Host, group.UID)).
				Set("Content-Type", "application/json").
				Send(tpl.GroupUpdateBody{SyncAt: &syncAt}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			text, err := res.Text()
			assert.Nil(err)
			assert.False(strings.Contains(text, `"id"`))

			json := tpl.GroupRes{}
			res.JSON(&json)
			assert.NotNil(json.Result)
			assert.Equal(group.UID, json.Result.UID)
			assert.Equal(syncAt, json.Result.SyncAt)
		})

		t.Run(`should 400`, func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Put(fmt.Sprintf("%s/v1/groups/%s", tt.Host, group.UID)).
				Set("Content-Type", "application/json").
				Send(tpl.GroupUpdateBody{}).
				End()
			assert.Nil(err)
			assert.Equal(400, res.StatusCode)
			res.Content() // close http client
		})
	})

	t.Run(`"DELETE /v1/groups/:uid"`, func(t *testing.T) {
		product, err := createProduct(tt)
		assert.Nil(t, err)

		label, err := createLabel(tt, product.Name)
		assert.Nil(t, err)

		module, err := createModule(tt, product.Name)
		assert.Nil(t, err)

		setting, err := createSetting(tt, product.Name, module.Name, "x", "y")
		assert.Nil(t, err)

		group, _, err := createGroupWithUsers(tt, 10)
		assert.Nil(t, err)

		t.Run(`should work`, func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Post(fmt.Sprintf("%s/v1/products/%s/labels/%s:assign", tt.Host, product.Name, label.Name)).
				Set("Content-Type", "application/json").
				Send(tpl.UsersGroupsBody{
					Groups: []string{group.UID},
				}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)
			res.Content() // close http client

			var count int64
			assert.Nil(tt.DB.Table(`group_label`).Where("group_id = ?", group.ID).Count(&count).Error)
			assert.Equal(int64(1), count)

			res, err = request.Post(fmt.Sprintf("%s/v1/products/%s/modules/%s/settings/%s:assign", tt.Host, product.Name, module.Name, setting.Name)).
				Set("Content-Type", "application/json").
				Send(tpl.UsersGroupsBody{
					Groups: []string{group.UID},
					Value:  "x",
				}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)
			res.Content() // close http client

			assert.Nil(tt.DB.Table(`group_setting`).Where("group_id = ?", group.ID).Count(&count).Error)
			assert.Equal(int64(1), count)

			assert.Nil(tt.DB.Table(`user_group`).Where("group_id = ?", group.ID).Count(&count).Error)
			assert.Equal(int64(10), count)

			res, err = request.Delete(fmt.Sprintf("%s/v1/groups/%s", tt.Host, group.UID)).
				Set("Content-Type", "application/json").
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.BoolRes{}
			res.JSON(&json)
			assert.True(json.Result)

			assert.Nil(tt.DB.Table(`group_label`).Where("group_id = ?", group.ID).Count(&count).Error)
			assert.Equal(int64(0), count)
			assert.Nil(tt.DB.Table(`group_setting`).Where("group_id = ?", group.ID).Count(&count).Error)
			assert.Equal(int64(0), count)
			assert.Nil(tt.DB.Table(`user_group`).Where("group_id = ?", group.ID).Count(&count).Error)
			assert.Equal(int64(0), count)
			assert.Nil(tt.DB.Table(`urbs_group`).Where("id = ?", group.ID).Count(&count).Error)
			assert.Equal(int64(0), count)
		})

		t.Run("should work idempotent", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Delete(fmt.Sprintf("%s/v1/groups/%s", tt.Host, group.UID)).
				Set("Content-Type", "application/json").
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.BoolRes{}
			res.JSON(&json)
			assert.True(json.Result)
		})
	})

	t.Run(`"POST /v1/groups/:uid/members:batch"`, func(t *testing.T) {
		group, err := createGroup(tt)
		assert.Nil(t, err)

		user := tpl.RandUID()
		users, err := createUsers(tt, 5)
		assert.Nil(t, err)

		t.Run("should work", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Post(fmt.Sprintf("%s/v1/groups/%s/members:batch", tt.Host, group.UID)).
				Set("Content-Type", "application/json").
				Send(tpl.UsersBody{Users: schema.GetUsersUID(users)}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.BoolRes{}
			res.JSON(&json)
			assert.True(json.Result)

			var count int64
			assert.Nil(tt.DB.Table(`user_group`).Where("group_id = ?", group.ID).Count(&count).Error)
			assert.Equal(int64(5), count)
		})

		t.Run("should work with duplicate user", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Post(fmt.Sprintf("%s/v1/groups/%s/members:batch", tt.Host, group.UID)).
				Set("Content-Type", "application/json").
				Send(tpl.UsersBody{Users: append(schema.GetUsersUID(users), user)}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.BoolRes{}
			res.JSON(&json)
			assert.True(json.Result)

			var count int64
			assert.Nil(tt.DB.Table(`user_group`).Where("group_id = ?", group.ID).Count(&count).Error)
			assert.Equal(int64(6), count)
		})

		t.Run("should work when user not exists", func(t *testing.T) {
			assert := assert.New(t)

			u := tpl.RandUID()
			res, err := request.Post(fmt.Sprintf("%s/v1/groups/%s/members:batch", tt.Host, group.UID)).
				Set("Content-Type", "application/json").
				Send(tpl.UsersBody{Users: []string{u}}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.BoolRes{}
			res.JSON(&json)
			assert.True(json.Result)

			var count int64
			assert.Nil(tt.DB.Table(`user_group`).Where("group_id = ?", group.ID).Count(&count).Error)
			assert.Equal(int64(7), count)
			assert.Nil(tt.DB.Table(`urbs_user`).Where("uid = ?", u).Count(&count).Error)
			assert.Equal(int64(1), count)
		})
	})

	t.Run(`"GET /v1/groups/:uid/members"`, func(t *testing.T) {
		group, users, err := createGroupWithUsers(tt, 10)
		assert.Nil(t, err)

		t.Run("should work", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Get(fmt.Sprintf("%s/v1/groups/%s/members", tt.Host, group.UID)).
				End()

			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			text, err := res.Text()
			assert.Nil(err)
			assert.True(strings.Contains(text, users[0].UID))
			assert.False(strings.Contains(text, `"id"`))

			json := tpl.GroupMembersRes{}
			res.JSON(&json)
			assert.NotNil(json.Result)
			assert.Equal(10, len(json.Result))
			assert.Equal("", json.NextPageToken)
		})

		t.Run("should work with pagination", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Get(fmt.Sprintf("%s/v1/groups/%s/members?pageSize=5&pageToken=", tt.Host, group.UID)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.GroupMembersRes{}
			res.JSON(&json)
			assert.NotNil(json.Result)
			assert.Equal(5, len(json.Result))
			assert.NotEqual("", json.NextPageToken)

			res, err = request.Get(fmt.Sprintf("%s/v1/groups/%s/members?pageSize=5&pageToken=%s", tt.Host, group.UID, json.NextPageToken)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json = tpl.GroupMembersRes{}
			res.JSON(&json)
			assert.NotNil(json.Result)
			assert.Equal(5, len(json.Result))
			assert.Equal("", json.NextPageToken)
		})
	})

	t.Run(`"DELETE /groups/:uid/members"`, func(t *testing.T) {
		group, users, err := createGroupWithUsers(tt, 10)
		assert.Nil(t, err)

		t.Run("should work with user query", func(t *testing.T) {
			assert := assert.New(t)

			var count int64
			assert.Nil(tt.DB.Table(`user_group`).Where("group_id = ? and user_id = ?", group.ID, users[0].ID).Count(&count).Error)
			assert.Equal(int64(1), count)

			res, err := request.Delete(fmt.Sprintf("%s/v1/groups/%s/members?user=%s", tt.Host, group.UID, users[0].UID)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.BoolRes{}
			res.JSON(&json)
			assert.True(json.Result)

			assert.Nil(tt.DB.Table(`user_group`).Where("group_id = ? and user_id = ?", group.ID, users[0].ID).Count(&count).Error)
			assert.Equal(int64(0), count)
		})

		t.Run("should work idempotent", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Delete(fmt.Sprintf("%s/v1/groups/%s/members?user=%s", tt.Host, group.UID, users[0].UID)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.BoolRes{}
			res.JSON(&json)
			assert.True(json.Result)

			var count int64
			assert.Nil(tt.DB.Table(`user_group`).Where("group_id = ? and user_id = ?", group.ID, users[0].ID).Count(&count).Error)
			assert.Equal(int64(0), count)
		})

		t.Run("should work with syncLt query", func(t *testing.T) {
			assert := assert.New(t)

			var count int64
			assert.Nil(tt.DB.Table(`user_group`).Where("group_id = ?", group.ID).Count(&count).Error)
			assert.Equal(int64(9), count)

			time.Sleep(time.Second)
			syncAt := time.Now().Unix()
			assert.True(syncAt > group.SyncAt)

			res, err := request.Put(fmt.Sprintf("%s/v1/groups/%s", tt.Host, group.UID)).
				Set("Content-Type", "application/json").
				Send(tpl.GroupUpdateBody{SyncAt: &syncAt}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)
			res.Content() // close http client

			res, err = request.Post(fmt.Sprintf("%s/v1/groups/%s/members:batch", tt.Host, group.UID)).
				Set("Content-Type", "application/json").
				Send(tpl.UsersBody{Users: []string{tpl.RandUID(), users[1].UID}}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)
			res.Content() // close http client

			assert.Nil(tt.DB.Table(`user_group`).Where("group_id = ?", group.ID).Count(&count).Error)
			assert.Equal(int64(10), count)

			res, err = request.Delete(fmt.Sprintf("%s/v1/groups/%s/members?syncLt=%d", tt.Host, group.UID, syncAt)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.BoolRes{}
			res.JSON(&json)
			assert.True(json.Result)
			assert.Nil(tt.DB.Table(`user_group`).Where("group_id = ?", group.ID).Count(&count).Error)
			assert.Equal(int64(2), count)
		})
	})

	t.Run(`"GET /v1/groups/:uid/labels"`, func(t *testing.T) {
		product, err := createProduct(tt)
		assert.Nil(t, err)

		label1, err := createLabel(tt, product.Name)
		assert.Nil(t, err)

		label2, err := createLabel(tt, product.Name)
		assert.Nil(t, err)

		group, err := createGroup(tt)
		assert.Nil(t, err)

		t.Run("should work", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Get(fmt.Sprintf("%s/v1/groups/%s/labels", tt.Host, group.UID)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.LabelsInfoRes{}
			_, err = res.JSON(&json)

			assert.Nil(err)
			assert.Equal(0, len(json.Result))
			assert.Equal(0, json.TotalSize)

			res, err = request.Post(fmt.Sprintf("%s/v1/products/%s/labels/%s:assign", tt.Host, product.Name, label1.Name)).
				Set("Content-Type", "application/json").
				Send(tpl.UsersGroupsBody{
					Groups: []string{group.UID},
				}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)
			res.Content() // close http client

			res, err = request.Get(fmt.Sprintf("%s/v1/groups/%s/labels", tt.Host, group.UID)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			text, err := res.Text()
			assert.Nil(err)
			assert.False(strings.Contains(text, `"id"`))

			json = tpl.LabelsInfoRes{}
			_, err = res.JSON(&json)

			assert.Nil(err)
			assert.Equal(1, len(json.Result))
			assert.Equal(1, json.TotalSize)
			assert.Equal("", json.NextPageToken)
			assert.Equal(service.IDToHID(label1.ID, "label"), json.Result[0].HID)
			assert.Equal(product.Name, json.Result[0].Product)
			assert.Equal(label1.Name, json.Result[0].Name)
			assert.Equal(0, len(json.Result[0].Clients))
			assert.Equal(0, len(json.Result[0].Channels))
			assert.True(json.Result[0].CreatedAt.After(time2020))

			res, err = request.Post(fmt.Sprintf("%s/v1/products/%s/labels/%s:assign", tt.Host, product.Name, label2.Name)).
				Set("Content-Type", "application/json").
				Send(tpl.UsersGroupsBody{
					Groups: []string{group.UID},
				}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)
			res.Content() // close http client

			res, err = request.Get(fmt.Sprintf("%s/v1/groups/%s/labels", tt.Host, group.UID)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			text, err = res.Text()
			assert.Nil(err)
			assert.False(strings.Contains(text, `"id"`))

			json = tpl.LabelsInfoRes{}
			_, err = res.JSON(&json)

			assert.Nil(err)
			assert.Equal(2, len(json.Result))
			assert.Equal(2, json.TotalSize)
			assert.Equal("", json.NextPageToken)
			assert.Equal(service.IDToHID(label2.ID, "label"), json.Result[0].HID)
		})
	})

	t.Run(`"DELETE /v1/groups/:uid/labels/:hid"`, func(t *testing.T) {
		product, err := createProduct(tt)
		assert.Nil(t, err)

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
			res.Content() // close http client

			var count int64
			assert.Nil(tt.DB.Table(`group_label`).Where("group_id = ?", group.ID).Count(&count).Error)
			assert.Equal(int64(1), count)

			res, err = request.Delete(fmt.Sprintf("%s/v1/groups/%s/labels/%s", tt.Host, group.UID, service.IDToHID(label.ID, "label"))).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.BoolRes{}
			_, err = res.JSON(&json)
			assert.Nil(err)
			assert.True(json.Result)

			assert.Nil(tt.DB.Table(`group_label`).Where("group_id = ?", group.ID).Count(&count).Error)
			assert.Equal(int64(0), count)
		})

		t.Run("should work idempotent", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Delete(fmt.Sprintf("%s/v1/groups/%s/labels/%s", tt.Host, group.UID, service.IDToHID(label.ID, "label"))).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.BoolRes{}
			_, err = res.JSON(&json)
			assert.Nil(err)
			assert.True(json.Result)
		})
	})

	t.Run(`"GET /v1/groups/:uid/settings"`, func(t *testing.T) {
		product, err := createProduct(tt)
		assert.Nil(t, err)

		module, err := createModule(tt, product.Name)
		assert.Nil(t, err)

		setting1, err := createSetting(tt, product.Name, module.Name, "a", "b")
		assert.Nil(t, err)

		setting2, err := createSetting(tt, product.Name, module.Name, "x", "y")
		assert.Nil(t, err)

		group, err := createGroup(tt)
		assert.Nil(t, err)

		t.Run("should work", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Get(fmt.Sprintf("%s/v1/groups/%s/settings", tt.Host, group.UID)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.MySettingsRes{}
			_, err = res.JSON(&json)

			assert.Nil(err)
			assert.Equal(0, len(json.Result))
			assert.Equal(0, json.TotalSize)

			res, err = request.Post(fmt.Sprintf("%s/v1/products/%s/modules/%s/settings/%s:assign", tt.Host, product.Name, module.Name, setting1.Name)).
				Set("Content-Type", "application/json").
				Send(tpl.UsersGroupsBody{
					Groups: []string{group.UID},
					Value:  "a",
				}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)
			res.Content() // close http client

			var count int64
			assert.Nil(tt.DB.Table(`group_setting`).Where("group_id = ?", group.ID).Count(&count).Error)
			assert.Equal(int64(1), count)

			res, err = request.Get(fmt.Sprintf("%s/v1/groups/%s/settings", tt.Host, group.UID)).
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
			assert.Equal(service.IDToHID(setting1.ID, "setting"), data.HID)
			assert.Equal(module.Name, data.Module)
			assert.Equal(setting1.Name, data.Name)
			assert.Equal("a", data.Value)
			assert.Equal("", data.LastValue)
			assert.True(data.CreatedAt.After(time2020))

			res, err = request.Post(fmt.Sprintf("%s/v1/products/%s/modules/%s/settings/%s:assign", tt.Host, product.Name, module.Name, setting2.Name)).
				Set("Content-Type", "application/json").
				Send(tpl.UsersGroupsBody{
					Groups: []string{group.UID},
					Value:  "x",
				}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)
			res.Content() // close http clients

			res, err = request.Get(fmt.Sprintf("%s/v1/groups/%s/settings", tt.Host, group.UID)).
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
			assert.Equal(service.IDToHID(setting2.ID, "setting"), data.HID)
			assert.Equal(module.Name, data.Module)
			assert.Equal(setting2.Name, data.Name)
			assert.Equal("x", data.Value)
			assert.Equal("", data.LastValue)
			assert.True(data.CreatedAt.After(time2020))
		})
	})

	t.Run(`"PUT /v1/groups/:uid/settings/:hid+:rollback"`, func(t *testing.T) {
		product, err := createProduct(tt)
		assert.Nil(t, err)

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
			res.Content() // close http client

			res, err = request.Get(fmt.Sprintf("%s/v1/groups/%s/settings", tt.Host, group.UID)).
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
			assert.Equal(service.IDToHID(setting.ID, "setting"), data.HID)
			assert.Equal(module.Name, data.Module)
			assert.Equal(setting.Name, data.Name)
			assert.Equal("x", data.Value)
			assert.Equal("", data.LastValue)
			assert.True(data.CreatedAt.After(time2020))

			res, err = request.Post(fmt.Sprintf("%s/v1/products/%s/modules/%s/settings/%s:assign", tt.Host, product.Name, module.Name, setting.Name)).
				Set("Content-Type", "application/json").
				Send(tpl.UsersGroupsBody{
					Groups: []string{group.UID},
					Value:  "y",
				}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)
			res.Content() // close http client

			res, err = request.Get(fmt.Sprintf("%s/v1/groups/%s/settings", tt.Host, group.UID)).
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

			data = json.Result[0]
			assert.Equal(service.IDToHID(setting.ID, "setting"), data.HID)
			assert.Equal(module.Name, data.Module)
			assert.Equal(setting.Name, data.Name)
			assert.Equal("y", data.Value)
			assert.Equal("x", data.LastValue)
			assert.True(data.CreatedAt.After(time2020))

			res, err = request.Put(fmt.Sprintf("%s/v1/groups/%s/settings/%s:rollback", tt.Host, group.UID, service.IDToHID(setting.ID, "setting"))).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)
			res.Content() // close http client

			json2 := tpl.BoolRes{}
			res.JSON(&json2)
			assert.True(json2.Result)

			res, err = request.Get(fmt.Sprintf("%s/v1/groups/%s/settings", tt.Host, group.UID)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json = tpl.MySettingsRes{}
			_, err = res.JSON(&json)
			data = json.Result[0]
			assert.Equal(service.IDToHID(setting.ID, "setting"), data.HID)
			assert.Equal("x", data.Value)
			assert.Equal("x", data.LastValue)
		})
	})

	t.Run(`"DELETE /v1/groups/:uid/settings/:hid"`, func(t *testing.T) {
		product, err := createProduct(tt)
		assert.Nil(t, err)

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
			res.Content() // close http client

			var count int64
			assert.Nil(tt.DB.Table(`group_setting`).Where("group_id = ?", group.ID).Count(&count).Error)
			assert.Equal(int64(1), count)

			res, err = request.Delete(fmt.Sprintf("%s/v1/groups/%s/settings/%s", tt.Host, group.UID, service.IDToHID(setting.ID, "setting"))).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.BoolRes{}
			_, err = res.JSON(&json)
			assert.Nil(err)
			assert.True(json.Result)

			assert.Nil(tt.DB.Table(`group_setting`).Where("group_id = ?", group.ID).Count(&count).Error)
			assert.Equal(int64(0), count)
		})

		t.Run("should work idempotent", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Delete(fmt.Sprintf("%s/v1/groups/%s/settings/%s", tt.Host, group.UID, service.IDToHID(setting.ID, "setting"))).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.BoolRes{}
			_, err = res.JSON(&json)
			assert.Nil(err)
			assert.True(json.Result)
		})
	})
}
