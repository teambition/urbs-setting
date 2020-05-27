package model

import (
	"context"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/teambition/urbs-setting/src/schema"
	"github.com/teambition/urbs-setting/src/service"
	"github.com/teambition/urbs-setting/src/tpl"
	"github.com/teambition/urbs-setting/src/util"
)

// LabelRule ...
type LabelRule struct {
	*Model
}

// ApplyRules ...
func (m *LabelRule) ApplyRules(ctx context.Context, userID int64, excludeLabels []int64) (int, error) {
	rules := []schema.LabelRule{}
	// 不把 excludeLabels 放入查询条件，从而尽量复用查询缓存
	sd := m.DB.From(schema.TableLabelRule).
		Where(goqu.C("kind").Eq("userPercent")).Order(goqu.C("updated_at").Desc()).Limit(200)
	err := sd.Executor().ScanStructsContext(ctx, &rules)
	if err != nil {
		return 0, err
	}
	ids := make([]interface{}, 0)
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
		sd := m.DB.Insert(schema.TableUserLabel).Cols("user_id", "label_id", "rls").
			FromQuery(goqu.From(goqu.T(schema.TableLabelRule).As("t1")).
				Select(goqu.V(userID), goqu.I("t1.label_id"), goqu.I("t1.id")).
				Where(goqu.I("t1.id").In(ids...))).
			OnConflict(goqu.DoNothing())
		rowsAffected, err := service.DeResult(sd.Executor().ExecContext(ctx))
		if err != nil {
			return 0, err
		}

		if rowsAffected > 0 {
			util.Go(5*time.Second, func(gctx context.Context) {
				m.tryIncreaseLabelsStatus(gctx, labelIDs, 1)
			})
		}
	}
	return len(ids), nil
}

// Acquire ...
func (m *LabelRule) Acquire(ctx context.Context, labelRuleID int64) (*schema.LabelRule, error) {
	labelRule := &schema.LabelRule{}
	if err := m.findOneByID(ctx, schema.TableLabelRule, labelRuleID, labelRule); err != nil {
		return nil, err
	}
	return labelRule, nil
}

// Find ...
func (m *LabelRule) Find(ctx context.Context, productID, labelID int64) ([]schema.LabelRule, error) {
	labelRules := make([]schema.LabelRule, 0)
	sd := m.DB.From(schema.TableLabelRule).
		Where(goqu.C("product_id").Eq(productID), goqu.C("label_id").Eq(labelID)).
		Order(goqu.C("id").Desc()).Limit(10)

	err := sd.Executor().ScanStructsContext(ctx, &labelRules)
	if err != nil {
		return nil, err
	}
	return labelRules, nil
}

// Create ...
func (m *LabelRule) Create(ctx context.Context, labelRule *schema.LabelRule) error {
	_, err := m.createOne(ctx, schema.TableLabelRule, labelRule)
	return err
}

// Update ...
func (m *LabelRule) Update(ctx context.Context, labelRuleID int64, changed map[string]interface{}) (*schema.LabelRule, error) {
	labelRule := &schema.LabelRule{}
	if _, err := m.updateByID(ctx, schema.TableLabelRule, labelRuleID, goqu.Record(changed)); err != nil {
		return nil, err
	}
	if err := m.findOneByID(ctx, schema.TableLabelRule, labelRuleID, labelRule); err != nil {
		return nil, err
	}
	return labelRule, nil
}

// Delete ...
func (m *LabelRule) Delete(ctx context.Context, id int64) (int64, error) {
	return m.deleteByID(ctx, schema.TableLabelRule, id)
}
