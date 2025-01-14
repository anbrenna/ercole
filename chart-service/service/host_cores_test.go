// Copyright (c) 2022 Sorint.lab S.p.A.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package service

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ercole-io/ercole/v2/chart-service/dto"
	"github.com/ercole-io/ercole/v2/utils"
)

func TestGetHostsHistory_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := ChartService{
		Database: db,
	}

	location := ""
	environment := ""
	olderThan := utils.MAX_TIME
	newerThan := utils.MIN_TIME
	expectedRes := []dto.HostCores{
		{
			Date:  utils.P("2020-04-15T00:00:00Z"),
			Cores: 1,
		},
	}

	db.EXPECT().GetHostCores(location, environment, olderThan, newerThan).
		Return(expectedRes, nil).Times(1)

	res, err := as.GetHostCores(location, environment, olderThan, newerThan)

	require.NoError(t, err)
	assert.Equal(t, expectedRes, res)
}
