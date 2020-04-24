package schema

// schema 模块不要引入官方库以外的其它模块或内部模块
import (
	"time"
)

// Product 详见 ./sql/schema.sql table `urbs_product`
// 产品线
type Product struct {
	ID        int64      `gorm:"column:id" json:"-"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deletedAt"` // 删除时间，用于灰度管理
	OfflineAt *time.Time `gorm:"column:offline_at" json:"offlineAt"` // 下线时间，用于灰度管理
	Name      string     `gorm:"column:name" json:"name"`            // varchar(63) 产品线名称，表内唯一
	Desc      string     `gorm:"column:description" json:"desc"`     // varchar(1022) 产品线描述
	Status    int64      `gorm:"column:status" json:"status"`        // -1 下线弃用，未使用
}

// TableName retuns table name
func (Product) TableName() string {
	return "urbs_product"
}
