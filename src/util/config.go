package util

// util 模块不要引入其它内部模块
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	yaml "gopkg.in/yaml.v2"
)

var once sync.Once

// ReadConfig 指定配置文件并解析; 未指定配置文件则通过环境变量获取
func ReadConfig(v interface{}, path ...string) {
	once.Do(func() {
		filePath, err := getConfigFilePath(path...)
		if err != nil {
			panic(err)
		}

		ext := filepath.Ext(filePath)

		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			panic(err)
		}

		err = parseConfig(data, ext, v)
		if err != nil {
			panic(err)
		}
	})
}

func getConfigFilePath(path ...string) (string, error) {
	// 优先使用的环境变量
	filePath := os.Getenv("CONFIG_FILE_PATH")

	// 或使用指定的路径
	if filePath == "" && len(path) > 0 {
		filePath = path[0]
	}

	if filePath == "" {
		return "", fmt.Errorf("config file not specified")
	}

	return filePath, nil
}

type unmarshaler func(data []byte, v interface{}) error

func parseConfig(data []byte, ext string, v interface{}) error {
	ext = strings.TrimLeft(ext, ".")

	var unmarshal unmarshaler

	switch ext {
	case "json":
		unmarshal = json.Unmarshal
	case "yaml", "yml":
		unmarshal = yaml.Unmarshal
	default:
		return fmt.Errorf("not supported config ext: %s", ext)
	}

	return unmarshal(data, v)
}
