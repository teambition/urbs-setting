package service

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/doug-martin/goqu/v9"
	mysqlDialect "github.com/doug-martin/goqu/v9/dialect/mysql"
	_ "github.com/go-sql-driver/mysql" // go-sql-driver
	"github.com/teambition/urbs-setting/src/conf"
	"github.com/teambition/urbs-setting/src/logging"
	"github.com/teambition/urbs-setting/src/util"
)

func init() {
	util.DigProvide(NewDB)
	goqu.RegisterDialect("default", mysqlDialect.DialectOptions()) // make mysql dialect as default too.
}

// SQL ...
type SQL struct {
	db   *sql.DB
	DB   *goqu.Database
	RdDB *goqu.Database
}

// DBStats ...
func (s *SQL) DBStats() sql.DBStats {
	return s.db.Stats()
}

// NewDB ...
func NewDB() *SQL {
	db := connectDB(conf.Config.MySQL)
	rdDB := db
	if conf.Config.MySQLRd.Host != "" {
		rdDB = connectDB(conf.Config.MySQLRd)
	}

	dialect := goqu.Dialect("mysql")
	return &SQL{
		db:   db,
		DB:   dialect.DB(db),
		RdDB: dialect.DB(rdDB),
	}
}

func connectDB(cfg conf.SQL) *sql.DB {
	if cfg.MaxIdleConns <= 0 {
		cfg.MaxIdleConns = 8
	}

	if cfg.MaxOpenConns <= 0 {
		cfg.MaxOpenConns = 64
	}

	if cfg.User == "" || cfg.Password == "" || cfg.Host == "" {
		logging.Panicf("Invalid SQL DB config %s:%s@(%s)/%s", cfg.User, cfg.Password, cfg.Host, cfg.Database)
	}

	parameters, err := url.ParseQuery(cfg.Parameters)
	if err != nil {
		logging.Panicf("Invalid SQL DB parameters %s", cfg.Parameters)
	}
	// 强制使用
	parameters.Set("collation", "utf8mb4_general_ci")
	parameters.Set("parseTime", "true")

	// https://github.com/go-sql-driver/mysql#parameters
	url := fmt.Sprintf(`%s:%s@(%s)/%s?%s`, cfg.User, cfg.Password, cfg.Host, cfg.Database, parameters.Encode())
	db, err := sql.Open("mysql", url)
	if err == nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err = db.PingContext(ctx)
		cancel()
	}
	if err != nil {
		url = strings.Replace(url, cfg.Password, cfg.Password[0:4]+"***", 1)
		logging.Panicf("SQL DB connect failed %s, with config %s", err, url)
	}

	// SetMaxIdleCons 设置连接池中的最大闲置连接数。
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	// SetMaxOpenCons 设置数据库的最大连接数量。
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	// SetConnMaxLifetiment 设置连接的最大可复用时间。
	// db.SetConnMaxLifetime(time.Hour)
	return db
}

// DeResult ...
func DeResult(re sql.Result, err error) (int64, error) {
	if err != nil {
		return 0, err
	}
	return re.RowsAffected()
}
