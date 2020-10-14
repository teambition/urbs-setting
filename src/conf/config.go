package conf

import (
	"context"
	"time"

	otgo "github.com/open-trust/ot-go-lib"
	"github.com/teambition/gear"
	"github.com/teambition/urbs-setting/src/util"
)

func init() {
	p := &Config
	util.ReadConfig(p)
	if err := p.Validate(); err != nil {
		panic(err)
	}
	p.GlobalCtx = gear.ContextWithSignal(context.Background())
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

// OpenTrust ...
type OpenTrust struct {
	OTID             otgo.OTID `json:"otid" yaml:"otid"`
	LegacyOTID       otgo.OTID `json:"legacy_otid" yaml:"legacy_otid"`
	PrivateKeys      []string  `json:"private_keys" yaml:"private_keys"`
	DomainPublicKeys []string  `json:"domain_public_keys" yaml:"domain_public_keys"`
}

// ConfigTpl ...
type ConfigTpl struct {
	GlobalCtx        context.Context
	SrvAddr          string    `json:"addr" yaml:"addr"`
	CertFile         string    `json:"cert_file" yaml:"cert_file"`
	KeyFile          string    `json:"key_file" yaml:"key_file"`
	Logger           Logger    `json:"logger" yaml:"logger"`
	MySQL            SQL       `json:"mysql" yaml:"mysql"`
	MySQLRd          SQL       `json:"mysql_read" yaml:"mysql_read"`
	CacheLabelExpire string    `json:"cache_label_expire" yaml:"cache_label_expire"`
	Channels         []string  `json:"channels" yaml:"channels"`
	Clients          []string  `json:"clients" yaml:"clients"`
	HIDKey           string    `json:"hid_key" yaml:"hid_key"`
	AuthKeys         []string  `json:"auth_keys" yaml:"auth_keys"`
	OpenTrust        OpenTrust `json:"open_trust" yaml:"open_trust"`
	cacheLabelExpire int64     // seconds, default to 60 seconds
}

// Validate 用于完成基本的配置验证和初始化工作。业务相关的配置验证建议放到相关代码中实现，如 mysql 的配置。
func (c *ConfigTpl) Validate() error {
	du, err := time.ParseDuration(c.CacheLabelExpire)
	if err != nil {
		return err
	}
	if du < time.Minute {
		du = time.Minute
	}
	c.cacheLabelExpire = int64(du / time.Second)
	return nil
}

// IsCacheLabelExpired 判断用户缓存的 labels 是否超过有效期
func (c *ConfigTpl) IsCacheLabelExpired(now, activeAt int64) bool {
	return now-activeAt > c.cacheLabelExpire
}

// Config ...
var Config ConfigTpl
