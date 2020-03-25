package model

import (
	"context"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/teambition/urbs-setting/src/schema"
	"github.com/teambition/urbs-setting/src/tpl"
)

// Module ...
type Module struct {
	DB *gorm.DB
}

// FindByName 根据 productID 和 name 返回 module 数据
func (m *Module) FindByName(ctx context.Context, productID int64, name, selectStr string) (*schema.Module, error) {
	var err error
	module := &schema.Module{}
	db := m.DB.Where("product_id = ? and name = ?", productID, name)

	if selectStr == "" {
		err = db.First(module).Error
	} else {
		err = db.Select(selectStr).First(module).Error
	}

	if err == nil {
		return module, nil
	}

	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}
	return nil, err
}

// Find 根据条件查找 modules
func (m *Module) Find(ctx context.Context, productID int64, pg tpl.Pagination) ([]schema.Module, error) {
	modules := make([]schema.Module, 0)
	cursor := pg.TokenToID()
	err := m.DB.Where("`product_id` = ? and `id` >= ?", productID, cursor).
		Order("`id`").Limit(pg.PageSize + 1).Find(&modules).Error
	return modules, err
}

// Count 计算 product modules 总数
func (m *Module) Count(ctx context.Context, productID int64) (int, error) {
	count := 0
	err := m.DB.Model(&schema.Module{}).Where("product_id = ?", productID).Count(&count).Error
	return count, err
}

// FindIDsByProductID 根据 productID 查找未下线模块 ID 数组
func (m *Module) FindIDsByProductID(ctx context.Context, productID int64) ([]int64, error) {
	modules := make([]schema.Module, 0)
	err := m.DB.Where("`product_id` = ? and `offline_at` is null", productID).Select("`id`").
		Limit(10000).Find(&modules).Error
	ids := make([]int64, len(modules))
	if err == nil {
		for i, m := range modules {
			ids[i] = m.ID
		}
	}
	return ids, err
}

// Create ...
func (m *Module) Create(ctx context.Context, module *schema.Module) error {
	return m.DB.Create(module).Error
}

// Update 更新指定功能模块
func (m *Module) Update(ctx context.Context, moduleID int64, changed map[string]interface{}) (*schema.Module, error) {
	module := &schema.Module{ID: moduleID}
	if len(changed) > 0 {
		if err := m.DB.Model(module).UpdateColumns(changed).Error; err != nil {
			return nil, err
		}
	}

	if err := m.DB.First(module).Error; err != nil {
		return nil, err
	}
	return module, nil
}

// Offline 标记模块下线
func (m *Module) Offline(ctx context.Context, moduleID int64) error {
	now := time.Now().UTC()
	db := m.DB.Model(&schema.Module{ID: moduleID}).UpdateColumns(schema.Module{
		OfflineAt: &now,
		Status:    -1,
	})
	if db.Error == nil {
		var settingIDs []int64
		if err := db.Model(&schema.Setting{}).Where("module_id = ?", moduleID).Pluck("id", &settingIDs).Error; err == nil {
			db.Model(&schema.Setting{}).Where("`id` in ( ? )", settingIDs).UpdateColumns(schema.Setting{
				OfflineAt: &now,
				Status:    -1,
			})
			go deleteUserAndGroupSettings(db, settingIDs)
		}
	}
	return db.Error
}
