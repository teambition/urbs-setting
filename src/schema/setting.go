package schema

// schema 模块不要引入官方库以外的其它模块或内部模块
import (
	"time"
)

// Setting 详见 ./sql/schema.sql table `setting`
// 功能模块的配置项
type Setting struct {
	ID        int64      `gorm:"column:id" json:"id"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updated_at"`
	OfflineAt *time.Time `gorm:"column:offline_at" json:"offline_at"` // 计划下线时间，用于灰度管理
	ModuleID  int64      `gorm:"column:module_id" json:"module_id"`   // 配置项所从属的功能模块 ID
	Name      string     `gorm:"column:name" json:"name"`             // varchar(63) 配置项名称，功能模块内唯一
	Desc      string     `gorm:"column:desc" json:"desc"`             // varchar(1022) 配置项描述信息
	Channels  string     `gorm:"column:channels" json:"channels"`     // varchar(255) 配置项适用的版本通道，未配置表示都适用
	Clients   string     `gorm:"column:clients" json:"clients"`       // varchar(255) 配置项适用的客户端类型，未配置表示都适用
	Values    string     `gorm:"column:values" json:"values"`         // varchar(1022) 配置项可选值集合
	Status    int64      `gorm:"column:status" json:"status"`         // -1 下线弃用，0 未使用，大于 0 为被使用计数
}