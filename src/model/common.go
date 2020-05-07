package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/teambition/urbs-setting/src/logging"
	"github.com/teambition/urbs-setting/src/schema"
	"github.com/teambition/urbs-setting/src/service"
	"github.com/teambition/urbs-setting/src/util"
)

func init() {
	util.DigProvide(NewModels)
}

// Model ...
type Model struct {
	DB *gorm.DB
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
	m := &Model{DB: sql.DB}
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
func (ms *Models) ApplyLabelRulesAndRefreshUserLabels(ctx context.Context, userID int64, now time.Time, force bool) (*schema.User, error) {
	user, labelIDs, ok, err := ms.User.RefreshLabels(ctx, userID, now.Unix(), force)
	if ok {
		hit, err := ms.LabelRule.ApplyRules(ctx, userID, labelIDs)
		if err != nil {
			return nil, err
		}
		if hit > 0 {
			// refresh label again
			user, labelIDs, ok, err = ms.User.RefreshLabels(ctx, userID, now.Unix(), true)
		}
	}

	if elapsed := time.Now().UTC().Sub(now) / time.Millisecond; elapsed > 100 {
		logging.Warningf("ApplyLabelRulesAndRefreshUserLabels: userID %d, consumed %d ms, refreshed %v, start %v\n",
			userID, elapsed, ok, now)
	}
	return user, err
}

// TryApplyLabelRulesAndRefreshUserLabels ...
func (ms *Models) TryApplyLabelRulesAndRefreshUserLabels(ctx context.Context, userID int64, now time.Time, force bool) *schema.User {
	user, err := ms.ApplyLabelRulesAndRefreshUserLabels(ctx, userID, now, force)
	if err != nil {
		logging.Errf("ApplyLabelRulesAndRefreshUserLabels: userID %d, error %v", userID, err)
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
	if err := ms.SettingRule.ApplyRules(ctx, productID, userID); err != nil {
		logging.Errf("%s error: %v", key, err)
	}
}

// ***** 以下为多个 model 可能共用的接口 *****

func (m *Model) lock(ctx context.Context, key string, expire time.Duration) error {
	now := time.Now().UTC()
	lock := &schema.Lock{Name: key, ExpireAt: now.Add(expire)}
	err := m.DB.Create(lock).Error
	if err != nil {
		l := &schema.Lock{}
		if e := m.DB.Where("`name` = ?", key).First(l).Error; e == nil {
			if l.ExpireAt.Before(now) {
				m.unlock(ctx, key) // 释放失效、异常的锁
				err = m.DB.Create(lock).Error
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
	if err := m.DB.Where("`name` = ?", key).Delete(&schema.Lock{}).Error; err != nil {
		logging.Errf("unlock: key %s, error %v", key, err)
	}
}

const refreshLabelStatusSQL = "select sum(t2.`status`) as Status " +
	"from `group_label` t1, `urbs_group` t2 " +
	"where t1.`label_id` = ? and t1.`group_id` = t2.`id` " +
	"group by t1.`label_id`"

// tryRefreshLabelStatus 更新指定 label 的 Status（灰度标签灰度进度，被作用的用户数，非精确）值
// 比如用户因为属于 n 个群组而被重复设置灰度标签
func (m *Model) tryRefreshLabelStatus(ctx context.Context, labelID int64) {
	if err := m.refreshLabelStatus(ctx, labelID); err != nil {
		logging.Errf("tryRefreshLabelStatus: labelID %d, error %v", labelID, err)
	}
}
func (m *Model) refreshLabelStatus(ctx context.Context, labelID int64) error {
	key := fmt.Sprintf("refreshLabelStatus:%d", labelID)
	if err := m.lock(ctx, key, time.Minute); err != nil {
		return err
	}
	defer m.unlock(ctx, key)

	count := int64(0)
	err := m.DB.Model(&schema.UserLabel{}).Where("`label_id` = ?", labelID).Count(&count).Error
	if err != nil {
		return err
	}

	row := m.DB.Raw(refreshLabelStatusSQL, labelID).Row()
	var status int64
	if err = row.Scan(&status); err != nil && err != sql.ErrNoRows {
		return err
	}

	label := &schema.Label{ID: labelID}
	err = m.DB.Model(label).UpdateColumn("status", count+status).Error
	return err
}

const refreshSettingStatus = "select sum(t2.`status`) as Status " +
	"from `group_setting` t1, `urbs_group` t2 " +
	"where t1.`setting_id` = ? and t1.`group_id` = t2.`id` " +
	"group by t1.`setting_id`"

// tryRefreshSettingStatus 更新指定 setting 的 Status（配置项灰度进度，被作用的用户数，非精确）值
// 比如用户因为属于 n 个群组而被重复设置配置项
func (m *Model) tryRefreshSettingStatus(ctx context.Context, settingID int64) {
	if err := m.refreshSettingStatus(ctx, settingID); err != nil {
		logging.Errf("tryRefreshSettingStatus: settingID %d, error %v", settingID, err)
	}
}
func (m *Model) refreshSettingStatus(ctx context.Context, settingID int64) error {
	key := fmt.Sprintf("refreshSettingStatus:%d", settingID)
	if err := m.lock(ctx, key, time.Minute); err != nil {
		return err
	}
	defer m.unlock(ctx, key)

	count := int64(0)
	err := m.DB.Model(&schema.UserSetting{}).Where("`setting_id` = ?", settingID).Count(&count).Error
	if err != nil {
		return err
	}

	row := m.DB.Raw(refreshSettingStatus, settingID).Row()
	var status int64
	if err = row.Scan(&status); err != nil && err != sql.ErrNoRows {
		return err
	}

	setting := &schema.Setting{ID: settingID}
	err = m.DB.Model(setting).UpdateColumn("status", count+status).Error
	return err
}

// tryRefreshGroupStatus 更新指定 group 的 Status（成员数量统计）值
func (m *Model) tryRefreshGroupStatus(ctx context.Context, groupID int64) {
	if err := m.refreshGroupStatus(ctx, groupID); err != nil {
		logging.Errf("tryRefreshGroupStatus: groupID %d, error %v", groupID, err)
	}
}
func (m *Model) refreshGroupStatus(ctx context.Context, groupID int64) error {
	key := fmt.Sprintf("refreshGroupStatus:%d", groupID)
	if err := m.lock(ctx, key, time.Minute); err != nil {
		return err
	}
	defer m.unlock(ctx, key)

	count := int64(0)
	err := m.DB.Model(&schema.UserGroup{}).Where("`group_id` = ?", groupID).Count(&count).Error
	if err != nil {
		return err
	}
	group := &schema.Group{ID: groupID}
	err = m.DB.Model(group).UpdateColumn("status", count).Error
	return err
}

// tryRefreshModuleStatus 更新指定 module 的 Status（功能模块的配置项统计）值
func (m *Model) tryRefreshModuleStatus(ctx context.Context, moduleID int64) {
	if err := m.refreshModuleStatus(ctx, moduleID); err != nil {
		logging.Errf("tryRefreshModuleStatus: moduleID %d, error %v", moduleID, err)
	}
}
func (m *Model) refreshModuleStatus(ctx context.Context, moduleID int64) error {
	key := fmt.Sprintf("refreshModuleStatus:%d", moduleID)
	if err := m.lock(ctx, key, time.Minute); err != nil {
		return err
	}
	defer m.unlock(ctx, key)

	count := int64(0)
	err := m.DB.Model(&schema.Setting{}).Where("`module_id` = ? and `offline_at` is null", moduleID).Count(&count).Error
	if err != nil {
		return err
	}
	module := &schema.Module{ID: moduleID}
	err = m.DB.Model(module).UpdateColumn("status", count).Error
	return err
}

// tryIncreaseStatisticStatus 加减指定 key 的统计值
func (m *Model) tryIncreaseStatisticStatus(ctx context.Context, key schema.StatisticKey, delta int) {
	if err := m.increaseStatisticStatus(ctx, key, delta); err != nil {
		logging.Errf("tryIncreaseStatisticStatus: key %s, delta: %d, error %v", key, delta, err)
	}
}
func (m *Model) increaseStatisticStatus(ctx context.Context, key schema.StatisticKey, delta int) error {
	exp := gorm.Expr("`status` + ?", delta)
	if delta < 0 {
		exp = gorm.Expr("`status` - ?", -delta)
	} else if delta == 0 {
		return nil
	}
	const sql = "insert ignore into `urbs_statistic` (`name`, `status`) values (?, ?) " +
		"on duplicate key update `status` = ?"

	return m.DB.Exec(sql, key, 1, exp).Error
}

func (m *Model) updateStatisticStatus(ctx context.Context, key schema.StatisticKey, status int64) error {
	const sql = "insert ignore into `urbs_statistic` (`name`, `status`) values (?, ?) " +
		"on duplicate key update `status` = ?"

	return m.DB.Exec(sql, key, status, status).Error
}

// updateStatisticStatus 更新指定 key 的 JSON 值
func (m *Model) updateStatisticValue(ctx context.Context, key schema.StatisticKey, value string) error {
	const sql = "insert ignore into `urbs_statistic` (`name`, `value`) values (?, ?) " +
		"on duplicate key update `value` = ?"

	return m.DB.Exec(sql, key, value, value).Error
}

// tryRefreshUsersTotalSize 更新用户总数
func (m *Model) tryRefreshUsersTotalSize(ctx context.Context) {
	if err := m.refreshUsersTotalSize(ctx); err != nil {
		logging.Errf("refreshUsersTotalSize: error %v", err)
	}
}
func (m *Model) refreshUsersTotalSize(ctx context.Context) error {
	key := string(schema.UsersTotalSize)
	if err := m.lock(ctx, key, 5*time.Minute); err != nil {
		return err
	}
	defer m.unlock(ctx, key)

	count := int64(0)
	err := m.DB.Model(&schema.User{}).Count(&count).Error
	if err != nil {
		return err
	}

	return m.updateStatisticStatus(ctx, schema.UsersTotalSize, count)
}

// tryRefreshGroupsTotalSize 更新群组总数
func (m *Model) tryRefreshGroupsTotalSize(ctx context.Context) {
	if err := m.refreshGroupsTotalSize(ctx); err != nil {
		logging.Errf("refreshGroupsTotalSize: error %v", err)
	}
}
func (m *Model) refreshGroupsTotalSize(ctx context.Context) error {
	key := string(schema.GroupsTotalSize)
	if err := m.lock(ctx, key, time.Minute); err != nil {
		return err
	}
	defer m.unlock(ctx, key)

	count := int64(0)
	err := m.DB.Model(&schema.Group{}).Count(&count).Error
	if err != nil {
		return err
	}

	return m.updateStatisticStatus(ctx, schema.GroupsTotalSize, count)
}

// tryRefreshProductsTotalSize 更新产品总数
func (m *Model) tryRefreshProductsTotalSize(ctx context.Context) {
	if err := m.refreshProductsTotalSize(ctx); err != nil {
		logging.Errf("refreshProductsTotalSize: error %v", err)
	}
}
func (m *Model) refreshProductsTotalSize(ctx context.Context) error {
	key := string(schema.ProductsTotalSize)
	if err := m.lock(ctx, key, time.Minute); err != nil {
		return err
	}
	defer m.unlock(ctx, key)

	count := int64(0)
	err := m.DB.Model(&schema.Product{}).Where("`offline_at` is null").Count(&count).Error
	if err != nil {
		return err
	}

	return m.updateStatisticStatus(ctx, schema.ProductsTotalSize, count)
}

// tryRefreshLabelsTotalSize 更新标签总数
func (m *Model) tryRefreshLabelsTotalSize(ctx context.Context) {
	if err := m.refreshLabelsTotalSize(ctx); err != nil {
		logging.Errf("refreshLabelsTotalSize: error %v", err)
	}
}
func (m *Model) refreshLabelsTotalSize(ctx context.Context) error {
	key := string(schema.LabelsTotalSize)
	if err := m.lock(ctx, key, time.Minute); err != nil {
		return err
	}
	defer m.unlock(ctx, key)

	count := int64(0)
	err := m.DB.Model(&schema.Label{}).Where("`offline_at` is null").Count(&count).Error
	if err != nil {
		return err
	}

	return m.updateStatisticStatus(ctx, schema.LabelsTotalSize, count)
}

// tryRefreshModulesTotalSize 更新模块总数
func (m *Model) tryRefreshModulesTotalSize(ctx context.Context) {
	if err := m.refreshModulesTotalSize(ctx); err != nil {
		logging.Errf("refreshModulesTotalSize: error %v", err)
	}
}
func (m *Model) refreshModulesTotalSize(ctx context.Context) error {
	key := string(schema.ModulesTotalSize)
	if err := m.lock(ctx, key, time.Minute); err != nil {
		return err
	}
	defer m.unlock(ctx, key)

	count := int64(0)
	err := m.DB.Model(&schema.Module{}).Where("`offline_at` is null").Count(&count).Error
	if err != nil {
		return err
	}

	return m.updateStatisticStatus(ctx, schema.ModulesTotalSize, count)
}

// tryRefreshSettingsTotalSize 更新配置项总数
func (m *Model) tryRefreshSettingsTotalSize(ctx context.Context) {
	if err := m.refreshSettingsTotalSize(ctx); err != nil {
		logging.Errf("refreshSettingsTotalSize: error %v", err)
	}
}
func (m *Model) refreshSettingsTotalSize(ctx context.Context) error {
	key := string(schema.SettingsTotalSize)
	if err := m.lock(ctx, key, time.Minute); err != nil {
		return err
	}
	defer m.unlock(ctx, key)

	count := int64(0)
	err := m.DB.Model(&schema.Setting{}).Where("`offline_at` is null").Count(&count).Error
	if err != nil {
		return err
	}

	return m.updateStatisticStatus(ctx, schema.SettingsTotalSize, count)
}

func (m *Model) tryIncreaseLabelsStatus(ctx context.Context, labelIDs []int64, delta int) {
	if err := m.increaseLabelsStatus(ctx, labelIDs, delta); err != nil {
		logging.Errf("increaseLabelsStatus: labelIDs [%v], delta: %d, error %v", labelIDs, delta, err)
	}
}
func (m *Model) increaseLabelsStatus(ctx context.Context, labelIDs []int64, delta int) error {
	exp := gorm.Expr("`status` + ?", delta)
	if delta < 0 {
		exp = gorm.Expr("`status` - ?", -delta)
	} else if delta == 0 {
		return nil
	}
	return m.DB.Model(&schema.Label{}).Where("`id` in ( ? )", labelIDs).Update("status", exp).Error
}

func (m *Model) tryIncreaseSettingsStatus(ctx context.Context, settingIDs []int64, delta int) {
	if err := m.increaseSettingsStatus(ctx, settingIDs, delta); err != nil {
		logging.Errf("increaseSettingsStatus: settingIDs [%v], delta: %d, error %v", settingIDs, delta, err)
	}
}
func (m *Model) increaseSettingsStatus(ctx context.Context, settingIDs []int64, delta int) error {
	exp := gorm.Expr("`status` + ?", delta)
	if delta < 0 {
		exp = gorm.Expr("`status` - ?", -delta)
	} else if delta == 0 {
		return nil
	}
	return m.DB.Model(&schema.Setting{}).Where("`id` in ( ? )", settingIDs).Update("status", exp).Error
}

func (m *Model) tryIncreaseModulesStatus(ctx context.Context, moduleIDs []int64, delta int) {
	if err := m.increaseModulesStatus(ctx, moduleIDs, delta); err != nil {
		logging.Errf("increaseModulesStatus: moduleIDs [%v], delta: %d, error %v", moduleIDs, delta, err)
	}
}
func (m *Model) increaseModulesStatus(ctx context.Context, moduleIDs []int64, delta int) error {
	exp := gorm.Expr("`status` + ?", delta)
	if delta < 0 {
		exp = gorm.Expr("`status` - ?", -delta)
	} else if delta == 0 {
		return nil
	}
	return m.DB.Model(&schema.Module{}).Where("`id` in ( ? )", moduleIDs).Update("status", exp).Error
}

func (m *Model) tryDeleteUserAndGroupLabels(ctx context.Context, labelIDs []int64) {
	var err error
	if len(labelIDs) > 0 {
		if err = m.DB.Exec("delete from `user_label` where `label_id` in ( ? )", labelIDs).Error; err == nil {
			err = m.DB.Exec("delete from `group_label` where `label_id` in ( ? )", labelIDs).Error
		}
	}
	if err != nil {
		logging.Errf("deleteUserAndGroupLabels with label_id [%v] error: %v", labelIDs, err)
	}
}

func (m *Model) tryDeleteLabelsRules(ctx context.Context, labelIDs []int64) {
	var err error
	if len(labelIDs) > 0 {
		err = m.DB.Exec("delete from `label_rule` where `label_id` in ( ? )", labelIDs).Error
	}
	if err != nil {
		logging.Errf("deleteLabelsRules with label_id [%v] error: %v", labelIDs, err)
	}
}

func (m *Model) tryDeleteUserAndGroupSettings(ctx context.Context, settingIDs []int64) {
	var err error
	if len(settingIDs) > 0 {
		if err = m.DB.Exec("delete from `user_setting` where `setting_id` in ( ? )", settingIDs).Error; err == nil {
			err = m.DB.Exec("delete from `group_setting` where `setting_id` in ( ? )", settingIDs).Error
		}
	}
	if err != nil {
		logging.Errf("deleteUserAndGroupSettings with setting_id [%v] error: %v", settingIDs, err)
	}
}

func (m *Model) tryDeleteSettingsRules(ctx context.Context, settingIDs []int64) {
	var err error
	if len(settingIDs) > 0 {
		err = m.DB.Exec("delete from `setting_rule` where `setting_id` in ( ? )", settingIDs).Error
	}
	if err != nil {
		logging.Errf("deleteSettingsRules with setting_id [%v] error: %v", settingIDs, err)
	}
}
