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
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/utils"
)

func TestSearchClusters_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	expectedRes := []dto.Cluster{
		{
			CPU:                         0,
			Environment:                 "PROD",
			Hostname:                    "fb-canvas-b9b1d8fa8328fe972b1e031621e8a6c9",
			HostnameAgentVirtualization: "fb-canvas-b9b1d8fa8328fe972b1e031621e8a6c9",
			Location:                    "Italy",
			Name:                        "not_in_cluster",
			VirtualizationNodes:         []string{"aspera-b1fe49e8501c9ef031e5acff4b5e69a9"},
			PhysicalServerModelNames:    []string{"model name"},
			Sockets:                     0,
			Type:                        "unknown",
			ID:                          utils.Str2oid("5e8c234b24f648a08585bd3d"),
		},
		{
			CPU:                         140,
			Environment:                 "PROD",
			Hostname:                    "test-virt",
			HostnameAgentVirtualization: "test-virt",
			Location:                    "Italy",
			Name:                        "Puzzait",
			VirtualizationNodes:         []string{"s157-cb32c10a56c256746c337e21b3f82402"},
			PhysicalServerModelNames:    []string{"new model name"},
			Sockets:                     10,
			Type:                        "vmware",
			ID:                          utils.Str2oid("5e8c234b24f648a08585bd41"),
		},
	}

	db.EXPECT().SearchClusters(
		"full", []string{"foo", "bar", "foobarx"}, "CPU",
		true, 1, 1,
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	).Return(expectedRes, nil).Times(1)

	res, err := as.SearchClusters(
		"full", "foo bar foobarx", "CPU",
		true, 1, 1,
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	)

	require.NoError(t, err)
	assert.Equal(t, expectedRes, res)
}

func TestSearchClusterNames_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	expectedRes := []dto.Cluster{
		{
			Name: "not_in_cluster",
		},
		{
			Name: "Puzzait",
		},
	}

	db.EXPECT().SearchClusters(
		"clusternames", []string{"foo", "bar", "foobarx"}, "CPU",
		true, 1, 1,
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	).Return(expectedRes, nil).Times(1)

	res, err := as.SearchClusters(
		"clusternames", "foo bar foobarx", "CPU",
		true, 1, 1,
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	)

	require.NoError(t, err)
	assert.Equal(t, expectedRes, res)
}

func TestSearchClusters_Fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().SearchClusters(
		"full", []string{"foo", "bar", "foobarx"}, "CPU",
		true, 1, 1,
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	).Return(nil, aerrMock).Times(1)

	res, err := as.SearchClusters(
		"full", "foo bar foobarx", "CPU",
		true, 1, 1,
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	)

	require.Nil(t, res)
	assert.Equal(t, aerrMock, err)
}

func TestGetClusterXLSX(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Log: logger.NewLogger("TEST"),
	}

	cluster := &dto.Cluster{
		ID:                          utils.Str2oid("5eb0222a45d85f4193704944"),
		CPU:                         140,
		CreatedAt:                   utils.P("2020-05-04T14:09:46.608Z"),
		Environment:                 "PROD",
		FetchEndpoint:               "???",
		Hostname:                    "test-virt",
		HostnameAgentVirtualization: "test-virt",
		Location:                    "Italy",
		Name:                        "Puzzait",
		Sockets:                     10,
		Type:                        "vmware",
		VirtualizationNodes:         []string{"s157-cb32c10a56c256746c337e21b3f82402"},
		VirtualizationNodesCount:    1,
		VirtualizationNodesStats: []dto.VirtualizationNodesStat{
			{
				TotalVMsCount:                   2,
				TotalVMsWithErcoleAgentCount:    1,
				TotalVMsWithoutErcoleAgentCount: 1,
				VirtualizationNode:              "s157-cb32c10a56c256746c337e21b3f82402"}},

		VMs: []dto.VM{
			{
				CappedCPU:               false,
				Hostname:                "test-virt",
				Name:                    "test-virt",
				VirtualizationNode:      "s157-cb32c10a56c256746c337e21b3f82402",
				PhysicalServerModelName: "HP ProLiant DL380 Gen9",
				IsErcoleInstalled:       false,
			},

			{
				CappedCPU:               true,
				Hostname:                "test-db",
				Name:                    "test-db",
				VirtualizationNode:      "s157-cb32c10a56c256746c337e21b3f82402",
				PhysicalServerModelName: "HP ProLiant DL380 Gen9",
				IsErcoleInstalled:       false,
			},
		},
		VMsCount:            2,
		VMsErcoleAgentCount: 1,
	}

	var clusterName = "pippo"
	var olderThan = utils.P("2019-11-05T14:02:03Z")

	db.EXPECT().
		GetCluster(clusterName, olderThan).
		Return(cluster, nil)

	xlsx, err := as.GetClusterXLSX(clusterName, olderThan)
	assert.NoError(t, err)

	sheet := "VMs"

	i := -1
	columns := []string{"A", "B", "C", "D", "E"}
	nextVal := func() string {
		i += 1
		cell := columns[i%5] + strconv.Itoa(i/5+1)
		return xlsx.GetCellValue(sheet, cell)
	}

	assert.Equal(t, "Physical Hosts", nextVal())
	assert.Equal(t, "Hostname", nextVal())
	assert.Equal(t, "VirtualizationNode", nextVal())
	assert.Equal(t, "PhysicalServerModelName", nextVal())
	assert.Equal(t, "CappedCPU", nextVal())

	assert.Equal(t, "test-virt", nextVal())
	assert.Equal(t, "test-virt", nextVal())
	assert.Equal(t, "s157-cb32c10a56c256746c337e21b3f82402", nextVal())
	assert.Equal(t, "HP ProLiant DL380 Gen9", nextVal())
	assert.Equal(t, "false", nextVal())

	assert.Equal(t, "test-db", nextVal())
	assert.Equal(t, "test-db", nextVal())
	assert.Equal(t, "s157-cb32c10a56c256746c337e21b3f82402", nextVal())
	assert.Equal(t, "HP ProLiant DL380 Gen9", nextVal())
	assert.Equal(t, "true", nextVal())

	assert.Equal(t, "", nextVal())
	assert.Equal(t, "", nextVal())
	assert.Equal(t, "", nextVal())
	assert.Equal(t, "", nextVal())
	assert.Equal(t, "", nextVal())
}

func TestSearchClustersAsXLSX_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Database: db,
	}

	data := []dto.Cluster{
		{
			Name:                     "Puzzait",
			Type:                     "VMWare/VMWare",
			CPU:                      140,
			Sockets:                  10,
			VirtualizationNodes:      []string{"s157-cb32c10a56c256746c337e21b3f82402"},
			PhysicalServerModelNames: []string{"model name"},
			VMsCount:                 2,
			VMsErcoleAgentCount:      2,
		},
	}

	db.EXPECT().SearchClusters(
		"full", []string{}, "",
		false, -1, -1,
		"Italy", "TST", utils.P("2019-12-05T14:02:03Z"),
	).Return(data, nil).Times(1)

	filter := dto.GlobalFilter{
		Location:    "Italy",
		Environment: "TST",
		OlderThan:   utils.P("2019-12-05T14:02:03Z"),
	}

	actual, err := as.SearchClustersAsXLSX(filter)
	require.NoError(t, err)
	assert.Equal(t, "Puzzait", actual.GetCellValue("Hypervisor", "A2"))
	assert.Equal(t, "VMWare/VMWare", actual.GetCellValue("Hypervisor", "B2"))
	assert.Equal(t, "140", actual.GetCellValue("Hypervisor", "C2"))
	assert.Equal(t, "10", actual.GetCellValue("Hypervisor", "D2"))
	assert.Equal(t, "[s157-cb32c10a56c256746c337e21b3f82402]", actual.GetCellValue("Hypervisor", "E2"))
	assert.Equal(t, "[model name]", actual.GetCellValue("Hypervisor", "F2"))
	assert.Equal(t, "2", actual.GetCellValue("Hypervisor", "G2"))
	assert.Equal(t, "2", actual.GetCellValue("Hypervisor", "H2"))
}
