package model

import (
	"context"
	"database/sql"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/teambition/gear"
	"github.com/teambition/urbs-setting/src/schema"
	"github.com/teambition/urbs-setting/src/service"
	"github.com/teambition/urbs-setting/src/tpl"
)

// Setting ...
type Setting struct {
	*Model
}

// FindByName 根据 moduleID 和 name 返回 setting 数据
func (m *Setting) FindByName(ctx context.Context, moduleID int64, name, selectStr string) (*schema.Setting, error) {
	var err error
	setting := &schema.Setting{}
	db := m.DB.Where("`module_id` = ? and `name` = ?", moduleID, name)

	if selectStr == "" {
		err = db.First(setting).Error
	} else {
		err = db.Select(selectStr).First(setting).Error
	}

	if err == nil {
		return setting, nil
	}

	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}
	return nil, err
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
	setting, err := m.FindByName(ctx, moduleID, settingName, "`id`, `offline_at`")
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
	setting := &schema.Setting{ID: settingID}
	if err := m.DB.First(setting).Error; err != nil {
		return nil, err
	}
	if setting.OfflineAt != nil {
		return nil, gear.ErrNotFound.WithMsgf("setting %d was offline", settingID)
	}
	return setting, nil
}

const listSettingsByProductIDSQL = "select t1.`id`, t1.`created_at`, t1.`updated_at`, t1.`name`, t1.`description`, t1.`channels`, t1.`clients`, t1.`vals`, t1.`status`, t1.`rls`, t2.`name` as module " +
	"from `urbs_setting` t1, `urbs_module` t2 " +
	"where t2.`product_id` = ? and t2.`offline_at` is null and t2.`id` = t1.`module_id` and t1.`id` <= ? and t1.`offline_at` is null " +
	"order by t1.`id` desc " +
	"limit ?"

const countSettingsByProductIDSQL = "select count(t1.`id`) " +
	"from `urbs_setting` t1, `urbs_module` t2 " +
	"where t2.`product_id` = ? and t2.`offline_at` is null and t2.`id` = t1.`module_id` and t1.`offline_at` is null"

const searchSettingsByProductIDSQL = "select t1.`id`, t1.`created_at`, t1.`updated_at`, t1.`name`, t1.`description`, t1.`channels`, t1.`clients`, t1.`vals`, t1.`status`, t1.`rls`, t2.`name` as module " +
	"from `urbs_setting` t1, `urbs_module` t2 " +
	"where t2.`product_id` = ? and t2.`offline_at` is null and t2.`id` = t1.`module_id` and t1.`id` <= ? and t1.`offline_at` is null and t1.`name` like ? " +
	"order by t1.`id` desc " +
	"limit ?"

const countSearchSettingsByProductIDSQL = "select count(t1.`id`) " +
	"from `urbs_setting` t1, `urbs_module` t2 " +
	"where t2.`product_id` = ? and t2.`offline_at` is null and t2.`id` = t1.`module_id` and t1.`offline_at` is null and t1.`name` like ? "

// FindByProductID 根据条件查找 settings
func (m *Setting) FindByProductID(ctx context.Context, product string, productID int64, pg tpl.Pagination) ([]tpl.SettingInfo, int, error) {
	data := make([]tpl.SettingInfo, 0)
	cursor := pg.TokenToID()
	total := 0

	if pg.Q == "" {
		if err := m.DB.Raw(countSettingsByProductIDSQL, productID).Row().Scan(&total); err != nil && err != sql.ErrNoRows {
			return nil, 0, err
		}
	} else {
		if err := m.DB.Raw(countSearchSettingsByProductIDSQL, productID, pg.Q).Row().Scan(&total); err != nil && err != sql.ErrNoRows {
			return nil, 0, err
		}
	}

	var err error
	var rows *sql.Rows
	if pg.Q == "" {
		rows, err = m.DB.Raw(listSettingsByProductIDSQL, productID, cursor, pg.PageSize+1).Rows()
	} else {
		rows, err = m.DB.Raw(searchSettingsByProductIDSQL, productID, cursor, pg.Q, pg.PageSize+1).Rows()
	}
	defer rows.Close()

	if err != nil {
		return nil, 0, err
	}

	for rows.Next() {
		var module string
		var s schema.Setting
		if err := rows.Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt, &s.Name, &s.Desc, &s.Channels, &s.Clients, &s.Values, &s.Status, &s.Release, &module); err != nil {
			return nil, 0, err
		}

		data = append(data, tpl.SettingInfoFrom(s, product, module))
	}

	return data, total, err
}

// Find 根据条件查找 settings
func (m *Setting) Find(ctx context.Context, moduleID int64, pg tpl.Pagination) ([]schema.Setting, int, error) {
	settings := make([]schema.Setting, 0)
	cursor := pg.TokenToID()
	db := m.DB.Where("`id` <= ? and `module_id` = ? and `offline_at` is null", cursor, moduleID)
	dbc := m.DB.Where("`module_id` = ? and `offline_at` is null", moduleID)
	if pg.Q != "" {
		db = m.DB.Where("`id` <= ? and `module_id` = ? and `offline_at` is null and `name` like ?", cursor, moduleID, pg.Q)
		dbc = m.DB.Where("`module_id` = ? and `offline_at` is null and `name` like ?", moduleID, pg.Q)
	}

	total := 0
	err := dbc.Model(&schema.Setting{}).Count(&total).Error
	if err == nil {
		err = db.Order("`id` desc").Limit(pg.PageSize + 1).Find(&settings).Error
	}
	if err != nil {
		return nil, 0, err
	}
	return settings, total, nil
}

// Create ...
func (m *Setting) Create(ctx context.Context, setting *schema.Setting) error {
	err := m.DB.Create(setting).Error
	if err == nil {
		go m.tryIncreaseModulesStatus(ctx, []int64{setting.ModuleID}, 1)
		go m.tryIncreaseStatisticStatus(ctx, schema.SettingsTotalSize, 1)
	}
	return err
}

// Update 更新指定功能模块配置项
func (m *Setting) Update(ctx context.Context, settingID int64, changed map[string]interface{}) (*schema.Setting, error) {
	setting := &schema.Setting{ID: settingID}
	if len(changed) > 0 {
		if err := m.DB.Model(setting).UpdateColumns(changed).Error; err != nil {
			return nil, err
		}
	}

	if err := m.DB.First(setting).Error; err != nil {
		return nil, err
	}
	return setting, nil
}

// Offline 标记配置项下线，同时真删除用户和群组的配置项值
func (m *Setting) Offline(ctx context.Context, moduleID, settingID int64) error {
	now := time.Now().UTC()
	res := m.DB.Model(&schema.Setting{ID: settingID}).UpdateColumns(schema.Setting{
		OfflineAt: &now,
		Status:    -1,
	})
	if res.RowsAffected > 0 {
		go m.tryDeleteSettingsRules(ctx, []int64{settingID})
		go m.tryDeleteUserAndGroupSettings(ctx, []int64{settingID})
		go m.tryIncreaseModulesStatus(ctx, []int64{moduleID}, -1)
		go m.tryIncreaseStatisticStatus(ctx, schema.SettingsTotalSize, -1)
	}
	return res.Error
}

const batchAddUserSettingSQL = "insert ignore into `user_setting` (`user_id`, `setting_id`, `value`, `rls`) " +
	"select `urbs_user`.`id`, ?, ?, ? from `urbs_user` where `urbs_user`.`uid` in ( ? ) " +
	"on duplicate key update `last_value` = `user_setting`.`value`, `value` = ?, `rls` = ?"
const batchAddGroupSettingSQL = "insert ignore into `group_setting` (`group_id`, `setting_id`, `value`, `rls`) " +
	"select `urbs_group`.`id`, ?, ?, ? from `urbs_group` where `urbs_group`.`uid` in ( ? ) " +
	"on duplicate key update `last_value` = `group_setting`.`value`, `value` = ?, `rls` = ?"
const checkAddUserSettingSQL = "select t2.`uid` " +
	"from `user_setting` t1, `urbs_user` t2 " +
	"where t1.`setting_id` = ? and t1.`rls` = ? and t1.`user_id` = t2.`id` " +
	"order by t1.`id` desc limit 1000"
const checkAddGroupSettingSQL = "select t2.`uid` " +
	"from `group_setting` t1, `urbs_group` t2 " +
	"where t1.`setting_id` = ? and t1.`rls` = ? and t1.`group_id` = t2.`id` " +
	"order by t1.`id` desc limit 1000"

// Assign 把标签批量分配给用户或群组，如果用户或群组不存在则忽略，如果已经分配，则把原值保存到 last_value 并更新值
func (m *Setting) Assign(ctx context.Context, settingID int64, value string, users, groups []string) (*tpl.SettingReleaseInfo, error) {
	var err error
	rowsAffected := int64(0)
	release, err := m.AcquireRelease(ctx, settingID)
	if err != nil {
		return nil, err
	}

	releaseInfo := &tpl.SettingReleaseInfo{Release: release, Value: value, Users: []string{}, Groups: []string{}}

	if len(users) > 0 {
		res := m.DB.Exec(batchAddUserSettingSQL, settingID, value, release, users, value, release)
		rowsAffected += res.RowsAffected
		err = res.Error
		if err == nil && res.RowsAffected > 0 {
			rows, err := m.DB.Raw(checkAddUserSettingSQL, settingID, release).Rows()

			if err != nil {
				rows.Close()
				return nil, err
			}

			for rows.Next() {
				var uid string
				if err := rows.Scan(&uid); err != nil {
					rows.Close()
					return nil, err
				}
				releaseInfo.Users = append(releaseInfo.Users, uid)
			}
			rows.Close()
		}
	}
	if err == nil && len(groups) > 0 {
		res := m.DB.Exec(batchAddGroupSettingSQL, settingID, value, release, groups, value, release)
		rowsAffected += res.RowsAffected
		err = res.Error
		if err == nil && res.RowsAffected > 0 {
			rows, err := m.DB.Raw(checkAddGroupSettingSQL, settingID, release).Rows()

			if err != nil {
				rows.Close()
				return nil, err
			}

			for rows.Next() {
				var uid string
				if err := rows.Scan(&uid); err != nil {
					rows.Close()
					return nil, err
				}
				releaseInfo.Groups = append(releaseInfo.Groups, uid)
			}
			rows.Close()
		}
	}

	if rowsAffected > 0 {
		go m.tryRefreshSettingStatus(ctx, settingID)
	}
	return releaseInfo, err
}

// Delete 对配置项进行物理删除
func (m *Setting) Delete(ctx context.Context, settingID int64) error {
	res := m.DB.Delete(&schema.Setting{ID: settingID})
	return res.Error
}

// RemoveUserSetting 删除用户的 setting
func (m *Setting) RemoveUserSetting(ctx context.Context, userID, settingID int64) (int64, error) {
	res := m.DB.Where("`user_id` = ? and `setting_id` = ?", userID, settingID).Delete(&schema.UserSetting{})
	if res.RowsAffected > 0 {
		go m.tryIncreaseSettingsStatus(ctx, []int64{settingID}, -1)
	}
	return res.RowsAffected, res.Error
}

const rollbackUserSettingSQL = "update `user_setting` set `value` = `user_setting`.`last_value` where `user_id` = ? and `setting_id` = ?"

// RollbackUserSetting 回滚用户的 setting
func (m *Setting) RollbackUserSetting(ctx context.Context, userID, settingID int64) error {
	err := m.DB.Exec(rollbackUserSettingSQL, userID, settingID).Error
	return err
}

// RemoveGroupSetting 删除群组的 setting
func (m *Setting) RemoveGroupSetting(ctx context.Context, groupID, settingID int64) (int64, error) {
	res := m.DB.Where("`group_id` = ? and `setting_id` = ?", groupID, settingID).Delete(&schema.GroupSetting{})
	if res.RowsAffected > 0 {
		go m.tryRefreshSettingStatus(ctx, settingID)
	}
	return res.RowsAffected, res.Error
}

const rollbackGroupSettingSQL = "update `group_setting` set `value` = `group_setting`.`last_value` where `group_id` = ? and `setting_id` = ?"

// RollbackGroupSetting 回滚群组的 setting
func (m *Setting) RollbackGroupSetting(ctx context.Context, groupID, settingID int64) error {
	err := m.DB.Exec(rollbackGroupSettingSQL, groupID, settingID).Error
	return err
}

// Recall 撤销指定批次的用户或群组的配置项
func (m *Setting) Recall(ctx context.Context, settingID, release int64) error {
	rowsAffected := int64(0)
	res := m.DB.Where("`setting_id` = ? and `rls` = ?", settingID, release).Delete(&schema.GroupSetting{})
	rowsAffected += res.RowsAffected

	if res.Error == nil {
		res = m.DB.Where("`setting_id` = ? and `rls` = ?", settingID, release).Delete(&schema.UserSetting{})
		rowsAffected += res.RowsAffected
	}
	if rowsAffected > 0 {
		go m.tryRefreshSettingStatus(ctx, settingID)
	}
	return res.Error
}

// AcquireRelease ...
func (m *Setting) AcquireRelease(ctx context.Context, settingID int64) (int64, error) {
	setting := &schema.Setting{ID: settingID}
	if err := m.DB.Model(setting).UpdateColumn("rls", gorm.Expr("`rls` + ?", 1)).Error; err != nil {
		return 0, err
	}
	// MySQL 不支持 RETURNING，并发操作分配时 release 可能不准确，不过真实场景下基本不可能并发操作
	if err := m.DB.Select("`id`, `rls`").First(setting).Error; err != nil {
		return 0, err
	}
	return setting.Release, nil
}

const listSettingUsersSQL = "select t1.`id`, t1.`created_at`, t1.`rls`, t2.`uid`, t1.`value`, t1.`last_value` " +
	"from `user_setting` t1, `urbs_user` t2 " +
	"where t1.`setting_id` = ? and t1.`id` <= ? and t1.`user_id` = t2.`id` " +
	"order by t1.`id` desc " +
	"limit ?"

const countSettingUsersSQL = "select count(t2.`id`) " +
	"from `user_setting` t1, `urbs_user` t2  " +
	"where t1.`setting_id` = ? and t1.`user_id` = t2.`id`"

const searchSettingUsersSQL = "select t1.`id`, t1.`created_at`, t1.`rls`, t2.`uid`, t1.`value`, t1.`last_value` " +
	"from `user_setting` t1, `urbs_user` t2 " +
	"where t1.`setting_id` = ? and t1.`id` <= ? and t1.`user_id` = t2.`id` and t2.`uid` like ? " +
	"order by t1.`id` desc " +
	"limit ?"

const countSearchSettingUsersSQL = "select count(t2.`id`) " +
	"from `user_setting` t1, `urbs_user` t2 " +
	"where t1.`setting_id` = ? and t1.`user_id` = t2.`id` and t2.`uid` like ?"

// ListUsers ...
func (m *Setting) ListUsers(ctx context.Context, settingID int64, pg tpl.Pagination) ([]tpl.SettingUserInfo, int, error) {
	data := []tpl.SettingUserInfo{}
	cursor := pg.TokenToID()
	total := 0

	if pg.Q == "" {
		if err := m.DB.Raw(countSettingUsersSQL, settingID).Row().Scan(&total); err != nil && err != sql.ErrNoRows {
			return nil, 0, err
		}
	} else {
		if err := m.DB.Raw(countSearchSettingUsersSQL, settingID, pg.Q).Row().Scan(&total); err != nil && err != sql.ErrNoRows {
			return nil, 0, err
		}
	}

	var err error
	var rows *sql.Rows
	if pg.Q == "" {
		rows, err = m.DB.Raw(listSettingUsersSQL, settingID, cursor, pg.PageSize+1).Rows()
	} else {
		rows, err = m.DB.Raw(searchSettingUsersSQL, settingID, cursor, pg.Q, pg.PageSize+1).Rows()
	}
	defer rows.Close()

	if err != nil {
		return nil, 0, err
	}

	for rows.Next() {
		info := tpl.SettingUserInfo{}
		if err := rows.Scan(&info.ID, &info.AssignedAt, &info.Release, &info.User, &info.Value, &info.LastValue); err != nil {
			return nil, 0, err
		}
		info.SettingHID = service.IDToHID(settingID, "setting")
		data = append(data, info)
	}
	return data, total, err
}

const listSettingGroupsSQL = "select t1.`id`, t1.`created_at`, t1.`rls`, t2.`uid`, t2.`kind`, t2.`description`, t2.`status`, t1.`value`, t1.`last_value` " +
	"from `group_setting` t1, `urbs_group` t2 " +
	"where t1.`setting_id` = ? and t1.`id` <= ? and t1.`group_id` = t2.`id` " +
	"order by t1.`id` desc " +
	"limit ?"

const countSettingGroupsSQL = "select count(t2.`id`) " +
	"from `group_setting` t1, `urbs_group` t2  " +
	"where t1.`setting_id` = ? and t1.`group_id` = t2.`id`"

const searchSettingGroupsSQL = "select t1.`id`, t1.`created_at`, t1.`rls`, t2.`uid`, t2.`kind`, t2.`description`, t2.`status`, t1.`value`, t1.`last_value` " +
	"from `group_setting` t1, `urbs_group` t2 " +
	"where t1.`setting_id` = ? and t1.`id` <= ? and t1.`group_id` = t2.`id` and t2.`uid` like ? " +
	"order by t1.`id` desc " +
	"limit ?"

const countSearchSettingGroupsSQL = "select count(t2.`id`) " +
	"from `group_setting` t1, `urbs_group` t2 " +
	"where t1.`setting_id` = ? and t1.`group_id` = t2.`id` and t2.`uid` like ?"

// ListGroups ...
func (m *Setting) ListGroups(ctx context.Context, settingID int64, pg tpl.Pagination) ([]tpl.SettingGroupInfo, int, error) {
	data := []tpl.SettingGroupInfo{}
	cursor := pg.TokenToID()
	total := 0

	if pg.Q == "" {
		if err := m.DB.Raw(countSettingGroupsSQL, settingID).Row().Scan(&total); err != nil && err != sql.ErrNoRows {
			return nil, 0, err
		}
	} else {
		if err := m.DB.Raw(countSearchSettingGroupsSQL, settingID, pg.Q).Row().Scan(&total); err != nil && err != sql.ErrNoRows {
			return nil, 0, err
		}
	}

	var err error
	var rows *sql.Rows
	if pg.Q == "" {
		rows, err = m.DB.Raw(listSettingGroupsSQL, settingID, cursor, pg.PageSize+1).Rows()
	} else {
		rows, err = m.DB.Raw(searchSettingGroupsSQL, settingID, cursor, pg.Q, pg.PageSize+1).Rows()
	}
	defer rows.Close()

	if err != nil {
		return nil, 0, err
	}

	for rows.Next() {
		info := tpl.SettingGroupInfo{}
		if err := rows.Scan(&info.ID, &info.AssignedAt, &info.Release, &info.Group, &info.Kind, &info.Desc, &info.Status, &info.Value, &info.LastValue); err != nil {
			return nil, 0, err
		}
		info.SettingHID = service.IDToHID(settingID, "setting")
		data = append(data, info)
	}
	return data, total, err
}
