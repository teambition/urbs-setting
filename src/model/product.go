package model

import (
	"context"
	"database/sql"
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

// Acquire ...
func (m *Product) Acquire(ctx context.Context, productName string) (*schema.Product, error) {
	product, err := m.FindByName(ctx, productName, "")
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, gear.ErrNotFound.WithMsgf("product %s not found", productName)
	}
	if product.DeletedAt != nil {
		return nil, gear.ErrNotFound.WithMsgf("product %s was deleted", productName)
	}
	if product.OfflineAt != nil {
		return nil, gear.ErrNotFound.WithMsgf("product %s was offline", productName)
	}
	return product, nil
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
func (m *Product) Find(ctx context.Context, pg tpl.Pagination) ([]schema.Product, int, error) {
	products := make([]schema.Product, 0)
	cursor := pg.TokenToID()
	db := m.DB.Where("`id` <= ? and `deleted_at` is null and `offline_at` is null", cursor)
	if pg.Q != "" {
		db = m.DB.Where("`id` <= ? and `deleted_at` is null and `offline_at` is null and `name` like ?", cursor, pg.Q)
	}

	total := 0
	err := db.Model(&schema.Product{}).Count(&total).Error
	if err == nil {
		err = db.Order("`id` desc").Limit(pg.PageSize + 1).Find(&products).Error
	}
	if err != nil {
		return nil, 0, err
	}
	return products, total, nil
}

// Create ...
func (m *Product) Create(ctx context.Context, product *schema.Product) error {
	err := m.DB.Create(product).Error
	if err == nil {
		go m.tryIncreaseStatisticStatus(ctx, schema.ProductsTotalSize, 1)
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
		go m.tryIncreaseStatisticStatus(ctx, schema.ProductsTotalSize, -1)
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
			go m.tryDeleteLabelsRules(ctx, labelIDs)
			go m.tryDeleteUserAndGroupLabels(ctx, labelIDs)
			go m.tryRefreshLabelsTotalSize(ctx)
		}

		var moduleIDs []int64
		err = m.DB.Model(&schema.Module{}).Where("`product_id` = ?", productID).Pluck("id", &moduleIDs).Error
		if err == nil {
			err = m.DB.Model(&schema.Module{}).Where("`id` in ( ? )", moduleIDs).UpdateColumns(schema.Setting{
				OfflineAt: &now,
				Status:    -1,
			}).Error
			go m.tryRefreshModulesTotalSize(ctx)

			// 逐个处理，可能数据量太大不适合一次性批量处理
			for _, moduleID := range moduleIDs {
				var settingIDs []int64
				if err := m.DB.Model(&schema.Setting{}).Where("`module_id` = ?", moduleID).Pluck("id", &settingIDs).Error; err == nil {
					m.DB.Model(&schema.Setting{}).Where("`id` in ( ? )", settingIDs).UpdateColumns(schema.Setting{
						OfflineAt: &now,
						Status:    -1,
					})
					go m.tryDeleteSettingsRules(ctx, settingIDs)
					go m.tryDeleteUserAndGroupSettings(ctx, settingIDs)
				}
			}
			go m.tryRefreshSettingsTotalSize(ctx)
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

	if err := m.DB.Raw(productLabelStatisticsSQL, productID).Row().Scan(&ignoreID, &n, &s, &r); err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	res := &tpl.ProductStatistics{Labels: n, Status: s, Release: r}

	var moduleIDs []int64
	if err := m.DB.Model(&schema.Module{}).Where("`product_id` = ? and `offline_at` is null", productID).Pluck("id", &moduleIDs).Error; err != nil {
		return nil, err
	}
	res.Modules = int64(len(moduleIDs))
	if err := m.DB.Raw(productSettingStatisticsSQL, productID, moduleIDs).Row().Scan(&ignoreID, &n, &s, &r); err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	res.Settings = n
	res.Status += s
	res.Release += r
	return res, nil
}
