package service

import (
	"github.com/ercole-io/ercole/v2/chart-service/dto"
	"time"
)

func (as *ChartService) GetHostCores(location, environment string, olderThan, newerThan time.Time) ([]dto.HostCores, error) {
	out, err := as.Database.GetHostCores(location, environment, olderThan, newerThan)
	if err != nil {
		return nil, err
	}

	return out, err
}