package model

import (
	"context"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/teambition/gear"
	"github.com/teambition/urbs-setting/src/schema"
	"github.com/teambition/urbs-setting/src/service"
	"github.com/teambition/urbs-setting/src/tpl"
	"github.com/teambition/urbs-setting/src/util"
)

// Group ...
type Group struct {
	*Model
}

// FindByUID 根据 uid 返回 user 数据
func (m *Group) FindByUID(ctx context.Context, uid string, selectStr string) (*schema.Group, error) {
	group := &schema.Group{}
	ok, err := m.findOneByCols(ctx, schema.TableGroup, goqu.Ex{"uid": uid}, selectStr, group)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}
	return group, nil
}

// Acquire ...
func (m *Group) Acquire(ctx context.Context, uid string) (*schema.Group, error) {
	group, err := m.FindByUID(ctx, uid, "")
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, gear.ErrNotFound.WithMsgf("group %s not found", uid)
	}
	return group, nil
}

// AcquireID ...
func (m *Group) AcquireID(ctx context.Context, uid string) (int64, error) {
	group, err := m.FindByUID(ctx, uid, "id, uid")
	if err != nil {
		return 0, err
	}
	if group == nil {
		return 0, gear.ErrNotFound.WithMsgf("group %s not found", uid)
	}
	return group.ID, nil
}

// Find 根据条件查找 groups
func (m *Group) Find(ctx context.Context, kind string, pg tpl.Pagination) ([]schema.Group, int, error) {
	groups := make([]schema.Group, 0)
	cursor := pg.TokenToID()
	sdc := m.RdDB.From(schema.TableGroup)
	sd := m.RdDB.From(schema.TableGroup).Where(goqu.C("id").Lte(cursor))
	if kind != "" {
		sdc = sdc.Where(goqu.C("kind").Eq(kind))
		sd = sd.Where(goqu.C("kind").Eq(kind))
	}
	if pg.Q != "" {
		sdc = sdc.Where(goqu.C("uid").ILike(pg.Q))
		sd = sd.Where(goqu.C("uid").ILike(pg.Q))
	}
	sd = sd.Order(goqu.C("id").Desc()).Limit(uint(pg.PageSize + 1))

	total, err := sdc.CountContext(ctx)
	if err != nil {
		return nil, 0, err
	}
	err = sd.Executor().ScanStructsContext(ctx, &groups)
	if err != nil {
		return nil, 0, err
	}
	return groups, int(total), nil
}

// FindLabels 根据群组 ID 返回其 labels 数据。TODO：支持更多筛选条件和分页
func (m *Group) FindLabels(ctx context.Context, groupID int64, pg tpl.Pagination) ([]tpl.MyLabel, int, error) {
	data := make([]tpl.MyLabel, 0)
	cursor := pg.TokenToID()

	sdc := m.RdDB.Select().
		From(
			goqu.T(schema.TableGroupLabel).As("t1"),
			goqu.T(schema.TableLabel).As("t2"),
			goqu.T(schema.TableProduct).As("t3")).
		Where(
			goqu.I("t1.group_id").Eq(groupID),
			goqu.I("t1.label_id").Eq(goqu.I("t2.id")))

	sd := m.RdDB.Select(
		goqu.I("t1.rls"),
		goqu.I("t1.created_at").As("assigned_at"),
		goqu.I("t2.id"),
		goqu.I("t2.name"),
		goqu.I("t2.description"),
		goqu.I("t3.name").As("product")).
		From(
			goqu.T(schema.TableGroupLabel).As("t1"),
			goqu.T(schema.TableLabel).As("t2"),
			goqu.T(schema.TableProduct).As("t3")).
		Where(
			goqu.I("t1.group_id").Eq(groupID),
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

// FindSettings 根据 Group ID, updateGt, productName 返回其 settings 数据。
func (m *Group) FindSettings(ctx context.Context, groupID, productID, moduleID, settingID int64, pg tpl.Pagination, channel, client string) ([]tpl.MySetting, int, error) {
	data := []tpl.MySetting{}
	cursor := pg.TokenToID()

	sdc := m.RdDB.Select().
		From(
			goqu.T(schema.TableGroupSetting).As("t1"),
			goqu.T(schema.TableSetting).As("t2"),
			goqu.T(schema.TableModule).As("t3"),
			goqu.T(schema.TableProduct).As("t4")).
		Where(goqu.I("t1.group_id").Eq(groupID))

	sd := m.RdDB.Select(
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
			goqu.T(schema.TableGroupSetting).As("t1"),
			goqu.T(schema.TableSetting).As("t2"),
			goqu.T(schema.TableModule).As("t3"),
			goqu.T(schema.TableProduct).As("t4")).
		Where(
			goqu.I("t1.group_id").Eq(groupID),
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

	if channel != "" {
		sdc = sdc.Where(goqu.L("FIND_IN_SET(?, ?)", channel, goqu.I("t2.channels")))
		sd = sd.Where(goqu.L("FIND_IN_SET(?, ?)", channel, goqu.I("t2.channels")))
	}

	if client != "" {
		sdc = sdc.Where(goqu.L("FIND_IN_SET(?, ?)", client, goqu.I("t2.clients")))
		sd = sd.Where(goqu.L("FIND_IN_SET(?, ?)", client, goqu.I("t2.clients")))
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

// BatchAdd 批量添加群组
func (m *Group) BatchAdd(ctx context.Context, groups []tpl.GroupBody) error {
	if len(groups) == 0 {
		return nil
	}
	syncAt := time.Now().UTC().Unix()
	vals := make([][]interface{}, len(groups))
	for i, g := range groups {
		vals[i] = goqu.Vals{g.UID, g.Kind, syncAt, g.Desc}
	}

	sd := m.DB.Insert(schema.TableGroup).Cols("uid", "kind", "sync_at", "description").
		Vals(vals...).OnConflict(goqu.DoNothing())
	rowsAffected, err := service.DeResult(sd.Executor().ExecContext(ctx))
	if rowsAffected > 0 {
		util.Go(10*time.Second, func(gctx context.Context) {
			m.tryRefreshGroupsTotalSize(gctx)
		})
	}
	return err
}

// Update 更新指定群组
func (m *Group) Update(ctx context.Context, groupID int64, changed map[string]interface{}) (*schema.Group, error) {
	group := &schema.Group{}
	if _, err := m.updateByID(ctx, schema.TableGroup, groupID, goqu.Record(changed)); err != nil {
		return nil, err
	}
	if err := m.findOneByID(ctx, schema.TableGroup, groupID, group); err != nil {
		return nil, err
	}
	return group, nil
}

// Delete 删除指定群组
func (m *Group) Delete(ctx context.Context, groupID int64) error {
	_, err := m.deleteByCols(ctx, schema.TableGroupLabel, goqu.Ex{"group_id": groupID})
	if err == nil {
		_, err = m.deleteByCols(ctx, schema.TableGroupSetting, goqu.Ex{"group_id": groupID})
	}
	if err == nil {
		_, err = m.deleteByCols(ctx, schema.TableUserGroup, goqu.Ex{"group_id": groupID})
	}

	if err == nil {
		var rowsAffected int64
		rowsAffected, err = m.deleteByID(ctx, schema.TableGroup, groupID)
		if rowsAffected > 0 {
			util.Go(5*time.Second, func(gctx context.Context) {
				m.tryIncreaseStatisticStatus(gctx, schema.GroupsTotalSize, -1)
			})
		}
	}
	return err
}

// BatchAddMembers 批量添加群组成员，已存在则更新 sync_at
func (m *Group) BatchAddMembers(ctx context.Context, group *schema.Group, users []string) error {
	if len(users) == 0 {
		return nil
	}

	sd := m.DB.Insert(schema.TableUserGroup).Cols("user_id", "group_id", "sync_at").
		FromQuery(goqu.From(goqu.T(schema.TableUser).As("t1")).
			Select(goqu.I("t1.id"), goqu.V(group.ID), goqu.V(group.SyncAt)).
			Where(goqu.I("t1.uid").In(tpl.StrSliceToInterface(users)...))).
		OnConflict(goqu.DoUpdate("sync_at", goqu.C("sync_at").Set(goqu.V(group.SyncAt))))

	rowsAffected, err := service.DeResult(sd.Executor().ExecContext(ctx))
	if rowsAffected > 0 {
		util.Go(10*time.Second, func(gctx context.Context) {
			m.tryRefreshGroupStatus(gctx, group.ID)
		})
	}

	return err
}

// FindMembers 根据条件查找群组成员
func (m *Group) FindMembers(ctx context.Context, groupID int64, pg tpl.Pagination) ([]tpl.GroupMember, int, error) {
	data := []tpl.GroupMember{}
	cursor := pg.TokenToID()

	sdc := m.RdDB.Select().
		From(
			goqu.T(schema.TableUserGroup).As("t1"),
			goqu.T(schema.TableUser).As("t2")).
		Where(
			goqu.I("t1.group_id").Eq(groupID),
			goqu.I("t1.user_id").Eq(goqu.I("t2.id")))

	sd := m.RdDB.Select(
		goqu.I("t1.id"),
		goqu.I("t2.uid"),
		goqu.I("t1.created_at"),
		goqu.I("t1.sync_at")).
		From(
			goqu.T(schema.TableUserGroup).As("t1"),
			goqu.T(schema.TableUser).As("t2")).
		Where(
			goqu.I("t1.group_id").Eq(groupID),
			goqu.I("t1.id").Lte(cursor),
			goqu.I("t1.user_id").Eq(goqu.I("t2.id")))

	if pg.Q != "" {
		sdc = sdc.Where(goqu.I("t2.uid").ILike(pg.Q))
		sd = sd.Where(goqu.I("t2.uid").ILike(pg.Q))
	}

	sd = sd.Order(goqu.I("t1.id").Desc()).Limit(uint(pg.PageSize + 1))

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
		member := tpl.GroupMember{}
		if err := scanner.ScanStruct(&member); err != nil {
			return nil, 0, err
		}
		data = append(data, member)
	}

	if err := scanner.Err(); err != nil {
		return nil, 0, err
	}
	return data, int(total), nil
}

// FindIDsByUser 根据 userID 查找加入的 Group ID 数组
func (m *Group) FindIDsByUser(ctx context.Context, userID int64) ([]int64, error) {
	ids := make([]int64, 0)
	sd := m.RdDB.From(schema.TableUserGroup).Where(goqu.C("user_id").Eq(userID)).Limit(1000)
	if err := sd.PluckContext(ctx, &ids, "group_id"); err != nil {
		return nil, err
	}

	return ids, nil
}

// RemoveMembers 删除群组的成员
func (m *Group) RemoveMembers(ctx context.Context, groupID, userID int64, syncLt int64) error {

	sd := m.DB.Delete(schema.TableUserGroup).Where(goqu.C("group_id").Eq(groupID))
	if syncLt > 0 {
		sd = sd.Where(goqu.C("sync_at").Lt(syncLt))
	} else if userID > 0 {
		sd = sd.Where(goqu.C("user_id").Eq(userID))
	} else {
		return nil
	}

	res, err := service.DeResult(sd.Executor().ExecContext(ctx))
	if res > 0 {
		util.Go(10*time.Second, func(gctx context.Context) {
			m.tryRefreshGroupStatus(gctx, groupID)
		})
	}
	return err
}
