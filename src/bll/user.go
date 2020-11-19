package bll

import (
	"context"
	"strings"
	"time"

	"github.com/teambition/urbs-setting/src/conf"
	"github.com/teambition/urbs-setting/src/logging"
	"github.com/teambition/urbs-setting/src/model"
	"github.com/teambition/urbs-setting/src/schema"
	"github.com/teambition/urbs-setting/src/tpl"
	"github.com/teambition/urbs-setting/src/util"
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

	readCtx := context.WithValue(ctx, model.ReadDB, true)
	productID, err := b.ms.Product.AcquireID(readCtx, product)
	if err != nil {
		return res
	}

	user, err := b.ms.User.Acquire(readCtx, uid)
	if err != nil {
		if strings.HasPrefix(uid, "anon-") {
			if labels, err := b.ms.LabelRule.ApplyRulesToAnonymous(ctx, uid, productID, schema.RuleUserPercent); err == nil {
				res.Result = labels
			}
		}
		return res
	}

	activeAt := user.GetCache(product).ActiveAt
	// user 上缓存的 labels 过期，则刷新获取最新，RefreshUser 要考虑并发场景
	if activeAt == 0 {
		if user = b.ms.TryApplyLabelRulesAndRefreshUserLabels(ctx, productID, product, user.ID, now, true); user == nil {
			return res
		}
	} else if conf.Config.IsCacheLabelExpired(now.Unix()-5, activeAt) { // 提前 5s 异步处理
		if conf.Config.IsCacheLabelDoubleExpired(now.Unix(), activeAt) { // 大于等于 2 倍过期时间的缓存，同步等待结果。
			if user = b.ms.TryApplyLabelRulesAndRefreshUserLabels(ctx, productID, product, user.ID, now, false); user == nil {
				return res
			}
		} else {
			util.Go(10*time.Second, func(gctx context.Context) {
				b.ms.TryApplyLabelRulesAndRefreshUserLabels(gctx, productID, product, user.ID, now, false)
			})
		}
	}
	userCache := user.GetCache(product)

	res.Result = userCache.Labels
	res.Timestamp = userCache.ActiveAt
	return res
}

// RefreshCachedLabels ...
func (b *User) RefreshCachedLabels(ctx context.Context, product, uid string) (*schema.User, error) {
	user, err := b.ms.User.Acquire(ctx, uid)
	if err != nil {
		return nil, err
	}
	readCtx := context.WithValue(ctx, model.ReadDB, true)
	var productID int64 = 0
	if product != "" {
		productID, err = b.ms.Product.AcquireID(readCtx, product)
		if err != nil {
			return nil, err
		}
	}
	if user, err = b.ms.ApplyLabelRulesAndRefreshUserLabels(ctx, productID, product, user.ID, time.Now().UTC(), true); err != nil {
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
func (b *User) ListSettings(ctx context.Context, req tpl.MySettingsQueryURL) (*tpl.MySettingsRes, error) {
	readCtx := context.WithValue(ctx, model.ReadDB, true)
	userID, err := b.ms.User.AcquireID(readCtx, req.UID)
	if err != nil {
		return nil, err
	}

	var productID int64
	var moduleID int64
	var settingID int64

	if req.Product != "" {
		productID, err = b.ms.Product.AcquireID(readCtx, req.Product)
		if err != nil {
			return nil, err
		}
	}
	if productID > 0 && req.Module != "" {
		moduleID, err = b.ms.Module.AcquireID(readCtx, productID, req.Module)
		if err != nil {
			return nil, err
		}
	}

	if moduleID > 0 && req.Setting != "" {
		settingID, err = b.ms.Setting.AcquireID(readCtx, moduleID, req.Setting)
		if err != nil {
			return nil, err
		}
	}

	pg := req.Pagination
	settings, total, err := b.ms.User.FindSettings(ctx, userID, productID, moduleID, settingID, pg, req.Channel, req.Client)
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
func (b *User) ListSettingsUnionAll(ctx context.Context, req tpl.MySettingsQueryURL) (*tpl.MySettingsRes, error) {
	res := &tpl.MySettingsRes{Result: []tpl.MySetting{}}
	readCtx := context.WithValue(ctx, model.ReadDB, true)

	var productID int64
	var moduleID int64
	var settingID int64
	productID, err := b.ms.Product.AcquireID(readCtx, req.Product)
	if err != nil {
		return nil, err
	}

	user, err := b.ms.User.Acquire(readCtx, req.UID)
	if err != nil {
		if strings.HasPrefix(req.UID, "anon-") {
			if settings, err := b.ms.SettingRule.ApplyRulesToAnonymous(ctx, req.UID, productID, req.Channel, req.Client, schema.RuleUserPercent); err == nil {
				for i := range settings {
					settings[i].Product = req.Product
				}
				res.Result = settings
			}
		}
		return res, nil
	}

	if req.Module != "" {
		moduleID, err = b.ms.Module.AcquireID(readCtx, productID, req.Module)
		if err != nil {
			return nil, err
		}
	}
	if req.Setting != "" {
		settingID, err = b.ms.Setting.AcquireID(readCtx, moduleID, req.Setting)
		if err != nil {
			return nil, err
		}
	}

	groupIDs, err := b.ms.Group.FindIDsByUser(readCtx, user.ID)
	if err != nil {
		return nil, err
	}

	pg := req.Pagination
	settings, err := b.ms.User.FindSettingsUnionAll(readCtx, groupIDs, user.ID, productID, moduleID, settingID, pg, req.Channel, req.Client)
	if err != nil {
		return nil, err
	}
	for i := range settings {
		settings[i].Product = req.Product
	}
	if pg.PageToken == "" { // 请求首页时尝试应用 SettingRules
		util.Go(10*time.Second, func(gctx context.Context) {
			b.ms.TryApplySettingRules(gctx, productID, user.ID)
		})
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
	user, _ := b.ms.User.FindByUID(context.WithValue(ctx, model.ReadDB, true), uid, "id")
	return user != nil
}

// BatchAdd ...
func (b *User) BatchAdd(ctx context.Context, users []string) error {
	return b.ms.User.BatchAdd(ctx, users)
}

// ApplyRules ...
func (b *User) ApplyRules(ctx context.Context, product string, body *tpl.ApplyRulesBody) error {
	readCtx := context.WithValue(ctx, model.ReadDB, true)
	productID, err := b.ms.Product.AcquireID(readCtx, product)
	if err != nil {
		return err
	}
	for _, UID := range body.Users {
		userID, err := b.ms.User.AcquireID(ctx, UID)
		if err != nil {
			logging.Warningf("newUserAcquireID: userID %d, error %v", userID, err)
			continue
		}
		_, err = b.ms.LabelRule.ApplyRules(ctx, productID, userID, []int64{}, body.Kind)
		if err != nil {
			logging.Warningf("newUserApplyLabelRules: userID %d, error %v", userID, err)
			continue
		}
		err = b.ms.SettingRule.ApplyRules(ctx, productID, userID, body.Kind)
		if err != nil {
			logging.Warningf("newUserApplySettingRules: userID %d, error %v", userID, err)
		}
	}
	return err
}
