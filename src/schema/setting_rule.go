package schema

// schema 模块不要引入官方库以外的其它模块或内部模块
import (
	"time"
)

// TableSettingRule is a table name in db.
const TableSettingRule = "setting_rule"

// SettingRule 详见 ./sql/schema.sql table `setting_rule`
// 灰度标签发布规则
type SettingRule struct {
	ID        int64     `db:"id" goqu:"skipinsert"`
	CreatedAt time.Time `db:"created_at" goqu:"skipinsert"`
	UpdatedAt time.Time `db:"updated_at" goqu:"skipinsert"`
	ProductID int64     `db:"product_id"` // 所从属的产品线 ID，与灰度标签的产品线一致
	SettingID int64     `db:"setting_id"` // 规则所指向的灰度标签 ID
	Kind      string    `db:"kind"`       // 规则类型，目前支持 "userPercent"
	Rule      string    `db:"rule"`       // varchar(1022)，规则值，JSON string，对于 percent 类，其格式为 {"value": percent}
	Value     string    `db:"value"`      // varchar(255)，配置值
	Release   int64     `db:"rls"`        // 标签发布（被设置）计数批次
}

// TableName retuns table name
func (SettingRule) TableName() string {
	return "setting_rule"
}

// ToPercent retuns table name
func (l SettingRule) ToPercent() int {
	return ToPercentRule(l.Kind, l.Rule).Rule.Value
}
