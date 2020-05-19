package bll

import (
	"context"
	"time"

	"github.com/teambition/urbs-setting/src/conf"
	"github.com/teambition/urbs-setting/src/model"
	"github.com/teambition/urbs-setting/src/schema"
	"github.com/teambition/urbs-setting/src/tpl"
)

// User ...
type User struct {
	ms *model.Models
}

// List 返回用户列表
func (b *User) List(ctx context.Context, pg tpl.Pagination) (*tpl.UsersRes, error) {
	users, total, err := b.ms.User.Find(ctx, pg)
	if err != nil {
		return nil, err
	}
	res := &tpl.UsersRes{Result: users}
	res.TotalSize = total

	if res.TotalSize == 0 && pg.Q == "" {
		statistic, _ := b.ms.Statistic.FindByKey(ctx, schema.UsersTotalSize)
		if statistic != nil {
			res.TotalSize = int(statistic.Status)
		}
	}
	if len(res.Result) > pg.PageSize {
		res.NextPageToken = tpl.IDToPageToken(res.Result[pg.PageSize].ID)
		res.Result = res.Result[:pg.PageSize]
	}
	return res, nil
}

// ListCachedLabels ... 该接口不返回错误
func (b *User) ListCachedLabels(ctx context.Context, uid, product string) *tpl.CacheLabelsInfoRes {
	now := time.Now().UTC()
	res := &tpl.CacheLabelsInfoRes{Result: []schema.UserCacheLabel{}, Timestamp: now.Unix()}
	if product == "" {
		return res
	}

	user, err := b.ms.User.Acquire(ctx, uid)
	if err != nil {
		return res
	}

	// user 上缓存的 labels 过期，则刷新获取最新，RefreshUser 要考虑并发场景
	if user.ActiveAt == 0 {
		if user = b.ms.TryApplyLabelRulesAndRefreshUserLabels(ctx, user.ID, now, true); user == nil {
			return res
		}
	} else if conf.Config.IsCacheLabelExpired(now.Unix()-5, user.ActiveAt) {
		// 提前 5s 异步处理
		go b.ms.TryApplyLabelRulesAndRefreshUserLabels(ctx, user.ID, now, false)
	}

	res.Result = user.GetLabels(product)
	res.Timestamp = user.ActiveAt
	return res
}

// RefreshCachedLabels ...
func (b *User) RefreshCachedLabels(ctx context.Context, uid string) (*schema.User, error) {
	user, err := b.ms.User.Acquire(ctx, uid)
	if err != nil {
		return nil, err
	}

	if user, err = b.ms.ApplyLabelRulesAndRefreshUserLabels(ctx, user.ID, time.Now().UTC(), true); err != nil {
		return nil, err
	}
	return user, nil
}

// ListLabels ...
func (b *User) ListLabels(ctx context.Context, uid string, pg tpl.Pagination) (*tpl.MyLabelsRes, error) {
	user, err := b.ms.User.Acquire(ctx, uid)
	if err != nil {
		return nil, err
	}

	labels, total, err := b.ms.User.FindLabels(ctx, user.ID, pg)
	if err != nil {
		return nil, err
	}

	res := &tpl.MyLabelsRes{Result: labels}
	res.TotalSize = total
	if len(res.Result) > pg.PageSize {
		res.NextPageToken = tpl.IDToPageToken(res.Result[pg.PageSize].ID)
		res.Result = res.Result[:pg.PageSize]
	}
	return res, nil
}

// ListSettings ...
func (b *User) ListSettings(ctx context.Context, uid string, pg tpl.Pagination) (*tpl.MySettingsRes, error) {
	user, err := b.ms.User.Acquire(ctx, uid)
	if err != nil {
		return nil, err
	}

	settings, total, err := b.ms.User.FindSettings(ctx, user.ID, pg)
	if err != nil {
		return nil, err
	}

	res := &tpl.MySettingsRes{Result: settings}
	res.TotalSize = total
	if len(res.Result) > pg.PageSize {
		res.NextPageToken = tpl.IDToPageToken(res.Result[pg.PageSize].ID)
		res.Result = res.Result[:pg.PageSize]
	}
	return res, nil
}

// ListSettingsUnionAll ...
func (b *User) ListSettingsUnionAll(ctx context.Context, uid, productName, channel, client string, pg tpl.Pagination) (*tpl.MySettingsRes, error) {
	res := &tpl.MySettingsRes{Result: []tpl.MySetting{}}
	user, err := b.ms.User.Acquire(ctx, uid)
	if err != nil {
		return res, nil
	}

	productID, err := b.ms.Product.AcquireID(ctx, productName)
	if err != nil {
		return nil, err
	}

	moduleIDs, err := b.ms.Module.FindIDsByProductID(ctx, productID)
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
	for i := range settings {
		settings[i].Product = productName
	}
	if pg.PageToken == "" { // 请求首页时尝试应用 SettingRules
		go b.ms.TryApplySettingRules(ctx, productID, user.ID)
	}

	res.Result = settings
	if len(res.Result) > pg.PageSize {
		res.NextPageToken = tpl.TimeToPageToken(res.Result[pg.PageSize].AssignedAt)
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
