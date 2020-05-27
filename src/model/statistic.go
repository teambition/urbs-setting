package model

import (
	"context"

	"github.com/doug-martin/goqu/v9"
	"github.com/teambition/urbs-setting/src/schema"
)

// Statistic ...
type Statistic struct {
	*Model
}

// FindByKey ...
func (m *Statistic) FindByKey(ctx context.Context, key schema.StatisticKey) (*schema.Statistic, error) {
	statistic := &schema.Statistic{}
	ok, err := m.findOneByCols(ctx, schema.TableStatistic, goqu.Ex{"name": key}, "", statistic)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}
	return statistic, nil
}
