package bll

import (
	"context"

	"github.com/teambition/gear"
	"github.com/teambition/urbs-setting/src/model"
	"github.com/teambition/urbs-setting/src/schema"
	"github.com/teambition/urbs-setting/src/tpl"
)

// Label ...
type Label struct {
	ms *model.Models
}

// List 返回产品下的标签列表
func (b *Label) List(ctx context.Context, productName string, pg tpl.Pagination) (*tpl.LabelsInfoRes, error) {
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
	labels, err := b.ms.Label.Find(ctx, product.ID, pg)
	if err != nil {
		return nil, err
	}
	total, err := b.ms.Label.Count(ctx, product.ID)
	if err != nil {
		return nil, err
	}

	labelInfos := tpl.LabelsInfoFrom(labels, productName)
	res := &tpl.LabelsInfoRes{Result: labelInfos}
	res.TotalSize = total
	if len(res.Result) > pg.PageSize {
		res.NextPageToken = tpl.IDToPageToken(res.Result[pg.PageSize].ID)
		res.Result = res.Result[:pg.PageSize]
	}
	return res, nil
}

// Create 创建标签
func (b *Label) Create(ctx context.Context, productName, labelName, desc string) (*tpl.LabelInfoRes, error) {
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

	label := schema.Label{ProductID: product.ID, Name: labelName, Desc: desc}
	if err = b.ms.Label.Create(ctx, &label); err != nil {
		return nil, err
	}
	return &tpl.LabelInfoRes{Result: tpl.LabelInfoFrom(label, productName)}, nil
}

// Offline 下线标签
func (b *Label) Offline(ctx context.Context, productName, labelName string) (*tpl.BoolRes, error) {
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
	label, err := b.ms.Label.FindByName(ctx, product.ID, labelName, "id, `offline_at`")
	if err != nil {
		return nil, err
	}
	if label == nil {
		return nil, gear.ErrNotFound.WithMsgf("label %s not found", labelName)
	}
	if label.OfflineAt == nil {
		if err = b.ms.Label.Offline(ctx, label.ID); err != nil {
			return nil, err
		}
		res.Result = true
	}
	return res, nil
}

// Assign 把标签批量分配给用户或群组
func (b *Label) Assign(ctx context.Context, productName, labelName string, users, groups []string) error {
	product, err := b.ms.Product.FindByName(ctx, productName, "id, `offline_at`")
	if err != nil {
		return err
	}
	if product == nil {
		return gear.ErrNotFound.WithMsgf("product %s not found", productName)
	}
	if product.OfflineAt != nil {
		return gear.ErrNotFound.WithMsgf("product %s was offline", productName)
	}

	label, err := b.ms.Label.FindByName(ctx, product.ID, labelName, "id, `offline_at`")
	if err != nil {
		return err
	}
	if label == nil {
		return gear.ErrNotFound.WithMsgf("label %s not found", labelName)
	}
	if label.OfflineAt != nil {
		return gear.ErrNotFound.WithMsgf("label %s was offline", labelName)
	}
	return b.ms.Label.Assign(ctx, label.ID, users, groups)
}

// Delete 物理删除标签
func (b *Label) Delete(ctx context.Context, productName, labelName string) (*tpl.BoolRes, error) {
	product, err := b.ms.Product.FindByName(ctx, productName, "id")
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, gear.ErrNotFound.WithMsgf("product %s not found", productName)
	}

	label, err := b.ms.Label.FindByName(ctx, product.ID, labelName, "id, `offline_at`")
	if err != nil {
		return nil, err
	}

	res := &tpl.BoolRes{Result: false}
	if label != nil {
		if label.OfflineAt == nil {
			return nil, gear.ErrConflict.WithMsgf("label %s is not offline", labelName)
		}

		if err = b.ms.Label.Delete(ctx, label.ID); err != nil {
			return nil, err
		}
		res.Result = true
	}
	return res, nil
}
