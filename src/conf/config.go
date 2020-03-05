package conf

import (
	"time"

	"github.com/teambition/urbs-setting/src/util"
)

func init() {
	p := &Config
	util.ReadConfig(p)
}

// Logger logger config
type Logger struct {
	Level string `json:"level" yaml:"level"`
}

// SQL ...
type SQL struct {
	Host         string `json:"host" yaml:"host"`
	User         string `json:"user" yaml:"user"`
	Password     string `json:"password" yaml:"password"`
	Database     string `json:"database" yaml:"database"`
	Parameters   string `json:"parameters" yaml:"parameters"`
	MaxIdleConns int    `json:"max_idle_conns" yaml:"max_idle_conns"`
	MaxOpenConns int    `json:"max_open_conns" yaml:"max_open_conns"`
}

// ConfigTpl ...
type ConfigTpl struct {
	SrvAddr          string        `json:"addr" yaml:"addr"`
	CertFile         string        `json:"cert_file" yaml:"cert_file"`
	KeyFile          string        `json:"key_file" yaml:"key_file"`
	Logger           Logger        `json:"logger" yaml:"logger"`
	MySQL            SQL           `json:"mysql" yaml:"mysql"`
	CacheLabelExpire time.Duration `json:"cache_label_expire" yaml:"cache_label_expire"`
	Channels         []string      `json:"channels" yaml:"channels"`
	Clients          []string      `json:"clients" yaml:"clients"`
	HIDKey           string        `json:"hid_key" yaml:"hid_key"`
	AuthKeys         []string      `json:"auth_keys" yaml:"auth_keys"`
}

// IsCacheLabelExpired 判断用户缓存的 labels 是否超过有效期
func (c *ConfigTpl) IsCacheLabelExpired(now, activeAt int64) bool {
	return now-activeAt > int64(c.CacheLabelExpire/time.Second)
}

// Config ...
var Config ConfigTpl
