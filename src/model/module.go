package model

import (
	"context"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/teambition/urbs-setting/src/schema"
)

// Module ...
type Module struct {
	DB *gorm.DB
}

// FindByName 根据 productID 和 name 返回 module 数据
func (m *Module) FindByName(ctx context.Context, productID int64, name, selectStr string) (*schema.Module, error) {
	var err error
	module := &schema.Module{ProductID: productID, Name: name}
	if selectStr == "" {
		err = m.DB.Take(module).Error
	} else {
		err = m.DB.Select(selectStr).Take(module).Error
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
func (m *Module) Find(ctx context.Context, productID int64) ([]schema.Module, error) {
	modules := make([]schema.Module, 0)
	err := m.DB.Where("`product_id` is ?", productID).Order("`status`, `created_at`").Limit(1000).Find(modules).Error
	return modules, err
}

// Create ...
func (m *Module) Create(ctx context.Context, module *schema.Module) error {
	db := m.DB.Create(module)
	if db.Error != nil {
		return db.Error
	}

	return db.Take(module).Error
}

// Offline 标记模块下线
func (m *Module) Offline(ctx context.Context, moduleID int64) error {
	now := time.Now().UTC()
	db := m.DB.Model(&schema.Module{ID: moduleID}).Update(schema.Module{
		OfflineAt: &now,
		Status:    -1,
	})
	if db.Error == nil {
		var settingIDs []int64
		if err := db.Model(&schema.Setting{ModuleID: moduleID}).Pluck("id", &settingIDs).Error; err == nil {
			db.Model(&schema.Setting{}).Where("`id` in (?)", settingIDs).Update(schema.Setting{
				OfflineAt: &now,
				Status:    -1,
			})
			go deleteUserAndGroupSettings(db.DB(), settingIDs)
		}
	}
	return db.Error
}
