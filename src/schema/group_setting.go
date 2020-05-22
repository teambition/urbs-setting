package schema

// schema 模块不要引入官方库以外的其它模块或内部模块
import (
	"time"
)

// GroupSetting 详见 ./sql/schema.sql table `group_setting`
// 记录群组对某功能模块配置项值，将作用于群组所有成员
type GroupSetting struct {
	ID        int64     `gorm:"column:id" db:"id"`
	CreatedAt time.Time `gorm:"column:created_at" db:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" db:"updated_at"`
	GroupID   int64     `gorm:"column:group_id" db:"group_id"`     // 群组内部 ID
	SettingID int64     `gorm:"column:setting_id" db:"setting_id"` // 配置项内部 ID
	Value     string    `gorm:"column:value" db:"value"`           // varchar(255)，配置值
	LastValue string    `gorm:"column:last_value" db:"last_value"` // varchar(255)，上一次配置值
	Release   int64     `gorm:"column:rls" db:"rls"`               // 配置项被设置计数批次
}
