package model

import (
	"context"
	"database/sql"

	"github.com/jinzhu/gorm"
)

// Healthz ...
type Healthz struct {
	DB *gorm.DB
}

// DBStats ...
func (m *Healthz) DBStats(ctx context.Context) sql.DBStats {
	return m.DB.DB().Stats()
}
