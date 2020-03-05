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

// ListLablesInCache ...
func (b *User) ListLablesInCache(ctx context.Context, uid, product string) (*tpl.CacheLabelsInfoRes, error) {
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
		user, err = b.ms.User.RefreshLabels(ctx, user.ID, now)
		if err != nil {
			return res, nil
		}
	}

	labels := user.GetLabels()
	if result, ok := labels[product]; ok {
		res.Result = result
		res.Timestamp = user.ActiveAt
	}
	return res, nil
}

// ListLables ...
func (b *User) ListLables(ctx context.Context, uid, product string) (*tpl.LabelsInfoRes, error) {
	user, err := b.ms.User.FindByUID(ctx, uid, "id")
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, gear.ErrNotFound.WithMsgf("user %s not found", uid)
	}

	labels, err := b.ms.User.FindLables(ctx, user.ID, product)
	if err != nil {
		return nil, err
	}
	res := &tpl.LabelsInfoRes{Result: labels}
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

// RemoveSetting ...
func (b *User) RemoveSetting(ctx context.Context, uid string, settingID int64) error {
	user, _ := b.ms.User.FindByUID(ctx, uid, "id")
	if user == nil {
		return gear.ErrNotFound.WithMsgf("User not found: %s", uid)
	}

	return b.ms.User.RemoveSetting(ctx, user.ID, settingID)
}
