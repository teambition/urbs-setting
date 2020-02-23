package conf

import (
	"time"

	"github.com/teambition/urbs-setting/src/util"
)

func init() {
	p := &Config
	util.ReadConfig(p, "/etc/urbs-setting/config.yml")
	p.init()
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
	MaxIdleConns int    `json:"max_idle_conns" yaml:"max_idle_conns"`
	MaxOpenConns int    `json:"max_open_conns" yaml:"max_open_conns"`
}

// ConfigTpl ...
type ConfigTpl struct {
	SrvAddr           string   `json:"addr" yaml:"addr"`
	CertFile          string   `json:"cert_file" yaml:"cert_file"`
	KeyFile           string   `json:"key_file" yaml:"key_file"`
	Logger            Logger   `json:"logger" yaml:"logger"`
	MySQL             SQL      `json:"mysql" yaml:"mysql"`
	CacheLabelExpireS string   `json:"cache_label_expire" yaml:"cache_label_expire"`
	CacheLabelExpire  int64    `json:"-" yaml:"-"` // 从 cache_label_expire 解析的值，seconds
	Channels          []string `json:"channels" yaml:"channels"`
	Clients           []string `json:"clients" yaml:"clients"`
	HIDKey            string   `json:"hid_key" yaml:"hid_key"`
}

func (c *ConfigTpl) init() {
	i, err := time.ParseDuration(c.CacheLabelExpireS)
	if err != nil || i < time.Second*10 {
		i = time.Second * 10
	}
	c.CacheLabelExpire = int64(i / time.Second)
}

// Config ...
var Config ConfigTpl
