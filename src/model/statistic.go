package model

import (
	"context"

	"github.com/teambition/urbs-setting/src/schema"
)

// Statistic ...
type Statistic struct {
	*Model
}

// FindByKey ...
func (m *Statistic) FindByKey(ctx context.Context, key schema.StatisticKey) (*schema.Statistic, error) {
	statistic := &schema.Statistic{}
	if err := m.DB.Where("`key` =  ?", key).First(&statistic).Error; err != nil {
		return nil, err
	}
	return statistic, nil
}

// FindByKeys ...
func (m *Statistic) FindByKeys(ctx context.Context, keys []schema.StatisticKey) ([]schema.Statistic, error) {
	statistics := make([]schema.Statistic, 0)
	err := m.DB.Where("`key` in ( ? )", keys).Find(&statistics).Error
	return statistics, err
}
