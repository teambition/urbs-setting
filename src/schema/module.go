package schema

import (
	"time"
)

// Module 详见 ./sql/schema.sql table `module`
// 产品线的功能模块
type Module struct {
	ID        int64      `gorm:"column:id" json:"id"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updated_at"`
	OfflineAt *time.Time `gorm:"column:offline_at" json:"offline_at"` // 计划下线时间，用于灰度管理
	ProductID int64      `gorm:"column:product_id" json:"product_id"` // 所从属的产品线 ID
	Name      string     `gorm:"column:name" json:"name"`             // varchar(63) 功能模块名称，产品线内唯一
	Desc      string     `gorm:"column:desc" json:"desc"`             // varchar(1023) 功能模块描述
	Status    int64      `gorm:"column:status" json:"status"`         // -1 弃用，0 未使用，大于 0 为有效配置项数
}
