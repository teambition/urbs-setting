package model

import (
	"context"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/teambition/urbs-setting/src/schema"
	"github.com/teambition/urbs-setting/src/service"
	"github.com/teambition/urbs-setting/src/tpl"
)

// Group ...
type Group struct {
	DB *gorm.DB
}

// FindByUID 根据 uid 返回 user 数据
func (m *Group) FindByUID(ctx context.Context, uid string, selectStr string) (*schema.Group, error) {
	var err error
	group := &schema.Group{}
	db := m.DB.Where("uid = ?", uid)

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

// Find 根据条件查找 groups
func (m *Group) Find(ctx context.Context, kind string, pg tpl.Pagination) ([]schema.Group, error) {
	groups := make([]schema.Group, 0)
	pageToken := pg.TokenToID()
	db := m.DB.Where("`id` >= ?", pageToken)
	if kind != "" {
		db = m.DB.Where("`id` >= ? and kind = ?", pageToken, kind)
	}

	err := db.Order("`id`").Limit(pg.PageSize + 1).Find(&groups).Error
	return groups, err
}

// Count 计算 group 总数
func (m *Group) Count(ctx context.Context, kind string) (int, error) {

	count := 0
	db := m.DB.Model(&schema.Group{})
	if kind != "" {
		db = db.Where("kind = ?", kind)
	}
	err := db.Count(&count).Error
	return count, err
}

const groupLabelsSQL = "select t2.`id`, t2.`created_at`, t2.`updated_at`, t2.`offline_at`, t2.`name`, " +
	"t2.`description`, t2.`status`, t2.`channels`, t2.`clients`, t3.`name` as `product` " +
	"from `group_label` t1, `urbs_label` t2, `urbs_product` t3 " +
	"where t1.`group_id` = ? and t1.`id` >= ? and t1.`label_id` = t2.`id` and t2.`product_id` = t3.id " +
	"order by t1.`id` asc " +
	"limit ?"

// FindLables 根据群组 ID 返回其 labels 数据。TODO：支持更多筛选条件和分页
func (m *Group) FindLables(ctx context.Context, groupID int64, pg tpl.Pagination) ([]tpl.LabelInfo, error) {
	data := []tpl.LabelInfo{}
	pageToken := pg.TokenToID()

	rows, err := m.DB.Raw(groupLabelsSQL, groupID, pageToken, pg.PageSize+1).Rows()
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var clients string
		var channels string
		labelInfo := tpl.LabelInfo{}
		if err := rows.Scan(&labelInfo.ID, &labelInfo.CreatedAt, &labelInfo.UpdatedAt, &labelInfo.OfflineAt,
			&labelInfo.Name, &labelInfo.Desc, &labelInfo.Status, &channels, &clients, &labelInfo.Product); err != nil {
			return nil, err
		}
		labelInfo.Channels = tpl.StringToSlice(channels)
		labelInfo.Clients = tpl.StringToSlice(clients)
		labelInfo.HID = service.IDToHID(labelInfo.ID, "label")
		data = append(data, labelInfo)
	}

	return data, nil
}

// CountLabels 计算 group labels 总数
func (m *Group) CountLabels(ctx context.Context, groupID int64) (int, error) {

	count := 0
	err := m.DB.Model(&schema.GroupLabel{}).Where("group_id = ?", groupID).Count(&count).Error
	return count, err
}

const groupSettingsSQL = "select t1.`created_at`, t1.`updated_at`, t1.`value`, t1.`last_value`, " +
	"t2.`id`, t2.`name`, t3.`name` as `module` " +
	"from `group_setting` t1, `urbs_setting` t2, `urbs_module` t3 " +
	"where t1.`group_id` = ? and t1.`updated_at` >= ? and t1.`setting_id` = t2.`id` and t2.`module_id` in ( ? ) and t2.`module_id` = t3.`id` " +
	"order by t1.`updated_at` asc " +
	"limit ?"

// FindSettings 根据 Group ID, updateGt, productName 返回其 settings 数据。
func (m *Group) FindSettings(ctx context.Context, groupID int64, moduleIDs []int64, pg tpl.Pagination) ([]tpl.MySetting, error) {
	data := []tpl.MySetting{}
	updatedAt := pg.TokenToTime()

	rows, err := m.DB.Raw(groupSettingsSQL, groupID, updatedAt, moduleIDs, pg.PageSize+1).Rows()
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		mySetting := tpl.MySetting{}
		if err := rows.Scan(&mySetting.CreatedAt, &mySetting.UpdatedAt, &mySetting.Value, &mySetting.LastValue,
			&mySetting.ID, &mySetting.Name, &mySetting.Module); err != nil {
			return nil, err
		}
		mySetting.HID = service.IDToHID(mySetting.ID, "setting")
		data = append(data, mySetting)
	}

	return data, nil
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
	return nil
}

const batchAddGroupMemberSQL = "insert ignore into `user_group` (`user_id`, `group_id`, `sync_at`) " +
	"select `urbs_user`.id, ?, ? from `urbs_user` where `urbs_user`.uid in ( ? ) " +
	"on duplicate key update `sync_at` = values(`sync_at`)"

// BatchAddMembers 批量添加群组成员，已存在则更新 sync_at
func (m *Group) BatchAddMembers(ctx context.Context, group *schema.Group, users []string) error {
	if len(users) == 0 {
		return nil
	}

	return m.DB.Exec(batchAddGroupMemberSQL, group.ID, group.SyncAt, users).Error
}

// CountMembers 计算成员总数
func (m *Group) CountMembers(ctx context.Context, groupID int64) (int, error) {

	count := 0
	err := m.DB.Model(&schema.UserGroup{}).Where("group_id = ?", groupID).Count(&count).Error
	return count, err
}

const groupMembersSQL = "select t1.`id`, t2.`uid`, t1.`created_at`, t1.`sync_at` " +
	"from `user_group` t1, `urbs_user` t2 " +
	"where t1.`group_id` = ? and t1.`id` >= ? and t1.`user_id` = t2.`id` " +
	"order by t1.`id` asc " +
	"limit ?"

// FindMembers 根据条件查找群组成员
func (m *Group) FindMembers(ctx context.Context, groupID int64, pg tpl.Pagination) ([]tpl.GroupMember, error) {
	data := []tpl.GroupMember{}
	pageToken := pg.TokenToID()

	rows, err := m.DB.Raw(groupMembersSQL, groupID, pageToken, pg.PageSize+1).Rows()
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		member := tpl.GroupMember{}
		if err := rows.Scan(&member.ID, &member.User, &member.CreatedAt, &member.SyncAt); err != nil {
			return nil, err
		}
		data = append(data, member)
	}

	return data, nil
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
		err = m.DB.Where("group_id = ? and sync_at < ?", groupID, syncLt).Delete(&schema.UserGroup{}).Error
	}
	if err == nil && userID > 0 {
		err = m.DB.Where("user_id = ? and group_id = ?", userID, groupID).Delete(&schema.UserGroup{}).Error
	}
	return err
}

// RemoveLable 删除群组的 label
func (m *Group) RemoveLable(ctx context.Context, groupID, lableID int64) error {
	return m.DB.Where("group_id = ? and label_id = ?", groupID, lableID).Delete(&schema.GroupLabel{}).Error
}

// RemoveSetting 删除群组的 setting
func (m *Group) RemoveSetting(ctx context.Context, groupID, settingID int64) error {
	return m.DB.Where("group_id = ? and setting_id = ?", groupID, settingID).Delete(&schema.GroupSetting{}).Error
}
