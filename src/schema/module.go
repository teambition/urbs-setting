package schema

// schema 模块不要引入官方库以外的其它模块或内部模块
import (
	"time"
)

// TableModule is a table name in db.
const TableModule = "urbs_module"

// Module 详见 ./sql/schema.sql table `urbs_module`
// 产品线的功能模块
type Module struct {
	ID        int64      `db:"id" json:"-" goqu:"skipinsert"`
	CreatedAt time.Time  `db:"created_at" json:"createdAt" goqu:"skipinsert"`
	UpdatedAt time.Time  `db:"updated_at" json:"updatedAt" goqu:"skipinsert"`
	OfflineAt *time.Time `db:"offline_at" json:"offlineAt"` // 计划下线时间，用于灰度管理
	ProductID int64      `db:"product_id"`                  // 所从属的产品线 ID
	Name      string     `db:"name" json:"name"`            // varchar(63) 功能模块名称，产品线内唯一
	Desc      string     `db:"description" json:"desc"`     // varchar(1022) 功能模块描述
	Status    int64      `db:"status" json:"status"`        // -1 下线弃用，有效配置项计数（被动异步计算，非精确值）
}

// TableName retuns table name
func (Module) TableName() string {
	return "urbs_module"
}
