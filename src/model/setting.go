package model

import (
	"context"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/teambition/urbs-setting/src/schema"
)

// Setting ...
type Setting struct {
	DB *gorm.DB
}

// FindByName 根据 moduleID 和 name 返回 setting 数据
func (m *Setting) FindByName(ctx context.Context, moduleID int64, name, selectStr string) (*schema.Setting, error) {
	var err error
	setting := &schema.Setting{ModuleID: moduleID, Name: name}
	if selectStr == "" {
		err = m.DB.Take(setting).Error
	} else {
		err = m.DB.Select(selectStr).Take(setting).Error
	}

	if err == nil {
		return setting, nil
	}

	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}
	return nil, err
}

// Find 根据条件查找 settings
func (m *Setting) Find(ctx context.Context, moduleID int64) ([]schema.Setting, error) {
	settings := make([]schema.Setting, 0)
	err := m.DB.Where("`module_id` is ?", moduleID).Order("`status`, `created_at`").Limit(1000).Find(settings).Error
	return settings, err
}

// Create ...
func (m *Setting) Create(ctx context.Context, setting *schema.Setting) error {
	db := m.DB.Create(setting)
	if db.Error != nil {
		return db.Error
	}

	return db.Take(setting).Error
}

// Offline 标记配置项下线，同时真删除用户和群组的配置项值
func (m *Setting) Offline(ctx context.Context, settingID int64) error {
	now := time.Now().UTC()
	db := m.DB.Model(&schema.Setting{ID: settingID}).Update(schema.Setting{
		OfflineAt: &now,
		Status:    -1,
	})
	if db.Error == nil {
		go deleteUserAndGroupSettings(db.DB(), []int64{settingID})
	}
	return db.Error
}

const batchAddUserSettingSQL = "insert ignore into `user_setting` (`user_id`, `setting_id`, `value`) " +
	"select `user`.id, ?, ? from `user` where `user`.uid in ( ? ) " +
	"on duplicate key update `last_value` = values(`value`), `value` = ?, `updated_at` = current_timestamp()"
const batchAddGroupSettingSQL = "insert ignore into `group_setting` (`group_id`, `setting_id`, `value`) " +
	"select `group`.id, ?, ? from `group` where `group`.uid in ( ? ) " +
	"on duplicate key update `last_value` = values(`value`), `value` = ?, `updated_at` = current_timestamp()"

// Assign 把标签批量分配给用户或群组，如果用户或群组不存在则忽略，如果已经分配，则把原值保存到 last_value 并更新值
func (m *Setting) Assign(ctx context.Context, settingID int64, value string, users, groups []string) error {
	var err error
	if len(users) > 0 {
		err = m.DB.Exec(batchAddUserSettingSQL, settingID, value, users, value).Error
	}
	if err == nil && len(groups) > 0 {
		err = m.DB.Exec(batchAddUserSettingSQL, settingID, value, groups, value).Error
	}

	return err
}
