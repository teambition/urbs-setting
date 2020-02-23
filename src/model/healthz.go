package model

import (
	"context"
	"database/sql"

	"github.com/jinzhu/gorm"
)

type healthz struct {
	DB *gorm.DB
}

func (m *healthz) DBStats(ctx context.Context) sql.DBStats {
	return m.DB.DB().Stats()
}
