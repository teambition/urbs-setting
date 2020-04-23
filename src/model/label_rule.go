package model

import (
	"context"

	"github.com/teambition/urbs-setting/src/schema"
	"github.com/teambition/urbs-setting/src/tpl"
)

// LabelRule ...
type LabelRule struct {
	*Model
}

const applyLabelRulesSQL = "insert ignore into `user_label` (`user_id`, `label_id`, `rls`) " +
	"select ?, t1.`label_id`, t1.`rls` from `label_rule` t1 where t1.`id` in ( ? )"

// ApplyRules ...
func (m *LabelRule) ApplyRules(ctx context.Context, userID int64, excludeLabels []int64) (int, error) {
	rules := []schema.LabelRule{}
	// 不把 excludeLabels 放入查询条件，从而尽量复用查询缓存
	err := m.DB.Order("`updated_at` desc").Limit(200).Find(&rules).Error
	if err != nil {
		return 0, err
	}
	ids := make([]int64, 0)
	labelIDs := make([]int64, 0)
	for _, rule := range rules {
		if tpl.Int64SliceHas(excludeLabels, rule.LabelID) {
			continue
		}

		p := rule.ToPercent()

		if p > 0 && (int((userID+rule.CreatedAt.Unix())%100) <= p) {
			// 百分比规则无效或者用户不在百分比区间内
			ids = append(ids, rule.ID)
			labelIDs = append(labelIDs, rule.LabelID)
		}
	}

	if len(ids) > 0 {
		if err = m.DB.Exec(applyLabelRulesSQL, userID, ids).Error; err != nil {
			return 0, err
		}

		go m.increaseLabelsStatus(ctx, labelIDs, 1)
	}
	return len(ids), nil
}

// Acquire ...
func (m *LabelRule) Acquire(ctx context.Context, labelRuleID int64) (*schema.LabelRule, error) {
	labelRule := &schema.LabelRule{ID: labelRuleID}
	if err := m.DB.First(labelRule).Error; err != nil {
		return nil, err
	}
	return labelRule, nil
}

// Find ...
func (m *LabelRule) Find(ctx context.Context, productID, labelID int64) ([]schema.LabelRule, error) {
	labelRules := make([]schema.LabelRule, 0)
	err := m.DB.Where("`product_id` = ? and `label_id` = ?", productID, labelID).
		Order("`id`").Limit(10).Find(&labelRules).Error
	return labelRules, err
}

// Create ...
func (m *LabelRule) Create(ctx context.Context, labelRule *schema.LabelRule) error {
	err := m.DB.Create(labelRule).Error
	return err
}

// Update ...
func (m *LabelRule) Update(ctx context.Context, labelRuleID int64, changed map[string]interface{}) (*schema.LabelRule, error) {
	labelRule := &schema.LabelRule{ID: labelRuleID}
	if len(changed) > 0 {
		if err := m.DB.Model(labelRule).UpdateColumns(changed).Error; err != nil {
			return nil, err
		}
	}

	if err := m.DB.First(labelRule).Error; err != nil {
		return nil, err
	}
	return labelRule, nil
}

// Delete ...
func (m *LabelRule) Delete(ctx context.Context, labelRuleID int64) error {
	res := m.DB.Delete(&schema.LabelRule{ID: labelRuleID})
	return res.Error
}
