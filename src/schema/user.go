package schema

import (
	"encoding/json"
	"time"

	"github.com/teambition/urbs-setting/src/conf"
)

// User 详见 ./sql/schema.sql table `user`
// 记录用户外部唯一 ID，uid 和最近活跃时间
// 缓存用户当前全部 label，根据 active_at 和 cache_label_expire 刷新
// labels 格式：TODO
type User struct {
	ID        int64     `gorm:"column:id" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UID       string    `gorm:"column:uid" json:"uid"`             // varchar(63)，用户外部ID，表内唯一， 如 Teambition user id
	ActiveAt  int64     `gorm:"column:active_at" json:"active_at"` // 最近活跃时间戳，1970 以来的秒数，但不及时更新
	Labels    string    `gorm:"column:labels" json:"labels"`       // varchar(2047)，缓存用户当前被设置的 labels
}

// UserLabelInfo 用于在 User 数据上缓存 labels
type UserLabelInfo struct {
	Product string `json:"p"`
	Label   string `json:"l"`
	Client  string `json:"cl,omitempty"`
	Channel string `json:"ch,omitempty"`
}

// IsStale 判断用户缓存的 labels 是否超过有效期
func (u *User) IsStale(now int64) bool {
	return now-u.ActiveAt > conf.Config.CacheLabelExpire
}

// GetLabels 从 user 上读取结构化的 labels 数据
func (u *User) GetLabels() []UserLabelInfo {
	data := []UserLabelInfo{}
	if u.Labels == "" {
		return data
	}
	_ = json.Unmarshal([]byte(u.Labels), &data)
	return data
}

// PutLabels 把结构化的 labels 数据转成字符串设置在 user.Labels 上
func (u *User) PutLabels(labels []UserLabelInfo) error {
	data, err := json.Marshal(labels)
	if err == nil {
		u.Labels = string(data)
	}
	return err
}
