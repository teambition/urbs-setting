package bll

import (
	"context"

	"github.com/teambition/gear"
	"github.com/teambition/urbs-setting/src/model"
	"github.com/teambition/urbs-setting/src/schema"
	"github.com/teambition/urbs-setting/src/tpl"
)

// Product ...
type Product struct {
	ms *model.Models
}

// List 返回产品列表，TODO：支持分页
func (b *Product) List(ctx context.Context) (*tpl.ProductsRes, error) {
	products, err := b.ms.Product.Find(ctx)
	if err != nil {
		return nil, err
	}
	res := &tpl.ProductsRes{Result: products}
	return res, nil
}

// Create 创建产品
func (b *Product) Create(ctx context.Context, name, desc string) (*tpl.ProductRes, error) {
	product := &schema.Product{Name: name, Desc: desc}
	if err := b.ms.Product.Create(ctx, product); err != nil {
		return nil, err
	}
	res := &tpl.ProductRes{Result: *product}
	return res, nil
}

// Offline 下线产品
func (b *Product) Offline(ctx context.Context, productName string) (*tpl.BoolRes, error) {
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

	res := &tpl.BoolRes{Result: false}
	if product.OfflineAt == nil {
		if err = b.ms.Product.Offline(ctx, product.ID); err != nil {
			return nil, err
		}
		res.Result = true
	}
	return res, nil
}

// Delete 逻辑删除产品
func (b *Product) Delete(ctx context.Context, productName string) (*tpl.BoolRes, error) {
	product, err := b.ms.Product.FindByName(ctx, productName, "id, `offline_at`, `deleted_at`")
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, gear.ErrNotFound.WithMsgf("product %s not found", productName)
	}
	if product.OfflineAt != nil {
		return nil, gear.ErrBadRequest.WithMsgf("product %s is not offline", productName)
	}

	res := &tpl.BoolRes{Result: false}
	if product.DeletedAt == nil {
		if err = b.ms.Product.Delete(ctx, product.ID); err != nil {
			return nil, err
		}
		res.Result = true
	}
	return res, nil
}
