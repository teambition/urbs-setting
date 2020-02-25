package bll

import (
	"context"
	"strings"
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

// GetLablesInCache ...
func (b *User) GetLablesInCache(ctx context.Context, uid, product, client, channel string) (*tpl.CacheLabelsInfoRes, error) {
	user, err := b.ms.User.FindByUID(ctx, uid, "id, `uid`, `active_at`, `labels`")
	if err != nil {
		return nil, err
	}

	res := &tpl.CacheLabelsInfoRes{Result: []schema.UserCacheLabel{}}
	if user == nil {
		return res, nil // user 不存在，返回空
	}

	now := time.Now().UTC().Unix()
	if conf.Config.IsCacheLabelExpired(now, user.ActiveAt) {
		// user 上缓存的 labels 过期，则刷新获取最新，RefreshUser 要考虑并发场景
		user.Labels, err = b.ms.User.RefreshLabels(ctx, user.ID, now)
		if err != nil {
			return res, nil
		}
	}

	labels := user.GetLabels()
	res.Result = make([]schema.UserCacheLabel, 0, len(labels))
	for _, l := range labels {
		if l.Product == product {
			if client != "" && l.Clients != "" && !strings.Contains(l.Clients, client) {
				continue
			}
			if channel != "" && l.Channels != "" && !strings.Contains(l.Channels, channel) {
				continue
			}
			r := l
			res.Result = append(res.Result, r)
		}
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

// RemoveSetting ...
func (b *User) RemoveSetting(ctx context.Context, uid string, settingID int64) error {
	user, _ := b.ms.User.FindByUID(ctx, uid, "id")
	if user == nil {
		return gear.ErrNotFound.WithMsgf("User not found: %s", uid)
	}

	return b.ms.User.RemoveSetting(ctx, user.ID, settingID)
}
