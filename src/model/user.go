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
func (m *User) RefreshLabels(ctx context.Context, id int64, now int64, force bool) (*schema.User, error) {
	user := &schema.User{ID: id}
	err := m.DB.Transaction(func(tx *gorm.DB) error {
		// 指定 id 的记录被锁住，如果表中无符合记录的数据则排他锁不生效
		if err := tx.Set("gorm:query_option", "FOR UPDATE").First(user).Error; err != nil {
			return err
		}

		if !force && !conf.Config.IsCacheLabelExpired(now, user.ActiveAt) {
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
			set[ignoreID] = struct{}{}

			label.Channels = tpl.StringToSlice(channels)
			label.Clients = tpl.StringToSlice(clients)
			arr, ok := data[product]
			if !ok {
				arr = make([]schema.UserCacheLabel, 0)
			}
			data[product] = append(arr, label)
		}
		_ = user.PutLabels(data)
		user.ActiveAt = time.Now().UTC().Unix()

		return tx.Model(&schema.User{ID: id}).UpdateColumns(map[string]interface{}{
			"labels": user.Labels, "active_at": user.ActiveAt}).Error // 返回 nil 提交事务，否则回滚
	})

	return user, err
}

const userSettingsWithGroupSQL = "(select t1.`created_at`, t1.`updated_at`, t1.`value`, t1.`last_value`, " +
	"t2.`id`, t2.`name`, t3.`name` as `module`, t2.`channels`, t2.`clients` " +
	"from `user_setting` t1, `urbs_setting` t2, `urbs_module` t3 " +
	"where t1.`user_id` = ? and t1.`updated_at` <= ? and t1.`setting_id` = t2.`id` and t2.`module_id` in ( ? ) and t2.`module_id` = t3.`id` " +
	"order by t1.`updated_at` desc limit ? ) " +
	"union all " +
	"(select t1.`created_at`, t1.`updated_at`, t1.`value`, t1.`last_value`, " +
	"t2.`id`, t2.`name`, t3.`name` as `module`, t2.`channels`, t2.`clients` " +
	"from `group_setting` t1, `urbs_setting` t2, `urbs_module` t3 " +
	"where t1.`group_id` in ( ? ) and t1.`updated_at` <= ? and t1.`setting_id` = t2.`id` and t2.`module_id` in ( ? ) and t2.`module_id` = t3.`id` " +
	"order by t1.`updated_at` desc limit ? ) " +
	"order by `updated_at` desc"

// FindSettingsUnionAll 根据用户 ID, updateGt, productName 返回其 settings 数据。
func (m *User) FindSettingsUnionAll(ctx context.Context, userID int64, groupIDs []int64, moduleIDs []int64, pg tpl.Pagination, channel, client string) ([]tpl.MySetting, error) {
	data := []tpl.MySetting{}
	cursor := pg.TokenToTime(time.Now().Add(time.Minute * 10))
	size := pg.PageSize + 1

	for {
		rows, err := m.DB.Raw(userSettingsWithGroupSQL, userID, cursor, moduleIDs, size,
			groupIDs, cursor, moduleIDs, size).Rows()

		if err != nil {
			rows.Close()
			return nil, err
		}

		set := make(map[int64]struct{})
		count := 0
		for rows.Next() {
			count++
			var clients string
			var channels string
			mySetting := tpl.MySetting{}
			if err := rows.Scan(&mySetting.CreatedAt, &mySetting.UpdatedAt, &mySetting.Value, &mySetting.LastValue,
				&mySetting.ID, &mySetting.Name, &mySetting.Module, &channels, &clients); err != nil {
				rows.Close()
				return nil, err
			}

			if _, ok := set[mySetting.ID]; ok {
				continue // 去重
			}
			set[mySetting.ID] = struct{}{}

			if channel != "" && channels != "" {
				if !tpl.StringSliceHas(tpl.StringToSlice(channels), channel) {
					continue // channel 不匹配
				}
			}
			if client != "" && clients != "" {
				if !tpl.StringSliceHas(tpl.StringToSlice(clients), client) {
					continue // client 不匹配
				}
			}

			mySetting.HID = service.IDToHID(mySetting.ID, "setting")
			data = append(data, mySetting)
		}
		rows.Close()

		if count < size {
			break // no data to select
		}
		if len(data) >= size {
			break // get enough
		}
		// select next page
		cursor = data[len(data)-1].UpdatedAt.Add(-time.Millisecond)
	}

	return data, nil
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
	cursor := pg.TokenToID()

	rows, err := m.DB.Raw(userLabelsSQL, userID, cursor, pg.PageSize+1).Rows()
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

const userSettingsSQL = "select t1.`created_at`, t1.`updated_at`, t1.`value`, t1.`last_value`, " +
	"t2.`id`, t2.`name`, t3.`name` as `module` " +
	"from `user_setting` t1, `urbs_setting` t2, `urbs_module` t3 " +
	"where t1.`user_id` = ? and t1.`id` >= ? and t1.`setting_id` = t2.`id` and t2.`module_id` in ( ? ) and t2.`module_id` = t3.`id` " +
	"order by t1.`id` asc " +
	"limit ?"

// FindSettings 根据用户 ID, moduleIDs 返回 settings 数据。
func (m *User) FindSettings(ctx context.Context, userID int64, moduleIDs []int64, pg tpl.Pagination) ([]tpl.MySetting, error) {
	data := []tpl.MySetting{}
	cursor := pg.TokenToID()

	rows, err := m.DB.Raw(userSettingsSQL, userID, cursor, moduleIDs, pg.PageSize+1).Rows()
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

const rollbackUserSettingSQL = "update `user_setting` set `value` = `user_setting`.`last_value` where user_id = ? and setting_id = ?"

// RollbackSetting 回滚用户的 setting
func (m *User) RollbackSetting(ctx context.Context, userID, settingID int64) error {
	err := m.DB.Exec(rollbackUserSettingSQL, userID, settingID).Error
	return err
}
