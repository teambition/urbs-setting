package model

import (
	"context"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/teambition/gear"
	"github.com/teambition/urbs-setting/src/schema"
	"github.com/teambition/urbs-setting/src/service"
	"github.com/teambition/urbs-setting/src/tpl"
	"github.com/teambition/urbs-setting/src/util"
)

// Label ...
type Label struct {
	*Model
}

// FindByName 根据 productID 和 name 返回 label 数据
func (m *Label) FindByName(ctx context.Context, productID int64, name, selectStr string) (*schema.Label, error) {
	label := &schema.Label{}
	ok, err := m.findOneByCols(ctx, schema.TableLabel, goqu.Ex{"product_id": productID, "name": name}, selectStr, label)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}
	return label, nil
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
	label, err := m.FindByName(ctx, productID, labelName, "id, offline_at")
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
	label := &schema.Label{}
	if err := m.findOneByID(ctx, schema.TableLabel, labelID, label); err != nil {
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

	sdc := m.RdDB.Select().
		From(goqu.T(schema.TableLabel)).
		Where(
			goqu.C("product_id").Eq(productID),
			goqu.C("offline_at").IsNull())

	sd := m.RdDB.Select().
		From(goqu.T(schema.TableLabel)).
		Where(
			goqu.C("product_id").Eq(productID),
			goqu.C("id").Lte(cursor),
			goqu.C("offline_at").IsNull())

	if pg.Q != "" {
		sdc = sdc.Where(goqu.C("name").ILike(pg.Q))
		sd = sd.Where(goqu.C("name").ILike(pg.Q))
	}

	sd = sd.Order(goqu.C("id").Desc()).Limit(uint(pg.PageSize + 1))

	total, err := sdc.CountContext(ctx)
	if err != nil {
		return nil, 0, err
	}

	if err = sd.Executor().ScanStructsContext(ctx, &labels); err != nil {
		return nil, 0, err
	}

	return labels, int(total), nil
}

// Create ...
func (m *Label) Create(ctx context.Context, label *schema.Label) error {
	rowsAffected, err := m.createOne(ctx, schema.TableLabel, label)
	if rowsAffected > 0 {
		util.Go(5*time.Second, func(gctx context.Context) {
			m.tryIncreaseStatisticStatus(gctx, schema.LabelsTotalSize, 1)
		})
	}
	return err
}

// Update 更新指定环境标签
func (m *Label) Update(ctx context.Context, labelID int64, changed map[string]interface{}) (*schema.Label, error) {
	label := &schema.Label{}
	if _, err := m.updateByID(ctx, schema.TableLabel, labelID, goqu.Record(changed)); err != nil {
		return nil, err
	}
	if err := m.findOneByID(ctx, schema.TableLabel, labelID, label); err != nil {
		return nil, err
	}
	return label, nil
}

// Offline 标记 label 下线，同时真删除用户和群组的 labels
func (m *Label) Offline(ctx context.Context, labelID int64) error {
	return m.offlineLabels(ctx, goqu.Ex{"id": labelID, "offline_at": nil})
}

// Assign 把标签批量分配给用户或群组，如果用户或群组不存在则忽略
func (m *Label) Assign(ctx context.Context, labelID int64, users, groups []string) (*tpl.LabelReleaseInfo, error) {
	var err error
	totalRowsAffected := int64(0)
	release, err := m.AcquireRelease(ctx, labelID)
	if err != nil {
		return nil, err
	}

	releaseInfo := &tpl.LabelReleaseInfo{Release: release, Users: []string{}, Groups: []string{}}
	if len(users) > 0 {
		sd := m.DB.Insert(schema.TableUserLabel).Cols("user_id", "label_id", "rls").
			FromQuery(goqu.From(goqu.T(schema.TableUser).As("t1")).
				Select(goqu.I("t1.id"), goqu.V(labelID), goqu.V(release)).
				Where(goqu.I("t1.uid").In(tpl.StrSliceToInterface(users)...))).
			OnConflict(goqu.DoUpdate("rls", goqu.C("rls").Set(goqu.V(release))))

		rowsAffected, err := service.DeResult(sd.Executor().ExecContext(ctx))
		if err != nil {
			return nil, err
		}

		totalRowsAffected += rowsAffected
		if rowsAffected > 0 {
			sd := m.DB.Select(goqu.I("t2.uid")).
				From(
					goqu.T(schema.TableUserLabel).As("t1"),
					goqu.T(schema.TableUser).As("t2")).
				Where(
					goqu.I("t1.label_id").Eq(goqu.V(labelID)),
					goqu.I("t1.rls").Eq(goqu.V(release)),
					goqu.I("t1.user_id").Eq(goqu.I("t2.id"))).
				Order(goqu.I("t1.id").Desc()).Limit(1000)

			if err := sd.Executor().ScanValsContext(ctx, &releaseInfo.Users); err != nil {
				return nil, err
			}
		}
	}

	if len(groups) > 0 {
		sd := m.DB.Insert(schema.TableGroupLabel).Cols("group_id", "label_id", "rls").
			FromQuery(goqu.From(goqu.T(schema.TableGroup).As("t1")).
				Select(goqu.I("t1.id"), goqu.V(labelID), goqu.V(release)).
				Where(goqu.I("t1.uid").In(tpl.StrSliceToInterface(groups)...))).
			OnConflict(goqu.DoUpdate("rls", goqu.C("rls").Set(goqu.V(release))))

		rowsAffected, err := service.DeResult(sd.Executor().ExecContext(ctx))
		if err != nil {
			return nil, err
		}

		totalRowsAffected += rowsAffected
		if rowsAffected > 0 {
			sd := m.DB.Select(goqu.I("t2.uid")).
				From(
					goqu.T(schema.TableGroupLabel).As("t1"),
					goqu.T(schema.TableGroup).As("t2")).
				Where(
					goqu.I("t1.label_id").Eq(goqu.V(labelID)),
					goqu.I("t1.rls").Eq(goqu.V(release)),
					goqu.I("t1.group_id").Eq(goqu.I("t2.id"))).
				Order(goqu.I("t1.id").Desc()).Limit(1000)

			if err := sd.Executor().ScanValsContext(ctx, &releaseInfo.Groups); err != nil {
				return nil, err
			}
		}
	}

	if totalRowsAffected > 0 {
		util.Go(10*time.Second, func(gctx context.Context) {
			m.tryRefreshLabelStatus(gctx, labelID)
		})
	}
	return releaseInfo, err
}

// Delete 对标签进行物理删除
func (m *Label) Delete(ctx context.Context, id int64) error {
	_, err := m.deleteByID(ctx, schema.TableLabel, id)
	return err
}

// Cleanup 清除产品环境标签下所有的用户、群组和百分比规则
func (m *Label) Cleanup(ctx context.Context, id int64) error {
	_, err := m.deleteByCols(ctx, schema.TableLabelRule, goqu.Ex{"label_id": id})
	if err != nil {
		return err
	}
	_, err = m.deleteByCols(ctx, schema.TableGroupLabel, goqu.Ex{"label_id": id})
	if err != nil {
		return err
	}
	_, err = m.deleteByCols(ctx, schema.TableUserLabel, goqu.Ex{"label_id": id})
	if err != nil {
		return err
	}
	_, err = m.updateByID(ctx, schema.TableLabel, id, goqu.Record{"status": 0})
	return err
}

// RemoveUserLabel 删除用户的 label
func (m *Label) RemoveUserLabel(ctx context.Context, userID, labelID int64) (int64, error) {
	rowsAffected, err := m.deleteByCols(ctx, schema.TableUserLabel, goqu.Ex{"user_id": userID, "label_id": labelID})
	if rowsAffected > 0 {
		util.Go(5*time.Second, func(gctx context.Context) {
			m.tryIncreaseLabelsStatus(gctx, []int64{labelID}, -1)
		})
	}
	return rowsAffected, err
}

// RemoveGroupLabel 删除群组的 label
func (m *Label) RemoveGroupLabel(ctx context.Context, groupID, labelID int64) (int64, error) {
	rowsAffected, err := m.deleteByCols(ctx, schema.TableGroupLabel, goqu.Ex{"group_id": groupID, "label_id": labelID})
	if rowsAffected > 0 {
		util.Go(10*time.Second, func(gctx context.Context) {
			m.tryRefreshLabelStatus(gctx, labelID)
		})
	}
	return rowsAffected, err
}

// Recall 撤销指定批次的用户或群组的环境标签
func (m *Label) Recall(ctx context.Context, labelID, release int64) error {
	totalRowsAffected := int64(0)
	rowsAffected, err := m.deleteByCols(ctx, schema.TableGroupLabel, goqu.Ex{"label_id": labelID, "rls": release})
	if err != nil {
		return err
	}
	totalRowsAffected += rowsAffected

	rowsAffected, err = m.deleteByCols(ctx, schema.TableUserLabel, goqu.Ex{"label_id": labelID, "rls": release})
	if err != nil {
		return err
	}
	totalRowsAffected += rowsAffected
	if totalRowsAffected > 0 {
		util.Go(10*time.Second, func(gctx context.Context) {
			m.tryRefreshLabelStatus(gctx, labelID)
		})
	}
	return nil
}

// AcquireRelease ...
func (m *Label) AcquireRelease(ctx context.Context, labelID int64) (int64, error) {
	label := &schema.Label{}
	if _, err := m.updateByID(ctx, schema.TableLabel, labelID, goqu.Record{
		"rls": goqu.L("rls + ?", 1),
	}); err != nil {
		return 0, err
	}

	// MySQL 不支持 RETURNING，并发操作分配时 release 可能不准确，不过真实场景下基本不可能并发操作
	if err := m.findOneByID(ctx, schema.TableLabel, labelID, label); err != nil {
		return 0, err
	}
	return label.Release, nil
}

// ListUsers ...
func (m *Label) ListUsers(ctx context.Context, labelID int64, pg tpl.Pagination) ([]tpl.LabelUserInfo, int, error) {
	data := []tpl.LabelUserInfo{}
	cursor := pg.TokenToID()

	sdc := m.RdDB.Select().
		From(
			goqu.T(schema.TableUserLabel).As("t1"),
			goqu.T(schema.TableUser).As("t2")).
		Where(
			goqu.I("t1.label_id").Eq(labelID),
			goqu.I("t1.user_id").Eq(goqu.I("t2.id")))

	sd := m.RdDB.Select(
		goqu.I("t1.id"),
		goqu.I("t1.created_at").As("assigned_at"),
		goqu.I("t1.rls"),
		goqu.I("t2.uid")).
		From(
			goqu.T(schema.TableUserLabel).As("t1"),
			goqu.T(schema.TableUser).As("t2")).
		Where(
			goqu.I("t1.label_id").Eq(labelID),
			goqu.I("t1.id").Lte(cursor),
			goqu.I("t1.user_id").Eq(goqu.I("t2.id")))

	if pg.Q != "" {
		sdc = sdc.Where(goqu.I("t2.uid").ILike(pg.Q))
		sd = sd.Where(goqu.I("t2.uid").ILike(pg.Q))
	}

	sd = sd.Order(goqu.I("t1.id").Desc()).Limit(uint(pg.PageSize + 1))

	total, err := sdc.CountContext(ctx)
	if err != nil {
		return nil, 0, err
	}
	scanner, err := sd.Executor().ScannerContext(ctx)
	if err != nil {
		return nil, 0, err
	}
	defer scanner.Close()

	for scanner.Next() {
		info := tpl.LabelUserInfo{}
		if err := scanner.ScanStruct(&info); err != nil {
			return nil, 0, err
		}
		info.LabelHID = service.IDToHID(labelID, "label")
		data = append(data, info)
	}

	if err := scanner.Err(); err != nil {
		return nil, 0, err
	}

	return data, int(total), err
}

// ListGroups ...
func (m *Label) ListGroups(ctx context.Context, labelID int64, pg tpl.Pagination) ([]tpl.LabelGroupInfo, int, error) {
	data := []tpl.LabelGroupInfo{}
	cursor := pg.TokenToID()
	sdc := m.RdDB.Select().
		From(
			goqu.T(schema.TableGroupLabel).As("t1"),
			goqu.T(schema.TableGroup).As("t2")).
		Where(
			goqu.I("t1.label_id").Eq(labelID),
			goqu.I("t1.group_id").Eq(goqu.I("t2.id")))

	sd := m.RdDB.Select(
		goqu.I("t1.id"),
		goqu.I("t1.created_at").As("assigned_at"),
		goqu.I("t1.rls"),
		goqu.I("t2.uid"),
		goqu.I("t2.kind"),
		goqu.I("t2.description"),
		goqu.I("t2.status")).
		From(
			goqu.T(schema.TableGroupLabel).As("t1"),
			goqu.T(schema.TableGroup).As("t2")).
		Where(
			goqu.I("t1.label_id").Eq(labelID),
			goqu.I("t1.id").Lte(cursor),
			goqu.I("t1.group_id").Eq(goqu.I("t2.id")))

	if pg.Q != "" {
		sdc = sdc.Where(goqu.I("t2.uid").ILike(pg.Q))
		sd = sd.Where(goqu.I("t2.uid").ILike(pg.Q))
	}

	sd = sd.Order(goqu.I("t1.id").Desc()).Limit(uint(pg.PageSize + 1))

	total, err := sdc.CountContext(ctx)
	if err != nil {
		return nil, 0, err
	}
	scanner, err := sd.Executor().ScannerContext(ctx)
	if err != nil {
		return nil, 0, err
	}
	defer scanner.Close()

	for scanner.Next() {
		info := tpl.LabelGroupInfo{}
		if err := scanner.ScanStruct(&info); err != nil {
			return nil, 0, err
		}
		info.LabelHID = service.IDToHID(labelID, "label")
		data = append(data, info)
	}

	if err := scanner.Err(); err != nil {
		return nil, 0, err
	}

	return data, int(total), err
}
