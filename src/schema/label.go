package schema

// schema 模块不要引入官方库以外的其它模块或内部模块
import (
	"time"
)

// Label 详见 ./sql/schema.sql table `urbs_label`
// 灰度标签
type Label struct {
	ID        int64      `gorm:"column:id"`
	CreatedAt time.Time  `gorm:"column:created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at"`
	OfflineAt *time.Time `gorm:"column:offline_at"`  // 计划下线时间，用于灰度管理
	ProductID int64      `gorm:"column:product_id"`  // 所从属的产品线 ID
	Name      string     `gorm:"column:name"`        // varchar(63) 灰度标签名称，产品线内唯一
	Desc      string     `gorm:"column:description"` // varchar(1022) 灰度标签描述
	Channels  string     `gorm:"column:channels"`    // varchar(255) 标签适用的版本通道，未配置表示都适用
	Clients   string     `gorm:"column:clients"`     // varchar(255) 标签适用的客户端类型，未配置表示都适用
	Status    int64      `gorm:"column:status"`      // -1 下线弃用，0 未使用，大于 0 为被使用计数
}

// TableName retuns table name
func (Label) TableName() string {
	return "urbs_label"
}
