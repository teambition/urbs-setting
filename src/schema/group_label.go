package schema

// schema 模块不要引入官方库以外的其它模块或内部模块
import (
	"time"
)

// GroupLabel 详见 ./sql/schema.sql table `group_label`
// 记录群组被设置的灰度标签，将作用于群组所有成员
type GroupLabel struct {
	ID        int64     `gorm:"column:id" db:"id"`
	CreatedAt time.Time `gorm:"column:created_at" db:"created_at"`
	GroupID   int64     `gorm:"column:group_id" db:"group_id"` // 群组内部 ID
	LabelID   int64     `gorm:"column:label_id" db:"label_id"` // 灰度标签内部 ID
	Release   int64     `gorm:"column:rls" db:"rls"`           // 标签被设置计数批次
}
