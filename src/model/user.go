package model

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/teambition/urbs-setting/src/conf"
	"github.com/teambition/urbs-setting/src/schema"
	"github.com/teambition/urbs-setting/src/service"
	"github.com/teambition/urbs-setting/src/tpl"
)

// User ...
type User struct {
	DB *gorm.DB
}

// FindByUID 根据 uid 返回 user 数据
func (m *User) FindByUID(ctx context.Context, uid string, selectStr string) (*schema.User, error) {
	var err error
	user := &schema.User{}
	db := m.DB.Where("uid = ?", uid)

	if selectStr == "" {
		err = db.First(user).Error
	} else {
		err = db.Select(selectStr).First(user).Error
	}

	if err == nil {
		return user, nil
	}

	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}
	return nil, err
}

const userCacheLabelsSQL = "(select t2.`id`, t1.`created_at`, t2.`name`, t3.`name` as `p`, t2.`channels`, t2.`clients` " +
	"from `user_label` t1, `urbs_label` t2, `urbs_product` t3 " +
	"where t1.`user_id` = ? and t1.`label_id` = t2.`id` and t2.`product_id` = t3.`id` " +
	"order by t1.`id` desc limit 200) " +
	"union all " + // 期望能基于 id distinct
	"(select t3.`id`, t2.`created_at`, t3.`name`, t4.`name` as `p`, t3.`channels`, t3.`clients` " +
	"from `user_group` t1, `group_label` t2, `urbs_label` t3, `urbs_product` t4 " +
	"where t1.`user_id` = ? and t1.`group_id` = t2.`group_id` and t2.`label_id` = t3.`id` and t3.`product_id` = t4.`id` " +
	"order by t2.`id` desc limit 200)" +
	"order by `created_at` desc"

// RefreshLabels 更新 user 上的 labels 缓存，包括通过 group 关系获得的 labels
func (m *User) RefreshLabels(ctx context.Context, id int64, now int64) (*schema.User, error) {
	user := &schema.User{ID: id}
	err := m.DB.Transaction(func(tx *gorm.DB) error {
		// 指定 id 的记录被锁住，如果表中无符合记录的数据则排他锁不生效
		if err := tx.Set("gorm:query_option", "FOR UPDATE").First(user).Error; err != nil {
			return err
		}

		if !conf.Config.IsCacheLabelExpired(now, user.ActiveAt) {
			// 已被其它请求更新
			return nil
		}

		data := make(schema.UserCacheLabels)
		rows, err := tx.Raw(userCacheLabelsSQL, id, id).Rows()
		defer rows.Close()

		if err != nil {
			return err
		}

		var ignoreTime time.Time
		set := make(map[int64]struct{})
		for rows.Next() {
			var ignoreID int64
			var product string
			var clients string
			var channels string
			label := schema.UserCacheLabel{}
			// ScanRows 扫描一行记录到 user
			if err := rows.Scan(&ignoreID, &ignoreTime, &label.Label, &product, &channels, &clients); err != nil {
				return err
			}
			if _, ok := set[ignoreID]; ok {
				continue // 去重
			}

			label.Channels = tpl.StringToSlice(channels)
			label.Clients = tpl.StringToSlice(clients)
			set[ignoreID] = struct{}{}
			arr, ok := data[product]
			if !ok {
				arr = make([]schema.UserCacheLabel, 0)
			}
			data[product] = append(arr, label)
		}
		_ = user.PutLabels(data)
		user.ActiveAt = time.Now().UTC().Unix()

		return tx.Model(&schema.User{ID: id}).Updates(map[string]interface{}{
			"labels": user.Labels, "active_at": user.ActiveAt}).Error // 返回 nil 提交事务，否则回滚
	})

	return user, err
}

const userLabelsSQL = "select t2.`id`, t2.`created_at`, t2.`updated_at`, t2.`offline_at`, t2.`name`, " +
	"t2.`description`, t2.`status`, t2.`channels`, t2.`clients`, t3.`name` as `product` " +
	"from `user_label` t1, `urbs_label` t2, `urbs_product` t3 " +
	"where t1.`user_id` = ? and t1.`id` >= ? and t1.`label_id` = t2.`id` and t2.`product_id` = t3.id " +
	"order by t1.`id` asc " +
	"limit ?"

// FindLables 根据用户 ID 返回其 labels 数据。
func (m *User) FindLables(ctx context.Context, userID int64, pg tpl.Pagination) ([]tpl.LabelInfo, error) {
	data := []tpl.LabelInfo{}
	pageToken := pg.TokenToID()

	rows, err := m.DB.Raw(userLabelsSQL, userID, pageToken, pg.PageSize+1).Rows()
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		labelInfo := tpl.LabelInfo{}
		var clients string
		var channels string
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

// CountLabels 计算 user labels 总数
func (m *User) CountLabels(ctx context.Context, userID int64) (int, error) {

	count := 0
	err := m.DB.Model(&schema.UserLabel{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

// BatchAdd 批量添加用户
// uids 经过了 `^[0-9A-Za-z._-]{3,63}$` 正则验证
func (m *User) BatchAdd(ctx context.Context, uids []string) error {
	if len(uids) == 0 {
		return nil
	}
	var buf bytes.Buffer
	fmt.Fprint(&buf, "insert ignore into `urbs_user` (`uid`) values")
	for _, uid := range uids {
		fmt.Fprintf(&buf, " ('%s'),", uid)
	}
	b := buf.Bytes()
	b[len(b)-1] = ';'
	return m.DB.Exec(string(b)).Error
}

// RemoveLable 删除用户的 label
func (m *User) RemoveLable(ctx context.Context, userID, lableID int64) error {
	return m.DB.Where("user_id = ? and label_id = ?", userID, lableID).Delete(&schema.UserLabel{}).Error
}

// RemoveSetting 删除用户的 setting
func (m *User) RemoveSetting(ctx context.Context, userID, settingID int64) error {
	return m.DB.Where("user_id = ? and setting_id = ?", userID, settingID).Delete(&schema.UserSetting{}).Error
}
