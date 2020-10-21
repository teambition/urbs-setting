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

func TestLabelAPIsV2(t *testing.T) {
	tt, cleanup := SetUpTestTools()
	defer cleanup()

	product, err := createProduct(tt)
	assert.Nil(t, err)

	t.Run(`POST "/v2/products/:product/labels/:label+:assign"`, func(t *testing.T) {
		label, err := createLabel(tt, product.Name)
		assert.Nil(t, err)

		users, err := createUsers(tt, 3)
		assert.Nil(t, err)

		group, err := createGroup(tt)
		assert.Nil(t, err)

		t.Run("should work", func(t *testing.T) {
			assert := assert.New(t)

			res, err := request.Post(fmt.Sprintf("%s/v2/products/%s/labels/%s:assign", tt.Host, product.Name, label.Name)).
				Set("Content-Type", "application/json").
				Send(tpl.UsersGroupsBodyV2{
					Users: schema.GetUsersUID(users[0:2]),
					Groups: []*tpl.GroupKindUID{
						{UID: group.UID, Kind: dto.GroupOrgKind},
					},
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
			res, err := request.Post(fmt.Sprintf("%s/v2/products/%s/labels/%s:assign", tt.Host, product.Name, label.Name)).
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
}
