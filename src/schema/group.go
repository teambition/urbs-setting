package schema

// schema 模块不要引入官方库以外的其它模块或内部模块
import (
	"time"
)

// Group 详见 ./sql/schema.sql table `group`
// 用户群组
type Group struct {
	ID        int64     `gorm:"column:id" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
	SyncAt    int64     `gorm:"column:sync_at" json:"sync_at"` // 群组成员同步时间点
	UID       string    `gorm:"column:uid" json:"uid"`         // varchar(63)，群组外部ID，表内唯一， 如 Teambition organization id
	Desc      string    `gorm:"column:desc" json:"desc"`       // varchar(1022)，群组描述
}
