package bll

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/teambition/urbs-setting/src/model"
	"github.com/teambition/urbs-setting/src/schema"
	"github.com/teambition/urbs-setting/src/service"
	"github.com/teambition/urbs-setting/src/tpl"
)

func TestUsers(t *testing.T) {
	user := &User{ms: model.NewModels(service.NewDB())}
	product := &Product{ms: model.NewModels(service.NewDB())}

	t.Run("newUserPercent should work with apply setting rule", func(t *testing.T) {
		assert := assert.New(t)
		ctx := context.Background()
		uid1 := tpl.RandUID()

		user.BatchAdd(ctx, []string{uid1})
		dbUser, err := user.ms.User.FindByUID(context.WithValue(ctx, model.ReadDB, true), uid1, "id")
		assert.Nil(err)
		assert.True(dbUser.ID > 0)

		userIntID := dbUser.ID

		productName := tpl.RandName()
		productRes, err := product.Create(ctx, productName, productName)

		settingRule := &schema.SettingRule{
			ProductID: productRes.Result.ID,
			SettingID: 10000,
			Kind:      schema.RuleNewUserPercent,
			Rule:      `{"value": 100 }`,
			Value:     "a",
		}
		assert.Equal(100, settingRule.ToPercent())
		err = user.ms.SettingRule.Create(ctx, settingRule)
		assert.Nil(err)

		body := &tpl.ApplyRulesBody{
			Kind: schema.RuleNewUserPercent,
		}
		body.Users = []string{uid1}

		user.ApplyRules(context.Background(), productName, body)
		time.Sleep(100 * time.Millisecond)

		us := &schema.UserSetting{}
		_, err = user.ms.User.DB.ScanStruct(us, "select * from `user_setting` where `user_id` = ? limit 1", userIntID)
		assert.Nil(err, err)
		assert.Equal("a", us.Value)
		assert.Equal(settingRule.SettingID, us.SettingID)

		_, err = user.ms.User.DB.Exec("delete from `user_setting` where `user_id` = ?", userIntID)
		assert.Nil(err)

		_, err = user.ms.User.DB.Exec("delete from `setting_rule` where `setting_id` = ?", settingRule.SettingID)
		assert.Nil(err)
	})

	t.Run("newUserPercent should work with apply label rule", func(t *testing.T) {
		assert := assert.New(t)
		ctx := context.Background()
		uid1 := tpl.RandUID()

		user.BatchAdd(ctx, []string{uid1})
		dbUser, _ := user.ms.User.FindByUID(context.WithValue(ctx, model.ReadDB, true), uid1, "id")
		userIntID := dbUser.ID

		productName := tpl.RandName()
		productRes, err := product.Create(ctx, productName, productName)

		labelRule := &schema.LabelRule{
			ProductID: productRes.Result.ID,
			LabelID:   10001,
			Kind:      schema.RuleNewUserPercent,
			Rule:      `{"value": 100 }`,
		}
		assert.Equal(100, labelRule.ToPercent())
		err = user.ms.LabelRule.Create(ctx, labelRule)
		assert.Nil(err)

		body := &tpl.ApplyRulesBody{
			Kind: schema.RuleNewUserPercent,
		}
		body.Users = []string{uid1}

		user.ApplyRules(context.Background(), productName, body)
		time.Sleep(100 * time.Millisecond)

		ul := &schema.UserLabel{}
		_, err = user.ms.User.DB.ScanStruct(ul, "select * from `user_label` where `user_id` = ? limit 1", userIntID)
		assert.Nil(err, err)
		assert.Equal(labelRule.LabelID, ul.LabelID)

		_, err = user.ms.User.DB.Exec("delete from `user_label` where `user_id` = ?", userIntID)
		assert.Nil(err)

		_, err = user.ms.User.DB.Exec("delete from `label_rule` where `label_id` = ?", labelRule.LabelID)
		assert.Nil(err)
	})
}

func TestChildLabelUserPercent(t *testing.T) {

	user := &User{ms: model.NewModels(service.NewDB())}
	product := &Product{ms: model.NewModels(service.NewDB())}

	t.Run("childLabelUserPercen should work with apply rule", func(t *testing.T) {
		assert := assert.New(t)
		require := require.New(t)
		ctx := context.Background()

		uid1 := tpl.RandUID()
		user.BatchAdd(ctx, []string{uid1})
		userObj, err := user.ms.User.Acquire(ctx, uid1)
		assert.Nil(err)

		productName := tpl.RandName()
		productRes, err := product.Create(ctx, productName, productName)

		label := &schema.Label{
			ProductID: productRes.Result.ID,
			Name:      tpl.RandName(),
		}
		err = user.ms.Label.Create(ctx, label)
		assert.Nil(err)
		labelRes, err := user.ms.Label.Acquire(ctx, productRes.Result.ID, label.Name)
		assert.Nil(err)

		labelRule := &schema.LabelRule{
			ProductID: productRes.Result.ID,
			LabelID:   labelRes.ID,
			Kind:      schema.RuleUserPercent,
			Rule:      `{"value": 100 }`,
		}
		assert.Equal(100, labelRule.ToPercent())
		err = user.ms.LabelRule.Create(ctx, labelRule)
		assert.Nil(err)

		label2 := &schema.Label{
			ProductID: productRes.Result.ID,
			Name:      labelRes.Name + "-gray",
		}
		err = user.ms.Label.Create(ctx, label2)
		assert.Nil(err)
		labelRes2, err := user.ms.Label.Acquire(ctx, productRes.Result.ID, label2.Name)
		assert.Nil(err)

		labelRule2 := &schema.LabelRule{
			ProductID: productRes.Result.ID,
			LabelID:   labelRes2.ID,
			Kind:      schema.RuleChildLabelUserPercent,
			Rule:      `{"value": 100 }`,
		}
		assert.Equal(100, labelRule.ToPercent())
		err = user.ms.LabelRule.Create(ctx, labelRule2)
		assert.Nil(err)

		userRes, err := user.ms.ApplyLabelRulesAndRefreshUserLabels(ctx, productRes.Result.ID, productName, userObj.ID, time.Now().UTC(), true)
		require.Nil(err)
		userLabels := userRes.GetLabels(productName)
		assert.True(len(userLabels) == 1)
		assert.Equal(labelRes.Name, userLabels[0].Label)

		userRes, err = user.ms.ApplyLabelRulesAndRefreshUserLabels(ctx, productRes.Result.ID, productName, userObj.ID, time.Now().UTC(), true)
		assert.Nil(err)
		userLabels = userRes.GetLabels(productName)
		assert.True(len(userLabels) == 2)
		assert.Equal(labelRes2.Name, userLabels[0].Label)
		assert.Equal(labelRes.Name, userLabels[1].Label)
	})
}
