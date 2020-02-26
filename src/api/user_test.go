package api

import (
	"fmt"
	"testing"

	"github.com/DavidCai1993/request"
	"github.com/stretchr/testify/assert"
	"github.com/teambition/urbs-setting/src/tpl"
)

func createUsers(appHost string, count int) ([]string, error) {
	uids := make([]string, count)
	for i := 0; i < count; i++ {
		uids[i] = tpl.RandName()
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
}
