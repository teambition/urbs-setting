package schema

// schema 模块不要引入官方库以外的其它模块或内部模块
import (
	"time"
)

// TableLabelRule is a table name in db.
const TableLabelRule = "label_rule"

// LabelRule 详见 ./sql/schema.sql table `label_rule`
// 环境标签发布规则
type LabelRule struct {
	ID        int64     `db:"id" goqu:"skipinsert"`
	CreatedAt time.Time `db:"created_at" goqu:"skipinsert"`
	UpdatedAt time.Time `db:"updated_at" goqu:"skipinsert"`
	ProductID int64     `db:"product_id"` // 所从属的产品线 ID，与环境标签的产品线一致
	LabelID   int64     `db:"label_id"`   // 规则所指向的环境标签 ID
	Kind      string    `db:"kind"`       // 规则类型，目前支持 "userPercent"
	Rule      string    `db:"rule"`       // varchar(1022)，规则值，JSON string，对于 percent 类，其格式为 {"value": percent}
	Release   int64     `db:"rls"`        // 标签发布（被设置）计数批次
}

// TableName retuns table name
func (LabelRule) TableName() string {
	return "label_rule"
}

// ToPercent retuns table name
func (l LabelRule) ToPercent() int {
	return ToPercentRule(l.Kind, l.Rule).Rule.Value
}
