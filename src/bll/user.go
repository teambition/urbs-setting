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

// ListCachedLables ...
func (b *User) ListCachedLables(ctx context.Context, uid, product string) (*tpl.CacheLabelsInfoRes, error) {
	user, err := b.ms.User.FindByUID(ctx, uid, "id, `uid`, `active_at`, `labels`")
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC().Unix()
	res := &tpl.CacheLabelsInfoRes{Result: []schema.UserCacheLabel{}, Timestamp: now}
	if user == nil {
		return res, nil // user 不存在，返回空
	}

	if conf.Config.IsCacheLabelExpired(now, user.ActiveAt) {
		// user 上缓存的 labels 过期，则刷新获取最新，RefreshUser 要考虑并发场景
		user, err = b.ms.User.RefreshLabels(ctx, user.ID, now, false)
		if err != nil {
			return res, nil
		}
	}

	labels := user.GetLabels(product)
	if len(labels) > 0 {
		res.Result = labels
		res.Timestamp = user.ActiveAt
	}
	return res, nil
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
		res.NextPageToken = tpl.TimeToPageToken(res.Result[pg.PageSize].UpdatedAt)
		res.Result = res.Result[:pg.PageSize]
	}
	return res, nil
}

// ListSettingsWithGroup ...
func (b *User) ListSettingsWithGroup(ctx context.Context, uid, productName string, pg tpl.Pagination) (*tpl.MySettingsRes, error) {
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

	groupIDs, err := b.ms.Group.FindIDsByUserID(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	settings, err := b.ms.User.FindSettingsWithGroup(ctx, user.ID, groupIDs, moduleIDs, pg)
	if err != nil {
		return nil, err
	}

	res := &tpl.MySettingsRes{Result: settings}
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
