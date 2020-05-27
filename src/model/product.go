package model

import (
	"context"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/teambition/gear"
	"github.com/teambition/urbs-setting/src/schema"
	"github.com/teambition/urbs-setting/src/tpl"
	"github.com/teambition/urbs-setting/src/util"
)

// Product ...
type Product struct {
	*Model
}

// FindByName 根据 name 返回 product 数据
func (m *Product) FindByName(ctx context.Context, name, selectStr string) (*schema.Product, error) {
	var err error
	product := &schema.Product{}
	ok, err := m.findOneByCols(ctx, schema.TableProduct, goqu.Ex{"name": name}, selectStr, product)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}
	return product, nil
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
	product, err := m.FindByName(ctx, productName, "id, offline_at, deleted_at")
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
	sdc := m.DB.Select().
		From(goqu.T(schema.TableProduct)).
		Where(
			goqu.C("deleted_at").IsNull(),
			goqu.C("offline_at").IsNull())

	sd := m.DB.Select().
		From(goqu.T(schema.TableProduct)).
		Where(
			goqu.C("id").Lte(cursor),
			goqu.C("deleted_at").IsNull(),
			goqu.C("offline_at").IsNull())

	if pg.Q != "" {
		sdc = sdc.Where(goqu.C("name").Like(pg.Q))
		sd = sd.Where(goqu.C("name").Like(pg.Q))
	}

	sd = sd.Order(goqu.C("id").Desc()).Limit(uint(pg.PageSize + 1))

	total, err := sdc.CountContext(ctx)
	if err != nil {
		return nil, 0, err
	}

	if err = sd.Executor().ScanStructsContext(ctx, &products); err != nil {
		return nil, 0, err
	}
	return products, int(total), nil
}

// Create ...
func (m *Product) Create(ctx context.Context, product *schema.Product) error {
	rowsAffected, err := m.createOne(ctx, schema.TableProduct, product)
	if rowsAffected > 0 {
		util.Go(5*time.Second, func(gctx context.Context) {
			m.tryIncreaseStatisticStatus(gctx, schema.ProductsTotalSize, 1)
		})
	}
	return err
}

// Update 更新指定功能模块
func (m *Product) Update(ctx context.Context, productID int64, changed map[string]interface{}) (*schema.Product, error) {
	product := &schema.Product{}
	if _, err := m.updateByID(ctx, schema.TableProduct, productID, goqu.Record(changed)); err != nil {
		return nil, err
	}
	if err := m.findOneByID(ctx, schema.TableProduct, productID, product); err != nil {
		return nil, err
	}
	return product, nil
}

// Offline 下线产品
func (m *Product) Offline(ctx context.Context, productID int64) error {
	now := time.Now().UTC()
	rowsAffected, err := m.updateByCols(ctx, schema.TableProduct,
		goqu.Ex{"id": productID, "offline_at": nil},
		goqu.Record{"offline_at": &now, "status": -1},
	)
	if rowsAffected > 0 {
		util.Go(5*time.Second, func(gctx context.Context) {
			m.tryIncreaseStatisticStatus(gctx, schema.ProductsTotalSize, -1)
		})

		err = m.offlineLabels(ctx, goqu.Ex{"product_id": productID, "offline_at": nil})
		if err == nil {
			err = m.offlineModules(ctx, goqu.Ex{"product_id": productID, "offline_at": nil})
		}

		if err == nil {
			util.Go(20*time.Second, func(gctx context.Context) {
				m.tryRefreshModulesTotalSize(gctx)
				m.tryRefreshSettingsTotalSize(gctx)
			})
		}
	}
	return err
}

// Delete 对产品进行逻辑删除
func (m *Product) Delete(ctx context.Context, productID int64) error {
	now := time.Now().UTC()
	_, err := m.updateByID(ctx, schema.TableProduct, productID, goqu.Record{"deleted_at": &now})
	return err
}

// Statistics 返回产品的统计数据
func (m *Product) Statistics(ctx context.Context, productID int64) (*tpl.ProductStatistics, error) {
	res := &tpl.ProductStatistics{}
	sd := m.DB.Select(
		goqu.COUNT("id").As("labels"),
		goqu.L("IFNULL(SUM(`status`), 0)").As("status"),
		goqu.L("IFNULL(SUM(`rls`), 0)").As("release")).
		From(goqu.T(schema.TableLabel)).
		Where(
			goqu.C("product_id").Eq(productID),
			goqu.C("offline_at").IsNull())

	if _, err := sd.Executor().ScanStructContext(ctx, res); err != nil {
		return nil, err
	}

	moduleIDs := make([]int64, 0)
	sd = m.DB.Select("id").
		From(goqu.T(schema.TableModule)).
		Where(
			goqu.C("product_id").Eq(productID),
			goqu.C("offline_at").IsNull())
	if err := sd.Executor().ScanValsContext(ctx, &moduleIDs); err != nil {
		return nil, err
	}

	if len(moduleIDs) > 0 {
		res.Modules = int64(len(moduleIDs))
		sd = m.DB.Select(
			goqu.COUNT("id").As("settings"),
			goqu.L("IFNULL(SUM(`status`), 0)").As("status"),
			goqu.L("IFNULL(SUM(`rls`), 0)").As("release")).
			From(goqu.T(schema.TableSetting)).
			Where(
				goqu.C("module_id").In(tpl.Int64SliceToInterface(moduleIDs)...),
				goqu.C("offline_at").IsNull())

		res2 := &tpl.ProductStatistics{}
		if _, err := sd.Executor().ScanStructContext(ctx, res2); err != nil {
			return nil, err
		}

		res.Settings = res2.Settings
		res.Status += res2.Status
		res.Release += res2.Release
	}
	return res, nil
}
