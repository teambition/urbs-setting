package service

import (
	"database/sql"
	"fmt"
	"net/url"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // go-sql-driver
	"github.com/teambition/urbs-setting/src/conf"
	"github.com/teambition/urbs-setting/src/logging"
	"github.com/teambition/urbs-setting/src/util"
)

func init() {
	util.DigProvide(NewDB)
}

// SQL ...
type SQL struct {
	DB *gorm.DB
}

// DBStats ...
func (s *SQL) DBStats() sql.DBStats {
	return s.DB.DB().Stats()
}

// NewDB ...
func NewDB() *SQL {
	cfg := conf.Config.MySQL

	if cfg.MaxIdleConns <= 0 {
		cfg.MaxIdleConns = 8
	}

	if cfg.MaxOpenConns <= 0 {
		cfg.MaxOpenConns = 64
	}

	if cfg.User == "" || cfg.Password == "" || cfg.Host == "" {
		logging.Panic(logging.SrvLog("Invalid SQL DB config %s:%s@(%s)/%s", cfg.User, cfg.Password, cfg.Host, cfg.Database))
	}

	parameters, err := url.ParseQuery(cfg.Parameters)
	if err != nil {
		logging.Panic(logging.SrvLog("Invalid SQL DB parameters %s", cfg.Parameters))
	}
	// 强制使用
	parameters.Set("collation", "utf8mb4_general_ci")
	parameters.Set("parseTime", "true")

	// https://github.com/go-sql-driver/mysql#parameters
	url := fmt.Sprintf(`%s:%s@(%s)/%s?%s`, cfg.User, cfg.Password, cfg.Host, cfg.Database, parameters.Encode())
	db, err := gorm.Open("mysql", url)

	if err != nil {
		url = strings.Replace(url, cfg.Password, cfg.Password[0:4]+"***", 1)
		logging.Panic(logging.SrvLog("SQL DB connect failed %s, with config %s", err, url))
	}

	// 表名使用单数。
	// https://the.agilesql.club/2019/05/should-i-pluralize-table-names-is-it-person-persons-people-or-people/
	db.SingularTable(true)
	// SetMaxIdleCons 设置连接池中的最大闲置连接数。
	db.DB().SetMaxIdleConns(cfg.MaxIdleConns)
	// SetMaxOpenCons 设置数据库的最大连接数量。
	db.DB().SetMaxOpenConns(cfg.MaxOpenConns)
	// SetConnMaxLifetiment 设置连接的最大可复用时间。
	// db.DB().SetConnMaxLifetime(time.Hour)

	return &SQL{
		DB: db,
	}
}
