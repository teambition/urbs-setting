package api

import (
	"fmt"
	"testing"

	"github.com/DavidCai1993/request"
	"github.com/stretchr/testify/assert"
	"github.com/teambition/urbs-setting/src/dto"
	"github.com/teambition/urbs-setting/src/schema"
	"github.com/teambition/urbs-setting/src/tpl"
)

func TestSettingAPIsV2(t *testing.T) {
	tt, cleanup := SetUpTestTools()
	defer cleanup()

	product, err := createProduct(tt)
	assert.Nil(t, err)

	t.Run(`POST "/v2/products/:product/modules/:module/settings/:setting+:assign"`, func(t *testing.T) {
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

			res, err := request.Post(fmt.Sprintf("%s/v2/products/%s/modules/%s/settings/%s:assign", tt.Host, product.Name, module.Name, setting.Name)).
				Set("Content-Type", "application/json").
				Send(tpl.UsersGroupsBodyV2{
					Users:  schema.GetUsersUID(users[0:2]),
					Groups: []*tpl.GroupKindUID{{UID: group.UID, Kind: dto.GroupOrgKind}},
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
			res, err := request.Post(fmt.Sprintf("%s/v2/products/%s/modules/%s/settings/%s:assign", tt.Host, product.Name, module.Name, setting.Name)).
				Set("Content-Type", "application/json").
				Send(tpl.UsersGroupsBodyV2{
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
}
