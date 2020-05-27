package schema

// schema 模块不要引入官方库以外的其它模块或内部模块
import (
	"time"
)

// TableLabel is a table name in db.
const TableLabel = "urbs_label"

// Label 详见 ./sql/schema.sql table `urbs_label`
// 环境标签
type Label struct {
	ID        int64      `db:"id" goqu:"skipinsert"`
	CreatedAt time.Time  `db:"created_at" goqu:"skipinsert"`
	UpdatedAt time.Time  `db:"updated_at" goqu:"skipinsert"`
	OfflineAt *time.Time `db:"offline_at"`  // 计划下线时间，用于灰度管理
	ProductID int64      `db:"product_id"`  // 所从属的产品线 ID
	Name      string     `db:"name"`        // varchar(63) 环境标签名称，产品线内唯一
	Desc      string     `db:"description"` // varchar(1022) 环境标签描述
	Channels  string     `db:"channels"`    // varchar(255) 标签适用的版本通道，未配置表示都适用
	Clients   string     `db:"clients"`     // varchar(255) 标签适用的客户端类型，未配置表示都适用
	Status    int64      `db:"status"`      // -1 下线弃用，使用用户计数（被动异步计算，非精确值）
	Release   int64      `db:"rls"`         // 标签发布（被设置）计数
}

// TableName retuns table name
func (Label) TableName() string {
	return "urbs_label"
}
