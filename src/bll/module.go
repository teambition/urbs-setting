package bll

import (
	"context"

	"github.com/teambition/gear"
	"github.com/teambition/urbs-setting/src/model"
	"github.com/teambition/urbs-setting/src/schema"
	"github.com/teambition/urbs-setting/src/tpl"
)

// Module ...
type Module struct {
	ms *model.Models
}

// List 返回产品下的功能模块列表，TODO：支持分页
func (b *Module) List(ctx context.Context, productName string, pg tpl.Pagination) (*tpl.ModulesRes, error) {
	productID, err := b.ms.Product.AcquireID(ctx, productName)
	if err != nil {
		return nil, err
	}
	modules, total, err := b.ms.Module.Find(ctx, productID, pg)
	if err != nil {
		return nil, err
	}

	res := &tpl.ModulesRes{Result: modules}
	res.TotalSize = total
	if len(res.Result) > pg.PageSize {
		res.NextPageToken = tpl.IDToPageToken(res.Result[pg.PageSize].ID)
		res.Result = res.Result[:pg.PageSize]
	}
	return res, nil
}

// Create 创建功能模块
func (b *Module) Create(ctx context.Context, productName, moduleName, desc string) (*tpl.ModuleRes, error) {
	productID, err := b.ms.Product.AcquireID(ctx, productName)
	if err != nil {
		return nil, err
	}

	module := &schema.Module{ProductID: productID, Name: moduleName, Desc: desc}
	if err = b.ms.Module.Create(ctx, module); err != nil {
		return nil, err
	}
	return &tpl.ModuleRes{Result: *module}, nil
}

// Update ...
func (b *Module) Update(ctx context.Context, productName, moduleName string, body tpl.ModuleUpdateBody) (*tpl.ModuleRes, error) {
	productID, err := b.ms.Product.AcquireID(ctx, productName)
	if err != nil {
		return nil, err
	}

	module, err := b.ms.Module.Acquire(ctx, productID, moduleName)
	if err != nil {
		return nil, err
	}

	module, err = b.ms.Module.Update(ctx, module.ID, body.ToMap())
	if err != nil {
		return nil, err
	}
	return &tpl.ModuleRes{Result: *module}, nil
}

// Offline 下线功能模块
func (b *Module) Offline(ctx context.Context, productName, moduleName string) (*tpl.BoolRes, error) {
	productID, err := b.ms.Product.AcquireID(ctx, productName)
	if err != nil {
		return nil, err
	}

	res := &tpl.BoolRes{Result: false}
	module, err := b.ms.Module.FindByName(ctx, productID, moduleName, "id, `offline_at`")
	if err != nil {
		return nil, err
	}
	if module == nil {
		return nil, gear.ErrNotFound.WithMsgf("module %s not found", moduleName)
	}
	if module.OfflineAt == nil {
		if err = b.ms.Module.Offline(ctx, module.ID); err != nil {
			return nil, err
		}
		res.Result = true
	}
	return res, nil
}
