package schema

import "time"

// schema 模块不要引入官方库以外的其它模块或内部模块

// Statistic 详见 ./sql/schema.sql table `urbs_statistic`
// 内部统计
type Statistic struct {
	ID        int64     `gorm:"column:id"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
	Name      string    `gorm:"column:name"` // varchar(255) 锁键，表内唯一
	Status    int64     `gorm:"column:status"`
	Value     string    `gorm:"column:value"` // varchar(8190) json value
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
