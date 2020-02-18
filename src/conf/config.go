package conf

import (
	"github.com/teambition/urbs-setting/src/util"
)

func init() {
	util.ReadConfig(&Config, "/etc/urbs-setting/config.yml")
}

// Logger logger config
type Logger struct {
	Level string `json:"level" yaml:"level"`
}

// ConfigTpl ...
type ConfigTpl struct {
	Logger   Logger `json:"logger" yaml:"logger"`
	SrvAddr  string `json:"addr" yaml:"addr"`
	CertFile string `json:"cert_file" yaml:"cert_file"`
	KeyFile  string `json:"key_file" yaml:"key_file"`
}

// Config ...
var Config ConfigTpl
