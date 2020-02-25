package model

import (
	"context"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/teambition/urbs-setting/src/schema"
)

// Product ...
type Product struct {
	DB *gorm.DB
}

// FindByName 根据 name 返回 product 数据
func (m *Product) FindByName(ctx context.Context, name, selectStr string) (*schema.Product, error) {
	var err error
	product := &schema.Product{Name: name}
	if selectStr == "" {
		err = m.DB.Take(product).Error
	} else {
		err = m.DB.Select(selectStr).Take(product).Error
	}

	if err == nil {
		return product, nil
	}

	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}
	return nil, err
}

// Find 根据条件查找 products
func (m *Product) Find(ctx context.Context) ([]schema.Product, error) {
	products := make([]schema.Product, 0)
	err := m.DB.Where("`deleted_at` is null").Order("`status`, `created_at`").Limit(1000).Find(products).Error
	return products, err
}

// Create ...
func (m *Product) Create(ctx context.Context, product *schema.Product) error {
	db := m.DB.Create(product)
	if db.Error != nil {
		return db.Error
	}

	return db.Take(product).Error
}

// Offline 下线产品
func (m *Product) Offline(ctx context.Context, productID int64) error {
	now := time.Now().UTC()
	db := m.DB.Model(&schema.Product{ID: productID}).Update(schema.Product{
		OfflineAt: &now,
		Status:    -1,
	})
	if db.Error == nil {
		var labelIDs []int64
		if err := db.Model(&schema.Label{ProductID: productID}).Pluck("id", &labelIDs).Error; err == nil {
			db.Model(&schema.Label{}).Where("`id` in (?)", labelIDs).Update(schema.Label{
				OfflineAt: &now,
				Status:    -1,
			})
			go deleteUserAndGroupLabels(db.DB(), labelIDs)
		}

		var moduleIDs []int64
		if err := db.Model(&schema.Module{ProductID: productID}).Pluck("id", &moduleIDs).Error; err == nil {
			db.Model(&schema.Module{}).Where("`id` in (?)", moduleIDs).Update(schema.Setting{
				OfflineAt: &now,
				Status:    -1,
			})

			// 逐个处理，可能数据量太大不适合一次性批量处理
			for _, moduleID := range moduleIDs {
				var settingIDs []int64
				if err := db.Model(&schema.Setting{ModuleID: moduleID}).Pluck("id", &settingIDs).Error; err == nil {
					db.Model(&schema.Setting{}).Where("`id` in (?)", settingIDs).Update(schema.Setting{
						OfflineAt: &now,
						Status:    -1,
					})
					go deleteUserAndGroupSettings(db.DB(), settingIDs)
				}
			}
		}
	}
	return db.Error
}

// Delete 对产品进行逻辑删除
func (m *Product) Delete(ctx context.Context, productID int64) error {
	now := time.Now().UTC()
	return m.DB.Model(&schema.Product{ID: productID}).Update(schema.Product{
		DeletedAt: &now,
	}).Error
}
