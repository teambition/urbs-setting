package schema

// schema 模块不要引入官方库以外的其它模块或内部模块
import (
	"time"
)

// TableUserGroup is a table name in db.
const TableUserGroup = "user_group"

// UserGroup 详见 ./sql/schema.sql table `user_group`
// 记录用户从属的群组，用户可以归属到多个群组
// 用户能从所归属的群组继承灰度标签和功能项配置，也就是基于群组进行灰度
type UserGroup struct {
	ID        int64     `db:"id" goqu:"skipinsert"`
	CreatedAt time.Time `db:"created_at" goqu:"skipinsert"`
	SyncAt    int64     `db:"sync_at"`  // 归属关系同步时间戳，1970 以来的秒数，应该与 group.sync_at 相等
	UserID    int64     `db:"user_id"`  // 用户内部 ID
	GroupID   int64     `db:"group_id"` // 群组内部 ID
}
