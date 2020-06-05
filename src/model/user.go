package model

import (
	"context"
	"database/sql"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/teambition/gear"
	"github.com/teambition/urbs-setting/src/conf"
	"github.com/teambition/urbs-setting/src/schema"
	"github.com/teambition/urbs-setting/src/service"
	"github.com/teambition/urbs-setting/src/tpl"
	"github.com/teambition/urbs-setting/src/util"
)

// User ...
type User struct {
	*Model
}

// FindByUID 根据 uid 返回 user 数据
func (m *User) FindByUID(ctx context.Context, uid string, selectStr string) (*schema.User, error) {
	user := &schema.User{}
	ok, err := m.findOneByCols(ctx, schema.TableUser, goqu.Ex{"uid": uid}, selectStr, user)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}
	return user, nil
}

// Acquire ...
func (m *User) Acquire(ctx context.Context, uid string) (*schema.User, error) {
	user, err := m.FindByUID(ctx, uid, "")
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, gear.ErrNotFound.WithMsgf("user %s not found", uid)
	}
	return user, nil
}

// AcquireID ...
func (m *User) AcquireID(ctx context.Context, uid string) (int64, error) {
	user, err := m.FindByUID(ctx, uid, "id")
	if err != nil {
		return 0, err
	}
	if user == nil {
		return 0, gear.ErrNotFound.WithMsgf("user %s not found", uid)
	}
	return user.ID, nil
}

// Find 根据条件查找 products
func (m *User) Find(ctx context.Context, pg tpl.Pagination) ([]schema.User, int, error) {
	users := make([]schema.User, 0)
	cursor := pg.TokenToID()
	sdc := m.DB.From(schema.TableUser)
	sd := m.DB.From(schema.TableUser).Where(goqu.C("id").Lte(cursor))

	var total int64
	var err error
	if pg.Q != "" {
		sdc = sdc.Where(goqu.C("uid").ILike(pg.Q))
		total, err = sdc.CountContext(ctx)
		if err != nil {
			return nil, 0, err
		}

		sd = sd.Where(goqu.C("uid").ILike(pg.Q))
	}

	sd = sd.Order(goqu.C("id").Desc()).Limit(uint(pg.PageSize + 1))
	err = sd.Executor().ScanStructsContext(ctx, &users)
	if err != nil {
		return nil, 0, err
	}
	return users, int(total), nil
}

// RefreshLabels 更新 user 上的 labels 缓存，包括通过 group 关系获得的 labels
func (m *User) RefreshLabels(ctx context.Context, id int64, now int64, force bool) (*schema.User, []int64, bool, error) {
	user := &schema.User{}
	labelIDs := make([]int64, 0)
	refreshed := false
	tx, err := m.DB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return nil, nil, false, err
	}

	err = tx.Wrap(func() error {
		// 指定 id 的记录被锁住，如果表中无符合记录的数据则排他锁不生效
		sd := tx.From(schema.TableUser).Where(goqu.C("id").Eq(id)).
			ForUpdate(exp.Wait).Order(goqu.C("id").Asc()).Limit(1)

		ok, err := sd.Executor().ScanStructContext(ctx, user)
		if err != nil {
			return err
		}
		if !ok {
			return gear.ErrNotFound.WithMsgf("user %d not found for RefreshLabels", id)
		}

		if !force && !conf.Config.IsCacheLabelExpired(now-5, user.ActiveAt) {
			// 已被其它请求更新
			return nil
		}

		data := make(schema.UserCacheLabels)

		sd = tx.Select(
			goqu.I("t1.created_at"),
			goqu.I("t2.id"),
			goqu.I("t2.name"),
			goqu.I("t2.channels"),
			goqu.I("t2.clients"),
			goqu.I("t3.name").As("product")).
			From(
				goqu.T(schema.TableUserLabel).As("t1"),
				goqu.T(schema.TableLabel).As("t2"),
				goqu.T(schema.TableProduct).As("t3")).
			Where(
				goqu.I("t1.user_id").Eq(id),
				goqu.I("t1.label_id").Eq(goqu.I("t2.id")),
				goqu.I("t2.product_id").Eq(goqu.I("t3.id"))).
			Order(goqu.I("t1.id").Desc()).Limit(200)

		sd = sd.UnionAll(tx.Select(
			goqu.I("t2.created_at"),
			goqu.I("t3.id"),
			goqu.I("t3.name"),
			goqu.I("t3.channels"),
			goqu.I("t3.clients"),
			goqu.I("t4.name").As("product")).
			From(
				goqu.T(schema.TableUserGroup).As("t1"),
				goqu.T(schema.TableGroupLabel).As("t2"),
				goqu.T(schema.TableLabel).As("t3"),
				goqu.T(schema.TableProduct).As("t4")).
			Where(
				goqu.I("t1.user_id").Eq(id),
				goqu.I("t1.group_id").Eq(goqu.I("t2.group_id")),
				goqu.I("t2.label_id").Eq(goqu.I("t3.id")),
				goqu.I("t3.product_id").Eq(goqu.I("t4.id"))).
			Order(goqu.I("t2.id").Desc()).Limit(200)).
			Order(goqu.C("created_at").Desc())

		scanner, err := sd.Executor().ScannerContext(ctx)
		if err != nil {
			return err
		}

		set := make(map[int64]struct{})
		for scanner.Next() {
			myLabelInfo := schema.MyLabelInfo{}
			if err := scanner.ScanStruct(&myLabelInfo); err != nil {
				scanner.Close()
				return err
			}
			if _, ok := set[myLabelInfo.ID]; ok {
				continue // 去重
			}
			set[myLabelInfo.ID] = struct{}{}

			labelIDs = append(labelIDs, myLabelInfo.ID)
			arr, ok := data[myLabelInfo.Product]
			if !ok {
				arr = make([]schema.UserCacheLabel, 0)
			}
			data[myLabelInfo.Product] = append(arr, schema.UserCacheLabel{
				Label:    myLabelInfo.Name,
				Clients:  tpl.StringToSlice(myLabelInfo.Clients),
				Channels: tpl.StringToSlice(myLabelInfo.Channels),
			})
		}

		if err := scanner.Close(); err != nil {
			return err
		}
		if err := scanner.Err(); err != nil {
			return err
		}

		refreshed = true
		user.ActiveAt = time.Now().UTC().Unix()
		_ = user.PutLabels(data)
		_, err = service.DeResult(tx.Update(schema.TableUser).
			Where(goqu.C("id").Eq(id)).
			Set(goqu.Record{"labels": user.Labels, "active_at": user.ActiveAt}).
			Executor().ExecContext(ctx))
		return err
	})

	if err != nil {
		return nil, nil, false, err
	}
	return user, labelIDs, refreshed, nil
}

// FindSettingsUnionAll 根据用户 ID, updateGt, productName 返回其 settings 数据。
func (m *User) FindSettingsUnionAll(ctx context.Context, groupIDs []int64, userID, productID, moduleID, settingID int64, pg tpl.Pagination, channel, client string) ([]tpl.MySetting, error) {
	data := []tpl.MySetting{}
	cursor := pg.TokenToTimestamp(time.Now().Add(time.Minute * 10))
	set := make(map[int64]struct{})
	size := pg.PageSize + 1

	s := m.DB.Select(
		goqu.I("t1.rls"),
		goqu.I("t1.updated_at").As("assigned_at"),
		goqu.I("t1.value"),
		goqu.I("t1.last_value"),
		goqu.I("t2.id"),
		goqu.I("t2.name"),
		goqu.I("t2.description"),
		goqu.I("t2.channels"),
		goqu.I("t2.clients"),
		goqu.I("t3.name").As("module"))

	for i := 0; i < 7; i++ { // 分页补偿最多 7 次
		sd := s.From(
			goqu.T(schema.TableUserSetting).As("t1"),
			goqu.T(schema.TableSetting).As("t2"),
			goqu.T(schema.TableModule).As("t3")).
			Where(
				goqu.I("t1.user_id").Eq(userID),
				goqu.L("unix_timestamp(`t1`.`updated_at`)*1000").Lte(cursor))

		if settingID > 0 {
			sd = sd.Where(
				goqu.I("t1.setting_id").Eq(settingID),
				goqu.I("t1.setting_id").Eq(goqu.I("t2.id")))
		} else if moduleID > 0 {
			sd = sd.Where(
				goqu.I("t1.setting_id").Eq(goqu.I("t2.id")),
				goqu.I("t2.module_id").Eq(moduleID))
		} else {
			sd = sd.Where(goqu.I("t1.setting_id").Eq(goqu.I("t2.id")))
		}

		if pg.Q != "" {
			sd = sd.Where(goqu.I("t2.name").ILike(pg.Q))
		}

		sd = sd.Where(
			goqu.I("t2.module_id").Eq(goqu.I("t3.id")),
			goqu.I("t3.product_id").Eq(productID)).
			Order(goqu.I("t1.updated_at").Desc()).Limit(uint(size))

		if len(groupIDs) > 0 {
			gsd := s.From(
				goqu.T(schema.TableGroupSetting).As("t1"),
				goqu.T(schema.TableSetting).As("t2"),
				goqu.T(schema.TableModule).As("t3")).
				Where(
					goqu.I("t1.group_id").In(groupIDs),
					goqu.L("unix_timestamp(`t1`.`updated_at`)*1000").Lte(cursor))

			if settingID > 0 {
				gsd = gsd.Where(
					goqu.I("t1.setting_id").Eq(settingID),
					goqu.I("t1.setting_id").Eq(goqu.I("t2.id")))
			} else if moduleID > 0 {
				gsd = gsd.Where(
					goqu.I("t1.setting_id").Eq(goqu.I("t2.id")),
					goqu.I("t2.module_id").Eq(moduleID))
			} else {
				gsd = gsd.Where(goqu.I("t1.setting_id").Eq(goqu.I("t2.id")))
			}

			if pg.Q != "" {
				gsd = gsd.Where(goqu.I("t2.name").ILike(pg.Q))
			}

			gsd = gsd.Where(
				goqu.I("t2.module_id").Eq(goqu.I("t3.id")),
				goqu.I("t3.product_id").Eq(productID)).
				Order(goqu.I("t1.updated_at").Desc()).Limit(uint(size))

			sd = sd.UnionAll(gsd).Order(goqu.C("assigned_at").Desc())
		}

		scanner, err := sd.Executor().ScannerContext(ctx)
		if err != nil {
			return nil, err
		}

		count := 0
		for scanner.Next() {
			count++
			mySetting := tpl.MySetting{}
			if err := scanner.ScanStruct(&mySetting); err != nil {
				scanner.Close()
				return nil, err
			}

			if _, ok := set[mySetting.ID]; ok {
				continue // 去重
			}
			set[mySetting.ID] = struct{}{}

			if mySetting.Channels != "" {
				if !tpl.StringSliceHas(tpl.StringToSlice(mySetting.Channels), channel) {
					continue // channel 不匹配
				}
			}
			if mySetting.Clients != "" {
				if !tpl.StringSliceHas(tpl.StringToSlice(mySetting.Clients), client) {
					continue // client 不匹配
				}
			}

			mySetting.HID = service.IDToHID(mySetting.ID, "setting")
			data = append(data, mySetting)
		}

		scanner.Close()
		if err := scanner.Err(); err != nil {
			return nil, err
		}

		if count < size {
			break // no data to select
		}
		if len(data) >= size {
			break // get enough
		}
		nt := data[len(data)-1].AssignedAt.UTC()
		// select next page
		cursor = nt.Unix()*1000 + int64(nt.Nanosecond()/1000000) - 1
	}

	return data, nil
}

// FindLabels 根据用户 ID 返回其 labels 数据。
func (m *User) FindLabels(ctx context.Context, userID int64, pg tpl.Pagination) ([]tpl.MyLabel, int, error) {
	data := []tpl.MyLabel{}
	cursor := pg.TokenToID()

	sdc := m.DB.Select().
		From(
			goqu.T(schema.TableUserLabel).As("t1"),
			goqu.T(schema.TableLabel).As("t2"),
			goqu.T(schema.TableProduct).As("t3")).
		Where(
			goqu.I("t1.user_id").Eq(userID),
			goqu.I("t1.label_id").Eq(goqu.I("t2.id")))

	sd := m.DB.Select(
		goqu.I("t1.rls"),
		goqu.I("t1.created_at").As("assigned_at"),
		goqu.I("t2.id"),
		goqu.I("t2.name"),
		goqu.I("t2.description"),
		goqu.I("t3.name").As("product")).
		From(
			goqu.T(schema.TableUserLabel).As("t1"),
			goqu.T(schema.TableLabel).As("t2"),
			goqu.T(schema.TableProduct).As("t3")).
		Where(
			goqu.I("t1.user_id").Eq(userID),
			goqu.I("t1.id").Lte(cursor),
			goqu.I("t1.label_id").Eq(goqu.I("t2.id")))

	if pg.Q != "" {
		sdc = sdc.Where(goqu.I("t2.name").ILike(pg.Q))
		sd = sd.Where(goqu.I("t2.name").ILike(pg.Q))
	}

	sdc = sdc.Where(goqu.I("t2.product_id").Eq(goqu.I("t3.id")))
	sd = sd.Where(goqu.I("t2.product_id").Eq(goqu.I("t3.id"))).
		Order(goqu.I("t1.id").Desc()).Limit(uint(pg.PageSize + 1))

	total, err := sdc.CountContext(ctx)
	if err != nil {
		return nil, 0, err
	}
	scanner, err := sd.Executor().ScannerContext(ctx)
	if err != nil {
		return nil, 0, err
	}
	defer scanner.Close()

	for scanner.Next() {
		myLabel := tpl.MyLabel{}
		if err := scanner.ScanStruct(&myLabel); err != nil {
			return nil, 0, err
		}
		myLabel.HID = service.IDToHID(myLabel.ID, "label")
		data = append(data, myLabel)
	}

	if err := scanner.Err(); err != nil {
		return nil, 0, err
	}
	return data, int(total), nil
}

// FindSettings 根据用户 ID 返回 settings 数据。
func (m *User) FindSettings(ctx context.Context, userID, productID, moduleID, settingID int64, pg tpl.Pagination) ([]tpl.MySetting, int, error) {
	data := []tpl.MySetting{}
	cursor := pg.TokenToID()
	sdc := m.DB.Select().
		From(
			goqu.T(schema.TableUserSetting).As("t1"),
			goqu.T(schema.TableSetting).As("t2"),
			goqu.T(schema.TableModule).As("t3"),
			goqu.T(schema.TableProduct).As("t4")).
		Where(goqu.I("t1.user_id").Eq(userID))

	sd := m.DB.Select(
		goqu.I("t1.rls"),
		goqu.I("t1.updated_at").As("assigned_at"),
		goqu.I("t1.value"),
		goqu.I("t1.last_value"),
		goqu.I("t2.id"),
		goqu.I("t2.name"),
		goqu.I("t2.description"),
		goqu.I("t3.name").As("module"),
		goqu.I("t4.name").As("product")).
		From(
			goqu.T(schema.TableUserSetting).As("t1"),
			goqu.T(schema.TableSetting).As("t2"),
			goqu.T(schema.TableModule).As("t3"),
			goqu.T(schema.TableProduct).As("t4")).
		Where(
			goqu.I("t1.user_id").Eq(userID),
			goqu.I("t1.id").Lte(cursor))

	if settingID > 0 {
		sdc = sdc.Where(
			goqu.I("t1.setting_id").Eq(settingID),
			goqu.I("t1.setting_id").Eq(goqu.I("t2.id")))
		sd = sd.Where(
			goqu.I("t1.setting_id").Eq(settingID),
			goqu.I("t1.setting_id").Eq(goqu.I("t2.id")))
	} else if moduleID > 0 {
		sdc = sdc.Where(
			goqu.I("t1.setting_id").Eq(goqu.I("t2.id")),
			goqu.I("t2.module_id").Eq(moduleID))
		sd = sd.Where(
			goqu.I("t1.setting_id").Eq(goqu.I("t2.id")),
			goqu.I("t2.module_id").Eq(moduleID))
	} else {
		sdc = sdc.Where(goqu.I("t1.setting_id").Eq(goqu.I("t2.id")))
		sd = sd.Where(goqu.I("t1.setting_id").Eq(goqu.I("t2.id")))
	}

	if pg.Q != "" {
		sdc = sdc.Where(goqu.I("t2.name").ILike(pg.Q))
		sd = sd.Where(goqu.I("t2.name").ILike(pg.Q))
	}

	sdc = sdc.Where(goqu.I("t2.module_id").Eq(goqu.I("t3.id")))
	sd = sd.Where(goqu.I("t2.module_id").Eq(goqu.I("t3.id")))
	if productID > 0 {
		sdc = sdc.Where(goqu.I("t3.product_id").Eq(productID))
		sd = sd.Where(goqu.I("t3.product_id").Eq(productID))
	}
	sd = sd.Where(goqu.I("t3.product_id").Eq(goqu.I("t4.id"))).
		Order(goqu.I("t1.id").Desc()).Limit(uint(pg.PageSize + 1))

	total, err := sdc.CountContext(ctx)
	if err != nil {
		return nil, 0, err
	}

	scanner, err := sd.Executor().ScannerContext(ctx)
	if err != nil {
		return nil, 0, err
	}
	defer scanner.Close()

	for scanner.Next() {
		mySetting := tpl.MySetting{}
		if err := scanner.ScanStruct(&mySetting); err != nil {
			return nil, 0, err
		}
		mySetting.HID = service.IDToHID(mySetting.ID, "setting")
		data = append(data, mySetting)
	}

	if err := scanner.Err(); err != nil {
		return nil, 0, err
	}
	return data, int(total), nil
}

// BatchAdd 批量添加用户
// uids 经过了 `^[0-9A-Za-z._-]{3,63}$` 正则验证
func (m *User) BatchAdd(ctx context.Context, uids []string) error {
	if len(uids) == 0 {
		return nil
	}

	vals := make([][]interface{}, len(uids))
	for i := range uids {
		vals[i] = goqu.Vals{uids[i]}
	}

	sd := m.DB.Insert(schema.TableUser).Cols("uid").Vals(vals...).OnConflict(goqu.DoNothing())
	rowsAffected, err := service.DeResult(sd.Executor().ExecContext(ctx))
	if rowsAffected > 0 {
		util.Go(30*time.Second, func(gctx context.Context) {
			m.tryRefreshUsersTotalSize(gctx)
		})
	}
	return err
}
