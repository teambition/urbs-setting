package model

import (
	"context"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/teambition/urbs-setting/src/schema"
	"github.com/teambition/urbs-setting/src/tpl"
)

// Product ...
type Product struct {
	DB *gorm.DB
}

// FindByName 根据 name 返回 product 数据
func (m *Product) FindByName(ctx context.Context, name, selectStr string) (*schema.Product, error) {
	var err error
	product := &schema.Product{}
	db := m.DB.Unscoped().Where("name = ?", name)

	if selectStr == "" {
		err = db.First(product).Error
	} else {
		err = db.Select(selectStr).First(product).Error
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
func (m *Product) Find(ctx context.Context, pg tpl.Pagination) ([]schema.Product, error) {
	products := make([]schema.Product, 0)
	pageToken := pg.TokenToID()
	err := m.DB.Where("`id` >= ? and `deleted_at` is null", pageToken).
		Order("`id`").Limit(pg.PageSize + 1).Find(&products).Error
	return products, err
}

// Count 计算 products 总数
func (m *Product) Count(ctx context.Context) (int, error) {
	count := 0
	err := m.DB.Model(&schema.Product{}).Where("deleted_at is null").Count(&count).Error
	return count, err
}

// Create ...
func (m *Product) Create(ctx context.Context, product *schema.Product) error {
	return m.DB.Create(product).Error
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
		if err := db.Model(&schema.Label{}).Where("product_id = ?", productID).Pluck("id", &labelIDs).Error; err == nil {
			db.Model(&schema.Label{}).Where("`id` in ( ? )", labelIDs).Update(schema.Label{
				OfflineAt: &now,
				Status:    -1,
			})
			go deleteUserAndGroupLabels(db, labelIDs)
		}

		var moduleIDs []int64
		if err := db.Model(&schema.Module{}).Where("product_id = ?", productID).Pluck("id", &moduleIDs).Error; err == nil {
			db.Model(&schema.Module{}).Where("`id` in ( ? )", moduleIDs).Update(schema.Setting{
				OfflineAt: &now,
				Status:    -1,
			})

			// 逐个处理，可能数据量太大不适合一次性批量处理
			for _, moduleID := range moduleIDs {
				var settingIDs []int64
				if err := db.Model(&schema.Setting{}).Where("module_id = ?", moduleID).Pluck("id", &settingIDs).Error; err == nil {
					db.Model(&schema.Setting{}).Where("`id` in ( ? )", settingIDs).Update(schema.Setting{
						OfflineAt: &now,
						Status:    -1,
					})
					go deleteUserAndGroupSettings(db, settingIDs)
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
