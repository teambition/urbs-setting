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
func (m *Group) Find(ctx context.Context) ([]schema.Group, error) {
	groups := make([]schema.Group, 0)
	err := m.DB.Order("`created_at`").Limit(1000).Find(&groups).Error
	return groups, err
}

const groupLabelsSQL = "select t2.`id`, t2.`name`, t2.`desc`, t2.`channels`, t2.`clients`, t3.`name` as `product` " +
	"from `group_label` t1, `label` t2, `product` t3 " +
	"where t1.`group_id` = ? and t1.`label_id` = t2.`id` and t2.`product_id` = t3.id " +
	"order by t1.`created_at` desc " +
	"limit 1000"

// FindLables 根据群组 ID 返回其 labels 数据。TODO：支持更多筛选条件和分页
func (m *Group) FindLables(ctx context.Context, groupID int64, product string) ([]tpl.LabelInfo, error) {
	data := []tpl.LabelInfo{}
	rows, err := m.DB.Raw(groupLabelsSQL, groupID).Rows()
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		label := schema.Label{}
		labelInfo := tpl.LabelInfo{}
		if err := rows.Scan(&label.ID, &labelInfo.Name, &labelInfo.Desc, &labelInfo.Channels, &labelInfo.Clients, &labelInfo.Product); err != nil {
			return nil, err
		}
		labelInfo.HID = service.HIDer.HID(label)
		data = append(data, labelInfo)
	}

	return data, nil
}

// BatchAdd 批量添加群组
func (m *Group) BatchAdd(ctx context.Context, groups []tpl.GroupBody) error {
	if len(groups) == 0 {
		return nil
	}

	syncAt := time.Now().UTC().Unix()
	stmt, err := m.DB.DB().Prepare("insert ignore into `group` (`uid`, `sync_at`, `desc`) values (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, group := range groups {
		if _, err := stmt.Exec(group.UID, syncAt, group.Desc); err != nil {
			return err
		}
	}
	return nil
}

const batchAddGroupMemberSQL = "insert ignore into `user_group` (`user_id`, `group_id`, `sync_at`) " +
	"select `user`.id, ?, ? from `user` where `user`.uid in ( ? ) " +
	"on duplicate key update `sync_at` = values(`sync_at`)"

// BatchAddMembers 批量添加群组成员，已存在则更新 sync_at
func (m *Group) BatchAddMembers(ctx context.Context, group *schema.Group, users []string) error {
	if len(users) == 0 {
		return nil
	}

	return m.DB.Exec(batchAddGroupMemberSQL, group.ID, group.SyncAt, users).Error
}

const groupMembersSQL = "select t2.`uid`, t1.`created_at`, t1.`sync_at` " +
	"from `user_group` t1, `user` t2 " +
	"where t1.`group_id` = ? and t1.`user_id` = t2.`id` " +
	"order by t1.`sync_at` desc " +
	"limit 10000"

// FindMembers 根据条件查找群组成员
func (m *Group) FindMembers(ctx context.Context, groupID int64) ([]tpl.GroupMember, error) {
	data := []tpl.GroupMember{}
	rows, err := m.DB.Raw(groupMembersSQL, groupID).Rows()
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		member := tpl.GroupMember{}
		if err := rows.Scan(&member.User, &member.CreatedAt, &member.SyncAt); err != nil {
			return nil, err
		}
		data = append(data, member)
	}

	return data, nil
}

// RemoveMembers 删除群组的成员
func (m *Group) RemoveMembers(ctx context.Context, groupID, userID int64, syncLt int64) error {
	var err error
	if syncLt > 0 {
		err = m.DB.Where("group_id = ? and sync_at < ?", groupID, syncLt).Delete(&schema.UserGroup{}).Error
	}
	if err == nil && userID > 0 {
		err = m.DB.Where("group_id = ? and user_id = ?", groupID, userID).Delete(&schema.UserGroup{}).Error
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