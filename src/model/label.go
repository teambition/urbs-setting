package model

import (
	"context"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/teambition/urbs-setting/src/schema"
	"github.com/teambition/urbs-setting/src/tpl"
)

// Label ...
type Label struct {
	DB *gorm.DB
}

// FindByName 根据 productID 和 name 返回 label 数据
func (m *Label) FindByName(ctx context.Context, productID int64, name, selectStr string) (*schema.Label, error) {
	var err error
	label := &schema.Label{}

	db := m.DB.Where("product_id = ? and name = ?", productID, name)

	if selectStr == "" {
		err = db.First(label).Error
	} else {
		err = db.Select(selectStr).First(label).Error
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
func (m *Label) Find(ctx context.Context, productID int64, pg tpl.Pagination) ([]schema.Label, error) {
	labels := make([]schema.Label, 0)
	cursor := pg.TokenToID()
	err := m.DB.Where("`product_id` = ? and `id` >= ?", productID, cursor).
		Order("`id`").Limit(pg.PageSize + 1).Find(&labels).Error
	return labels, err
}

// Count 计算 product labels 总数
func (m *Label) Count(ctx context.Context, productID int64) (int, error) {
	count := 0
	err := m.DB.Model(&schema.Label{}).Where("product_id = ?", productID).Count(&count).Error
	return count, err
}

// Create ...
func (m *Label) Create(ctx context.Context, label *schema.Label) error {
	return m.DB.Create(label).Error
}

// Update 更新指定灰度标签
func (m *Label) Update(ctx context.Context, labelID int64, changed map[string]interface{}) (*schema.Label, error) {
	label := &schema.Label{ID: labelID}
	if len(changed) > 0 {
		if err := m.DB.Model(label).UpdateColumns(changed).Error; err != nil {
			return nil, err
		}
	}

	if err := m.DB.First(label).Error; err != nil {
		return nil, err
	}
	return label, nil
}

// Offline 标记 label 下线，同时真删除用户和群组的 labels
func (m *Label) Offline(ctx context.Context, labelID int64) error {
	now := time.Now().UTC()
	db := m.DB.Model(&schema.Label{ID: labelID}).UpdateColumns(schema.Label{
		OfflineAt: &now,
		Status:    -1,
	})
	if db.Error == nil {
		go deleteUserAndGroupLabels(db, []int64{labelID})
	}
	return db.Error
}

const batchAddUserLabelSQL = "insert ignore into `user_label` (`user_id`, `label_id`) " +
	"select `urbs_user`.id, ? from `urbs_user` where `urbs_user`.uid in ( ? )"
const batchAddGroupLabelSQL = "insert ignore into `group_label` (`group_id`, `label_id`) " +
	"select `urbs_group`.id, ? from `urbs_group` where `urbs_group`.uid in ( ? )"

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

// Delete 对标签进行物理删除
func (m *Label) Delete(ctx context.Context, labelID int64) error {
	return m.DB.Delete(&schema.Label{ID: labelID}).Error
}
