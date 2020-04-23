package schema

// schema 模块不要引入官方库以外的其它模块或内部模块
import (
	"time"
)

// LabelRule 详见 ./sql/schema.sql table `label_rule`
// 灰度标签发布规则
type LabelRule struct {
	ID        int64     `gorm:"column:id"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
	ProductID int64     `gorm:"column:product_id"` // 所从属的产品线 ID，与灰度标签的产品线一致
	LabelID   int64     `gorm:"column:label_id"`   // 规则所指向的灰度标签 ID
	Kind      string    `gorm:"column:kind"`       // 规则类型，目前支持 "user_percent"
	Rule      string    `gorm:"column:rule"`       // varchar(1022)，规则值，JSON string，对于 percent 类，其格式为 {"value": percent}
	Release   int64     `gorm:"column:rls"`        // 标签发布（被设置）计数批次
}

// TableName retuns table name
func (LabelRule) TableName() string {
	return "label_rule"
}

// ToPercent retuns table name
func (l LabelRule) ToPercent() int {
	return ToPercentRule(l.Kind, l.Rule).Rule.Value
}
