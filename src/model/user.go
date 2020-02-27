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
	"from `user_label` t1, `label` t2, `product` t3 " +
	"where t1.`user_id` = ? and t1.`label_id` = t2.`id` and t2.`product_id` = t3.`id` " +
	"order by t1.`created_at` desc limit 1000) " +
	"union all " + // 期望能基于 id distinct
	"(select t3.`id`, t2.`created_at`, t3.`name`, t4.`name` as `p`, t3.`channels`, t3.`clients` " +
	"from `user_group` t1, `group_label` t2, `label` t3, `product` t4 " +
	"where t1.`user_id` = ? and t1.`group_id` = t2.`group_id` and t2.`label_id` = t3.`id` and t3.`product_id` = t4.`id` " +
	"order by t2.`created_at` desc limit 1000)" +
	"order by `created_at` desc"

// RefreshLabels 更新 user 上的 labels 缓存，包括通过 group 关系获得的 labels
func (m *User) RefreshLabels(ctx context.Context, id int64, now int64) (string, error) {
	labels := ""

	err := m.DB.Transaction(func(tx *gorm.DB) error {
		user := &schema.User{ID: id}
		// 指定 id 的记录被锁住，如果表中无符合记录的数据则排他锁不生效
		if err := tx.Set("gorm:query_option", "FOR UPDATE").First(user).Error; err != nil {
			return err
		}

		if !conf.Config.IsCacheLabelExpired(now, user.ActiveAt) {
			// 已被其它请求更新
			labels = user.Labels
			return nil
		}

		data := make([]schema.UserCacheLabel, 0, 8)
		rows, err := tx.Raw(userCacheLabelsSQL, id, id).Rows()
		defer rows.Close()

		if err != nil {
			return err
		}

		var ignoreTime time.Time
		set := make(map[int64]struct{})
		for rows.Next() {
			var ignoreID int64
			label := schema.UserCacheLabel{}
			// ScanRows 扫描一行记录到 user
			if err := rows.Scan(&ignoreID, &ignoreTime, &label.Label, &label.Product, &label.Channels, &label.Clients); err != nil {
				return err
			}
			if _, ok := set[ignoreID]; ok {
				continue // 去重
			}

			set[ignoreID] = struct{}{}
			data = append(data, label)
		}
		user.PutLabels(data)
		user.ActiveAt = time.Now().UTC().Unix()
		labels = user.Labels

		return tx.Model(&schema.User{ID: id}).Updates(map[string]interface{}{
			"labels": user.Labels, "active_at": user.ActiveAt}).Error // 返回 nil 提交事务，否则回滚
	})

	return labels, err
}

const userLabelsSQL = "select t2.`id`, t2.`name`, t2.`desc`, t2.`channels`, t2.`clients`, t3.`name` as `product` " +
	"from `user_label` t1, `label` t2, `product` t3 " +
	"where t1.`user_id` = ? and t1.`label_id` = t2.`id` and t2.`product_id` = t3.id " +
	"order by t1.`created_at` desc " +
	"limit 1000"

// FindLables 根据用户 ID 返回其 labels 数据。TODO：支持更多筛选条件和分页
func (m *User) FindLables(ctx context.Context, userID int64, product string) ([]tpl.LabelInfo, error) {
	data := []tpl.LabelInfo{}
	rows, err := m.DB.Raw(userLabelsSQL, userID).Rows()
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

// BatchAdd 批量添加用户
func (m *User) BatchAdd(ctx context.Context, uids []string) error {
	if len(uids) == 0 {
		return nil
	}
	var buf bytes.Buffer
	fmt.Fprint(&buf, "insert ignore into `user` (`uid`) values")
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
