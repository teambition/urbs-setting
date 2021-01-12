package model

import (
	"context"
	"hash/crc32"
	"math/rand"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
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
func (m *LabelRule) ApplyRules(ctx context.Context, productID int64, userID int64, excludeLabels []int64, kind string) (int, error) {
	rules := []schema.LabelRule{}
	exps := []exp.Expression{goqu.C("kind").Eq(kind)}
	if productID > 0 {
		exps = append(exps, goqu.C("product_id").Eq(productID))
	}
	sd := m.RdDB.From(schema.TableLabelRule).Where(exps...).Order(goqu.C("updated_at").Desc()).Limit(200)
	err := sd.Executor().ScanStructsContext(ctx, &rules)
	if err != nil {
		return 0, err
	}
	// 不把 excludeLabels 放入查询条件，从而尽量复用查询缓存
	res, err := m.ComputeUserRule(ctx, userID, excludeLabels, rules)
	if err != nil {
		return 0, err
	}
	return res, nil
}

// ApplyRule ...
func (m *LabelRule) ApplyRule(ctx context.Context, productID int64, userID int64, labelID int64, kind string) (int, error) {
	rules := []schema.LabelRule{}
	exps := []exp.Expression{
		goqu.C("kind").Eq(kind),
		goqu.C("label_id").Eq(labelID),
		goqu.C("product_id").Eq(productID),
	}
	sd := m.RdDB.From(schema.TableLabelRule).Where(exps...).Order(goqu.C("updated_at").Desc()).Limit(200)
	err := sd.Executor().ScanStructsContext(ctx, &rules)
	if err != nil {
		return 0, err
	}
	res, err := m.ComputeUserRule(ctx, userID, []int64{}, rules)
	if err != nil {
		return 0, err
	}
	return res, nil
}

// ComputeUserRule ...
func (m *LabelRule) ComputeUserRule(ctx context.Context, userID int64, excludeLabels []int64, rules []schema.LabelRule) (int, error) {
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
		index := rand.Intn(len(ids))
		ids = []interface{}{ids[index]}
		labelIDs = []int64{labelIDs[index]}

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

// ApplyRulesToAnonymous ...
func (m *LabelRule) ApplyRulesToAnonymous(ctx context.Context, anonymousID string, productID int64, kind string) ([]schema.UserCacheLabel, error) {
	rules := []schema.LabelRule{}
	sd := m.RdDB.From(schema.TableLabelRule).
		Where(
			goqu.C("kind").Eq(kind),
			goqu.C("product_id").Eq(productID)).
		Order(goqu.C("updated_at").Desc()).Limit(200)
	err := sd.Executor().ScanStructsContext(ctx, &rules)
	if err != nil {
		return nil, err
	}

	anonID := int64(crc32.ChecksumIEEE([]byte(anonymousID)))
	labelIDs := make([]int64, 0)
	for _, rule := range rules {
		p := rule.ToPercent()
		if p > 0 && (int((anonID+rule.CreatedAt.Unix())%100) <= p) {
			// 百分比规则无效或者用户不在百分比区间内
			labelIDs = append(labelIDs, rule.LabelID)
		}
	}

	data := make([]schema.UserCacheLabel, 0)
	if len(labelIDs) > 0 {
		sd := m.RdDB.Select(
			goqu.I("t1.id"),
			goqu.I("t1.name"),
			goqu.I("t1.channels"),
			goqu.I("t1.clients")).
			From(goqu.T(schema.TableLabel).As("t1")).
			Where(goqu.I("t1.id").In(labelIDs))

		scanner, err := sd.Executor().ScannerContext(ctx)
		if err != nil {
			return nil, err
		}

		mp := make(map[int64]schema.UserCacheLabel)
		for scanner.Next() {
			myLabelInfo := schema.MyLabelInfo{}
			if err := scanner.ScanStruct(&myLabelInfo); err != nil {
				scanner.Close()
				return nil, err
			}

			mp[myLabelInfo.ID] = schema.UserCacheLabel{
				Label:    myLabelInfo.Name,
				Clients:  tpl.StringToSlice(myLabelInfo.Clients),
				Channels: tpl.StringToSlice(myLabelInfo.Channels),
			}
		}

		scanner.Close()
		if err := scanner.Err(); err != nil {
			return nil, err
		}

		for _, id := range labelIDs {
			if label, ok := mp[id]; ok {
				data = append(data, label)
			}
		}
	}

	return data, nil
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
	sd := m.RdDB.From(schema.TableLabelRule).
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
