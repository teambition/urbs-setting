package api

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/DavidCai1993/request"
	"github.com/stretchr/testify/assert"
	"github.com/teambition/urbs-setting/src/tpl"
)

func createGroup(appHost string) (string, error) {
	groupUID := tpl.RandUID()

	_, err := request.Post(fmt.Sprintf("%s/v1/groups:batch", appHost)).
		Set("Content-Type", "application/json").
		Send(tpl.GroupsBody{Groups: []tpl.GroupBody{
			tpl.GroupBody{UID: groupUID, Desc: groupUID},
		}}).
		End()

	if err != nil {
		return "", err
	}
	return groupUID, nil
}

func createGroupWithUsers(appHost string, count int) (group string, users []string, err error) {
	group = tpl.RandUID()

	users = make([]string, count)
	for i := 0; i < count; i++ {
		users[i] = tpl.RandUID()
	}

	_, err = request.Post(fmt.Sprintf("%s/v1/groups:batch", appHost)).
		Set("Content-Type", "application/json").
		Send(tpl.GroupsBody{Groups: []tpl.GroupBody{
			tpl.GroupBody{UID: group, Desc: group},
		}}).
		End()

	if err != nil {
		return "", nil, err
	}

	_, err = request.Post(fmt.Sprintf("%s/v1/groups/%s/members:batch", appHost, group)).
		Set("Content-Type", "application/json").
		Send(tpl.UsersBody{Users: users}).
		End()

	return group, users, nil
}

func TestGroupAPIs(t *testing.T) {
	tt, cleanup := SetUpTestTools()
	defer cleanup()

	uid1 := tpl.RandUID()
	uid2 := tpl.RandUID()

	user := tpl.RandUID()
	users, err := createUsers(tt.Host, 5)
	assert.Nil(t, err)

	t.Run(`"POST /v1/groups:batch"`, func(t *testing.T) {
		t.Run("should work", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Post(fmt.Sprintf("%s/v1/groups:batch", tt.Host)).
				Set("Content-Type", "application/json").
				Send(tpl.GroupsBody{Groups: []tpl.GroupBody{
					tpl.GroupBody{UID: uid1, Desc: "test"},
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
					tpl.GroupBody{UID: uid1, Desc: "test"},
					tpl.GroupBody{UID: uid2, Desc: "test"},
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

			json := tpl.GroupsRes{}
			res.JSON(&json)
			assert.NotNil(json.Result)
			assert.True(len(json.Result) > 0)
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

	t.Run(`"POST /groups/:uid/members:batch"`, func(t *testing.T) {
		t.Run("should work", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Post(fmt.Sprintf("%s/v1/groups/%s/members:batch", tt.Host, uid1)).
				Set("Content-Type", "application/json").
				Send(tpl.UsersBody{Users: users}).
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
				Send(tpl.UsersBody{Users: append(users, user)}).
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
			assert.True(strings.Contains(text, users[0]))
			assert.True(strings.Contains(text, user))

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
