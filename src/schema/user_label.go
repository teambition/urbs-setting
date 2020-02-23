package schema

import (
	"time"
)

// UserLabel 详见 ./sql/schema.sql table `user_label`
// 记录用户被分配的灰度标签，不同客户端不同大版本可能有不同的灰度标签
type UserLabel struct {
	ID        int64     `gorm:"column:id" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UserID    int64     `gorm:"column:user_id" json:"user_id"`   // 用户内部 ID
	LabelID   int64     `gorm:"column:label_id" json:"label_id"` // 灰度标签内部 ID
	Channel   string    `gorm:"column:channel" json:"channel"`   // varchar(31)，产品版本更新通道
	Client    string    `gorm:"column:client" json:"client"`     // varchar(31)，产品客户端类型
}
