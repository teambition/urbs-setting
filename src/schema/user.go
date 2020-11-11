package schema

// schema 模块不要引入官方库以外的其它模块或内部模块
import (
	"encoding/json"
	"time"
)

// TableUser is a table name in db.
const TableUser = "urbs_user"

// User 详见 ./sql/schema.sql table `urbs_user`
// 记录用户外部唯一 ID，uid 和最近活跃时间
// 缓存用户当前全部 label，根据 active_at 和 cache_label_expire 刷新
// labels 格式：TODO
type User struct {
	ID        int64     `db:"id" json:"-" goqu:"skipinsert"`
	CreatedAt time.Time `db:"created_at" json:"createdAt" goqu:"skipinsert"`
	UID       string    `db:"uid" json:"uid"`            // varchar(63)，用户外部ID，表内唯一， 如 Teambition user id
	ActiveAt  int64     `db:"active_at" json:"activeAt"` // 最近活跃时间戳，1970 以来的秒数，但不及时更新
	Labels    string    `db:"labels" json:"labels"`      // varchar(8190)，缓存用户当前被设置的 labels
}

// GetUsersUID 返回 users 数组的 uid 数组
func GetUsersUID(users []User) []string {
	uids := make([]string, len(users))
	for i, u := range users {
		uids[i] = u.UID
	}
	return uids
}

// TableName retuns table name
func (User) TableName() string {
	return "urbs_user"
}

// MyLabelInfo ...
type MyLabelInfo struct {
	ID        int64     `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	Name      string    `db:"name"`
	Channels  string    `db:"channels"`
	Clients   string    `db:"clients"`
	Product   string    `db:"product"`
}

// UserCache 用于在 User 数据上缓存数据
type UserCache struct {
	ActiveAt int64            `json:"activeAt"` // 最近活跃时间戳，1970 以来的秒数，但不及时更新
	Labels   []UserCacheLabel `json:"labels"`
}

// UserCacheLabel 用于在 User 数据上缓存 labels
type UserCacheLabel struct {
	Label    string   `json:"l"`
	Clients  []string `json:"cls,omitempty"`
	Channels []string `json:"chs,omitempty"`
}

// UserCacheLabelMap 用于在 User 数据上缓存
type UserCacheLabelMap map[string]*UserCache

// GetCache 从 user 上读取结构化的缓存数据
func (u *User) GetCache(product string) *UserCache {
	userCache := &UserCache{}
	if u.Labels == "" {
		return userCache
	}
	data := u.GetCacheMap()
	for k, ucl := range data {
		if k == product {
			return ucl
		}
	}
	return userCache
}

// GetLabels 从 user 上读取结构化的 labels 数据
func (u *User) GetLabels(product string) []UserCacheLabel {
	return u.GetCache(product).Labels
}

// GetCacheMap 从 user 上读取结构化的缓存数据
func (u *User) GetCacheMap() UserCacheLabelMap {
	data := make(UserCacheLabelMap)
	if u.Labels == "" {
		return data
	}
	_ = json.Unmarshal([]byte(u.Labels), &data)
	return data
}

// PutCacheMap 把结构化的 labels 数据转成字符串设置在 user.Labels 上
func (u *User) PutCacheMap(labels UserCacheLabelMap) error {
	data, err := json.Marshal(labels)
	if err == nil {
		u.Labels = string(data)
	}
	return err
}
