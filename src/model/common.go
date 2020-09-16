package model

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/teambition/gear"
	"github.com/teambition/urbs-setting/src/logging"
	"github.com/teambition/urbs-setting/src/schema"
	"github.com/teambition/urbs-setting/src/service"
	"github.com/teambition/urbs-setting/src/tpl"
	"github.com/teambition/urbs-setting/src/util"
)

func init() {
	util.DigProvide(NewModels)
}

type dbMode string

// ReadDB ...
const ReadDB dbMode = "ReadDB"

// Model ...
type Model struct {
	SQL  *service.SQL
	DB   *goqu.Database
	RdDB *goqu.Database
}

// Models ...
type Models struct {
	Model       *Model
	Healthz     *Healthz
	User        *User
	Group       *Group
	Product     *Product
	Label       *Label
	Module      *Module
	Setting     *Setting
	LabelRule   *LabelRule
	SettingRule *SettingRule
	Statistic   *Statistic
}

// NewModels ...
func NewModels(sql *service.SQL) *Models {
	m := &Model{SQL: sql, DB: sql.DB, RdDB: sql.RdDB}
	return &Models{
		Model:       m,
		Healthz:     &Healthz{m},
		User:        &User{m},
		Group:       &Group{m},
		Product:     &Product{m},
		Label:       &Label{m},
		Module:      &Module{m},
		Setting:     &Setting{m},
		LabelRule:   &LabelRule{m},
		SettingRule: &SettingRule{m},
		Statistic:   &Statistic{m},
	}
}

// ***** 以下为需要组合多个 model 接口能力而对外暴露的接口 *****

// ApplyLabelRulesAndRefreshUserLabels ...
func (ms *Models) ApplyLabelRulesAndRefreshUserLabels(ctx context.Context, productID int64, product string, userID int64, now time.Time, force bool) (*schema.User, error) {
	user, labelIDs, ok, err := ms.User.RefreshLabels(ctx, userID, now.Unix(), force)
	userProductLables := user.GetLabels(product)
	if ok && len(userProductLables) == 0 {
		hit, err := ms.LabelRule.ApplyRules(ctx, productID, userID, labelIDs, schema.RuleUserPercent)
		if err != nil {
			return nil, err
		}
		if hit > 0 {
			// refresh label again
			user, labelIDs, ok, err = ms.User.RefreshLabels(ctx, userID, now.Unix(), true)
		}
	} else if len(userProductLables) > 0 {
		pg := tpl.Pagination{PageSize: 200}
		pg.Q = userProductLables[0].Label + "%"
		labels, _, err := ms.Label.Find(ctx, productID, pg)
		if err != nil {
			return nil, err
		}
		sort.SliceStable(labels, func(i, j int) bool {
			return len(labels[i].Name) < len(labels[j].Name)
		})
		for _, item := range labels {
			if !strings.HasPrefix(item.Name, userProductLables[0].Label+"-") {
				continue
			}
			hit, err := ms.LabelRule.ApplyRule(ctx, productID, userID, item.ID, schema.RuleChildLabelUserPercent)
			if err != nil {
				return nil, err
			}
			if hit > 0 {
				user, labelIDs, ok, err = ms.User.RefreshLabels(ctx, userID, now.Unix(), true)
			}
			break
		}
	}

	if elapsed := time.Now().UTC().Sub(now) / time.Millisecond; elapsed > 200 {
		logging.Warningf("ApplyLabelRulesAndRefreshUserLabels: userID %d, consumed %d ms, refreshed %v, start %v\n",
			userID, elapsed, ok, now)
	}
	return user, err
}

// TryApplyLabelRulesAndRefreshUserLabels ...
func (ms *Models) TryApplyLabelRulesAndRefreshUserLabels(ctx context.Context, productID int64, product string, userID int64, now time.Time, force bool) *schema.User {
	user, err := ms.ApplyLabelRulesAndRefreshUserLabels(ctx, productID, product, userID, now, force)
	if err != nil {
		logging.Warningf("ApplyLabelRulesAndRefreshUserLabels: userID %d, error %v", userID, err)
		return nil
	}
	return user
}

// TryApplySettingRules ...
func (ms *Models) TryApplySettingRules(ctx context.Context, productID, userID int64) {
	key := fmt.Sprintf("TryApplySettingRules:%d:%d", productID, userID)
	if err := ms.Model.lock(ctx, key, 10*time.Minute); err != nil {
		return
	}

	// 此处不要释放锁，锁期不再执行对应 setting rule
	// defer ms.Model.unlock(ctx, key)
	if err := ms.SettingRule.ApplyRules(ctx, productID, userID, schema.RuleUserPercent); err != nil {
		logging.Warningf("%s error: %v", key, err)
	}
}

// ***** 以下为多个 model 可能共用的接口 *****
func (m *Model) findOneByID(ctx context.Context, table string, id int64, i interface{}) error {
	if id <= 0 || table == "" {
		return fmt.Errorf("invalid id %d or table %s for findOneByID", id, table)
	}

	db := m.DB
	if ctx.Value(ReadDB) != nil {
		db = m.RdDB
	}
	sd := db.From(table).Where(goqu.C("id").Eq(id)).Order(goqu.C("id").Asc()).Limit(1)

	ok, err := sd.Executor().ScanStructContext(ctx, i)
	if err != nil {
		return err
	}
	if !ok {
		return gear.ErrNotFound.WithMsgf("%s %d not found", table, id)
	}
	return nil
}

func (m *Model) findOneByCols(ctx context.Context, table string, cls goqu.Ex, selectStr string, i interface{}) (bool, error) {
	if len(cls) == 0 {
		return false, fmt.Errorf("invalid clause %v for findOneByCols", cls)
	}

	db := m.DB
	if ctx.Value(ReadDB) != nil {
		db = m.RdDB
	}
	sd := db.From(table).Where(cls).Order(goqu.C("id").Asc()).Limit(1)
	if selectStr != "" {
		sd = sd.Select(goqu.L(selectStr))
	}

	return sd.Executor().ScanStructContext(ctx, i)
}

func (m *Model) createOne(ctx context.Context, table string, obj interface{}) (int64, error) {
	if obj == nil {
		return 0, fmt.Errorf("invalid obj for createOne")
	}
	sd := m.DB.Insert(table).Rows(obj)
	res, err := sd.Executor().ExecContext(ctx)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	if id <= 0 || rowsAffected <= 0 {
		return 0, fmt.Errorf("createOne failed")
	}

	err = m.findOneByID(ctx, table, id, obj)
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil
}

func (m *Model) updateByID(ctx context.Context, table string, id int64, changed goqu.Record) (int64, error) {
	if id <= 0 || table == "" {
		return 0, fmt.Errorf("invalid id %d or table %s for updateByID", id, table)
	}
	return m.updateByCols(ctx, table, goqu.Ex{"id": id}, changed)
}

func (m *Model) updateByCols(ctx context.Context, table string, cls goqu.Ex, changed goqu.Record) (int64, error) {
	if len(cls) == 0 {
		return 0, fmt.Errorf("invalid clause %v for updateByCols", cls)
	}
	if len(changed) == 0 {
		return 0, nil
	}
	sd := m.DB.Update(table).Where(cls).Set(changed)
	return service.DeResult(sd.Executor().ExecContext(ctx))
}

func (m *Model) deleteByID(ctx context.Context, table string, id int64) (int64, error) {
	if id <= 0 || table == "" {
		return 0, fmt.Errorf("invalid id %d or table %s for deleteByID", id, table)
	}

	return m.deleteByCols(ctx, table, goqu.Ex{"id": id})
}

func (m *Model) deleteByCols(ctx context.Context, table string, cls goqu.Ex) (int64, error) {
	if len(cls) == 0 {
		return 0, fmt.Errorf("invalid clause %v for deleteByCols", cls)
	}

	sd := m.DB.Delete(table).Where(cls)
	return service.DeResult(sd.Executor().ExecContext(ctx))
}

func (m *Model) offlineLabels(ctx context.Context, cls goqu.Ex) error {
	ids := make([]int64, 0)
	sd := m.DB.Select("id").
		From(goqu.T(schema.TableLabel)).
		Where(cls)
	if err := sd.Executor().ScanValsContext(ctx, &ids); err != nil {
		return err
	}
	if len(ids) == 0 {
		return nil
	}

	now := time.Now().UTC()
	rowsAffected, err := m.updateByCols(ctx, schema.TableLabel, goqu.Ex{"id": ids}, goqu.Record{
		"offline_at": &now,
		"status":     -1,
	})
	if rowsAffected > 0 {
		util.Go(10*time.Second, func(gctx context.Context) {
			m.tryIncreaseStatisticStatus(gctx, schema.LabelsTotalSize, -int(rowsAffected))
			m.tryDeleteLabelsRules(gctx, ids)
			m.tryDeleteUserAndGroupLabels(gctx, ids)
		})
	}
	return err
}

func (m *Model) offlineSettingsInModule(ctx context.Context, moduleID int64, cls goqu.Ex) error {
	cls["module_id"] = moduleID
	ids := make([]int64, 0)
	sd := m.DB.Select("id").
		From(goqu.T(schema.TableSetting)).
		Where(cls)
	if err := sd.Executor().ScanValsContext(ctx, &ids); err != nil {
		return err
	}
	if len(ids) == 0 {
		return nil
	}

	now := time.Now().UTC()
	rowsAffected, err := m.updateByCols(ctx, schema.TableSetting, goqu.Ex{"id": ids}, goqu.Record{
		"offline_at": &now,
		"status":     -1,
	})
	if rowsAffected > 0 {
		util.Go(10*time.Second, func(gctx context.Context) {
			m.tryIncreaseStatisticStatus(gctx, schema.SettingsTotalSize, -int(rowsAffected))
			m.tryDeleteSettingsRules(gctx, ids)
			m.tryDeleteUserAndGroupSettings(gctx, ids)
			m.tryIncreaseModulesStatus(gctx, []int64{moduleID}, -1)
		})
	}
	return err
}

func (m *Model) offlineModules(ctx context.Context, cls goqu.Ex) error {
	ids := make([]int64, 0)
	sd := m.DB.Select("id").
		From(goqu.T(schema.TableModule)).
		Where(cls)
	if err := sd.Executor().ScanValsContext(ctx, &ids); err != nil {
		return err
	}
	if len(ids) == 0 {
		return nil
	}

	now := time.Now().UTC()
	rowsAffected, err := m.updateByCols(ctx, schema.TableModule, goqu.Ex{"id": ids}, goqu.Record{
		"offline_at": &now,
		"status":     -1,
	})
	if rowsAffected > 0 {
		util.Go(5*time.Second, func(gctx context.Context) {
			m.tryIncreaseStatisticStatus(gctx, schema.ModulesTotalSize, -int(rowsAffected))
		})
		for i := range ids {
			if err = m.offlineSettingsInModule(ctx, ids[i], goqu.Ex{"offline_at": nil}); err != nil {
				return err
			}
		}
	}
	return err
}

func (m *Model) lock(ctx context.Context, key string, expire time.Duration) error {
	now := time.Now().UTC()
	lock := &schema.Lock{Name: key, ExpireAt: now.Add(expire)}
	_, err := m.DB.Insert(schema.TableLock).Rows(lock).Executor().ExecContext(ctx)
	if err != nil {
		l := &schema.Lock{}
		sd := m.DB.From(schema.TableLock).Where(goqu.C("name").Eq(key)).Order(goqu.C("id").Asc()).Limit(1)
		ok, _ := sd.Executor().ScanStructContext(ctx, l)
		if ok {
			if l.ExpireAt.Before(now) {
				m.unlock(ctx, key) // 释放失效、异常的锁
				_, err = m.DB.Insert(schema.TableLock).Rows(lock).Executor().ExecContext(ctx)
			} else {
				lock = l
			}
		}
	}
	if err != nil {
		err = fmt.Errorf("%s locked, should expire at: %v, error: %s", key, lock.ExpireAt, err.Error())
	}
	return err
}

func (m *Model) unlock(ctx context.Context, key string) {
	sd := m.DB.Delete(schema.TableLock).Where(goqu.C("name").Eq(key))
	_, err := service.DeResult(sd.Executor().ExecContext(ctx))
	if err != nil {
		logging.Warningf("unlock: key %s, error %v", key, err)
	}
}

// tryRefreshLabelStatus 更新指定 label 的 Status（环境标签灰度进度，被作用的用户数，非精确）值
// 比如用户因为属于 n 个群组而被重复设置环境标签
func (m *Model) tryRefreshLabelStatus(ctx context.Context, labelID int64) {
	if err := m.refreshLabelStatus(ctx, labelID); err != nil {
		logging.Debugf("tryRefreshLabelStatus: labelID %d, error %v", labelID, err)
	}
}
func (m *Model) refreshLabelStatus(ctx context.Context, labelID int64) error {
	key := fmt.Sprintf("refreshLabelStatus:%d", labelID)
	if err := m.lock(ctx, key, time.Minute); err != nil {
		return err
	}
	defer m.unlock(ctx, key)

	sd := m.DB.From(schema.TableUserLabel).Where(goqu.C("label_id").Eq(labelID))
	count, err := sd.CountContext(ctx)
	if err != nil {
		return err
	}

	sd = m.DB.Select(
		goqu.L("IFNULL(SUM(`t2`.`status`), 0)").As("status")).
		From(
			goqu.T(schema.TableGroupLabel).As("t1"),
			goqu.T(schema.TableGroup).As("t2")).
		Where(goqu.I("t1.label_id").Eq(labelID),
			goqu.I("t1.group_id").Eq(goqu.I("t2.id")))

	status := int64(0)
	if _, err := sd.Executor().ScanValContext(ctx, &status); err != nil {
		return err
	}

	_, err = m.updateByID(ctx, schema.TableLabel, labelID, goqu.Record{"status": count + status})
	return err
}

// tryRefreshSettingStatus 更新指定 setting 的 Status（配置项灰度进度，被作用的用户数，非精确）值
// 比如用户因为属于 n 个群组而被重复设置配置项
func (m *Model) tryRefreshSettingStatus(ctx context.Context, settingID int64) {
	if err := m.refreshSettingStatus(ctx, settingID); err != nil {
		logging.Debugf("tryRefreshSettingStatus: settingID %d, error %v", settingID, err)
	}
}
func (m *Model) refreshSettingStatus(ctx context.Context, settingID int64) error {
	key := fmt.Sprintf("refreshSettingStatus:%d", settingID)
	if err := m.lock(ctx, key, time.Minute); err != nil {
		return err
	}
	defer m.unlock(ctx, key)

	sd := m.DB.From(schema.TableUserSetting).Where(goqu.C("setting_id").Eq(settingID))
	count, err := sd.CountContext(ctx)
	if err != nil {
		return err
	}

	sd = m.DB.Select(
		goqu.L("IFNULL(SUM(`t2`.`status`), 0)").As("status")).
		From(
			goqu.T(schema.TableGroupSetting).As("t1"),
			goqu.T(schema.TableGroup).As("t2")).
		Where(goqu.I("t1.setting_id").Eq(settingID),
			goqu.I("t1.group_id").Eq(goqu.I("t2.id")))

	status := int64(0)
	if _, err := sd.Executor().ScanValContext(ctx, &status); err != nil {
		return err
	}

	_, err = m.updateByID(ctx, schema.TableSetting, settingID, goqu.Record{"status": count + status})
	return err
}

// tryRefreshGroupStatus 更新指定 group 的 Status（成员数量统计）值
func (m *Model) tryRefreshGroupStatus(ctx context.Context, groupID int64) {
	if err := m.refreshGroupStatus(ctx, groupID); err != nil {
		logging.Debugf("tryRefreshGroupStatus: groupID %d, error %v", groupID, err)
	}
}
func (m *Model) refreshGroupStatus(ctx context.Context, groupID int64) error {
	key := fmt.Sprintf("refreshGroupStatus:%d", groupID)
	if err := m.lock(ctx, key, time.Minute); err != nil {
		return err
	}
	defer m.unlock(ctx, key)

	sd := m.DB.From(schema.TableUserGroup).Where(goqu.C("group_id").Eq(groupID))
	count, err := sd.CountContext(ctx)
	if err != nil {
		return err
	}

	_, err = m.updateByID(ctx, schema.TableGroup, groupID, goqu.Record{"status": count})
	return err
}

// tryRefreshModuleStatus 更新指定 module 的 Status（功能模块的配置项统计）值
func (m *Model) tryRefreshModuleStatus(ctx context.Context, moduleID int64) {
	if err := m.refreshModuleStatus(ctx, moduleID); err != nil {
		logging.Debugf("tryRefreshModuleStatus: moduleID %d, error %v", moduleID, err)
	}
}
func (m *Model) refreshModuleStatus(ctx context.Context, moduleID int64) error {
	key := fmt.Sprintf("refreshModuleStatus:%d", moduleID)
	if err := m.lock(ctx, key, time.Minute); err != nil {
		return err
	}
	defer m.unlock(ctx, key)

	sd := m.DB.From(schema.TableSetting).Where(goqu.C("module_id").Eq(moduleID), goqu.C("offline_at").IsNull())
	count, err := sd.CountContext(ctx)
	if err != nil {
		return err
	}

	_, err = m.updateByID(ctx, schema.TableModule, moduleID, goqu.Record{"status": count})
	return err
}

// tryIncreaseStatisticStatus 加减指定 key 的统计值
func (m *Model) tryIncreaseStatisticStatus(ctx context.Context, key schema.StatisticKey, delta int) {
	if err := m.increaseStatisticStatus(ctx, key, delta); err != nil {
		logging.Debugf("tryIncreaseStatisticStatus: key %s, delta: %d, error %v", key, delta, err)
	}
}
func (m *Model) increaseStatisticStatus(ctx context.Context, key schema.StatisticKey, delta int) error {
	exp := goqu.L("status + ?", delta)
	if delta < 0 {
		exp = goqu.L("status - ?", -delta)
	} else if delta == 0 {
		return nil
	}

	sd := m.DB.Insert(schema.TableStatistic).
		Rows(goqu.Record{"name": key, "status": 1}).
		OnConflict(goqu.DoUpdate("name", goqu.C("status").Set(exp)))

	_, err := service.DeResult(sd.Executor().ExecContext(ctx))
	return err
}

func (m *Model) updateStatisticStatus(ctx context.Context, key schema.StatisticKey, status int64) error {
	sd := m.DB.Insert(schema.TableStatistic).
		Rows(goqu.Record{"name": key, "status": status}).
		OnConflict(goqu.DoUpdate("name", goqu.C("status").Set(goqu.V(status))))

	_, err := service.DeResult(sd.Executor().ExecContext(ctx))
	return err
}

// updateStatisticStatus 更新指定 key 的 JSON 值
func (m *Model) updateStatisticValue(ctx context.Context, key schema.StatisticKey, value string) error {
	sd := m.DB.Insert(schema.TableStatistic).
		Rows(goqu.Record{"name": key, "value": value}).
		OnConflict(goqu.DoUpdate("name", goqu.C("value").Set(goqu.V(value))))

	_, err := service.DeResult(sd.Executor().ExecContext(ctx))
	return err
}

// tryRefreshUsersTotalSize 更新用户总数
func (m *Model) tryRefreshUsersTotalSize(ctx context.Context) {
	if err := m.refreshUsersTotalSize(ctx); err != nil {
		logging.Debugf("refreshUsersTotalSize: error %v", err)
	}
}
func (m *Model) refreshUsersTotalSize(ctx context.Context) error {
	key := string(schema.UsersTotalSize)
	if err := m.lock(ctx, key, 5*time.Minute); err != nil {
		return err
	}
	defer m.unlock(ctx, key)

	count, err := m.DB.From(schema.TableUser).CountContext(ctx)
	if err != nil {
		return err
	}

	return m.updateStatisticStatus(ctx, schema.UsersTotalSize, count)
}

// tryRefreshGroupsTotalSize 更新群组总数
func (m *Model) tryRefreshGroupsTotalSize(ctx context.Context) {
	if err := m.refreshGroupsTotalSize(ctx); err != nil {
		logging.Debugf("refreshGroupsTotalSize: error %v", err)
	}
}
func (m *Model) refreshGroupsTotalSize(ctx context.Context) error {
	key := string(schema.GroupsTotalSize)
	if err := m.lock(ctx, key, time.Minute); err != nil {
		return err
	}
	defer m.unlock(ctx, key)

	count, err := m.DB.From(schema.TableGroup).CountContext(ctx)
	if err != nil {
		return err
	}

	return m.updateStatisticStatus(ctx, schema.GroupsTotalSize, count)
}

// tryRefreshProductsTotalSize 更新产品总数
func (m *Model) tryRefreshProductsTotalSize(ctx context.Context) {
	if err := m.refreshProductsTotalSize(ctx); err != nil {
		logging.Debugf("refreshProductsTotalSize: error %v", err)
	}
}
func (m *Model) refreshProductsTotalSize(ctx context.Context) error {
	key := string(schema.ProductsTotalSize)
	if err := m.lock(ctx, key, time.Minute); err != nil {
		return err
	}
	defer m.unlock(ctx, key)

	count, err := m.DB.From(schema.TableProduct).Where(goqu.C("offline_at").IsNull()).CountContext(ctx)
	if err != nil {
		return err
	}

	return m.updateStatisticStatus(ctx, schema.ProductsTotalSize, count)
}

// tryRefreshLabelsTotalSize 更新标签总数
func (m *Model) tryRefreshLabelsTotalSize(ctx context.Context) {
	if err := m.refreshLabelsTotalSize(ctx); err != nil {
		logging.Debugf("refreshLabelsTotalSize: error %v", err)
	}
}
func (m *Model) refreshLabelsTotalSize(ctx context.Context) error {
	key := string(schema.LabelsTotalSize)
	if err := m.lock(ctx, key, time.Minute); err != nil {
		return err
	}
	defer m.unlock(ctx, key)

	count, err := m.DB.From(schema.TableLabel).Where(goqu.C("offline_at").IsNull()).CountContext(ctx)
	if err != nil {
		return err
	}

	return m.updateStatisticStatus(ctx, schema.LabelsTotalSize, count)
}

// tryRefreshModulesTotalSize 更新模块总数
func (m *Model) tryRefreshModulesTotalSize(ctx context.Context) {
	if err := m.refreshModulesTotalSize(ctx); err != nil {
		logging.Debugf("refreshModulesTotalSize: error %v", err)
	}
}
func (m *Model) refreshModulesTotalSize(ctx context.Context) error {
	key := string(schema.ModulesTotalSize)
	if err := m.lock(ctx, key, time.Minute); err != nil {
		return err
	}
	defer m.unlock(ctx, key)

	count, err := m.DB.From(schema.TableModule).Where(goqu.C("offline_at").IsNull()).CountContext(ctx)
	if err != nil {
		return err
	}

	return m.updateStatisticStatus(ctx, schema.ModulesTotalSize, count)
}

// tryRefreshSettingsTotalSize 更新配置项总数
func (m *Model) tryRefreshSettingsTotalSize(ctx context.Context) {
	if err := m.refreshSettingsTotalSize(ctx); err != nil {
		logging.Debugf("refreshSettingsTotalSize: error %v", err)
	}
}
func (m *Model) refreshSettingsTotalSize(ctx context.Context) error {
	key := string(schema.SettingsTotalSize)
	if err := m.lock(ctx, key, time.Minute); err != nil {
		return err
	}
	defer m.unlock(ctx, key)

	count, err := m.DB.From(schema.TableSetting).Where(goqu.C("offline_at").IsNull()).CountContext(ctx)
	if err != nil {
		return err
	}

	return m.updateStatisticStatus(ctx, schema.SettingsTotalSize, count)
}

func (m *Model) tryIncreaseLabelsStatus(ctx context.Context, labelIDs []int64, delta int) {
	if err := m.increaseLabelsStatus(ctx, labelIDs, delta); err != nil {
		logging.Debugf("increaseLabelsStatus: labelIDs [%v], delta: %d, error %v", labelIDs, delta, err)
	}
}
func (m *Model) increaseLabelsStatus(ctx context.Context, labelIDs []int64, delta int) error {
	if len(labelIDs) == 0 {
		return nil
	}

	exp := goqu.L("status + ?", delta)
	if delta < 0 {
		exp = goqu.L("status - ?", -delta)
	} else if delta == 0 {
		return nil
	}
	_, err := m.updateByCols(ctx, schema.TableLabel, goqu.Ex{"id": labelIDs}, goqu.Record{"status": exp})
	return err
}

func (m *Model) tryIncreaseSettingsStatus(ctx context.Context, settingIDs []int64, delta int) {
	if err := m.increaseSettingsStatus(ctx, settingIDs, delta); err != nil {
		logging.Debugf("increaseSettingsStatus: settingIDs [%v], delta: %d, error %v", settingIDs, delta, err)
	}
}
func (m *Model) increaseSettingsStatus(ctx context.Context, settingIDs []int64, delta int) error {
	if len(settingIDs) == 0 {
		return nil
	}

	exp := goqu.L("status + ?", delta)
	if delta < 0 {
		exp = goqu.L("status - ?", -delta)
	} else if delta == 0 {
		return nil
	}
	_, err := m.updateByCols(ctx, schema.TableSetting, goqu.Ex{"id": settingIDs}, goqu.Record{"status": exp})
	return err
}

func (m *Model) tryIncreaseModulesStatus(ctx context.Context, moduleIDs []int64, delta int) {
	if err := m.increaseModulesStatus(ctx, moduleIDs, delta); err != nil {
		logging.Debugf("increaseModulesStatus: moduleIDs [%v], delta: %d, error %v", moduleIDs, delta, err)
	}
}
func (m *Model) increaseModulesStatus(ctx context.Context, moduleIDs []int64, delta int) error {
	if len(moduleIDs) == 0 {
		return nil
	}

	exp := goqu.L("status + ?", delta)
	if delta < 0 {
		exp = goqu.L("status - ?", -delta)
	} else if delta == 0 {
		return nil
	}
	_, err := m.updateByCols(ctx, schema.TableModule, goqu.Ex{"id": moduleIDs}, goqu.Record{"status": exp})
	return err
}

func (m *Model) tryDeleteUserAndGroupLabels(ctx context.Context, labelIDs []int64) {
	var err error
	if len(labelIDs) > 0 {
		_, err = m.deleteByCols(ctx, schema.TableUserLabel, goqu.Ex{"label_id": labelIDs})
		if err == nil {
			_, err = m.deleteByCols(ctx, schema.TableGroupLabel, goqu.Ex{"label_id": labelIDs})
		}
	}
	if err != nil {
		logging.Warningf("deleteUserAndGroupLabels with label_id [%v] error: %v", labelIDs, err)
	}
}

func (m *Model) tryDeleteLabelsRules(ctx context.Context, labelIDs []int64) {
	var err error
	if len(labelIDs) > 0 {
		_, err = m.deleteByCols(ctx, schema.TableLabelRule, goqu.Ex{"label_id": labelIDs})
	}
	if err != nil {
		logging.Warningf("deleteLabelsRules with label_id [%v] error: %v", labelIDs, err)
	}
}

func (m *Model) tryDeleteUserAndGroupSettings(ctx context.Context, settingIDs []int64) {
	var err error
	if len(settingIDs) > 0 {
		_, err = m.deleteByCols(ctx, schema.TableUserSetting, goqu.Ex{"setting_id": settingIDs})
		if err == nil {
			_, err = m.deleteByCols(ctx, schema.TableGroupSetting, goqu.Ex{"setting_id": settingIDs})
		}
	}
	if err != nil {
		logging.Warningf("deleteUserAndGroupSettings with setting_id [%v] error: %v", settingIDs, err)
	}
}

func (m *Model) tryDeleteSettingsRules(ctx context.Context, settingIDs []int64) {
	var err error
	if len(settingIDs) > 0 {
		_, err = m.deleteByCols(ctx, schema.TableSettingRule, goqu.Ex{"setting_id": settingIDs})
	}
	if err != nil {
		logging.Warningf("deleteSettingsRules with setting_id [%v] error: %v", settingIDs, err)
	}
}
