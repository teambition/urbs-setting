package model

import (
	"context"
	"database/sql"
)

// Healthz ...
type Healthz struct {
	*Model
}

// DBStats ...
func (m *Healthz) DBStats(ctx context.Context) sql.DBStats {
	return m.SQL.DBStats()
}
