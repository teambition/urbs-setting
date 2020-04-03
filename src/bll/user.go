package bll

import (
	"context"
	"time"

	"github.com/teambition/gear"
	"github.com/teambition/urbs-setting/src/conf"
	"github.com/teambition/urbs-setting/src/model"
	"github.com/teambition/urbs-setting/src/schema"
	"github.com/teambition/urbs-setting/src/tpl"
)

// User ...
type User struct {
	ms *model.Models
}

// ListCachedLables ... 该接口不返回错误
func (b *User) ListCachedLables(ctx context.Context, uid, product string) *tpl.CacheLabelsInfoRes {
	now := time.Now().UTC().Unix()
	res := &tpl.CacheLabelsInfoRes{Result: []schema.UserCacheLabel{}, Timestamp: now}

	user, err := b.ms.User.FindByUID(ctx, uid, "id, `uid`, `active_at`, `labels`")
	if err != nil || user == nil {
		return res
	}

	// user 上缓存的 labels 过期，则刷新获取最新，RefreshUser 要考虑并发场景
	if user.ActiveAt == 0 {
		if user, err = b.ms.User.RefreshLabels(ctx, user.ID, now, true); err != nil {
			return res
		}
	} else if conf.Config.IsCacheLabelExpired(now-5, user.ActiveAt) {
		// 提前 5s 异步处理
		go b.ms.User.RefreshLabels(ctx, user.ID, now, false)
	}

	res.Result = user.GetLabels(product)
	res.Timestamp = user.ActiveAt
	return res
}

// RefreshCachedLables ...
func (b *User) RefreshCachedLables(ctx context.Context, uid string) error {
	user, err := b.ms.User.FindByUID(ctx, uid, "id, `uid`, `active_at`, `labels`")
	if err != nil {
		return err
	}
	if user == nil {
		return gear.ErrNotFound.WithMsgf("user %s not found", uid)
	}

	_, err = b.ms.User.RefreshLabels(ctx, user.ID, time.Now().UTC().Unix(), true)
	return err
}

// ListLables ...
func (b *User) ListLables(ctx context.Context, uid string, pg tpl.Pagination) (*tpl.LabelsInfoRes, error) {
	user, err := b.ms.User.FindByUID(ctx, uid, "id")
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, gear.ErrNotFound.WithMsgf("user %s not found", uid)
	}

	labels, err := b.ms.User.FindLables(ctx, user.ID, pg)
	if err != nil {
		return nil, err
	}
	total, err := b.ms.User.CountLabels(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	res := &tpl.LabelsInfoRes{Result: labels}
	res.TotalSize = total
	if len(res.Result) > pg.PageSize {
		res.NextPageToken = tpl.IDToPageToken(res.Result[pg.PageSize].ID)
		res.Result = res.Result[:pg.PageSize]
	}
	return res, nil
}

// ListSettings ...
func (b *User) ListSettings(ctx context.Context, uid, productName string, pg tpl.Pagination) (*tpl.MySettingsRes, error) {
	user, err := b.ms.User.FindByUID(ctx, uid, "id")
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, gear.ErrNotFound.WithMsgf("user %s not found", uid)
	}

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

	moduleIDs, err := b.ms.Module.FindIDsByProductID(ctx, product.ID)
	if err != nil {
		return nil, err
	}
	settings, err := b.ms.User.FindSettings(ctx, user.ID, moduleIDs, pg)
	if err != nil {
		return nil, err
	}

	res := &tpl.MySettingsRes{Result: settings}
	if len(res.Result) > pg.PageSize {
		res.NextPageToken = tpl.IDToPageToken(res.Result[pg.PageSize].ID)
		res.Result = res.Result[:pg.PageSize]
	}
	return res, nil
}

// ListSettingsUnionAll ...
func (b *User) ListSettingsUnionAll(ctx context.Context, uid, productName, channel, client string, pg tpl.Pagination) (*tpl.MySettingsRes, error) {
	res := &tpl.MySettingsRes{Result: []tpl.MySetting{}}
	user, err := b.ms.User.FindByUID(ctx, uid, "id")
	if err != nil {
		return nil, err
	}
	if user == nil {
		return res, nil
	}

	product, err := b.ms.Product.FindByName(ctx, productName, "id, `offline_at`, `deleted_at`")
	if err != nil {
		return nil, err
	}
	if product == nil {
		return res, gear.ErrNotFound.WithMsgf("product %s not found", productName)
	}
	if product.DeletedAt != nil {
		return nil, gear.ErrNotFound.WithMsgf("product %s was deleted", productName)
	}
	if product.OfflineAt != nil {
		return nil, gear.ErrNotFound.WithMsgf("product %s was offline", productName)
	}

	moduleIDs, err := b.ms.Module.FindIDsByProductID(ctx, product.ID)
	if err != nil {
		return nil, err
	}

	groupIDs, err := b.ms.Group.FindIDsByUserID(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	settings, err := b.ms.User.FindSettingsUnionAll(ctx, user.ID, groupIDs, moduleIDs, pg, channel, client)
	if err != nil {
		return nil, err
	}

	res.Result = settings
	if len(res.Result) > pg.PageSize {
		res.NextPageToken = tpl.TimeToPageToken(res.Result[pg.PageSize].UpdatedAt)
		res.Result = res.Result[:pg.PageSize]
	}
	return res, nil
}

// CheckExists ...
func (b *User) CheckExists(ctx context.Context, uid string) bool {
	user, _ := b.ms.User.FindByUID(ctx, uid, "id")
	return user != nil
}

// BatchAdd ...
func (b *User) BatchAdd(ctx context.Context, users []string) error {
	return b.ms.User.BatchAdd(ctx, users)
}

// RemoveLable ...
func (b *User) RemoveLable(ctx context.Context, uid string, lableID int64) error {
	user, _ := b.ms.User.FindByUID(ctx, uid, "id")
	if user == nil {
		return gear.ErrNotFound.WithMsgf("User not found: %s", uid)
	}

	return b.ms.User.RemoveLable(ctx, user.ID, lableID)
}

// RollbackSetting ...
func (b *User) RollbackSetting(ctx context.Context, uid string, settingID int64) error {
	user, _ := b.ms.User.FindByUID(ctx, uid, "id")
	if user == nil {
		return gear.ErrNotFound.WithMsgf("User not found: %s", uid)
	}

	return b.ms.User.RollbackSetting(ctx, user.ID, settingID)
}

// RemoveSetting ...
func (b *User) RemoveSetting(ctx context.Context, uid string, settingID int64) error {
	user, _ := b.ms.User.FindByUID(ctx, uid, "id")
	if user == nil {
		return gear.ErrNotFound.WithMsgf("User not found: %s", uid)
	}

	return b.ms.User.RemoveSetting(ctx, user.ID, settingID)
}
