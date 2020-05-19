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

// Label ...
type Label struct {
	*Model
}

// FindByName 根据 productID 和 name 返回 label 数据
func (m *Label) FindByName(ctx context.Context, productID int64, name, selectStr string) (*schema.Label, error) {
	var err error
	label := &schema.Label{}

	db := m.DB.Where("`product_id` = ? and `name` = ?", productID, name)

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

// Acquire ...
func (m *Label) Acquire(ctx context.Context, productID int64, labelName string) (*schema.Label, error) {
	label, err := m.FindByName(ctx, productID, labelName, "")
	if err != nil {
		return nil, err
	}
	if label == nil {
		return nil, gear.ErrNotFound.WithMsgf("label %s not found", labelName)
	}
	if label.OfflineAt != nil {
		return nil, gear.ErrNotFound.WithMsgf("label %s was offline", labelName)
	}
	return label, nil
}

// AcquireID ...
func (m *Label) AcquireID(ctx context.Context, productID int64, labelName string) (int64, error) {
	label, err := m.FindByName(ctx, productID, labelName, "`id`, `offline_at`")
	if err != nil {
		return 0, err
	}
	if label == nil {
		return 0, gear.ErrNotFound.WithMsgf("label %s not found", labelName)
	}
	if label.OfflineAt != nil {
		return 0, gear.ErrNotFound.WithMsgf("label %s was offline", labelName)
	}
	return label.ID, nil
}

// AcquireByID ...
func (m *Label) AcquireByID(ctx context.Context, labelID int64) (*schema.Label, error) {
	label := &schema.Label{ID: labelID}
	if err := m.DB.First(label).Error; err != nil {
		return nil, err
	}
	if label.OfflineAt != nil {
		return nil, gear.ErrNotFound.WithMsgf("label %d was offline", labelID)
	}
	return label, nil
}

// Find 根据条件查找 labels
func (m *Label) Find(ctx context.Context, productID int64, pg tpl.Pagination) ([]schema.Label, int, error) {
	labels := make([]schema.Label, 0)
	cursor := pg.TokenToID()
	db := m.DB.Where("`id` <= ? and `product_id` = ? and `offline_at` is null", cursor, productID)
	dbc := m.DB.Where("`product_id` = ? and `offline_at` is null", productID)
	if pg.Q != "" {
		db = m.DB.Where("`id` <= ? and `product_id` = ? and `offline_at` is null and `name` like ?", cursor, productID, pg.Q)
		dbc = m.DB.Where("`product_id` = ? and `offline_at` is null and `name` like ?", productID, pg.Q)
	}

	total := 0
	err := dbc.Model(&schema.Label{}).Count(&total).Error
	if err == nil {
		err = db.Order("`id` desc").Limit(pg.PageSize + 1).Find(&labels).Error
	}
	if err != nil {
		return nil, 0, err
	}
	return labels, total, nil
}

// Create ...
func (m *Label) Create(ctx context.Context, label *schema.Label) error {
	err := m.DB.Create(label).Error
	if err == nil {
		go m.tryIncreaseStatisticStatus(ctx, schema.LabelsTotalSize, 1)
	}
	return err
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
	res := m.DB.Model(&schema.Label{ID: labelID}).UpdateColumns(schema.Label{
		OfflineAt: &now,
		Status:    -1,
	})
	if res.RowsAffected > 0 {
		go m.tryDeleteLabelsRules(ctx, []int64{labelID})
		go m.tryDeleteUserAndGroupLabels(ctx, []int64{labelID})
		go m.tryIncreaseStatisticStatus(ctx, schema.LabelsTotalSize, -1)
	}
	return res.Error
}

const batchAddUserLabelSQL = "insert ignore into `user_label` (`user_id`, `label_id`, `rls`) " +
	"select `urbs_user`.`id`, ?, ? from `urbs_user` where `urbs_user`.`uid` in ( ? ) " +
	"on duplicate key update `rls` = ?"
const batchAddGroupLabelSQL = "insert ignore into `group_label` (`group_id`, `label_id`, `rls`) " +
	"select `urbs_group`.`id`, ?, ? from `urbs_group` where `urbs_group`.`uid` in ( ? ) " +
	"on duplicate key update `rls` = ?"
const checkAddUserLabelSQL = "select t2.`uid` " +
	"from `user_label` t1, `urbs_user` t2 " +
	"where t1.`label_id` = ? and t1.`rls` = ? and t1.`user_id` = t2.`id` " +
	"order by t1.`id` desc limit 1000"
const checkAddGroupLabelSQL = "select t2.`uid` " +
	"from `group_label` t1, `urbs_group` t2 " +
	"where t1.`label_id` = ? and t1.`rls` = ? and t1.`group_id` = t2.`id` " +
	"order by t1.`id` desc limit 1000"

// Assign 把标签批量分配给用户或群组，如果用户或群组不存在则忽略
func (m *Label) Assign(ctx context.Context, labelID int64, users, groups []string) (*tpl.LabelReleaseInfo, error) {
	var err error
	rowsAffected := int64(0)
	release, err := m.AcquireRelease(ctx, labelID)
	if err != nil {
		return nil, err
	}

	releaseInfo := &tpl.LabelReleaseInfo{Release: release, Users: []string{}, Groups: []string{}}
	if len(users) > 0 {
		res := m.DB.Exec(batchAddUserLabelSQL, labelID, release, users, release)
		rowsAffected += res.RowsAffected
		err = res.Error
		if err == nil && res.RowsAffected > 0 {
			rows, err := m.DB.Raw(checkAddUserLabelSQL, labelID, release).Rows()

			if err != nil {
				rows.Close()
				return nil, err
			}

			for rows.Next() {
				var uid string
				if err := rows.Scan(&uid); err != nil {
					rows.Close()
					return nil, err
				}
				releaseInfo.Users = append(releaseInfo.Users, uid)
			}
			rows.Close()
		}
	}

	if err == nil && len(groups) > 0 {
		res := m.DB.Exec(batchAddGroupLabelSQL, labelID, release, groups, release)
		rowsAffected += res.RowsAffected
		err = res.Error
		if err == nil && res.RowsAffected > 0 {
			rows, err := m.DB.Raw(checkAddGroupLabelSQL, labelID, release).Rows()

			if err != nil {
				rows.Close()
				return nil, err
			}

			for rows.Next() {
				var uid string
				if err := rows.Scan(&uid); err != nil {
					rows.Close()
					return nil, err
				}
				releaseInfo.Groups = append(releaseInfo.Groups, uid)
			}
			rows.Close()
		}
	}

	if rowsAffected > 0 {
		go m.tryRefreshLabelStatus(ctx, labelID)
	}
	return releaseInfo, err
}

// Delete 对标签进行物理删除
func (m *Label) Delete(ctx context.Context, labelID int64) error {
	res := m.DB.Delete(&schema.Label{ID: labelID})
	return res.Error
}

// RemoveUserLabel 删除用户的 label
func (m *Label) RemoveUserLabel(ctx context.Context, userID, labelID int64) (int64, error) {
	res := m.DB.Where("`user_id` = ? and `label_id` = ?", userID, labelID).Delete(&schema.UserLabel{})
	if res.RowsAffected > 0 {
		go m.tryIncreaseLabelsStatus(ctx, []int64{labelID}, -1)
	}
	return res.RowsAffected, res.Error
}

// RemoveGroupLabel 删除群组的 label
func (m *Label) RemoveGroupLabel(ctx context.Context, groupID, labelID int64) (int64, error) {
	res := m.DB.Where("`group_id` = ? and `label_id` = ?", groupID, labelID).Delete(&schema.GroupLabel{})
	if res.RowsAffected > 0 {
		go m.tryRefreshLabelStatus(ctx, labelID)
	}
	return res.RowsAffected, res.Error
}

// Recall 撤销指定批次的用户或群组的灰度标签
func (m *Label) Recall(ctx context.Context, labelID, release int64) error {
	rowsAffected := int64(0)
	res := m.DB.Where("`label_id` = ? and `rls` = ?", labelID, release).Delete(&schema.GroupLabel{})
	rowsAffected += res.RowsAffected

	if res.Error == nil {
		res = m.DB.Where("`label_id` = ? and `rls` = ?", labelID, release).Delete(&schema.UserLabel{})
		rowsAffected += res.RowsAffected
	}
	if rowsAffected > 0 {
		go m.tryRefreshLabelStatus(ctx, labelID)
	}
	return res.Error
}

// AcquireRelease ...
func (m *Label) AcquireRelease(ctx context.Context, labelID int64) (int64, error) {
	label := &schema.Label{ID: labelID}
	if err := m.DB.Model(label).UpdateColumn("rls", gorm.Expr("`rls` + ?", 1)).Error; err != nil {
		return 0, err
	}
	// MySQL 不支持 RETURNING，并发操作分配时 release 可能不准确，不过真实场景下基本不可能并发操作
	if err := m.DB.Select("`id`, `rls`").First(label).Error; err != nil {
		return 0, err
	}
	return label.Release, nil
}

const listLabelUsersSQL = "select t1.`id`, t1.`created_at`, t1.`rls`, t2.`uid` " +
	"from `user_label` t1, `urbs_user` t2 " +
	"where t1.`label_id` = ? and t1.`id` <= ? and t1.`user_id` = t2.`id` " +
	"order by t1.`id` desc " +
	"limit ?"

const countLabelUsersSQL = "select count(t2.`id`) " +
	"from `user_label` t1, `urbs_user` t2  " +
	"where t1.`label_id` = ? and t1.`user_id` = t2.`id`"

const searchLabelUsersSQL = "select t1.`id`, t1.`created_at`, t1.`rls`, t2.`uid` " +
	"from `user_label` t1, `urbs_user` t2 " +
	"where t1.`label_id` = ? and t1.`id` <= ? and t1.`user_id` = t2.`id` and t2.`uid` like ? " +
	"order by t1.`id` desc " +
	"limit ?"

const countSearchLabelUsersSQL = "select count(t2.`id`) " +
	"from `user_label` t1, `urbs_user` t2 " +
	"where t1.`label_id` = ? and t1.`user_id` = t2.`id` and t2.`uid` like ?"

// ListUsers ...
func (m *Label) ListUsers(ctx context.Context, labelID int64, pg tpl.Pagination) ([]tpl.LabelUserInfo, int, error) {
	data := []tpl.LabelUserInfo{}
	cursor := pg.TokenToID()
	total := 0

	if pg.Q == "" {
		if err := m.DB.Raw(countLabelUsersSQL, labelID).Row().Scan(&total); err != nil && err != sql.ErrNoRows {
			return nil, 0, err
		}
	} else {
		if err := m.DB.Raw(countSearchLabelUsersSQL, labelID, pg.Q).Row().Scan(&total); err != nil && err != sql.ErrNoRows {
			return nil, 0, err
		}
	}

	var err error
	var rows *sql.Rows
	if pg.Q == "" {
		rows, err = m.DB.Raw(listLabelUsersSQL, labelID, cursor, pg.PageSize+1).Rows()
	} else {
		rows, err = m.DB.Raw(searchLabelUsersSQL, labelID, cursor, pg.Q, pg.PageSize+1).Rows()
	}
	defer rows.Close()

	if err != nil {
		return nil, 0, err
	}

	for rows.Next() {
		info := tpl.LabelUserInfo{}
		if err := rows.Scan(&info.ID, &info.AssignedAt, &info.Release, &info.User); err != nil {
			return nil, 0, err
		}
		info.LabelHID = service.IDToHID(labelID, "label")
		data = append(data, info)
	}
	return data, total, err
}

const listLabelGroupsSQL = "select t1.`id`, t1.`created_at`, t1.`rls`, t2.`uid`, t2.`kind`, t2.`description`, t2.`status` " +
	"from `group_label` t1, `urbs_group` t2 " +
	"where t1.`label_id` = ? and t1.`id` <= ? and t1.`group_id` = t2.`id` " +
	"order by t1.`id` desc " +
	"limit ?"

const countLabelGroupsSQL = "select count(t2.`id`) " +
	"from `group_label` t1, `urbs_group` t2  " +
	"where t1.`label_id` = ? and t1.`group_id` = t2.`id`"

const searchLabelGroupsSQL = "select t1.`id`, t1.`created_at`, t1.`rls`, t2.`uid`, t2.`kind`, t2.`description`, t2.`status` " +
	"from `group_label` t1, `urbs_group` t2 " +
	"where t1.`label_id` = ? and t1.`id` <= ? and t1.`group_id` = t2.`id` and t2.`uid` like ? " +
	"order by t1.`id` desc " +
	"limit ?"

const countSearchLabelGroupsSQL = "select count(t2.`id`) " +
	"from `group_label` t1, `urbs_group` t2 " +
	"where t1.`label_id` = ? and t1.`group_id` = t2.`id` and t2.`uid` like ?"

// ListGroups ...
func (m *Label) ListGroups(ctx context.Context, labelID int64, pg tpl.Pagination) ([]tpl.LabelGroupInfo, int, error) {
	data := []tpl.LabelGroupInfo{}
	cursor := pg.TokenToID()
	total := 0

	if pg.Q == "" {
		if err := m.DB.Raw(countLabelGroupsSQL, labelID).Row().Scan(&total); err != nil && err != sql.ErrNoRows {
			return nil, 0, err
		}
	} else {
		if err := m.DB.Raw(countSearchLabelGroupsSQL, labelID, pg.Q).Row().Scan(&total); err != nil && err != sql.ErrNoRows {
			return nil, 0, err
		}
	}

	var err error
	var rows *sql.Rows
	if pg.Q == "" {
		rows, err = m.DB.Raw(listLabelGroupsSQL, labelID, cursor, pg.PageSize+1).Rows()
	} else {
		rows, err = m.DB.Raw(searchLabelGroupsSQL, labelID, cursor, pg.Q, pg.PageSize+1).Rows()
	}
	defer rows.Close()

	if err != nil {
		return nil, 0, err
	}

	for rows.Next() {
		info := tpl.LabelGroupInfo{}
		if err := rows.Scan(&info.ID, &info.AssignedAt, &info.Release, &info.Group, &info.Kind, &info.Desc, &info.Status); err != nil {
			return nil, 0, err
		}
		info.LabelHID = service.IDToHID(labelID, "label")
		data = append(data, info)
	}
	return data, total, err
}
