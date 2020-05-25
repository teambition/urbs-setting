package schema

import "time"

// schema 模块不要引入官方库以外的其它模块或内部模块

// TableStatistic is a table name in db.
const TableStatistic = "urbs_statistic"

// Statistic 详见 ./sql/schema.sql table `urbs_statistic`
// 内部统计
type Statistic struct {
	ID        int64     `db:"id" goqu:"skipinsert"`
	CreatedAt time.Time `db:"created_at" goqu:"skipinsert"`
	UpdatedAt time.Time `db:"updated_at" goqu:"skipinsert"`
	Name      string    `db:"name"` // varchar(255) 锁键，表内唯一
	Status    int64     `db:"status"`
	Value     string    `db:"value"` // varchar(8190) json value
}

// TableName retuns table name
func (Statistic) TableName() string {
	return "urbs_statistic"
}

// StatisticKey ...
type StatisticKey string

// UsersTotalSize ...
const (
	UsersTotalSize        StatisticKey = "UsersTotalSize"
	GroupsTotalSize       StatisticKey = "GroupsTotalSize"
	ProductsTotalSize     StatisticKey = "ProductsTotalSize"
	LabelsTotalSize       StatisticKey = "LabelsTotalSize"
	ModulesTotalSize      StatisticKey = "ModulesTotalSize"
	SettingsTotalSize     StatisticKey = "SettingsTotalSize"
	LabelRulesTotalSize   StatisticKey = "LabelRulesTotalSize"
	SettingRulesTotalSize StatisticKey = "SettingRulesTotalSize"
)
