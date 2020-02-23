package service

import (
	"database/sql"
	"fmt"

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

	if cfg.User == "" || cfg.Password == "" || cfg.Host == "" || cfg.Database == "" {
		logging.Panic(logging.SrvLog("Invalid SQL DB config %s:%s***@(%s)/%s", cfg.User, cfg.Password[:4], cfg.Host, cfg.Database))
	}

	// https://github.com/go-sql-driver/mysql#parameters
	url := fmt.Sprintf(`%s:%s@(%s)/%s?collation=utf8mb4_general_ci&parseTime=true&loc=UTC&readTimeout=10s&writeTimeout=10s&timeout=10s`, cfg.User, cfg.Password, cfg.Host, cfg.Database)
	db, err := gorm.Open("mysql", url)

	if err != nil {
		logging.Panic(logging.SrvLog("SQL DB connect failed with config %s:%s***@(%s)/%s, msg: %s", cfg.User, cfg.Password[:4], cfg.Host, cfg.Database, err))
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
