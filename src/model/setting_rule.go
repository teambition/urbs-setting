package model

import (
	"context"

	"github.com/teambition/urbs-setting/src/schema"
)

// SettingRule ...
type SettingRule struct {
	*Model
}

const applySettingRulesSQL = "insert ignore into `user_setting` (`user_id`, `setting_id`, `rls`, `value`) " +
	"select ?, t1.`setting_id`, t1.`rls`, t1.`value` from `setting_rule` t1 where t1.`id` in ( ? )"

const findSettingRulesResultSQL = "select t1.`setting_id` " +
	"from `setting_rule` t1 `user_setting` t2 " +
	"where t1.`id` in ( ? ) and t2.`user_id` = ? and t1.`setting_id` = t2.`setting_id` and t1.`rls` = t2.`rls`"

// ApplyRules ...
func (m *SettingRule) ApplyRules(ctx context.Context, productID, userID int64) error {
	rules := []schema.SettingRule{}
	// 不把 excludeSettings 放入查询条件，从而尽量复用查询缓存
	err := m.DB.Where("`product_id` = ?", productID).Order("`updated_at` desc").Limit(1000).Find(&rules).Error
	if err != nil {
		return err
	}
	ids := make([]int64, 0)

	for _, rule := range rules {
		p := rule.ToPercent()
		if p > 0 && (int((userID+rule.CreatedAt.Unix())%100) <= p) {
			// 百分比规则无效或者用户不在百分比区间内
			ids = append(ids, rule.ID)
		}
	}

	if len(ids) > 0 {
		res := m.DB.Exec(applySettingRulesSQL, userID, ids)
		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected > 0 {
			rows, err := m.DB.Raw(findSettingRulesResultSQL, ids, userID).Rows()
			defer rows.Close()

			if err != nil {
				return err
			}

			settingIDs := make([]int64, 0)
			for rows.Next() {
				var settingID int64
				if err := rows.Scan(&settingID); err != nil {
					return err
				}
				settingIDs = append(settingIDs, settingID)
			}
			if len(settingIDs) > 0 {
				m.increaseSettingsStatus(ctx, settingIDs, 1)
			}
		}
	}
	return nil
}

// Acquire ...
func (m *SettingRule) Acquire(ctx context.Context, settingRuleID int64) (*schema.SettingRule, error) {
	settingRule := &schema.SettingRule{ID: settingRuleID}
	if err := m.DB.First(settingRule).Error; err != nil {
		return nil, err
	}
	return settingRule, nil
}

// Find ...
func (m *SettingRule) Find(ctx context.Context, productID, settingID int64) ([]schema.SettingRule, error) {
	settingRules := make([]schema.SettingRule, 0)
	err := m.DB.Where("`product_id` = ? and `setting_id` = ?", productID, settingID).
		Order("`id`").Limit(10).Find(&settingRules).Error
	return settingRules, err
}

// Create ...
func (m *SettingRule) Create(ctx context.Context, settingRule *schema.SettingRule) error {
	err := m.DB.Create(settingRule).Error
	return err
}

// Update ...
func (m *SettingRule) Update(ctx context.Context, settingRuleID int64, changed map[string]interface{}) (*schema.SettingRule, error) {
	settingRule := &schema.SettingRule{ID: settingRuleID}
	if len(changed) > 0 {
		if err := m.DB.Model(settingRule).UpdateColumns(changed).Error; err != nil {
			return nil, err
		}
	}

	if err := m.DB.First(settingRule).Error; err != nil {
		return nil, err
	}
	return settingRule, nil
}

// Delete ...
func (m *SettingRule) Delete(ctx context.Context, settingRuleID int64) error {
	res := m.DB.Delete(&schema.SettingRule{ID: settingRuleID})
	return res.Error
}
