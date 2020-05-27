package schema

// schema 模块不要引入官方库以外的其它模块或内部模块
import (
	"time"
)

// TableGroupSetting is a table name in db.
const TableGroupSetting = "group_setting"

// GroupSetting 详见 ./sql/schema.sql table `group_setting`
// 记录群组对某功能模块配置项值，将作用于群组所有成员
type GroupSetting struct {
	ID        int64     `db:"id" goqu:"skipinsert"`
	CreatedAt time.Time `db:"created_at" goqu:"skipinsert"`
	UpdatedAt time.Time `db:"updated_at" goqu:"skipinsert"`
	GroupID   int64     `db:"group_id"`   // 群组内部 ID
	SettingID int64     `db:"setting_id"` // 配置项内部 ID
	Value     string    `db:"value"`      // varchar(255)，配置值
	LastValue string    `db:"last_value"` // varchar(255)，上一次配置值
	Release   int64     `db:"rls"`        // 配置项被设置计数批次
}
