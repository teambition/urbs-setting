package model

import (
	"context"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/teambition/gear"
	"github.com/teambition/urbs-setting/src/schema"
	"github.com/teambition/urbs-setting/src/service"
	"github.com/teambition/urbs-setting/src/tpl"
	"github.com/teambition/urbs-setting/src/util"
)

// Setting ...
type Setting struct {
	*Model
}

// FindByName 根据 moduleID 和 name 返回 setting 数据
func (m *Setting) FindByName(ctx context.Context, moduleID int64, name, selectStr string) (*schema.Setting, error) {
	var err error
	setting := &schema.Setting{}
	ok, err := m.findOneByCols(ctx, schema.TableSetting, goqu.Ex{"module_id": moduleID, "name": name}, selectStr, setting)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	return setting, nil
}

// Acquire ...
func (m *Setting) Acquire(ctx context.Context, moduleID int64, settingName string) (*schema.Setting, error) {
	setting, err := m.FindByName(ctx, moduleID, settingName, "")
	if err != nil {
		return nil, err
	}
	if setting == nil {
		return nil, gear.ErrNotFound.WithMsgf("setting %s not found", settingName)
	}
	if setting.OfflineAt != nil {
		return nil, gear.ErrNotFound.WithMsgf("setting %s was offline", settingName)
	}
	return setting, nil
}

// AcquireID ...
func (m *Setting) AcquireID(ctx context.Context, moduleID int64, settingName string) (int64, error) {
	setting, err := m.FindByName(ctx, moduleID, settingName, "id, offline_at")
	if err != nil {
		return 0, err
	}
	if setting == nil {
		return 0, gear.ErrNotFound.WithMsgf("setting %s not found", settingName)
	}
	if setting.OfflineAt != nil {
		return 0, gear.ErrNotFound.WithMsgf("setting %s was offline", settingName)
	}
	return setting.ID, nil
}

// AcquireByID ...
func (m *Setting) AcquireByID(ctx context.Context, settingID int64) (*schema.Setting, error) {
	setting := &schema.Setting{}
	if err := m.findOneByID(ctx, schema.TableSetting, settingID, setting); err != nil {
		return nil, err
	}
	if setting.OfflineAt != nil {
		return nil, gear.ErrNotFound.WithMsgf("setting %d was offline", settingID)
	}
	return setting, nil
}

// Find 根据条件查找 settings
func (m *Setting) Find(ctx context.Context, productID, moduleID int64, pg tpl.Pagination) ([]schema.Setting, int, error) {
	data := make([]schema.Setting, 0)
	cursor := pg.TokenToID()

	sdc := m.DB.Select().
		From(
			goqu.T(schema.TableSetting).As("t1"),
			goqu.T(schema.TableModule).As("t2")).
		Where(
			goqu.I("t2.product_id").Eq(productID),
			goqu.I("t2.offline_at").IsNull())

	sd := m.DB.Select(
		goqu.I("t1.id"),
		goqu.I("t1.created_at"),
		goqu.I("t1.updated_at"),
		goqu.I("t1.name"),
		goqu.I("t1.description"),
		goqu.I("t1.channels"),
		goqu.I("t1.clients"),
		goqu.I("t1.vals"),
		goqu.I("t1.status"),
		goqu.I("t1.rls"),
		goqu.I("t2.name").As("module")).
		From(
			goqu.T(schema.TableSetting).As("t1"),
			goqu.T(schema.TableModule).As("t2")).
		Where(
			goqu.I("t2.product_id").Eq(productID),
			goqu.I("t2.offline_at").IsNull())

	if moduleID > 0 {
		sdc = sdc.Where(goqu.I("t2.id").Eq(moduleID))
		sd = sd.Where(goqu.I("t2.id").Eq(moduleID))
	}

	sdc = sdc.Where(
		goqu.I("t2.id").Eq(goqu.I("t1.module_id")),
		goqu.I("t1.offline_at").IsNull())
	sd = sd.Where(
		goqu.I("t2.id").Eq(goqu.I("t1.module_id")),
		goqu.I("t1.id").Lte(cursor),
		goqu.I("t1.offline_at").IsNull())

	if pg.Q != "" {
		sdc = sdc.Where(goqu.I("t1.name").ILike(pg.Q))
		sd = sd.Where(goqu.I("t1.name").ILike(pg.Q))
	}

	sd = sd.Order(goqu.I("t1.id").Desc()).Limit(uint(pg.PageSize + 1))

	total, err := sdc.CountContext(ctx)
	if err != nil {
		return nil, 0, err
	}

	scanner, err := sd.Executor().ScannerContext(ctx)
	if err != nil {
		return nil, 0, err
	}
	defer scanner.Close()

	for scanner.Next() {
		setting := schema.Setting{}
		if err := scanner.ScanStruct(&setting); err != nil {
			return nil, 0, err
		}
		data = append(data, setting)
	}

	if err := scanner.Err(); err != nil {
		return nil, 0, err
	}
	return data, int(total), nil
}

// Create ...
func (m *Setting) Create(ctx context.Context, setting *schema.Setting) error {
	rowsAffected, err := m.createOne(ctx, schema.TableSetting, setting)
	if rowsAffected > 0 {
		util.Go(5*time.Second, func(gctx context.Context) {
			m.tryIncreaseModulesStatus(gctx, []int64{setting.ModuleID}, 1)
			m.tryIncreaseStatisticStatus(gctx, schema.SettingsTotalSize, 1)
		})
	}
	return err
}

// Update 更新指定功能模块配置项
func (m *Setting) Update(ctx context.Context, settingID int64, changed map[string]interface{}) (*schema.Setting, error) {
	setting := &schema.Setting{}
	if _, err := m.updateByID(ctx, schema.TableSetting, settingID, goqu.Record(changed)); err != nil {
		return nil, err
	}
	if err := m.findOneByID(ctx, schema.TableSetting, settingID, setting); err != nil {
		return nil, err
	}
	return setting, nil
}

// Offline 标记配置项下线，同时真删除用户和群组的配置项值
func (m *Setting) Offline(ctx context.Context, moduleID, settingID int64) error {
	return m.offlineSettingsInModule(ctx, moduleID, goqu.Ex{"id": settingID, "offline_at": nil})
}

// Assign 把标签批量分配给用户或群组，如果用户或群组不存在则忽略，如果已经分配，则把原值保存到 last_value 并更新值
func (m *Setting) Assign(ctx context.Context, settingID int64, value string, users, groups []string) (*tpl.SettingReleaseInfo, error) {
	totalRowsAffected := int64(0)
	release, err := m.AcquireRelease(ctx, settingID)
	if err != nil {
		return nil, err
	}

	releaseInfo := &tpl.SettingReleaseInfo{Release: release, Value: value, Users: []string{}, Groups: []string{}}
	if len(users) > 0 {
		sd := m.DB.Insert(schema.TableUserSetting).Cols("user_id", "setting_id", "value", "rls").
			FromQuery(goqu.From(goqu.T(schema.TableUser).As("t1")).
				Select(goqu.I("t1.id"), goqu.V(settingID), goqu.V(value), goqu.V(release)).
				Where(goqu.I("t1.uid").In(tpl.StrSliceToInterface(users)...))).
			OnConflict(goqu.DoUpdate("", goqu.Record{
				"last_value": goqu.T(schema.TableUserSetting).Col("value"),
				"value":      value,
				"rls":        release,
			}))

		rowsAffected, err := service.DeResult(sd.Executor().ExecContext(ctx))
		if err != nil {
			return nil, err
		}

		totalRowsAffected += rowsAffected
		if rowsAffected > 0 {
			sd := m.DB.Select(goqu.I("t2.uid")).
				From(
					goqu.T(schema.TableUserSetting).As("t1"),
					goqu.T(schema.TableUser).As("t2")).
				Where(
					goqu.I("t1.setting_id").Eq(goqu.V(settingID)),
					goqu.I("t1.rls").Eq(goqu.V(release)),
					goqu.I("t1.user_id").Eq(goqu.I("t2.id"))).
				Order(goqu.I("t1.id").Desc()).Limit(1000)

			if err := sd.Executor().ScanValsContext(ctx, &releaseInfo.Users); err != nil {
				return nil, err
			}
		}
	}

	if len(groups) > 0 {
		sd := m.DB.Insert(schema.TableGroupSetting).Cols("group_id", "setting_id", "value", "rls").
			FromQuery(goqu.From(goqu.T(schema.TableGroup).As("t1")).
				Select(goqu.I("t1.id"), goqu.V(settingID), goqu.V(value), goqu.V(release)).
				Where(goqu.I("t1.uid").In(tpl.StrSliceToInterface(groups)...))).
			OnConflict(goqu.DoUpdate("", goqu.Record{
				"last_value": goqu.T(schema.TableGroupSetting).Col("value"),
				"value":      value,
				"rls":        release,
			}))

		rowsAffected, err := service.DeResult(sd.Executor().ExecContext(ctx))
		if err != nil {
			return nil, err
		}

		totalRowsAffected += rowsAffected
		if rowsAffected > 0 {
			sd := m.DB.Select(goqu.I("t2.uid")).
				From(
					goqu.T(schema.TableGroupSetting).As("t1"),
					goqu.T(schema.TableGroup).As("t2")).
				Where(
					goqu.I("t1.setting_id").Eq(goqu.V(settingID)),
					goqu.I("t1.rls").Eq(goqu.V(release)),
					goqu.I("t1.group_id").Eq(goqu.I("t2.id"))).
				Order(goqu.I("t1.id").Desc()).Limit(1000)

			if err := sd.Executor().ScanValsContext(ctx, &releaseInfo.Groups); err != nil {
				return nil, err
			}
		}
	}

	if totalRowsAffected > 0 {
		util.Go(10*time.Second, func(gctx context.Context) {
			m.tryRefreshSettingStatus(gctx, settingID)
		})
	}
	return releaseInfo, err
}

// Delete 对配置项进行物理删除
func (m *Setting) Delete(ctx context.Context, id int64) error {
	_, err := m.deleteByID(ctx, schema.TableSetting, id)
	return err
}

// RemoveUserSetting 删除用户的 setting
func (m *Setting) RemoveUserSetting(ctx context.Context, userID, settingID int64) (int64, error) {
	rowsAffected, err := m.deleteByCols(ctx, schema.TableUserSetting,
		goqu.Ex{"user_id": userID, "setting_id": settingID})
	if rowsAffected > 0 {
		util.Go(5*time.Second, func(gctx context.Context) {
			m.tryIncreaseSettingsStatus(gctx, []int64{settingID}, -1)
		})
	}
	return rowsAffected, err
}

// RollbackUserSetting 回滚用户的 setting
func (m *Setting) RollbackUserSetting(ctx context.Context, userID, settingID int64) error {
	_, err := m.updateByCols(ctx, schema.TableUserSetting,
		goqu.Ex{"user_id": userID, "setting_id": settingID},
		goqu.Record{"value": goqu.T(schema.TableUserSetting).Col("last_value")})
	return err
}

// RemoveGroupSetting 删除群组的 setting
func (m *Setting) RemoveGroupSetting(ctx context.Context, groupID, settingID int64) (int64, error) {
	rowsAffected, err := m.deleteByCols(ctx, schema.TableGroupSetting,
		goqu.Ex{"group_id": groupID, "setting_id": settingID})
	if rowsAffected > 0 {
		util.Go(10*time.Second, func(gctx context.Context) {
			m.tryRefreshSettingStatus(gctx, settingID)
		})
	}
	return rowsAffected, err
}

// RollbackGroupSetting 回滚群组的 setting
func (m *Setting) RollbackGroupSetting(ctx context.Context, groupID, settingID int64) error {
	_, err := m.updateByCols(ctx, schema.TableGroupSetting,
		goqu.Ex{"group_id": groupID, "setting_id": settingID},
		goqu.Record{"value": goqu.T(schema.TableGroupSetting).Col("last_value")})
	return err
}

// Recall 撤销指定批次的用户或群组的配置项
func (m *Setting) Recall(ctx context.Context, settingID, release int64) error {
	totalRowsAffected := int64(0)
	rowsAffected, err := m.deleteByCols(ctx, schema.TableGroupSetting, goqu.Ex{"setting_id": settingID, "rls": release})
	if err != nil {
		return err
	}
	totalRowsAffected += rowsAffected

	rowsAffected, err = m.deleteByCols(ctx, schema.TableUserSetting, goqu.Ex{"setting_id": settingID, "rls": release})
	if err != nil {
		return err
	}
	totalRowsAffected += rowsAffected
	if totalRowsAffected > 0 {
		util.Go(10*time.Second, func(gctx context.Context) {
			m.tryRefreshSettingStatus(gctx, settingID)
		})
	}

	return nil
}

// AcquireRelease ...
func (m *Setting) AcquireRelease(ctx context.Context, settingID int64) (int64, error) {
	setting := &schema.Setting{}
	if _, err := m.updateByID(ctx, schema.TableSetting, settingID, goqu.Record{
		"rls": goqu.L("rls + ?", 1),
	}); err != nil {
		return 0, err
	}

	// MySQL 不支持 RETURNING，并发操作分配时 release 可能不准确，不过真实场景下基本不可能并发操作
	if err := m.findOneByID(ctx, schema.TableSetting, settingID, setting); err != nil {
		return 0, err
	}
	return setting.Release, nil
}

// ListUsers ...
func (m *Setting) ListUsers(ctx context.Context, settingID int64, pg tpl.Pagination) ([]tpl.SettingUserInfo, int, error) {
	data := []tpl.SettingUserInfo{}
	cursor := pg.TokenToID()

	sdc := m.DB.Select().
		From(
			goqu.T(schema.TableUserSetting).As("t1"),
			goqu.T(schema.TableUser).As("t2")).
		Where(
			goqu.I("t1.setting_id").Eq(settingID),
			goqu.I("t1.user_id").Eq(goqu.I("t2.id")))

	sd := m.DB.Select(
		goqu.I("t1.id"),
		goqu.I("t1.updated_at").As("assigned_at"),
		goqu.I("t1.rls"),
		goqu.I("t2.uid"),
		goqu.I("t1.value"),
		goqu.I("t1.last_value")).
		From(
			goqu.T(schema.TableUserSetting).As("t1"),
			goqu.T(schema.TableUser).As("t2")).
		Where(
			goqu.I("t1.setting_id").Eq(settingID),
			goqu.I("t1.id").Lte(cursor),
			goqu.I("t1.user_id").Eq(goqu.I("t2.id")))

	if pg.Q != "" {
		sdc = sdc.Where(goqu.I("t2.uid").ILike(pg.Q))
		sd = sd.Where(goqu.I("t2.uid").ILike(pg.Q))
	}

	sd = sd.Order(goqu.I("t1.id").Desc()).Limit(uint(pg.PageSize + 1))

	total, err := sdc.CountContext(ctx)
	if err != nil {
		return nil, 0, err
	}
	scanner, err := sd.Executor().ScannerContext(ctx)
	if err != nil {
		return nil, 0, err
	}
	defer scanner.Close()

	for scanner.Next() {
		info := tpl.SettingUserInfo{}
		if err := scanner.ScanStruct(&info); err != nil {
			return nil, 0, err
		}
		info.SettingHID = service.IDToHID(settingID, "setting")
		data = append(data, info)
	}

	if err := scanner.Err(); err != nil {
		return nil, 0, err
	}

	return data, int(total), err
}

// ListGroups ...
func (m *Setting) ListGroups(ctx context.Context, settingID int64, pg tpl.Pagination) ([]tpl.SettingGroupInfo, int, error) {
	data := []tpl.SettingGroupInfo{}
	cursor := pg.TokenToID()
	sdc := m.DB.Select().
		From(
			goqu.T(schema.TableGroupSetting).As("t1"),
			goqu.T(schema.TableGroup).As("t2")).
		Where(
			goqu.I("t1.setting_id").Eq(settingID),
			goqu.I("t1.group_id").Eq(goqu.I("t2.id")))

	sd := m.DB.Select(
		goqu.I("t1.id"),
		goqu.I("t1.updated_at").As("assigned_at"),
		goqu.I("t1.rls"),
		goqu.I("t2.uid"),
		goqu.I("t2.kind"),
		goqu.I("t2.description"),
		goqu.I("t2.status"),
		goqu.I("t1.value"),
		goqu.I("t1.last_value")).
		From(
			goqu.T(schema.TableGroupSetting).As("t1"),
			goqu.T(schema.TableGroup).As("t2")).
		Where(
			goqu.I("t1.setting_id").Eq(settingID),
			goqu.I("t1.id").Lte(cursor),
			goqu.I("t1.group_id").Eq(goqu.I("t2.id")))

	if pg.Q != "" {
		sdc = sdc.Where(goqu.I("t2.uid").ILike(pg.Q))
		sd = sd.Where(goqu.I("t2.uid").ILike(pg.Q))
	}

	sd = sd.Order(goqu.I("t1.id").Desc()).Limit(uint(pg.PageSize + 1))

	total, err := sdc.CountContext(ctx)
	if err != nil {
		return nil, 0, err
	}
	scanner, err := sd.Executor().ScannerContext(ctx)
	if err != nil {
		return nil, 0, err
	}
	defer scanner.Close()

	for scanner.Next() {
		info := tpl.SettingGroupInfo{}
		if err := scanner.ScanStruct(&info); err != nil {
			return nil, 0, err
		}
		info.SettingHID = service.IDToHID(settingID, "setting")
		data = append(data, info)
	}

	if err := scanner.Err(); err != nil {
		return nil, 0, err
	}

	return data, int(total), err
}
