package model

import (
	"context"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/teambition/urbs-setting/src/schema"
)

// Label ...
type Label struct {
	DB *gorm.DB
}

// FindByName 根据 productID 和 name 返回 label 数据
func (m *Label) FindByName(ctx context.Context, productID int64, name, selectStr string) (*schema.Label, error) {
	var err error
	label := &schema.Label{ProductID: productID, Name: name}
	if selectStr == "" {
		err = m.DB.Take(label).Error
	} else {
		err = m.DB.Select(selectStr).Take(label).Error
	}

	if err == nil {
		return label, nil
	}

	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}
	return nil, err
}

// Find 根据条件查找 labels
func (m *Label) Find(ctx context.Context, productID int64) ([]schema.Label, error) {
	labels := make([]schema.Label, 0)
	err := m.DB.Where("`product_id` is ?", productID).Order("`status`, `created_at`").Limit(1000).Find(labels).Error
	return labels, err
}

// Create ...
func (m *Label) Create(ctx context.Context, label *schema.Label) error {
	db := m.DB.Create(label)
	if db.Error != nil {
		return db.Error
	}

	return db.Take(label).Error
}

// Offline 标记 label 下线，同时真删除用户和群组的 labels
func (m *Label) Offline(ctx context.Context, labelID int64) error {
	now := time.Now().UTC()
	db := m.DB.Model(&schema.Label{ID: labelID}).Update(schema.Label{
		OfflineAt: &now,
		Status:    -1,
	})
	if db.Error == nil {
		go deleteUserAndGroupLabels(db.DB(), []int64{labelID})
	}
	return db.Error
}

const batchAddUserLabelSQL = "insert ignore into `user_label` (`user_id`, `label_id`) " +
	"select `user`.id, ? from `user` where `user`.uid in ( ? )"
const batchAddGroupLabelSQL = "insert ignore into `group_label` (`group_id`, `label_id`) " +
	"select `group`.id, ? from `group` where `group`.uid in ( ? )"

// Assign 把标签批量分配给用户或群组，如果用户或群组不存在则忽略
func (m *Label) Assign(ctx context.Context, labelID int64, users, groups []string) error {
	var err error
	if len(users) > 0 {
		err = m.DB.Exec(batchAddUserLabelSQL, labelID, users).Error
	}
	if err == nil && len(groups) > 0 {
		err = m.DB.Exec(batchAddGroupLabelSQL, labelID, groups).Error
	}

	return err
}
