package bll

import (
	"context"
	"time"

	"github.com/teambition/urbs-setting/src/model"
	"github.com/teambition/urbs-setting/src/schema"
	"github.com/teambition/urbs-setting/src/tpl"
)

// User ...
type User struct {
	ms *model.Models
}

// GetLables ...
func (b *User) GetLables(ctx context.Context, uid, product, client, channel string) (*tpl.LabelsResponse, error) {
	user, err := b.ms.User.FindByUID(ctx, uid, "id, `uid`, `active_at`, `labels`")
	if err != nil {
		return nil, err
	}

	res := &tpl.LabelsResponse{Result: []schema.UserLabelInfo{}}
	if user == nil {
		return res, nil // user 不存在，返回空
	}

	now := time.Now().UTC().Unix()
	if user.IsStale(now) {
		// user 上缓存的 labels 过期，则刷新获取最新，RefreshUser 要考虑并发场景
		user.Labels, err = b.ms.User.RefreshLabels(ctx, user.ID, now)
		if err != nil {
			return res, nil
		}
	}

	labels := user.GetLabels()
	res.Result = make([]schema.UserLabelInfo, 0, len(labels))
	for _, l := range labels {
		if l.Product == product {
			if client != "" && l.Client != "" && l.Client != client {
				continue
			}
			if channel != "" && l.Channel != "" && l.Channel != channel {
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
