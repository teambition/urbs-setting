package schema

import (
	"time"
)

// UserSetting 详见 ./sql/schema.sql table `user_setting`
// 记录用户对某功能模块配置项值
type UserSetting struct {
	ID        int64     `gorm:"column:id" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
	UserID    int64     `gorm:"column:user_id" json:"user_id"`       // 用户内部 ID
	SettingID int64     `gorm:"column:setting_id" json:"setting_id"` // 配置项内部 ID
	Channel   string    `gorm:"column:channel" json:"channel"`       // varchar(31)，产品版本更新通道
	Client    string    `gorm:"column:client" json:"client"`         // varchar(31)，产品客户端类型
	Value     string    `gorm:"column:value" json:"value"`           // varchar(255)，配置值
	LastValue string    `gorm:"column:last_value" json:"last_value"` // varchar(255)，上一次配置值
}
