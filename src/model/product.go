package model

import (
	"context"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/teambition/gear"
	"github.com/teambition/urbs-setting/src/schema"
	"github.com/teambition/urbs-setting/src/tpl"
)

// Product ...
type Product struct {
	*Model
}

// FindByName 根据 name 返回 product 数据
func (m *Product) FindByName(ctx context.Context, name, selectStr string) (*schema.Product, error) {
	var err error
	product := &schema.Product{}
	db := m.DB.Unscoped().Where("`name` = ?", name)

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

// AcquireID ...
func (m *Product) AcquireID(ctx context.Context, productName string) (int64, error) {
	product, err := m.FindByName(ctx, productName, "`id`, `offline_at`, `deleted_at`")
	if err != nil {
		return 0, err
	}
	if product == nil {
		return 0, gear.ErrNotFound.WithMsgf("product %s not found", productName)
	}
	if product.DeletedAt != nil {
		return 0, gear.ErrNotFound.WithMsgf("product %s was deleted", productName)
	}
	if product.OfflineAt != nil {
		return 0, gear.ErrNotFound.WithMsgf("product %s was offline", productName)
	}
	return product.ID, nil
}

// Find 根据条件查找 products
func (m *Product) Find(ctx context.Context, pg tpl.Pagination) ([]schema.Product, error) {
	products := make([]schema.Product, 0)
	cursor := pg.TokenToID()

	err := m.DB.Where("`id` >= ? and `deleted_at` is null", cursor).
		Order("`id`").Limit(pg.PageSize + 1).Find(&products).Error
	return products, err
}

// Count 计算 products 总数
func (m *Product) Count(ctx context.Context) (int, error) {
	count := 0
	err := m.DB.Model(&schema.Product{}).Where("`deleted_at` is null").Count(&count).Error
	return count, err
}

// Create ...
func (m *Product) Create(ctx context.Context, product *schema.Product) error {
	err := m.DB.Create(product).Error
	if err == nil {
		go m.increaseStatisticStatus(ctx, schema.ProductsTotalSize, 1)
	}
	return err
}

// Update 更新指定功能模块
func (m *Product) Update(ctx context.Context, productID int64, changed map[string]interface{}) (*schema.Product, error) {
	product := &schema.Product{ID: productID}
	if len(changed) > 0 {
		if err := m.DB.Model(product).UpdateColumns(changed).Error; err != nil {
			return nil, err
		}
	}

	if err := m.DB.First(product).Error; err != nil {
		return nil, err
	}
	return product, nil
}

// Offline 下线产品
func (m *Product) Offline(ctx context.Context, productID int64) error {
	now := time.Now().UTC()
	res := m.DB.Model(&schema.Product{ID: productID}).UpdateColumns(schema.Product{
		OfflineAt: &now,
		Status:    -1,
	})
	if res.RowsAffected > 0 {
		go m.increaseStatisticStatus(ctx, schema.ProductsTotalSize, -1)
	}

	err := res.Error
	if err == nil {
		var labelIDs []int64
		err = m.DB.Model(&schema.Label{}).Where("`product_id` = ?", productID).Pluck("id", &labelIDs).Error
		if err == nil {
			err = m.DB.Model(&schema.Label{}).Where("`id` in ( ? )", labelIDs).UpdateColumns(schema.Label{
				OfflineAt: &now,
				Status:    -1,
			}).Error
			go m.deleteLabelsRules(ctx, labelIDs)
			go m.deleteUserAndGroupLabels(ctx, labelIDs)
			go m.refreshLabelsTotalSize(ctx)
		}

		var moduleIDs []int64
		err = m.DB.Model(&schema.Module{}).Where("`product_id` = ?", productID).Pluck("id", &moduleIDs).Error
		if err == nil {
			err = m.DB.Model(&schema.Module{}).Where("`id` in ( ? )", moduleIDs).UpdateColumns(schema.Setting{
				OfflineAt: &now,
				Status:    -1,
			}).Error
			go m.refreshModulesTotalSize(ctx)

			// 逐个处理，可能数据量太大不适合一次性批量处理
			for _, moduleID := range moduleIDs {
				var settingIDs []int64
				if err := m.DB.Model(&schema.Setting{}).Where("`module_id` = ?", moduleID).Pluck("id", &settingIDs).Error; err == nil {
					m.DB.Model(&schema.Setting{}).Where("`id` in ( ? )", settingIDs).UpdateColumns(schema.Setting{
						OfflineAt: &now,
						Status:    -1,
					})
					go m.deleteSettingsRules(ctx, settingIDs)
					go m.deleteUserAndGroupSettings(ctx, settingIDs)
				}
			}
			go m.refreshSettingsTotalSize(ctx)
		}
	}
	return err
}

// Delete 对产品进行逻辑删除
func (m *Product) Delete(ctx context.Context, productID int64) error {
	now := time.Now().UTC()
	res := m.DB.Model(&schema.Product{ID: productID}).UpdateColumns(schema.Product{
		DeletedAt: &now,
	})
	return res.Error
}

const productLabelStatisticsSQL = "select `product_id`, count(`id`) as n, sum(`status`) as s, sum(`rls`) as r " +
	"from `urbs_label` " +
	"where `product_id` = ? and `offline_at` is null " +
	"group by `product_id`"
const productSettingStatisticsSQL = "select ? as `product_id`, count(`id`) as n, sum(`status`) as s, sum(`rls`) as r " +
	"from `urbs_setting` " +
	"where `module_id` in ( ? ) and `offline_at` is null " +
	"group by `product_id`"

// Statistics 返回产品的统计数据
func (m *Product) Statistics(ctx context.Context, productID int64) (*tpl.ProductStatistics, error) {
	var n int64
	var s int64
	var r int64
	var ignoreID int64

	if err := m.DB.Raw(productLabelStatisticsSQL, productID).Row().Scan(&ignoreID, &n, &s, &r); err != nil {
		return nil, err
	}

	res := &tpl.ProductStatistics{Labels: n, Status: s, Release: r}

	var moduleIDs []int64
	if err := m.DB.Model(&schema.Module{}).Where("`product_id` = ? and `offline_at` is null", productID).Pluck("id", &moduleIDs).Error; err != nil {
		return nil, err
	}
	res.Modules = int64(len(moduleIDs))
	if err := m.DB.Raw(productSettingStatisticsSQL, productID, moduleIDs).Row().Scan(&ignoreID, &n, &s, &r); err != nil {
		return nil, err
	}
	res.Settings = n
	res.Status += s
	res.Release += r
	return res, nil
}
