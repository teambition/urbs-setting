package model

import (
	"context"
	"database/sql"
	"time"

	"github.com/jinzhu/gorm"
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
	var err error
	group := &schema.Group{}
	db := m.DB.Where("`uid` = ?", uid)

	if selectStr == "" {
		err = db.First(group).Error
	} else {
		err = db.Select(selectStr).First(group).Error
	}

	if err == nil {
		return group, nil
	}

	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}
	return nil, err
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
	group, err := m.FindByUID(ctx, uid, "`id`")
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
	cursor := pg.TokenToID(true)
	db := m.DB.Where("`id` <= ?", cursor)
	if pg.Q != "" {
		db = m.DB.Where("`id` <= ? and `uid` like ?", cursor, pg.Q)
	}
	if kind != "" {
		db = m.DB.Where("`id` <= ? and `kind` = ?", cursor, kind)
		if pg.Q != "" {
			db = m.DB.Where("`id` <= ? and `kind` = ? and `uid` like ?", cursor, kind, pg.Q)
		}
	}

	total := 0
	err := db.Model(&schema.Group{}).Count(&total).Error
	if err == nil {
		err = db.Order("`id` desc").Limit(pg.PageSize + 1).Find(&groups).Error
	}
	if err != nil {
		return nil, 0, err
	}
	return groups, total, nil
}

const listGroupLabelsSQL = "select t2.`id`, t2.`created_at`, t2.`updated_at`, t2.`offline_at`, t2.`name`, " +
	"t2.`description`, t2.`status`, t2.`channels`, t2.`clients`, t3.`name` as `product` " +
	"from `group_label` t1, `urbs_label` t2, `urbs_product` t3 " +
	"where t1.`group_id` = ? and t1.`id` <= ? and t1.`label_id` = t2.`id` and t2.`product_id` = t3.`id` " +
	"order by t1.`id` desc " +
	"limit ?"

const countGroupLabelsSQL = "select count(t2.`id`) " +
	"from `group_label` t1, `urbs_label` t2  " +
	"where t1.`group_id` = ? and t1.`label_id` = t2.`id`"

const searchGroupLabelsSQL = "select t2.`id`, t2.`created_at`, t2.`updated_at`, t2.`offline_at`, t2.`name`, " +
	"t2.`description`, t2.`status`, t2.`channels`, t2.`clients`, t3.`name` as `product` " +
	"from `group_label` t1, `urbs_label` t2, `urbs_product` t3 " +
	"where t1.`group_id` = ? and t1.`id` <= ? and t1.`label_id` = t2.`id` and t2.`name` like ? and t2.`product_id` = t3.`id` " +
	"order by t1.`id` desc " +
	"limit ?"

const countSearchGroupLabelsSQL = "select count(t2.`id`) " +
	"from `group_label` t1, `urbs_label` t2 " +
	"where t1.`group_id` = ? and t1.`label_id` = t2.`id` and t2.`name` like ?"

// FindLabels 根据群组 ID 返回其 labels 数据。TODO：支持更多筛选条件和分页
func (m *Group) FindLabels(ctx context.Context, groupID int64, pg tpl.Pagination) ([]tpl.LabelInfo, int, error) {
	data := make([]tpl.LabelInfo, 0)
	cursor := pg.TokenToID(true)
	total := 0

	if pg.Q == "" {
		if err := m.DB.Raw(countGroupLabelsSQL, groupID).Row().Scan(&total); err != nil && err != sql.ErrNoRows {
			return nil, 0, err
		}
	} else {
		if err := m.DB.Raw(countSearchGroupLabelsSQL, groupID, pg.Q).Row().Scan(&total); err != nil && err != sql.ErrNoRows {
			return nil, 0, err
		}
	}

	var err error
	var rows *sql.Rows
	if pg.Q == "" {
		rows, err = m.DB.Raw(listGroupLabelsSQL, groupID, cursor, pg.PageSize+1).Rows()
	} else {
		rows, err = m.DB.Raw(searchGroupLabelsSQL, groupID, cursor, pg.Q, pg.PageSize+1).Rows()
	}
	defer rows.Close()

	if err != nil {
		return nil, 0, err
	}

	for rows.Next() {
		var clients string
		var channels string
		labelInfo := tpl.LabelInfo{}
		if err := rows.Scan(&labelInfo.ID, &labelInfo.CreatedAt, &labelInfo.UpdatedAt, &labelInfo.OfflineAt,
			&labelInfo.Name, &labelInfo.Desc, &labelInfo.Status, &channels, &clients, &labelInfo.Product); err != nil {
			return nil, 0, err
		}
		labelInfo.Channels = tpl.StringToSlice(channels)
		labelInfo.Clients = tpl.StringToSlice(clients)
		labelInfo.HID = service.IDToHID(labelInfo.ID, "label")
		data = append(data, labelInfo)
	}

	return data, total, nil
}

const listGroupSettingsSQL = "select t1.`created_at`, t1.`updated_at`, t1.`value`, t1.`last_value`, " +
	"t2.`id`, t2.`name`, t3.`name` as `module`, t4.`name` as `product` " +
	"from `group_setting` t1, `urbs_setting` t2, `urbs_module` t3, `urbs_product` t4 " +
	"where t1.`group_id` = ? and t1.`id` <= ? and t1.`setting_id` = t2.`id` and t2.`module_id` = t3.`id` and t3.`product_id` = t4.`id` " +
	"order by t1.`id` desc " +
	"limit ?"

const countGroupSettingsSQL = "select count(t2.`id`) " +
	"from `group_setting` t1, `urbs_setting` t2 " +
	"where t1.`group_id` = ? and t1.`setting_id` = t2.`id`"

const searchGroupSettingsSQL = "select t1.`created_at`, t1.`updated_at`, t1.`value`, t1.`last_value`, " +
	"t2.`id`, t2.`name`, t3.`name` as `module`, t4.`name` as `product` " +
	"from `group_setting` t1, `urbs_setting` t2, `urbs_module` t3, `urbs_product` t4 " +
	"where t1.`group_id` = ? and t1.`id` <= ? and t1.`setting_id` = t2.`id` and t2.`name` like ? and t2.`module_id` = t3.`id` and t3.`product_id` = t4.`id` " +
	"order by t1.`id` desc " +
	"limit ?"

const countSearchGroupSettingsSQL = "select count(t2.`id`) " +
	"from `group_setting` t1, `urbs_setting` t2 " +
	"where t1.`group_id` = ? and t1.`setting_id` = t2.`id` and t2.`name` like ?"

// FindSettings 根据 Group ID, updateGt, productName 返回其 settings 数据。
func (m *Group) FindSettings(ctx context.Context, groupID int64, pg tpl.Pagination) ([]tpl.MySetting, int, error) {
	data := []tpl.MySetting{}
	cursor := pg.TokenToID(true)
	total := 0

	if pg.Q == "" {
		if err := m.DB.Raw(countGroupSettingsSQL, groupID).Row().Scan(&total); err != nil && err != sql.ErrNoRows {
			return nil, 0, err
		}
	} else {
		if err := m.DB.Raw(countSearchGroupSettingsSQL, groupID, pg.Q).Row().Scan(&total); err != nil && err != sql.ErrNoRows {
			return nil, 0, err
		}
	}

	var err error
	var rows *sql.Rows
	if pg.Q == "" {
		rows, err = m.DB.Raw(listGroupSettingsSQL, groupID, cursor, pg.PageSize+1).Rows()
	} else {
		rows, err = m.DB.Raw(searchGroupSettingsSQL, groupID, cursor, pg.Q, pg.PageSize+1).Rows()
	}

	defer rows.Close()

	if err != nil {
		return nil, 0, err
	}

	for rows.Next() {
		mySetting := tpl.MySetting{}
		if err := rows.Scan(&mySetting.CreatedAt, &mySetting.UpdatedAt, &mySetting.Value, &mySetting.LastValue,
			&mySetting.ID, &mySetting.Name, &mySetting.Module, &mySetting.Product); err != nil {
			return nil, 0, err
		}
		mySetting.HID = service.IDToHID(mySetting.ID, "setting")
		data = append(data, mySetting)
	}

	return data, total, nil
}

// BatchAdd 批量添加群组
func (m *Group) BatchAdd(ctx context.Context, groups []tpl.GroupBody) error {
	if len(groups) == 0 {
		return nil
	}

	syncAt := time.Now().UTC().Unix()
	stmt, err := m.DB.DB().Prepare("insert ignore into `urbs_group` (`uid`, `kind`, `sync_at`, `description`) values (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, group := range groups {
		if _, err := stmt.Exec(group.UID, group.Kind, syncAt, group.Desc); err != nil {
			return err
		}
	}
	go m.refreshGroupsTotalSize(ctx)
	return nil
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
			go m.increaseStatisticStatus(ctx, schema.GroupsTotalSize, -1)
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
	go m.refreshGroupStatus(ctx, group.ID)
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
	cursor := pg.TokenToID(true)
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
	go m.refreshGroupStatus(ctx, groupID)
	return err
}
