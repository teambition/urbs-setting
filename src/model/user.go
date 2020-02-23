package model

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/teambition/urbs-setting/src/schema"
)

type user struct {
	DB *gorm.DB
}

// FindByUID 根据 uid 返回 user 数据
func (m *user) FindByUID(ctx context.Context, uid string, selectStr string) (*schema.User, error) {
	user := &schema.User{UID: uid}
	if selectStr == "" {
		selectStr = "id"
	}

	err := m.DB.Select(selectStr).Take(user).Error
	if err == nil {
		return user, nil
	}

	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}
	return nil, err
}

func (m *user) RefreshLabels(ctx context.Context, id int64, now int64) (string, error) {
	labels := ""
	err := m.DB.Transaction(func(tx *gorm.DB) error {
		user := &schema.User{ID: id}
		// 指定 id 的记录被锁住，如果表中无符合记录的数据则排他锁不生效
		if err := tx.Set("gorm:query_option", "FOR UPDATE").Take(user).Error; err != nil {
			return err
		}

		if !user.IsStale(now) {
			// 已被其它请求更新
			labels = user.Labels
			return nil
		}

		data := []schema.UserLabelInfo{}

		rows1, err := tx.Raw("select t1.`channel`, t1.`client`, t2.`name` as `label`, t3.`name` as `product` "+
			"from `user_label` t1, `label` t2, `product` t3 "+
			"where t1.`user_id` = ? and t1.`label_id` = t2.id and t2.`product_id` = t3.id "+
			"order by t1.`created_at` desc "+
			"limit 1000", id).Rows()
		defer rows1.Close()

		if err != nil {
			return err
		}
		for rows1.Next() {
			label := schema.UserLabelInfo{}
			// ScanRows 扫描一行记录到 user
			if err := rows1.Scan(&label.Channel, &label.Client, &label.Label, &label.Product); err != nil {
				return err
			}
			data = append(data, label)
		}

		rows2, err := tx.Raw("select t3.`name` as `label`, t4.`name` as `product` "+
			"from `user_group` t1, `group_label` t2, `label` t3, `product` t4 "+
			"where t1.`user_id` = ? and t1.`group_id` = t2.`group_id` and t2.`label_id` = t3.id and t3.`product_id` = t4.id "+
			"order by t2.`created_at` desc "+
			"limit 1000", id).Rows()
		defer rows1.Close()

		if err != nil {
			return err
		}
		for rows2.Next() {
			label := schema.UserLabelInfo{}
			// ScanRows 扫描一行记录到 user
			if err := rows2.Scan(&label.Label, &label.Product); err != nil {
				return err
			}
			data = append(data, label)
		}
		user.PutLabels(data)
		user.ActiveAt = now
		labels = user.Labels

		return tx.Model(user).Updates(user).Error // 返回 nil 提交事务，否则回滚
	})

	return labels, err
}
