package schema

import (
	"time"
)

// Product 详见 ./sql/schema.sql table `product`
// 产品线
type Product struct {
	ID        int64      `gorm:"column:id" json:"id"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"` // 删除时间，用于灰度管理
	Name      string     `gorm:"column:name" json:"name"`             // varchar(63) 产品线名称，表内唯一
	Desc      string     `gorm:"column:desc" json:"desc"`             // varchar(1023) 产品线描述
	Status    int64      `gorm:"column:status" json:"status"`         // -1 弃用，0 未使用，大于 0 为有效功能模块数
}
