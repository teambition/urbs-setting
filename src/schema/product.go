package schema

// schema 模块不要引入官方库以外的其它模块或内部模块
import (
	"time"
)

// TableProduct is a table name in db.
const TableProduct = "urbs_product"

// Product 详见 ./sql/schema.sql table `urbs_product`
// 产品线
type Product struct {
	ID        int64      `db:"id" json:"-" goqu:"skipinsert"`
	CreatedAt time.Time  `db:"created_at" json:"createdAt" goqu:"skipinsert"`
	UpdatedAt time.Time  `db:"updated_at" json:"updatedAt" goqu:"skipinsert"`
	DeletedAt *time.Time `db:"deleted_at" json:"deletedAt"` // 删除时间，用于灰度管理
	OfflineAt *time.Time `db:"offline_at" json:"offlineAt"` // 下线时间，用于灰度管理
	Name      string     `db:"name" json:"name"`            // varchar(63) 产品线名称，表内唯一
	Desc      string     `db:"description" json:"desc"`     // varchar(1022) 产品线描述
	Status    int64      `db:"status" json:"status"`        // -1 下线弃用，未使用
}

// TableName retuns table name
func (Product) TableName() string {
	return "urbs_product"
}
