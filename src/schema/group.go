package schema

// schema 模块不要引入官方库以外的其它模块或内部模块
import (
	"time"
)

// Group 详见 ./sql/schema.sql table `urbs_group`
// 用户群组
type Group struct {
	ID        int64     `gorm:"column:id" db:"id" json:"-"`
	CreatedAt time.Time `gorm:"column:created_at" db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at" db:"updated_at" json:"updatedAt"`
	SyncAt    int64     `gorm:"column:sync_at" db:"sync_at" json:"syncAt"`           // 群组成员同步时间点，1970 以来的秒数
	UID       string    `gorm:"column:uid" db:"uid" json:"uid"`                      // varchar(63)，群组外部ID，表内唯一， 如 Teambition organization id
	Kind      string    `gorm:"column:kind" db:"kind" json:"kind"`                   // varchar(63)，群组外部ID，表内唯一， 如 Teambition organization id
	Desc      string    `gorm:"column:description" db:"description" json:"desc"`     // varchar(1022)，群组描述
	Status    int64     `gorm:"column:status" db:"status" json:"status" db:"status"` // 成员计数（被动异步计算，非精确值）
}

// TableName retuns table name
func (Group) TableName() string {
	return "urbs_group"
}
