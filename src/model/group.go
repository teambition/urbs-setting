package model

import (
	"context"
	"database/sql"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/teambition/gear"
	"github.com/teambition/urbs-setting/src/schema"
	"github.com/teambition/urbs-setting/src/service"
	"github.com/teambition/urbs-setting/src/tpl"
)

// Group ...
type Group struct {
	*Model
}

// FindByUID 根据 uid 返回 user 数据
func (m *Group) FindByUID(ctx context.Context, uid string, selectStr string) (*schema.Group, error) {
	group := &schema.Group{}
	sd := m.GDB.From("urbs_group").Where(goqu.C("uid").Eq(uid)).Order(goqu.C("id").Asc()).Limit(1)

	if selectStr != "" {
		sd = sd.Select(goqu.L(selectStr))
	}

	ok, err := sd.Executor().ScanStructContext(ctx, group)
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
	group, err := m.FindByUID(ctx, uid, "`id`, `uid`")
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
	sdc := m.GDB.From("urbs_group")
	sd := m.GDB.From("urbs_group").Where(goqu.C("id").Lte(cursor))
	if kind != "" {
		sdc = sdc.Where(goqu.C("kind").Eq(kind))
		sd = sd.Where(goqu.C("kind").Eq(kind))
	}
	if pg.Q != "" {
		sdc = sdc.Where(goqu.C("uid").Like(pg.Q))
		sd = sd.Where(goqu.C("uid").Like(pg.Q))
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

	sdc := m.GDB.Select().
		From(
			goqu.T("group_label").As("t1"),
			goqu.T("urbs_label").As("t2"),
			goqu.T("urbs_product").As("t3")).
		Where(
			goqu.I("t1.group_id").Eq(groupID),
			goqu.I("t1.label_id").Eq(goqu.I("t2.id")))

	sd := m.GDB.Select(
		goqu.I("t1.rls"),
		goqu.I("t1.created_at").As("assigned_at"),
		goqu.I("t2.id"),
		goqu.I("t2.name"),
		goqu.I("t2.description"),
		goqu.I("t3.name").As("product")).
		From(
			goqu.T("group_label").As("t1"),
			goqu.T("urbs_label").As("t2"),
			goqu.T("urbs_product").As("t3")).
		Where(
			goqu.I("t1.group_id").Eq(groupID),
			goqu.I("t1.id").Lte(cursor),
			goqu.I("t1.label_id").Eq(goqu.I("t2.id")))

	if pg.Q != "" {
		sdc = sdc.Where(goqu.I("t2.name").Like(pg.Q))
		sd = sd.Where(goqu.I("t2.name").Like(pg.Q))
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
func (m *Group) FindSettings(ctx context.Context, groupID, productID, moduleID, settingID int64, pg tpl.Pagination) ([]tpl.MySetting, int, error) {
	data := []tpl.MySetting{}
	cursor := pg.TokenToID()

	sdc := m.GDB.Select().
		From(
			goqu.T("group_setting").As("t1"),
			goqu.T("urbs_setting").As("t2"),
			goqu.T("urbs_module").As("t3"),
			goqu.T("urbs_product").As("t4")).
		Where(goqu.I("t1.group_id").Eq(groupID))

	sd := m.GDB.Select(
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
			goqu.T("group_setting").As("t1"),
			goqu.T("urbs_setting").As("t2"),
			goqu.T("urbs_module").As("t3"),
			goqu.T("urbs_product").As("t4")).
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

	if pg.Q != "" {
		sdc = sdc.Where(goqu.I("t2.name").Like(pg.Q))
		sd = sd.Where(goqu.I("t2.name").Like(pg.Q))
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

	insertDataset := m.GDB.Insert("urbs_group").Cols("uid", "kind", "sync_at", "description").
		Vals(vals...).OnConflict(goqu.DoNothing())
	res, err := insertDataset.Executor().ExecContext(ctx)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if rowsAffected > 0 {
		go m.tryRefreshGroupsTotalSize(ctx)
	}
	return err
}

// Update 更新指定群组
func (m *Group) Update(ctx context.Context, groupID int64, changed map[string]interface{}) (*schema.Group, error) {
	group := &schema.Group{ID: groupID}
	if len(changed) > 0 {
		if err := m.DB.Model(group).UpdateColumns(changed).Error; err != nil {
			return nil, err
		}
	}

	if err := m.DB.First(group).Error; err != nil {
		return nil, err
	}
	return group, nil
}

// Delete 更新指定群组
func (m *Group) Delete(ctx context.Context, groupID int64) error {
	err := m.DB.Where("`group_id` = ?", groupID).Delete(&schema.GroupLabel{}).Error
	if err == nil {
		err = m.DB.Where("`group_id` = ?", groupID).Delete(&schema.GroupSetting{}).Error
	}
	if err == nil {
		err = m.DB.Where("`group_id` = ?", groupID).Delete(&schema.UserGroup{}).Error
	}
	if err == nil {
		res := m.DB.Where("`id` = ?", groupID).Delete(&schema.Group{})
		if res.RowsAffected > 0 {
			go m.tryIncreaseStatisticStatus(ctx, schema.GroupsTotalSize, -1)
		}
		err = res.Error
	}
	return err
}

const batchAddGroupMemberSQL = "insert ignore into `user_group` (`user_id`, `group_id`, `sync_at`) " +
	"select `urbs_user`.id, ?, ? from `urbs_user` where `urbs_user`.uid in ( ? ) " +
	"on duplicate key update `sync_at` = ?"

// BatchAddMembers 批量添加群组成员，已存在则更新 sync_at
func (m *Group) BatchAddMembers(ctx context.Context, group *schema.Group, users []string) error {
	if len(users) == 0 {
		return nil
	}

	err := m.DB.Exec(batchAddGroupMemberSQL, group.ID, group.SyncAt, users, group.SyncAt).Error
	go m.tryRefreshGroupStatus(ctx, group.ID)
	return err
}

const listGroupMembersSQL = "select t1.`id`, t2.`uid`, t1.`created_at`, t1.`sync_at` " +
	"from `user_group` t1, `urbs_user` t2 " +
	"where t1.`group_id` = ? and t1.`id` <= ? and t1.`user_id` = t2.`id` " +
	"order by t1.`id` desc " +
	"limit ?"

const countGroupMembersSQL = "select count(t2.`id`) " +
	"from `user_group` t1, `urbs_user` t2 " +
	"where t1.`group_id` = ? and t1.`user_id` = t2.`id`"

const searchGroupMembersSQL = "select t1.`id`, t1.`uid`, t1.`created_at`, t1.`sync_at` " +
	"from `user_group` t1, `urbs_user` t2 " +
	"where t1.`group_id` = ? and t1.`id` <= ? and t1.`user_id` = t2.`id` and t2.`uid` like ? " +
	"order by t1.`id` desc " +
	"limit ?"

const countSearchGroupMembersSQL = "select count(t2.`id`) " +
	"from `user_group` t1, `urbs_user` t2 " +
	"where t1.`group_id` = ? and t1.`user_id` = t2.`id` and t2.`uid` like ?"

// FindMembers 根据条件查找群组成员
func (m *Group) FindMembers(ctx context.Context, groupID int64, pg tpl.Pagination) ([]tpl.GroupMember, int, error) {
	data := []tpl.GroupMember{}
	cursor := pg.TokenToID()
	total := 0

	if pg.Q == "" {
		if err := m.DB.Raw(countGroupMembersSQL, groupID).Row().Scan(&total); err != nil && err != sql.ErrNoRows {
			return nil, 0, err
		}
	} else {
		if err := m.DB.Raw(countSearchGroupMembersSQL, groupID, pg.Q).Row().Scan(&total); err != nil && err != sql.ErrNoRows {
			return nil, 0, err
		}
	}

	var err error
	var rows *sql.Rows
	if pg.Q == "" {
		rows, err = m.DB.Raw(listGroupMembersSQL, groupID, cursor, pg.PageSize+1).Rows()
	} else {
		rows, err = m.DB.Raw(searchGroupMembersSQL, groupID, cursor, pg.Q, pg.PageSize+1).Rows()
	}

	defer rows.Close()

	if err != nil {
		return nil, 0, err
	}

	for rows.Next() {
		member := tpl.GroupMember{}
		if err := rows.Scan(&member.ID, &member.User, &member.CreatedAt, &member.SyncAt); err != nil {
			return nil, 0, err
		}
		data = append(data, member)
	}

	return data, total, nil
}

// FindIDsByUserID 根据 userID 查找加入的 Group ID 数组
func (m *Group) FindIDsByUserID(ctx context.Context, userID int64) ([]int64, error) {
	userGroups := make([]schema.UserGroup, 0)
	err := m.DB.Where("`user_id` = ?", userID).Select("`group_id`").
		Limit(1000).Find(&userGroups).Error
	ids := make([]int64, len(userGroups))
	if err == nil {
		for i, u := range userGroups {
			ids[i] = u.GroupID
		}
	}
	return ids, err
}

// RemoveMembers 删除群组的成员
func (m *Group) RemoveMembers(ctx context.Context, groupID, userID int64, syncLt int64) error {
	var err error
	if syncLt > 0 {
		err = m.DB.Where("`group_id` = ? and `sync_at` < ?", groupID, syncLt).Delete(&schema.UserGroup{}).Error
	}
	if err == nil && userID > 0 {
		err = m.DB.Where("`user_id` = ? and `group_id` = ?", userID, groupID).Delete(&schema.UserGroup{}).Error
	}
	go m.tryRefreshGroupStatus(ctx, groupID)
	return err
}
