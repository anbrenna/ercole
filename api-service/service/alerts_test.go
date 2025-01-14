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

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/api-service/dto/filter"
	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

func TestSearchAlerts_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	expectedRes := &dto.Pagination{
		Items: []map[string]interface{}{
			{
				"AffectedHosts": 12,
				"Code":          "NEW_SERVER",
				"Count":         12,
				"OldestAlert":   "2020-05-06T15:40:04.543+02:00",
				"Severity":      "INFO",
			},
			{
				"AffectedHosts": 12,
				"Code":          "NEW_SERVER",
				"Count":         12,
				"OldestAlert":   "2020-05-06T15:40:04.543+02:00",
				"Severity":      "INFO",
			},
		},
		Count:    2,
		PageSize: 25,
		Page:     0,
	}

	alertFilter := filter.Alert{
		Mode:     "aggregated-code-severity",
		Keywords: []string{"foo", "bar", "foobarx"},
		SortBy:   "AlertCode",
		SortDesc: true,
		Filter:   filter.Filter{Limit: 1, Page: 1},
		Severity: model.AlertSeverityCritical,
		Status:   model.AlertStatusNew,
		From:     utils.P("2019-11-05T14:02:03Z"),
		To:       utils.P("2020-04-07T14:02:03Z"),
	}
	db.EXPECT().SearchAlerts(alertFilter).Return(
		expectedRes,
		nil,
	).Times(1)

	res, err := as.SearchAlerts(alertFilter)

	require.NoError(t, err)
	assert.Equal(t, expectedRes, res)
}

func TestSearchAlerts_Fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	alertFilter := filter.Alert{
		Mode:     "aggregated-code-severity",
		Keywords: []string{"foo", "bar", "foobarx"},
		SortBy:   "AlertCode",
		SortDesc: true,
		Filter:   filter.Filter{Limit: 1, Page: 1},
		Severity: model.AlertSeverityCritical,
		Status:   model.AlertStatusNew,
		From:     utils.P("2019-11-05T14:02:03Z"),
		To:       utils.P("2019-12-05T14:02:03Z"),
	}
	db.EXPECT().SearchAlerts(alertFilter).Return(
		nil,
		aerrMock,
	).Times(1)

	res, err := as.SearchAlerts(alertFilter)

	require.Equal(t, aerrMock, err)
	assert.Nil(t, res)
}

func TestAcknowledgeAlerts(t *testing.T) {
	testCases := []struct {
		filter dto.AlertsFilter
		expErr error
	}{
		{
			filter: dto.AlertsFilter{},
			expErr: nil,
		},
		{
			filter: dto.AlertsFilter{},
			expErr: aerrMock,
		},
	}

	var count int64

	for _, tc := range testCases {
		mockCtrl := gomock.NewController(t)
		defer func() {
			mockCtrl.Finish()
		}()

		db := NewMockMongoDatabaseInterface(mockCtrl)
		as := APIService{
			Database: db,
		}

		db.EXPECT().CountAlertsNODATA(tc.filter).Return(count, nil)
		db.EXPECT().UpdateAlertsStatus(tc.filter, model.AlertStatusAck).Return(tc.expErr)

		actErr := as.AckAlerts(tc.filter)
		assert.Equal(t, tc.expErr, actErr)
	}
}

func TestAcknowledgeAlerts_FailAlertCodeNoData(t *testing.T) {
	a_ack := dto.AlertsFilter{
		AlertCode: utils.Str2ptr(model.AlertCodeNoData),
	}

	dataErr := utils.NewErrorf("%w: you are trying to ack alerts with code: %s",
		utils.ErrInvalidAck,
		model.AlertCodeNoData)

	mockCtrl := gomock.NewController(t)
	defer func() {
		mockCtrl.Finish()
	}()

	as := APIService{}

	actErr := as.AckAlerts(a_ack)
	require.Error(t, actErr, dataErr.Message)
}

func TestAcknowledgeAlerts_FailCountAlertsNoData(t *testing.T) {
	a_ack := dto.AlertsFilter{}

	var count int64

	mockCtrl := gomock.NewController(t)
	defer func() {
		mockCtrl.Finish()
	}()

	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().CountAlertsNODATA(a_ack).Return(count, aerrMock)

	actErr := as.AckAlerts(a_ack)
	require.Equal(t, aerrMock, actErr)
}

func TestAcknowledgeAlerts_FailAckAlertsNoData(t *testing.T) {
	a_ack := dto.AlertsFilter{}
	var count int64 = 10

	mockCtrl := gomock.NewController(t)
	defer func() {
		mockCtrl.Finish()
	}()

	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().CountAlertsNODATA(a_ack).Return(count, nil)

	actErr := as.AckAlerts(a_ack)
	require.Error(t, actErr)
}

func TestSearchAlertsAsXLSX_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Database: db,
	}

	data := []map[string]interface{}{
		{
			"_id":                     utils.Str2oid("5f1943c97238d4bb6c98ef82"),
			"alertAffectedTechnology": "Oracle/Database",
			"alertCategory":           "LICENSE",
			"alertCode":               "NEW_LICENSE",
			"alertSeverity":           "CRITICAL",
			"alertStatus":             "NEW",
			"date":                    utils.PDT("2020-07-23T10:01:13.746+02:00"),
			"description":             "A new Enterprise license has been enabled to ercsoldbx",
			"hostname":                "ercsoldbx",
			"otherInfo": map[string]interface{}{
				"hostname": "ercsoldbx",
			},
		},
		{
			"_id":                     utils.Str2oid("5f1943c97238d4bb6c98ef83"),
			"alertAffectedTechnology": "Oracle/Database",
			"alertCategory":           "LICENSE",
			"alertCode":               "NEW_OPTION",
			"alertSeverity":           "CRITICAL",
			"alertStatus":             "NEW",
			"date":                    utils.PDT("2020-07-23T10:01:13.746+02:00"),
			"description":             "The database ERCSOL19 on ercsoldbx has enabled new features (Diagnostics Pack) on server",
			"hostname":                "ercsoldbx",
			"otherInfo": map[string]interface{}{
				"dbname": "ERCSOL19",
				"features": []string{
					"Diagnostics Pack",
				},
				"hostname": "ercsoldbx",
			},
		},
	}

	db.EXPECT().GetAlerts("Italy", "TST", "NEW", utils.P("2020-06-10T11:54:59Z"), utils.P("2020-06-17T11:54:59Z"), utils.P("2019-12-05T14:02:03Z")).
		Return(data, nil).Times(1)

	filter := dto.GlobalFilter{
		Location:    "Italy",
		Environment: "TST",
		OlderThan:   utils.P("2019-12-05T14:02:03Z"),
	}

	from := utils.P("2020-06-10T11:54:59Z")
	to := utils.P("2020-06-17T11:54:59Z")

	actual, err := as.SearchAlertsAsXLSX("NEW", from, to, filter)
	require.NoError(t, err)
	assert.Equal(t, "LICENSE", actual.GetCellValue("Alerts", "A2"))
	assert.Equal(t, "2020-07-23 08:01:13.746 +0000 UTC", actual.GetCellValue("Alerts", "B2"))
	assert.Equal(t, "CRITICAL", actual.GetCellValue("Alerts", "C2"))
	assert.Equal(t, "ercsoldbx", actual.GetCellValue("Alerts", "D2"))
	assert.Equal(t, "NEW_LICENSE", actual.GetCellValue("Alerts", "E2"))
	assert.Equal(t, "A new Enterprise license has been enabled to ercsoldbx", actual.GetCellValue("Alerts", "F2"))

	assert.Equal(t, "LICENSE", actual.GetCellValue("Alerts", "A3"))
	assert.Equal(t, "2020-07-23 08:01:13.746 +0000 UTC", actual.GetCellValue("Alerts", "B3"))
	assert.Equal(t, "CRITICAL", actual.GetCellValue("Alerts", "C3"))
	assert.Equal(t, "ercsoldbx", actual.GetCellValue("Alerts", "D3"))
	assert.Equal(t, "NEW_OPTION", actual.GetCellValue("Alerts", "E3"))
	assert.Equal(t, "The database ERCSOL19 on ercsoldbx has enabled new features (Diagnostics Pack) on server", actual.GetCellValue("Alerts", "F3"))
}

func TestUpdateAlertsStatus_Success(t *testing.T) {
	testCases := []struct {
		filter dto.AlertsFilter
		expErr error
	}{
		{
			filter: dto.AlertsFilter{},
			expErr: nil,
		},
		{
			filter: dto.AlertsFilter{},
			expErr: aerrMock,
		},
	}

	for _, tc := range testCases {
		mockCtrl := gomock.NewController(t)
		defer func() {
			mockCtrl.Finish()
		}()

		db := NewMockMongoDatabaseInterface(mockCtrl)
		as := APIService{
			Database: db,
		}

		db.EXPECT().UpdateAlertsStatus(tc.filter, model.AlertStatusDismissed).Return(tc.expErr)

		actErr := as.UpdateAlertsStatus(tc.filter, model.AlertStatusDismissed)
		assert.Equal(t, tc.expErr, actErr)
	}
}
