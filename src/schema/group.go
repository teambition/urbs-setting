package schema

// schema 模块不要引入官方库以外的其它模块或内部模块
import (
	"time"
)

// TableGroup is a table name in db.
const TableGroup = "urbs_group"

// Group 详见 ./sql/schema.sql table `urbs_group`
// 用户群组
type Group struct {
	ID        int64     `db:"id" json:"-" goqu:"skipinsert"`
	CreatedAt time.Time `db:"created_at" json:"createdAt" goqu:"skipinsert"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt" goqu:"skipinsert"`
	SyncAt    int64     `db:"sync_at" json:"syncAt"`            // 群组成员同步时间点，1970 以来的秒数
	UID       string    `db:"uid" json:"uid"`                   // varchar(63)，群组外部ID，表内唯一， 如 Teambition organization id
	Kind      string    `db:"kind" json:"kind"`                 // varchar(63)，群组外部ID，表内唯一， 如 Teambition organization id
	Desc      string    `db:"description" json:"desc"`          // varchar(1022)，群组描述
	Status    int64     `db:"status" json:"status" db:"status"` // 成员计数（被动异步计算，非精确值）
}

// TableName retuns table name
func (Group) TableName() string {
	return "urbs_group"
}
