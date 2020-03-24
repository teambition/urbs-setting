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

func createGroup(tt *TestTools) (group schema.Group, err error) {
	uid := tpl.RandUID()
	_, err = request.Post(fmt.Sprintf("%s/v1/groups:batch", tt.Host)).
		Set("Content-Type", "application/json").
		Send(tpl.GroupsBody{Groups: []tpl.GroupBody{
			tpl.GroupBody{UID: uid, Kind: "org", Desc: uid},
		}}).
		End()

	if err == nil {
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

	_, err = request.Post(fmt.Sprintf("%s/v1/groups:batch", tt.Host)).
		Set("Content-Type", "application/json").
		Send(tpl.GroupsBody{Groups: []tpl.GroupBody{
			tpl.GroupBody{UID: groupUID, Kind: "org", Desc: groupUID},
		}}).
		End()

	if err == nil {
		_, err = request.Post(fmt.Sprintf("%s/v1/groups/%s/members:batch", tt.Host, groupUID)).
			Set("Content-Type", "application/json").
			Send(tpl.UsersBody{Users: userUIDs}).
			End()
	}

	if err == nil {
		err = tt.DB.Where("uid= ?", groupUID).First(&group).Error
	}

	if err == nil {
		err = tt.DB.Where("uid in ( ? )", userUIDs).Find(&users).Error
	}
	return
}

func TestGroupAPIs(t *testing.T) {
	tt, cleanup := SetUpTestTools()
	defer cleanup()

	uid1 := tpl.RandUID()
	uid2 := tpl.RandUID()

	user := tpl.RandUID()
	users, err := createUsers(tt, 5)
	assert.Nil(t, err)

	t.Run(`"POST /v1/groups:batch"`, func(t *testing.T) {
		t.Run("should work", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Post(fmt.Sprintf("%s/v1/groups:batch", tt.Host)).
				Set("Content-Type", "application/json").
				Send(tpl.GroupsBody{Groups: []tpl.GroupBody{
					tpl.GroupBody{UID: uid1, Kind: "org", Desc: "test"},
				}}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.BoolRes{}
			res.JSON(&json)

			assert.True(json.Result)
		})

		t.Run("should work with duplicate group", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Post(fmt.Sprintf("%s/v1/groups:batch", tt.Host)).
				Set("Content-Type", "application/json").
				Send(tpl.GroupsBody{Groups: []tpl.GroupBody{
					tpl.GroupBody{UID: uid1, Kind: "org", Desc: "test"},
					tpl.GroupBody{UID: uid2, Kind: "project", Desc: "test"},
				}}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.BoolRes{}
			res.JSON(&json)
			assert.True(json.Result)
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
			assert.NotNil(json.Result)
			assert.True(len(json.Result) > 0)
			assert.Equal("org", json.Result[0].Kind)

			res, err = request.Get(fmt.Sprintf("%s/v1/groups?kind=org", tt.Host)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			text, err = res.Text()
			assert.Nil(err)
			assert.True(strings.Contains(text, uid1))
			assert.True(strings.Contains(text, `"kind":"org"`))
			assert.False(strings.Contains(text, `"kind":"project"`))
			assert.False(strings.Contains(text, `"id"`))

			json2 := tpl.GroupsRes{}
			res.JSON(&json2)
			assert.True(json2.TotalSize < json.TotalSize)
			assert.True(len(json2.Result) > 0)
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

	t.Run(`"POST /v1/groups/:uid/members:batch"`, func(t *testing.T) {
		t.Run("should work", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Post(fmt.Sprintf("%s/v1/groups/%s/members:batch", tt.Host, uid1)).
				Set("Content-Type", "application/json").
				Send(tpl.UsersBody{Users: schema.GetUsersUID(users)}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.BoolRes{}
			res.JSON(&json)
			assert.True(json.Result)
		})

		t.Run("should work with duplicate user", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Post(fmt.Sprintf("%s/v1/groups/%s/members:batch", tt.Host, uid1)).
				Set("Content-Type", "application/json").
				Send(tpl.UsersBody{Users: append(schema.GetUsersUID(users), user)}).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.BoolRes{}
			res.JSON(&json)
			assert.True(json.Result)
		})
	})

	t.Run(`"GET /v1/groups/:uid/members"`, func(t *testing.T) {
		t.Run("should work", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Get(fmt.Sprintf("%s/v1/groups/%s/members", tt.Host, uid1)).
				End()

			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			text, err := res.Text()
			assert.Nil(err)
			assert.True(strings.Contains(text, users[0].UID))
			assert.True(strings.Contains(text, user))
			assert.False(strings.Contains(text, `"id"`))

			json := tpl.GroupMembersRes{}
			res.JSON(&json)
			assert.NotNil(json.Result)
			assert.Equal(6, len(json.Result))
		})
	})

	t.Run(`"DELETE /groups/:uid/members"`, func(t *testing.T) {
		t.Run("should work with user query", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Delete(fmt.Sprintf("%s/v1/groups/%s/members?user=%s", tt.Host, uid1, user)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.BoolRes{}
			res.JSON(&json)
			assert.True(json.Result)

			res, err = request.Get(fmt.Sprintf("%s/v1/groups/%s/members", tt.Host, uid1)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)
			json2 := tpl.GroupMembersRes{}
			res.JSON(&json2)
			assert.NotNil(json2.Result)
			assert.Equal(5, len(json2.Result))
		})

		t.Run("should work idempotent", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Delete(fmt.Sprintf("%s/v1/groups/%s/members?user=%s", tt.Host, uid1, user)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.BoolRes{}
			res.JSON(&json)
			assert.True(json.Result)
		})

		t.Run("should work with sync_lt query", func(t *testing.T) {
			assert := assert.New(t)

			syncLt := time.Now().UTC().Unix() + 60
			res, err := request.Delete(fmt.Sprintf("%s/v1/groups/%s/members?sync_lt=%d", tt.Host, uid1, syncLt)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)

			json := tpl.BoolRes{}
			res.JSON(&json)
			assert.True(json.Result)

			res, err = request.Get(fmt.Sprintf("%s/v1/groups/%s/members", tt.Host, uid1)).
				End()
			assert.Nil(err)
			assert.Equal(200, res.StatusCode)
			json2 := tpl.GroupMembersRes{}
			res.JSON(&json2)
			assert.NotNil(json2.Result)
			assert.Equal(0, len(json2.Result))
		})
	})

	// t.Run(`"GET /v1/groups/:uid/labels"`, func(t *testing.T) {
	// 	t.Run("should work", func(t *testing.T) {
	// 		assert := assert.New(t)
	// 	})
	// })
}
