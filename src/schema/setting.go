package schema

// schema 模块不要引入官方库以外的其它模块或内部模块
import (
	"time"
)

// TableSetting is a table name in db.
const TableSetting = "urbs_setting"

// Setting 详见 ./sql/schema.sql table `urbs_setting`
// 功能模块的配置项
type Setting struct {
	ID        int64      `db:"id" goqu:"skipinsert"`
	CreatedAt time.Time  `db:"created_at" goqu:"skipinsert"`
	UpdatedAt time.Time  `db:"updated_at" goqu:"skipinsert"`
	OfflineAt *time.Time `db:"offline_at"`               // 计划下线时间，用于灰度管理
	ModuleID  int64      `db:"module_id"`                // 配置项所从属的功能模块 ID
	Module    string     `db:"module" goqu:"skipinsert"` // 仅为查询方便追加字段，数据库中没有该字段
	Name      string     `db:"name"`                     // varchar(63) 配置项名称，功能模块内唯一
	Desc      string     `db:"description"`              // varchar(1022) 配置项描述信息
	Channels  string     `db:"channels"`                 // varchar(255) 配置项适用的版本通道，未配置表示都适用
	Clients   string     `db:"clients"`                  // varchar(255) 配置项适用的客户端类型，未配置表示都适用
	Values    string     `db:"vals"`                     // varchar(1022) 配置项可选值集合
	Status    int64      `db:"status"`                   // -1 下线弃用，使用用户计数（被动异步计算，非精确值）
	Release   int64      `db:"rls"`                      // 配置项发布（被设置）计数
}

// TableName retuns table name
func (Setting) TableName() string {
	return "urbs_setting"
}
