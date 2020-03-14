package schema

// schema 模块不要引入官方库以外的其它模块或内部模块
import (
	"time"
)

// Setting 详见 ./sql/schema.sql table `urbs_setting`
// 功能模块的配置项
type Setting struct {
	ID        int64      `gorm:"column:id"`
	CreatedAt time.Time  `gorm:"column:created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at"`
	OfflineAt *time.Time `gorm:"column:offline_at"`  // 计划下线时间，用于灰度管理
	ModuleID  int64      `gorm:"column:module_id"`   // 配置项所从属的功能模块 ID
	Name      string     `gorm:"column:name"`        // varchar(63) 配置项名称，功能模块内唯一
	Desc      string     `gorm:"column:description"` // varchar(1022) 配置项描述信息
	Channels  string     `gorm:"column:channels"`    // varchar(255) 配置项适用的版本通道，未配置表示都适用
	Clients   string     `gorm:"column:clients"`     // varchar(255) 配置项适用的客户端类型，未配置表示都适用
	Values    string     `gorm:"column:vals"`        // varchar(1022) 配置项可选值集合
	Status    int64      `gorm:"column:status"`      // -1 下线弃用，0 未使用，大于 0 为被使用计数
}

// TableName retuns table name
func (Setting) TableName() string {
	return "urbs_setting"
}
