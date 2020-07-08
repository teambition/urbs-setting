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

// Module ...
type Module struct {
	*Model
}

// FindByName 根据 productID 和 name 返回 module 数据
func (m *Module) FindByName(ctx context.Context, productID int64, name, selectStr string) (*schema.Module, error) {
	module := &schema.Module{}
	ok, err := m.findOneByCols(ctx, schema.TableModule, goqu.Ex{"product_id": productID, "name": name}, selectStr, module)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	return module, nil
}

// Acquire ...
func (m *Module) Acquire(ctx context.Context, productID int64, moduleName string) (*schema.Module, error) {
	module, err := m.FindByName(ctx, productID, moduleName, "")
	if err != nil {
		return nil, err
	}
	if module == nil {
		return nil, gear.ErrNotFound.WithMsgf("module %s not found", moduleName)
	}
	if module.OfflineAt != nil {
		return nil, gear.ErrNotFound.WithMsgf("module %s was offline", moduleName)
	}
	return module, nil
}

// AcquireID ...
func (m *Module) AcquireID(ctx context.Context, productID int64, moduleName string) (int64, error) {
	module, err := m.FindByName(ctx, productID, moduleName, "id, offline_at")
	if err != nil {
		return 0, err
	}
	if module == nil {
		return 0, gear.ErrNotFound.WithMsgf("module %s not found", moduleName)
	}
	if module.OfflineAt != nil {
		return 0, gear.ErrNotFound.WithMsgf("module %s was offline", moduleName)
	}
	return module.ID, nil
}

// Find 根据条件查找 modules
func (m *Module) Find(ctx context.Context, productID int64, pg tpl.Pagination) ([]schema.Module, int, error) {
	modules := make([]schema.Module, 0)
	cursor := pg.TokenToID()
	sdc := m.RdDB.Select().
		From(goqu.T(schema.TableModule)).
		Where(
			goqu.C("product_id").Eq(productID),
			goqu.C("offline_at").IsNull())

	sd := m.RdDB.Select().
		From(goqu.T(schema.TableModule)).
		Where(
			goqu.C("id").Lte(cursor),
			goqu.C("product_id").Eq(productID),
			goqu.C("offline_at").IsNull())

	if pg.Q != "" {
		sdc = sdc.Where(goqu.C("name").ILike(pg.Q))
		sd = sd.Where(goqu.C("name").ILike(pg.Q))
	}

	sd = sd.Order(goqu.C("id").Desc()).Limit(uint(pg.PageSize + 1))

	total, err := sdc.CountContext(ctx)
	if err != nil {
		return nil, 0, err
	}

	if err = sd.Executor().ScanStructsContext(ctx, &modules); err != nil {
		return nil, 0, err
	}

	return modules, int(total), nil
}

// Create ...
func (m *Module) Create(ctx context.Context, module *schema.Module) error {
	rowsAffected, err := m.createOne(ctx, schema.TableModule, module)
	if rowsAffected > 0 {
		util.Go(5*time.Second, func(gctx context.Context) {
			m.tryIncreaseStatisticStatus(gctx, schema.ModulesTotalSize, 1)
		})
	}
	return err
}

// Update 更新指定功能模块
func (m *Module) Update(ctx context.Context, moduleID int64, changed map[string]interface{}) (*schema.Module, error) {
	module := &schema.Module{}
	if _, err := m.updateByID(ctx, schema.TableModule, moduleID, goqu.Record(changed)); err != nil {
		return nil, err
	}
	if err := m.findOneByID(ctx, schema.TableModule, moduleID, module); err != nil {
		return nil, err
	}
	return module, nil
}

// Offline 标记模块下线
func (m *Module) Offline(ctx context.Context, moduleID int64) error {
	return m.offlineModules(ctx, goqu.Ex{"id": moduleID, "offline_at": nil})
}
