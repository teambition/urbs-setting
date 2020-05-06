package schema

import "time"

// schema 模块不要引入官方库以外的其它模块或内部模块

// Lock 详见 ./sql/schema.sql table `urbs_lock`
// 内部锁
type Lock struct {
	ID       int64     `gorm:"column:id"`
	Name     string    `gorm:"column:name"` // varchar(255) 锁键，表内唯一
	ExpireAt time.Time `gorm:"column:expire_at"`
}

// TableName retuns table name
func (Lock) TableName() string {
	return "urbs_lock"
}
