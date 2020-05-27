package bll

import (
	"context"
	"strings"

	"github.com/teambition/gear"
	"github.com/teambition/urbs-setting/src/model"
	"github.com/teambition/urbs-setting/src/schema"
	"github.com/teambition/urbs-setting/src/tpl"
)

// Label ...
type Label struct {
	ms *model.Models
}

// List 返回产品下的标签列表
func (b *Label) List(ctx context.Context, productName string, pg tpl.Pagination) (*tpl.LabelsInfoRes, error) {
	productID, err := b.ms.Product.AcquireID(ctx, productName)
	if err != nil {
		return nil, err
	}

	labels, total, err := b.ms.Label.Find(ctx, productID, pg)
	if err != nil {
		return nil, err
	}

	labelInfos := tpl.LabelsInfoFrom(labels, productName)
	res := &tpl.LabelsInfoRes{Result: labelInfos}
	res.TotalSize = total
	if len(res.Result) > pg.PageSize {
		res.NextPageToken = tpl.IDToPageToken(res.Result[pg.PageSize].ID)
		res.Result = res.Result[:pg.PageSize]
	}
	return res, nil
}

// Create 创建标签
func (b *Label) Create(ctx context.Context, productName string, body *tpl.LabelBody) (*tpl.LabelInfoRes, error) {
	productID, err := b.ms.Product.AcquireID(ctx, productName)
	if err != nil {
		return nil, err
	}

	label := &schema.Label{ProductID: productID, Name: body.Name, Desc: body.Desc}
	if body.Channels != nil {
		label.Channels = strings.Join(*body.Channels, ",")
	}
	if body.Clients != nil {
		label.Clients = strings.Join(*body.Clients, ",")
	}
	if err = b.ms.Label.Create(ctx, label); err != nil {
		return nil, err
	}
	return &tpl.LabelInfoRes{Result: tpl.LabelInfoFrom(*label, productName)}, nil
}

// Update ...
func (b *Label) Update(ctx context.Context, productName, labelName string, body tpl.LabelUpdateBody) (*tpl.LabelInfoRes, error) {
	productID, err := b.ms.Product.AcquireID(ctx, productName)
	if err != nil {
		return nil, err
	}

	label, err := b.ms.Label.Acquire(ctx, productID, labelName)
	label, err = b.ms.Label.Update(ctx, label.ID, body.ToMap())
	if err != nil {
		return nil, err
	}
	return &tpl.LabelInfoRes{Result: tpl.LabelInfoFrom(*label, productName)}, nil
}

// Offline 下线标签
func (b *Label) Offline(ctx context.Context, productName, labelName string) (*tpl.BoolRes, error) {
	productID, err := b.ms.Product.AcquireID(ctx, productName)
	if err != nil {
		return nil, err
	}

	res := &tpl.BoolRes{Result: false}
	label, err := b.ms.Label.FindByName(ctx, productID, labelName, "id, `offline_at`")
	if err != nil {
		return nil, err
	}
	if label == nil {
		return nil, gear.ErrNotFound.WithMsgf("label %s not found", labelName)
	}
	if label.OfflineAt == nil {
		if err = b.ms.Label.Offline(ctx, label.ID); err != nil {
			return nil, err
		}
		res.Result = true
	}
	return res, nil
}

// Assign 把标签批量分配给用户或群组
func (b *Label) Assign(ctx context.Context, productName, labelName string, users, groups []string) (*tpl.LabelReleaseInfo, error) {
	productID, err := b.ms.Product.AcquireID(ctx, productName)
	if err != nil {
		return nil, err
	}

	label, err := b.ms.Label.Acquire(ctx, productID, labelName)
	return b.ms.Label.Assign(ctx, label.ID, users, groups)
}

// Delete 物理删除标签
func (b *Label) Delete(ctx context.Context, productName, labelName string) (*tpl.BoolRes, error) {
	productID, err := b.ms.Product.AcquireID(ctx, productName)
	if err != nil {
		return nil, err
	}

	label, err := b.ms.Label.FindByName(ctx, productID, labelName, "id, `offline_at`")
	if err != nil {
		return nil, err
	}

	res := &tpl.BoolRes{Result: false}
	if label != nil {
		if label.OfflineAt == nil {
			return nil, gear.ErrConflict.WithMsgf("label %s is not offline", labelName)
		}

		if err = b.ms.Label.Delete(ctx, label.ID); err != nil {
			return nil, err
		}
		res.Result = true
	}
	return res, nil
}

// Recall 撤销指定批次的用户或群组的环境标签
func (b *Label) Recall(ctx context.Context, productName, labelName string, release int64) (*tpl.BoolRes, error) {
	productID, err := b.ms.Product.AcquireID(ctx, productName)
	if err != nil {
		return nil, err
	}

	label, err := b.ms.Label.Acquire(ctx, productID, labelName)
	if err != nil {
		return nil, err
	}

	res := &tpl.BoolRes{Result: false}
	if err = b.ms.Label.Recall(ctx, label.ID, release); err != nil {
		return nil, err
	}
	res.Result = true
	return res, nil
}

// CreateRule ...
func (b *Label) CreateRule(ctx context.Context, productName, labelName string, body tpl.LabelRuleBody) (*tpl.LabelRuleInfoRes, error) {
	productID, err := b.ms.Product.AcquireID(ctx, productName)
	if err != nil {
		return nil, err
	}

	label, err := b.ms.Label.Acquire(ctx, productID, labelName)
	if err != nil {
		return nil, err
	}

	labelRule := &schema.LabelRule{
		ProductID: productID,
		LabelID:   label.ID,
		Kind:      body.Kind,
		Rule:      body.ToRule(),
		Release:   0,
	}
	if err = b.ms.LabelRule.Create(ctx, labelRule); err != nil {
		return nil, err
	}
	// 创建成功再从 label 获取当前的 release 发布计数
	release, err := b.ms.Label.AcquireRelease(ctx, label.ID)
	if err != nil {
		return nil, err
	}

	changed := map[string]interface{}{"rls": release}
	labelRule, err = b.ms.LabelRule.Update(ctx, labelRule.ID, changed)
	if err != nil {
		return nil, err
	}
	return &tpl.LabelRuleInfoRes{Result: tpl.LabelRuleInfoFrom(*labelRule)}, nil
}

// ListRules ...
func (b *Label) ListRules(ctx context.Context, productName, labelName string) (*tpl.LabelRulesInfoRes, error) {
	productID, err := b.ms.Product.AcquireID(ctx, productName)
	if err != nil {
		return nil, err
	}

	label, err := b.ms.Label.Acquire(ctx, productID, labelName)
	if err != nil {
		return nil, err
	}

	labelRules, err := b.ms.LabelRule.Find(ctx, productID, label.ID)
	if err != nil {
		return nil, err
	}

	res := &tpl.LabelRulesInfoRes{Result: tpl.LabelRulesInfoFrom(labelRules)}
	res.TotalSize = len(labelRules)
	return res, nil
}

// UpdateRule ...
func (b *Label) UpdateRule(ctx context.Context, productName, labelName string, ruleID int64, body tpl.LabelRuleBody) (*tpl.LabelRuleInfoRes, error) {
	productID, err := b.ms.Product.AcquireID(ctx, productName)
	if err != nil {
		return nil, err
	}

	label, err := b.ms.Label.Acquire(ctx, productID, labelName)
	if err != nil {
		return nil, err
	}

	labelRule, err := b.ms.LabelRule.Acquire(ctx, ruleID)
	if err != nil {
		return nil, err
	}

	if labelRule.LabelID != label.ID || body.Kind != labelRule.Kind {
		return nil, gear.ErrNotFound.WithMsgf("label rule not matched!")
	}

	changed := map[string]interface{}{}
	rule := body.ToRule()
	if rule != labelRule.Rule {
		changed["rule"] = rule
	}

	if len(changed) > 0 {
		release, err := b.ms.Label.AcquireRelease(ctx, label.ID)
		if err != nil {
			return nil, err
		}
		changed["rls"] = release
		labelRule, err = b.ms.LabelRule.Update(ctx, labelRule.ID, changed)
		if err != nil {
			return nil, err
		}
	}

	return &tpl.LabelRuleInfoRes{Result: tpl.LabelRuleInfoFrom(*labelRule)}, nil
}

// DeleteRule ...
func (b *Label) DeleteRule(ctx context.Context, productName, labelName string, ruleID int64) (*tpl.BoolRes, error) {
	productID, err := b.ms.Product.AcquireID(ctx, productName)
	if err != nil {
		return nil, err
	}

	label, err := b.ms.Label.Acquire(ctx, productID, labelName)
	if err != nil {
		return nil, err
	}

	res := &tpl.BoolRes{Result: false}
	labelRule, err := b.ms.LabelRule.Acquire(ctx, ruleID)
	if err != nil {
		return res, nil
	}

	if labelRule.LabelID != label.ID {
		return nil, gear.ErrNotFound.WithMsgf("label rule not matched!")
	}

	rowsAffected, err := b.ms.LabelRule.Delete(ctx, labelRule.ID)
	if err != nil {
		return nil, err
	}
	res.Result = rowsAffected > 0
	return res, nil
}

// ListUsers 返回产品下环境标签的用户列表
func (b *Label) ListUsers(ctx context.Context, productName, labelName string, pg tpl.Pagination) (*tpl.LabelUsersInfoRes, error) {
	productID, err := b.ms.Product.AcquireID(ctx, productName)
	if err != nil {
		return nil, err
	}

	label, err := b.ms.Label.Acquire(ctx, productID, labelName)
	if err != nil {
		return nil, err
	}

	data, total, err := b.ms.Label.ListUsers(ctx, label.ID, pg)
	if err != nil {
		return nil, err
	}
	res := &tpl.LabelUsersInfoRes{Result: data}
	res.TotalSize = total
	if len(res.Result) > pg.PageSize {
		res.NextPageToken = tpl.IDToPageToken(res.Result[pg.PageSize].ID)
		res.Result = res.Result[:pg.PageSize]
	}
	return res, nil
}

// DeleteUser ...
func (b *Label) DeleteUser(ctx context.Context, productName, labelName, uid string) (*tpl.BoolRes, error) {
	productID, err := b.ms.Product.AcquireID(ctx, productName)
	if err != nil {
		return nil, err
	}

	label, err := b.ms.Label.Acquire(ctx, productID, labelName)
	if err != nil {
		return nil, err
	}

	user, err := b.ms.User.Acquire(ctx, uid)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := b.ms.Label.RemoveUserLabel(ctx, user.ID, label.ID)
	if err != nil {
		return nil, err
	}
	return &tpl.BoolRes{Result: rowsAffected > 0}, nil
}

// ListGroups 返回产品下环境标签的群组列表
func (b *Label) ListGroups(ctx context.Context, productName, labelName string, pg tpl.Pagination) (*tpl.LabelGroupsInfoRes, error) {
	productID, err := b.ms.Product.AcquireID(ctx, productName)
	if err != nil {
		return nil, err
	}

	label, err := b.ms.Label.Acquire(ctx, productID, labelName)
	if err != nil {
		return nil, err
	}

	data, total, err := b.ms.Label.ListGroups(ctx, label.ID, pg)
	if err != nil {
		return nil, err
	}
	res := &tpl.LabelGroupsInfoRes{Result: data}
	res.TotalSize = total
	if len(res.Result) > pg.PageSize {
		res.NextPageToken = tpl.IDToPageToken(res.Result[pg.PageSize].ID)
		res.Result = res.Result[:pg.PageSize]
	}
	return res, nil
}

// DeleteGroup ...
func (b *Label) DeleteGroup(ctx context.Context, productName, labelName, uid string) (*tpl.BoolRes, error) {
	productID, err := b.ms.Product.AcquireID(ctx, productName)
	if err != nil {
		return nil, err
	}

	label, err := b.ms.Label.Acquire(ctx, productID, labelName)
	if err != nil {
		return nil, err
	}

	group, err := b.ms.Group.Acquire(ctx, uid)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := b.ms.Label.RemoveGroupLabel(ctx, group.ID, label.ID)
	if err != nil {
		return nil, err
	}
	return &tpl.BoolRes{Result: rowsAffected > 0}, nil
}
