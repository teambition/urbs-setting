package schema

// schema 模块不要引入官方库以外的其它模块或内部模块
import (
	"time"
)

// UserLabel 详见 ./sql/schema.sql table `user_label`
// 记录用户被分配的灰度标签，不同客户端不同大版本可能有不同的灰度标签
type UserLabel struct {
	ID        int64     `gorm:"column:id"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UserID    int64     `gorm:"column:user_id"`  // 用户内部 ID
	LabelID   int64     `gorm:"column:label_id"` // 灰度标签内部 ID
	Release   int64     `gorm:"column:rls"`      // 标签被设置计数批次
}
