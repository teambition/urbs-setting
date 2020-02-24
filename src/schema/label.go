package schema

// schema 模块不要引入官方库以外的其它模块或内部模块
import (
	"time"
)

// Label 详见 ./sql/schema.sql table `label`
// 灰度标签
type Label struct {
	ID        int64      `gorm:"column:id" json:"id"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updated_at"`
	OfflineAt *time.Time `gorm:"column:offline_at" json:"offline_at"` // 计划下线时间，用于灰度管理
	ProductID int64      `gorm:"column:product_id" json:"product_id"` // 所从属的产品线 ID
	Name      string     `gorm:"column:name" json:"name"`             // varchar(63) 灰度标签名称，产品线内唯一
	Desc      string     `gorm:"column:desc" json:"desc"`             // varchar(1023) 灰度标签描述
	Status    int64      `gorm:"column:status" json:"status"`         // -1 弃用，0 未使用，大于 0 为被使用计数
}
