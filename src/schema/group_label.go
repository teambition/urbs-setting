package schema

import (
	"time"
)

// GroupLabel 详见 ./sql/schema.sql table `group_label`
// 记录群组被设置的灰度标签，将作用于群组所有成员
type GroupLabel struct {
	ID        int64     `gorm:"column:id" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	GroupID   int64     `gorm:"column:group_id" json:"group_id"` // 群组内部 ID
	LabelID   int64     `gorm:"column:label_id" json:"label_id"` // 灰度标签内部 ID
}
