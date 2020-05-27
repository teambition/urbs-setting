package schema

// schema 模块不要引入官方库以外的其它模块或内部模块
import (
	"time"
)

// TableGroupLabel is a table name in db.
const TableGroupLabel = "group_label"

// GroupLabel 详见 ./sql/schema.sql table `group_label`
// 记录群组被设置的环境标签，将作用于群组所有成员
type GroupLabel struct {
	ID        int64     `db:"id" goqu:"skipinsert"`
	CreatedAt time.Time `db:"created_at" goqu:"skipinsert"`
	GroupID   int64     `db:"group_id"` // 群组内部 ID
	LabelID   int64     `db:"label_id"` // 环境标签内部 ID
	Release   int64     `db:"rls"`      // 标签被设置计数批次
}
