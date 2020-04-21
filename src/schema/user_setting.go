package schema

// schema 模块不要引入官方库以外的其它模块或内部模块
import (
	"time"
)

// UserSetting 详见 ./sql/schema.sql table `user_setting`
// 记录用户对某功能模块配置项值
type UserSetting struct {
	ID        int64     `gorm:"column:id"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
	UserID    int64     `gorm:"column:user_id"`    // 用户内部 ID
	SettingID int64     `gorm:"column:setting_id"` // 配置项内部 ID
	Value     string    `gorm:"column:value"`      // varchar(255)，配置值
	LastValue string    `gorm:"column:last_value"` // varchar(255)，上一次配置值
	Release   int64     `gorm:"column:rls"`        // 配置项被设置计数批次
}
