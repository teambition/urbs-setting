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
	product, err := b.ms.Product.FindByName(ctx, productName, "id, `deleted_at`")
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, gear.ErrNotFound.WithMsgf("product %s not found", productName)
	}
	if product.DeletedAt != nil {
		return nil, gear.ErrNotFound.WithMsgf("product %s was deleted", productName)
	}
	modules, err := b.ms.Module.Find(ctx, product.ID, pg)
	if err != nil {
		return nil, err
	}

	total, err := b.ms.Module.Count(ctx, product.ID)
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
	product, err := b.ms.Product.FindByName(ctx, productName, "id, `offline_at`, `deleted_at`")
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

	module := &schema.Module{ProductID: product.ID, Name: moduleName, Desc: desc}
	if err = b.ms.Module.Create(ctx, module); err != nil {
		return nil, err
	}
	return &tpl.ModuleRes{Result: *module}, nil
}

// Update ...
func (b *Module) Update(ctx context.Context, productName, moduleName string, body tpl.ModuleUpdateBody) (*tpl.ModuleRes, error) {
	product, err := b.ms.Product.FindByName(ctx, productName, "id, `offline_at`, `deleted_at`")
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

	module, err := b.ms.Module.FindByName(ctx, product.ID, moduleName, "id, `offline_at`")
	if err != nil {
		return nil, err
	}
	if module == nil {
		return nil, gear.ErrNotFound.WithMsgf("label %s not found", moduleName)
	}
	if module.OfflineAt != nil {
		return nil, gear.ErrNotFound.WithMsgf("label %s was offline", moduleName)
	}

	module, err = b.ms.Module.Update(ctx, module.ID, body.ToMap())
	if err != nil {
		return nil, err
	}
	return &tpl.ModuleRes{Result: *module}, nil
}

// Offline 下线功能模块
func (b *Module) Offline(ctx context.Context, productName, moduleName string) (*tpl.BoolRes, error) {
	product, err := b.ms.Product.FindByName(ctx, productName, "id, `offline_at`, `deleted_at`")
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

	res := &tpl.BoolRes{Result: false}
	module, err := b.ms.Module.FindByName(ctx, product.ID, moduleName, "id, `offline_at`")
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
