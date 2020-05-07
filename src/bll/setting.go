package bll

import (
	"context"

	"github.com/teambition/gear"
	"github.com/teambition/urbs-setting/src/model"
	"github.com/teambition/urbs-setting/src/schema"
	"github.com/teambition/urbs-setting/src/tpl"
)

// Setting ...
type Setting struct {
	ms *model.Models
}

// List 返回产品下的功能模块配置项列表
func (b *Setting) List(ctx context.Context, productName, moduleName string, pg tpl.Pagination) (*tpl.SettingsInfoRes, error) {
	productID, err := b.ms.Product.AcquireID(ctx, productName)
	if err != nil {
		return nil, err
	}

	module, err := b.ms.Module.Acquire(ctx, productID, moduleName)
	if err != nil {
		return nil, err
	}

	settings, total, err := b.ms.Setting.Find(ctx, module.ID, pg)
	if err != nil {
		return nil, err
	}

	res := &tpl.SettingsInfoRes{Result: tpl.SettingsInfoFrom(settings, productName, moduleName)}
	res.TotalSize = total
	if len(res.Result) > pg.PageSize {
		res.NextPageToken = tpl.IDToPageToken(res.Result[pg.PageSize].ID)
		res.Result = res.Result[:pg.PageSize]
	}
	return res, nil
}

// ListByProduct 返回产品下的所有功能模块的配置项列表
func (b *Setting) ListByProduct(ctx context.Context, productName string, pg tpl.Pagination) (*tpl.SettingsInfoRes, error) {
	productID, err := b.ms.Product.AcquireID(ctx, productName)
	if err != nil {
		return nil, err
	}

	settingsInfo, total, err := b.ms.Setting.FindByProductID(ctx, productName, productID, pg)
	if err != nil {
		return nil, err
	}

	res := &tpl.SettingsInfoRes{Result: settingsInfo}
	res.TotalSize = total
	if len(res.Result) > pg.PageSize {
		res.NextPageToken = tpl.IDToPageToken(res.Result[pg.PageSize].ID)
		res.Result = res.Result[:pg.PageSize]
	}
	return res, nil
}

// Get 返回产品下指定功能模块配置项
func (b *Setting) Get(ctx context.Context, productName, moduleName, settingName string) (*tpl.SettingInfoRes, error) {
	productID, err := b.ms.Product.AcquireID(ctx, productName)
	if err != nil {
		return nil, err
	}

	module, err := b.ms.Module.Acquire(ctx, productID, moduleName)
	if err != nil {
		return nil, err
	}

	setting, err := b.ms.Setting.Acquire(ctx, module.ID, settingName)
	if err != nil {
		return nil, err
	}

	res := &tpl.SettingInfoRes{Result: tpl.SettingInfoFrom(*setting, productName, moduleName)}
	return res, nil
}

// Create 创建功能模块配置项
func (b *Setting) Create(ctx context.Context, productName, moduleName, settingName, desc string) (*tpl.SettingInfoRes, error) {
	productID, err := b.ms.Product.AcquireID(ctx, productName)
	if err != nil {
		return nil, err
	}

	module, err := b.ms.Module.Acquire(ctx, productID, moduleName)
	if err != nil {
		return nil, err
	}

	setting := &schema.Setting{ModuleID: module.ID, Name: settingName, Desc: desc}
	if err = b.ms.Setting.Create(ctx, setting); err != nil {
		return nil, err
	}
	return &tpl.SettingInfoRes{Result: tpl.SettingInfoFrom(*setting, productName, moduleName)}, nil
}

// Update ...
func (b *Setting) Update(ctx context.Context, productName, moduleName, settingName string, body tpl.SettingUpdateBody) (*tpl.SettingInfoRes, error) {
	productID, err := b.ms.Product.AcquireID(ctx, productName)
	if err != nil {
		return nil, err
	}

	module, err := b.ms.Module.Acquire(ctx, productID, moduleName)
	if err != nil {
		return nil, err
	}

	setting, err := b.ms.Setting.Acquire(ctx, module.ID, settingName)
	if err != nil {
		return nil, err
	}

	setting, err = b.ms.Setting.Update(ctx, setting.ID, body.ToMap())
	if err != nil {
		return nil, err
	}
	return &tpl.SettingInfoRes{Result: tpl.SettingInfoFrom(*setting, productName, moduleName)}, nil
}

// Offline 下线功能模块配置项
func (b *Setting) Offline(ctx context.Context, productName, moduleName, settingName string) (*tpl.BoolRes, error) {
	productID, err := b.ms.Product.AcquireID(ctx, productName)
	if err != nil {
		return nil, err
	}

	module, err := b.ms.Module.Acquire(ctx, productID, moduleName)
	if err != nil {
		return nil, err
	}

	res := &tpl.BoolRes{Result: false}
	setting, err := b.ms.Setting.FindByName(ctx, module.ID, settingName, "id, `offline_at`")
	if err != nil {
		return nil, err
	}
	if setting == nil {
		return nil, gear.ErrNotFound.WithMsgf("setting %s not found", settingName)
	}
	if setting.OfflineAt == nil {
		if err = b.ms.Setting.Offline(ctx, module.ID, setting.ID); err != nil {
			return nil, err
		}
		res.Result = true
	}
	return res, nil
}

// Assign 把配置项批量分配给用户或群组
func (b *Setting) Assign(ctx context.Context, productName, moduleName, settingName, value string, users, groups []string) (*tpl.SettingReleaseInfo, error) {
	productID, err := b.ms.Product.AcquireID(ctx, productName)
	if err != nil {
		return nil, err
	}

	module, err := b.ms.Module.Acquire(ctx, productID, moduleName)
	if err != nil {
		return nil, err
	}

	setting, err := b.ms.Setting.Acquire(ctx, module.ID, settingName)
	if err != nil {
		return nil, err
	}
	vals := tpl.StringToSlice(setting.Values)
	if value != "" && !tpl.StringSliceHas(vals, value) {
		return nil, gear.ErrBadRequest.WithMsgf("value %s is not in setting", value)
	}
	return b.ms.Setting.Assign(ctx, setting.ID, value, users, groups)
}

// Delete 物理删除配置项
func (b *Setting) Delete(ctx context.Context, productName, moduleName, settingName string) (*tpl.BoolRes, error) {
	productID, err := b.ms.Product.AcquireID(ctx, productName)
	if err != nil {
		return nil, err
	}

	module, err := b.ms.Module.Acquire(ctx, productID, moduleName)
	if err != nil {
		return nil, err
	}

	setting, err := b.ms.Setting.FindByName(ctx, module.ID, settingName, "id, `offline_at`")
	if err != nil {
		return nil, err
	}

	res := &tpl.BoolRes{Result: false}
	if setting != nil {
		if setting.OfflineAt == nil {
			return nil, gear.ErrConflict.WithMsgf("setting %s is not offline", settingName)
		}
		if err = b.ms.Setting.Delete(ctx, setting.ID); err != nil {
			return nil, err
		}
		res.Result = true
	}
	return res, nil
}

// Recall 撤销指定批次的用户或群组的配置项
func (b *Setting) Recall(ctx context.Context, productName, moduleName, settingName string, release int64) (*tpl.BoolRes, error) {
	productID, err := b.ms.Product.AcquireID(ctx, productName)
	if err != nil {
		return nil, err
	}

	module, err := b.ms.Module.Acquire(ctx, productID, moduleName)
	if err != nil {
		return nil, err
	}

	setting, err := b.ms.Setting.Acquire(ctx, module.ID, settingName)
	if err != nil {
		return nil, err
	}

	res := &tpl.BoolRes{Result: false}
	if err = b.ms.Setting.Recall(ctx, setting.ID, release); err != nil {
		return nil, err
	}
	res.Result = true
	return res, nil
}

// CreateRule ...
func (b *Setting) CreateRule(ctx context.Context, productName, moduleName, settingName string, body tpl.SettingRuleBody) (*tpl.SettingRuleInfoRes, error) {
	productID, err := b.ms.Product.AcquireID(ctx, productName)
	if err != nil {
		return nil, err
	}

	module, err := b.ms.Module.Acquire(ctx, productID, moduleName)
	if err != nil {
		return nil, err
	}

	setting, err := b.ms.Setting.Acquire(ctx, module.ID, settingName)
	if err != nil {
		return nil, err
	}
	vals := tpl.StringToSlice(setting.Values)
	if body.Value != "" && !tpl.StringSliceHas(vals, body.Value) {
		return nil, gear.ErrBadRequest.WithMsgf("value %s is not in setting", body.Value)
	}

	settingRule := &schema.SettingRule{
		ProductID: productID,
		SettingID: setting.ID,
		Kind:      body.Kind,
		Rule:      body.ToRule(),
		Value:     body.Value,
		Release:   0,
	}
	if err = b.ms.SettingRule.Create(ctx, settingRule); err != nil {
		return nil, err
	}
	// 创建成功再从 setting 获取当前的 release 发布计数
	release, err := b.ms.Setting.AcquireRelease(ctx, setting.ID)
	if err != nil {
		return nil, err
	}

	settingRule, err = b.ms.SettingRule.Update(ctx, settingRule.ID, map[string]interface{}{"rls": release})
	if err != nil {
		return nil, err
	}
	return &tpl.SettingRuleInfoRes{Result: tpl.SettingRuleInfoFrom(*settingRule)}, nil
}

// ListRules ...
func (b *Setting) ListRules(ctx context.Context, productName, moduleName, settingName string) (*tpl.SettingRulesInfoRes, error) {
	productID, err := b.ms.Product.AcquireID(ctx, productName)
	if err != nil {
		return nil, err
	}

	module, err := b.ms.Module.Acquire(ctx, productID, moduleName)
	if err != nil {
		return nil, err
	}

	setting, err := b.ms.Setting.Acquire(ctx, module.ID, settingName)
	if err != nil {
		return nil, err
	}

	settingRules, err := b.ms.SettingRule.Find(ctx, productID, setting.ID)
	if err != nil {
		return nil, err
	}

	res := &tpl.SettingRulesInfoRes{Result: tpl.SettingRulesInfoFrom(settingRules)}
	res.TotalSize = len(settingRules)
	return res, nil
}

// UpdateRule ...
func (b *Setting) UpdateRule(ctx context.Context, productName, moduleName, settingName string, ruleID int64, body tpl.SettingRuleBody) (*tpl.SettingRuleInfoRes, error) {
	productID, err := b.ms.Product.AcquireID(ctx, productName)
	if err != nil {
		return nil, err
	}

	module, err := b.ms.Module.Acquire(ctx, productID, moduleName)
	if err != nil {
		return nil, err
	}

	setting, err := b.ms.Setting.Acquire(ctx, module.ID, settingName)
	if err != nil {
		return nil, err
	}

	settingRule, err := b.ms.SettingRule.Acquire(ctx, ruleID)
	if err != nil {
		return nil, err
	}

	if settingRule.SettingID != setting.ID || body.Kind != settingRule.Kind {
		return nil, gear.ErrNotFound.WithMsgf("label rule not matched!")
	}

	changed := map[string]interface{}{}
	if body.Value != "" {
		vals := tpl.StringToSlice(setting.Values)
		if !tpl.StringSliceHas(vals, body.Value) {
			return nil, gear.ErrBadRequest.WithMsgf("value %s is not in setting", body.Value)
		}
		if body.Value != settingRule.Value {
			changed["value"] = body.Value
		}
	}

	rule := body.ToRule()
	if rule != settingRule.Rule {
		changed["rule"] = rule
	}

	if len(changed) > 0 {
		release, err := b.ms.Setting.AcquireRelease(ctx, setting.ID)
		if err != nil {
			return nil, err
		}
		changed["rls"] = release
		settingRule, err = b.ms.SettingRule.Update(ctx, settingRule.ID, changed)
		if err != nil {
			return nil, err
		}
	}

	return &tpl.SettingRuleInfoRes{Result: tpl.SettingRuleInfoFrom(*settingRule)}, nil
}

// DeleteRule ...
func (b *Setting) DeleteRule(ctx context.Context, productName, moduleName, settingName string, ruleID int64) (*tpl.BoolRes, error) {
	productID, err := b.ms.Product.AcquireID(ctx, productName)
	if err != nil {
		return nil, err
	}

	module, err := b.ms.Module.Acquire(ctx, productID, moduleName)
	if err != nil {
		return nil, err
	}

	setting, err := b.ms.Setting.Acquire(ctx, module.ID, settingName)
	if err != nil {
		return nil, err
	}

	res := &tpl.BoolRes{Result: false}
	settingRule, err := b.ms.SettingRule.Acquire(ctx, ruleID)
	if err != nil {
		return res, nil
	}

	if settingRule.SettingID != setting.ID {
		return nil, gear.ErrNotFound.WithMsgf("setting rule not matched!")
	}

	rowsAffected, err := b.ms.SettingRule.Delete(ctx, settingRule.ID)
	if err != nil {
		return nil, err
	}
	res.Result = rowsAffected > 0
	return res, nil
}

// ListUsers 返回产品下功能配置项的用户列表
func (b *Setting) ListUsers(ctx context.Context, productName, moduleName, settingName string, pg tpl.Pagination) (*tpl.SettingUsersInfoRes, error) {
	productID, err := b.ms.Product.AcquireID(ctx, productName)
	if err != nil {
		return nil, err
	}

	module, err := b.ms.Module.Acquire(ctx, productID, moduleName)
	if err != nil {
		return nil, err
	}

	setting, err := b.ms.Setting.Acquire(ctx, module.ID, settingName)
	if err != nil {
		return nil, err
	}

	data, total, err := b.ms.Setting.ListUsers(ctx, setting.ID, pg)
	if err != nil {
		return nil, err
	}
	res := &tpl.SettingUsersInfoRes{Result: data}
	res.TotalSize = total
	if len(res.Result) > pg.PageSize {
		res.NextPageToken = tpl.IDToPageToken(res.Result[pg.PageSize].ID)
		res.Result = res.Result[:pg.PageSize]
	}
	return res, nil
}

// ListGroups 返回产品下功能配置项的群组列表
func (b *Setting) ListGroups(ctx context.Context, productName, moduleName, settingName string, pg tpl.Pagination) (*tpl.SettingGroupsInfoRes, error) {
	productID, err := b.ms.Product.AcquireID(ctx, productName)
	if err != nil {
		return nil, err
	}

	module, err := b.ms.Module.Acquire(ctx, productID, moduleName)
	if err != nil {
		return nil, err
	}

	setting, err := b.ms.Setting.Acquire(ctx, module.ID, settingName)
	if err != nil {
		return nil, err
	}

	data, total, err := b.ms.Setting.ListGroups(ctx, setting.ID, pg)
	if err != nil {
		return nil, err
	}
	res := &tpl.SettingGroupsInfoRes{Result: data}
	res.TotalSize = total
	if len(res.Result) > pg.PageSize {
		res.NextPageToken = tpl.IDToPageToken(res.Result[pg.PageSize].ID)
		res.Result = res.Result[:pg.PageSize]
	}
	return res, nil
}
