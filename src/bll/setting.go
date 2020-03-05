package bll

import (
	"context"

	"github.com/teambition/gear"
	"github.com/teambition/urbs-setting/src/model"
	"github.com/teambition/urbs-setting/src/schema"
	"github.com/teambition/urbs-setting/src/tpl"
)

// Setting ...
type Setting struct {
	ms *model.Models
}

// List 返回产品下的功能模块配置项列表，TODO：支持分页
func (b *Setting) List(ctx context.Context, productName, moduleName string) (*tpl.SettingsRes, error) {
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

	module, err := b.ms.Module.FindByName(ctx, product.ID, moduleName, "id")
	if err != nil {
		return nil, err
	}
	if module == nil {
		return nil, gear.ErrNotFound.WithMsgf("module %s not found", moduleName)
	}

	settings, err := b.ms.Setting.Find(ctx, module.ID)
	if err != nil {
		return nil, err
	}

	res := &tpl.SettingsRes{Result: settings}
	return res, nil
}

// Create 创建功能模块配置项
func (b *Setting) Create(ctx context.Context, productName, moduleName, settingName, desc string) (*tpl.SettingRes, error) {
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
		return nil, gear.ErrNotFound.WithMsgf("module %s not found", moduleName)
	}
	if module.OfflineAt != nil {
		return nil, gear.ErrNotFound.WithMsgf("module %s was offline", moduleName)
	}

	setting := &schema.Setting{ModuleID: module.ID, Name: settingName, Desc: desc}
	if err = b.ms.Setting.Create(ctx, setting); err != nil {
		return nil, err
	}
	return &tpl.SettingRes{Result: *setting}, nil
}

// Offline 下线功能模块配置项
func (b *Setting) Offline(ctx context.Context, productName, moduleName, settingName string) (*tpl.BoolRes, error) {
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
		return nil, gear.ErrNotFound.WithMsgf("module %s not found", moduleName)
	}
	if module.OfflineAt != nil {
		return nil, gear.ErrNotFound.WithMsgf("module %s was offline", moduleName)
	}

	res := &tpl.BoolRes{Result: false}
	setting, err := b.ms.Setting.FindByName(ctx, module.ID, settingName, "id, `offline_at`")
	if err != nil {
		return nil, err
	}
	if setting == nil {
		return nil, gear.ErrNotFound.WithMsgf("setting %s not found", settingName)
	}
	if setting.OfflineAt == nil {
		if err = b.ms.Setting.Offline(ctx, setting.ID); err != nil {
			return nil, err
		}
		res.Result = true
	}
	return res, nil
}

// Assign 把配置项批量分配给用户或群组
func (b *Setting) Assign(ctx context.Context, productName, moduleName, settingName, value string, users, groups []string) error {
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

	module, err := b.ms.Module.FindByName(ctx, product.ID, moduleName, "id, `offline_at`")
	if err != nil {
		return err
	}
	if module == nil {
		return gear.ErrNotFound.WithMsgf("module %s not found", moduleName)
	}
	if module.OfflineAt != nil {
		return gear.ErrNotFound.WithMsgf("module %s was offline", moduleName)
	}

	setting, err := b.ms.Setting.FindByName(ctx, module.ID, settingName, "id, `offline_at`")
	if err != nil {
		return err
	}
	if setting == nil {
		return gear.ErrNotFound.WithMsgf("setting %s not found", settingName)
	}
	if setting.OfflineAt != nil {
		return gear.ErrNotFound.WithMsgf("setting %s was offline", settingName)
	}
	return b.ms.Setting.Assign(ctx, setting.ID, value, users, groups)
}

// Delete 物理删除配置项
func (b *Setting) Delete(ctx context.Context, productName, moduleName, settingName string) (*tpl.BoolRes, error) {
	product, err := b.ms.Product.FindByName(ctx, productName, "id")
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, gear.ErrNotFound.WithMsgf("product %s not found", productName)
	}

	module, err := b.ms.Module.FindByName(ctx, product.ID, moduleName, "id")
	if err != nil {
		return nil, err
	}
	if module == nil {
		return nil, gear.ErrNotFound.WithMsgf("module %s not found", moduleName)
	}

	setting, err := b.ms.Setting.FindByName(ctx, module.ID, settingName, "id, `offline_at`")
	if err != nil {
		return nil, err
	}

	res := &tpl.BoolRes{Result: false}
	if setting != nil {
		if setting.OfflineAt == nil {
			return nil, gear.ErrConflict.WithMsgf("setting %s is not offline", settingName)
		}
		if err = b.ms.Setting.Delete(ctx, setting.ID); err != nil {
			return nil, err
		}
		res.Result = true
	}
	return res, nil
}
