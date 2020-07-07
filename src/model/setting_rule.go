package model

import (
	"context"
	"hash/crc32"

	"github.com/doug-martin/goqu/v9"
	"github.com/teambition/urbs-setting/src/schema"
	"github.com/teambition/urbs-setting/src/service"
	"github.com/teambition/urbs-setting/src/tpl"
)

// SettingRule ...
type SettingRule struct {
	*Model
}

// ApplyRules ...
func (m *SettingRule) ApplyRules(ctx context.Context, productID, userID int64) error {
	rules := []schema.SettingRule{}
	sd := m.DB.From(schema.TableSettingRule).
		Where(goqu.C("product_id").Eq(productID), goqu.C("kind").Eq("userPercent")).
		Order(goqu.C("updated_at").Desc()).Limit(1000)
	err := sd.Executor().ScanStructsContext(ctx, &rules)
	if err != nil {
		return err
	}

	ids := make([]interface{}, 0)
	for _, rule := range rules {
		p := rule.ToPercent()
		if p > 0 && (int((userID+rule.CreatedAt.Unix())%100) <= p) {
			// 百分比规则无效或者用户不在百分比区间内
			ids = append(ids, rule.ID)
		}
	}

	if len(ids) > 0 {
		sd := m.DB.Insert(schema.TableUserSetting).Cols("user_id", "setting_id", "rls", "value").
			FromQuery(goqu.From(goqu.T(schema.TableSettingRule).As("t1")).
				Select(goqu.V(userID), goqu.I("t1.setting_id"), goqu.I("t1.rls"), goqu.I("t1.value")).
				Where(goqu.I("t1.id").In(ids...))).
			OnConflict(goqu.DoNothing())
		rowsAffected, err := service.DeResult(sd.Executor().ExecContext(ctx))
		if err != nil {
			return err
		}

		if rowsAffected > 0 {
			settingIDs := make([]int64, 0)
			sd := m.DB.Select(goqu.I("t1.setting_id")).
				From(
					goqu.T(schema.TableSettingRule).As("t1"),
					goqu.T(schema.TableUserSetting).As("t2")).
				Where(
					goqu.I("t1.id").In(ids...),
					goqu.I("t1.setting_id").Eq(goqu.I("t2.setting_id")),
					goqu.I("t1.rls").Eq(goqu.I("t2.rls")),
					goqu.I("t2.user_id").Eq(userID)).
				Limit(1000)
			if err := sd.Executor().ScanValsContext(ctx, &settingIDs); err != nil {
				return err
			}

			if len(settingIDs) > 0 {
				m.tryIncreaseSettingsStatus(ctx, settingIDs, 1)
			}
		}
	}
	return nil
}

// ApplyRulesToAnonymous ...
func (m *SettingRule) ApplyRulesToAnonymous(ctx context.Context, anonymousID string, productID int64, channel, client string) ([]tpl.MySetting, error) {
	rules := []schema.SettingRule{}
	sd := m.DB.From(schema.TableSettingRule).
		Where(goqu.C("product_id").Eq(productID), goqu.C("kind").Eq("userPercent")).
		Order(goqu.C("updated_at").Desc()).Limit(1000)
	err := sd.Executor().ScanStructsContext(ctx, &rules)
	if err != nil {
		return nil, err
	}

	anonID := int64(crc32.ChecksumIEEE([]byte(anonymousID)))
	ids := make([]interface{}, 0)
	for _, rule := range rules {
		p := rule.ToPercent()
		if p > 0 && (int((anonID+rule.CreatedAt.Unix())%100) <= p) {
			// 百分比规则无效或者用户不在百分比区间内
			ids = append(ids, rule.ID)
		}
	}

	data := make([]tpl.MySetting, 0)
	if len(ids) > 0 {
		sd := m.DB.Select(
			goqu.I("t1.rls"),
			goqu.I("t1.updated_at").As("assigned_at"),
			goqu.I("t1.value"),
			goqu.I("t2.id"),
			goqu.I("t2.name"),
			goqu.I("t2.description"),
			goqu.I("t2.channels"),
			goqu.I("t2.clients"),
			goqu.I("t3.name").As("module")).
			From(
				goqu.T(schema.TableSettingRule).As("t1"),
				goqu.T(schema.TableSetting).As("t2"),
				goqu.T(schema.TableModule).As("t3")).
			Where(
				goqu.I("t1.id").In(ids),
				goqu.I("t1.setting_id").Eq(goqu.I("t2.id")),
				goqu.I("t2.module_id").Eq(goqu.I("t3.id"))).
			Order(goqu.I("t1.updated_at").Desc())

		scanner, err := sd.Executor().ScannerContext(ctx)
		if err != nil {
			return nil, err
		}

		for scanner.Next() {
			mySetting := tpl.MySetting{}
			if err := scanner.ScanStruct(&mySetting); err != nil {
				scanner.Close()
				return nil, err
			}

			if mySetting.Channels != "" {
				if !tpl.StringSliceHas(tpl.StringToSlice(mySetting.Channels), channel) {
					continue // channel 不匹配
				}
			}
			if mySetting.Clients != "" {
				if !tpl.StringSliceHas(tpl.StringToSlice(mySetting.Clients), client) {
					continue // client 不匹配
				}
			}

			mySetting.HID = service.IDToHID(mySetting.ID, "setting")
			data = append(data, mySetting)
		}

		scanner.Close()
		if err := scanner.Err(); err != nil {
			return nil, err
		}
	}
	return data, nil
}

// Acquire ...
func (m *SettingRule) Acquire(ctx context.Context, settingRuleID int64) (*schema.SettingRule, error) {
	settingRule := &schema.SettingRule{}
	if err := m.findOneByID(ctx, schema.TableSettingRule, settingRuleID, settingRule); err != nil {
		return nil, err
	}
	return settingRule, nil
}

// Find ...
func (m *SettingRule) Find(ctx context.Context, productID, settingID int64) ([]schema.SettingRule, error) {
	settingRules := make([]schema.SettingRule, 0)
	sd := m.DB.From(schema.TableSettingRule).
		Where(goqu.C("product_id").Eq(productID), goqu.C("setting_id").Eq(settingID)).
		Order(goqu.C("id").Desc()).Limit(10)

	err := sd.Executor().ScanStructsContext(ctx, &settingRules)
	if err != nil {
		return nil, err
	}
	return settingRules, nil
}

// Create ...
func (m *SettingRule) Create(ctx context.Context, settingRule *schema.SettingRule) error {
	_, err := m.createOne(ctx, schema.TableSettingRule, settingRule)
	return err
}

// Update ...
func (m *SettingRule) Update(ctx context.Context, settingRuleID int64, changed map[string]interface{}) (*schema.SettingRule, error) {
	settingRule := &schema.SettingRule{}
	if _, err := m.updateByID(ctx, schema.TableSettingRule, settingRuleID, goqu.Record(changed)); err != nil {
		return nil, err
	}
	if err := m.findOneByID(ctx, schema.TableSettingRule, settingRuleID, settingRule); err != nil {
		return nil, err
	}
	return settingRule, nil
}

// Delete ...
func (m *SettingRule) Delete(ctx context.Context, id int64) (int64, error) {
	return m.deleteByID(ctx, schema.TableSettingRule, id)
}
