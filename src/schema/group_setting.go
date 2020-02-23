package schema

import (
	"time"
)

// GroupSetting 详见 ./sql/schema.sql table `group_setting`
// 记录群组对某功能模块配置项值，将作用于群组所有成员
type GroupSetting struct {
	ID        int64     `gorm:"column:id" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
	GroupID   int64     `gorm:"column:group_id" json:"group_id"`     // 群组内部 ID
	SettingID int64     `gorm:"column:setting_id" json:"setting_id"` // 配置项内部 ID
	Value     string    `gorm:"column:value" json:"value"`           // varchar(255)，配置值
	LastValue string    `gorm:"column:last_value" json:"last_value"` // varchar(255)，上一次配置值
}
