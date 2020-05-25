package schema

// schema 模块不要引入官方库以外的其它模块或内部模块
import (
	"time"
)

// TableUserSetting is a table name in db.
const TableUserSetting = "user_setting"

// UserSetting 详见 ./sql/schema.sql table `user_setting`
// 记录用户对某功能模块配置项值
type UserSetting struct {
	ID        int64     `db:"id" goqu:"skipinsert"`
	CreatedAt time.Time `db:"created_at" goqu:"skipinsert"`
	UpdatedAt time.Time `db:"updated_at" goqu:"skipinsert"`
	UserID    int64     `db:"user_id"`    // 用户内部 ID
	SettingID int64     `db:"setting_id"` // 配置项内部 ID
	Value     string    `db:"value"`      // varchar(255)，配置值
	LastValue string    `db:"last_value"` // varchar(255)，上一次配置值
	Release   int64     `db:"rls"`        // 配置项被设置计数批次
}
