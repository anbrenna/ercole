// Copyright (c) 2020 Sorint.lab S.p.A.
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

package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/amreo/ercole-services/config"
	"github.com/amreo/ercole-services/utils"
	gomock "github.com/golang/mock/gomock"
	"github.com/plandem/xlsx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearchAddms_JSONPaged(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	expectedRes := map[string]interface{}{
		"Content": []interface{}{
			map[string]interface{}{
				"Action":         "Run SQL Tuning Advisor on the SELECT statement with SQL_ID \"4ztz048yfq32s\".",
				"Benefit":        83.34,
				"CreatedAt":      utils.P("2020-04-07T08:52:59.872+02:00"),
				"Dbname":         "ERCOLE",
				"Environment":    "TST",
				"Finding":        "SQL statements consuming significant database time were found. These statements offer a good opportunity for performance improvement.",
				"Hostname":       "test-db",
				"Location":       "Germany",
				"Recommendation": "SQL Tuning",
				"_id":            utils.Str2oid("5e8c234b24f648a08585bd43"),
			},
			map[string]interface{}{
				"Action":         "Look at the \"Top SQL Statements\" finding for SQL statements consuming significant I/O on this segment. For example, the SELECT statement with SQL_ID \"4ztz048yfq32s\" is responsible for 100% of \"User I/O\" and \"Cluster\" waits for this segment.",
				"Benefit":        68.24,
				"CreatedAt":      utils.P("2020-04-07T08:52:59.872+02:00"),
				"Dbname":         "ERCOLE",
				"Environment":    "TST",
				"Finding":        "Individual database segments responsible for significant \"User I/O\" and \"Cluster\" waits were found.",
				"Hostname":       "test-db",
				"Location":       "Germany",
				"Recommendation": "Segment Tuning",
				"_id":            utils.Str2oid("5e8c234b24f648a08585bd43"),
			},
		},
		"Metadata": map[string]interface{}{
			"Empty":         false,
			"First":         true,
			"Last":          true,
			"Number":        0,
			"Size":          20,
			"TotalElements": 25,
			"TotalPages":    1,
		},
	}

	resFromService := []map[string]interface{}{
		expectedRes,
	}

	as.EXPECT().
		SearchAddms("foobar", "Benefit", true, 2, 3, "Italy", "TST", utils.P("2020-06-10T11:54:59Z")).
		Return(resFromService, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchAddms)
	req, err := http.NewRequest("GET", "/addms?search=foobar&sort-by=Benefit&sort-desc=true&page=2&size=3&location=Italy&environment=TST&older-than=2020-06-10T11%3A54%3A59Z", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(expectedRes), rr.Body.String())
}

func TestSearchAddms_JSONUnpaged(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	expectedRes := []map[string]interface{}{
		map[string]interface{}{
			"Action":         "Run SQL Tuning Advisor on the SELECT statement with SQL_ID \"4ztz048yfq32s\".",
			"Benefit":        83.34,
			"CreatedAt":      utils.P("2020-04-07T08:52:59.872+02:00"),
			"Dbname":         "ERCOLE",
			"Environment":    "TST",
			"Finding":        "SQL statements consuming significant database time were found. These statements offer a good opportunity for performance improvement.",
			"Hostname":       "test-db",
			"Location":       "Germany",
			"Recommendation": "SQL Tuning",
			"_id":            utils.Str2oid("5e8c234b24f648a08585bd43"),
		},
		map[string]interface{}{
			"Action":         "Look at the \"Top SQL Statements\" finding for SQL statements consuming significant I/O on this segment. For example, the SELECT statement with SQL_ID \"4ztz048yfq32s\" is responsible for 100% of \"User I/O\" and \"Cluster\" waits for this segment.",
			"Benefit":        68.24,
			"CreatedAt":      utils.P("2020-04-07T08:52:59.872+02:00"),
			"Dbname":         "ERCOLE",
			"Environment":    "TST",
			"Finding":        "Individual database segments responsible for significant \"User I/O\" and \"Cluster\" waits were found.",
			"Hostname":       "test-db",
			"Location":       "Germany",
			"Recommendation": "Segment Tuning",
			"_id":            utils.Str2oid("5e8c234b24f648a08585bd43"),
		},
	}

	as.EXPECT().
		SearchAddms("", "", false, -1, -1, "", "", utils.MAX_TIME).
		Return(expectedRes, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchAddms)
	req, err := http.NewRequest("GET", "/addms", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(expectedRes), rr.Body.String())
}

func TestSearchAddms_JSONUnprocessableEntity1(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchAddms)
	req, err := http.NewRequest("GET", "/addms?sort-desc=sdfdfsdfs", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchAddms_JSONUnprocessableEntity2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchAddms)
	req, err := http.NewRequest("GET", "/addms?page=sdfdfsdfs", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchAddms_JSONUnprocessableEntity3(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchAddms)
	req, err := http.NewRequest("GET", "/addms?size=sdfdfsdfs", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchAddms_JSONUnprocessableEntity4(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchAddms)
	req, err := http.NewRequest("GET", "/addms?older-than=sdfdfsdfs", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchAddms_JSONInternalServerError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	as.EXPECT().
		SearchAddms("", "", false, -1, -1, "", "", utils.MAX_TIME).
		Return(nil, aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchAddms)
	req, err := http.NewRequest("GET", "/addms", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestSearchAddms_XLSXSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Log: utils.NewLogger("TEST"),
	}

	expectedRes := []map[string]interface{}{
		map[string]interface{}{
			"Action":         "Run SQL Tuning Advisor on the SELECT statement with SQL_ID \"4ztz048yfq32s\".",
			"Benefit":        83.34,
			"CreatedAt":      utils.P("2020-04-07T08:52:59.872+02:00"),
			"Dbname":         "ERCOLE",
			"Environment":    "TST",
			"Finding":        "SQL statements consuming significant database time were found. These statements offer a good opportunity for performance improvement.",
			"Hostname":       "test-db",
			"Location":       "Germany",
			"Recommendation": "SQL Tuning",
			"_id":            utils.Str2oid("5e8c234b24f648a08585bd43"),
		},
		map[string]interface{}{
			"Action":         "Look at the \"Top SQL Statements\" finding for SQL statements consuming significant I/O on this segment. For example, the SELECT statement with SQL_ID \"4ztz048yfq32s\" is responsible for 100% of \"User I/O\" and \"Cluster\" waits for this segment.",
			"Benefit":        68.24,
			"CreatedAt":      utils.P("2020-04-07T08:52:59.872+02:00"),
			"Dbname":         "ERCOLE",
			"Environment":    "TST",
			"Finding":        "Individual database segments responsible for significant \"User I/O\" and \"Cluster\" waits were found.",
			"Hostname":       "test-db",
			"Location":       "Germany",
			"Recommendation": "Segment Tuning",
			"_id":            utils.Str2oid("5e8c234b24f648a08585bd43"),
		},
	}

	as.EXPECT().
		SearchAddms("foobar", "Benefit", true, -1, -1, "Germany", "TST", utils.P("2020-06-10T11:54:59Z")).
		Return(expectedRes, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchAddms)
	req, err := http.NewRequest("GET", "/addms?search=foobar&location=Germany&environment=TST&older-than=2020-06-10T11%3A54%3A59Z", nil)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	sp, err := xlsx.Open(rr.Body)
	require.NoError(t, err)
	sh := sp.SheetByName("Addm")
	require.NotNil(t, sh)
	assert.Equal(t, "Run SQL Tuning Advisor on the SELECT statement with SQL_ID \"4ztz048yfq32s\".", sh.Cell(0, 1).String())
	AssertXLSXFloat(t, 83.34, sh.Cell(1, 1))
	assert.Equal(t, "ERCOLE", sh.Cell(2, 1).String())
	assert.Equal(t, "TST", sh.Cell(3, 1).String())
	assert.Equal(t, "SQL statements consuming significant database time were found. These statements offer a good opportunity for performance improvement.", sh.Cell(4, 1).String())
	assert.Equal(t, "test-db", sh.Cell(5, 1).String())
	assert.Equal(t, "SQL Tuning", sh.Cell(6, 1).String())

	assert.Equal(t, "Look at the \"Top SQL Statements\" finding for SQL statements consuming significant I/O on this segment. For example, the SELECT statement with SQL_ID \"4ztz048yfq32s\" is responsible for 100% of \"User I/O\" and \"Cluster\" waits for this segment.", sh.Cell(0, 2).String())
	AssertXLSXFloat(t, 68.24, sh.Cell(1, 2))
	assert.Equal(t, "ERCOLE", sh.Cell(2, 2).String())
	assert.Equal(t, "TST", sh.Cell(3, 2).String())
	assert.Equal(t, "Individual database segments responsible for significant \"User I/O\" and \"Cluster\" waits were found.", sh.Cell(4, 2).String())
	assert.Equal(t, "test-db", sh.Cell(5, 2).String())
	assert.Equal(t, "Segment Tuning", sh.Cell(6, 2).String())
}

func TestSearchAddms_XLSXUnprocessableEntity1(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Log: utils.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchAddms)
	req, err := http.NewRequest("GET", "/addms?older-than=aasdasd", nil)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchAddms_XLSXInternalServerError1(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Log: utils.NewLogger("TEST"),
	}

	as.EXPECT().
		SearchAddms("", "Benefit", true, -1, -1, "", "", utils.MAX_TIME).
		Return(nil, aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchAddms)
	req, err := http.NewRequest("GET", "/addms", nil)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestSearchAddms_XLSXInternalServerError2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	expectedRes := []map[string]interface{}{
		map[string]interface{}{
			"Action":         "Run SQL Tuning Advisor on the SELECT statement with SQL_ID \"4ztz048yfq32s\".",
			"Benefit":        83.34,
			"CreatedAt":      utils.P("2020-04-07T08:52:59.872+02:00"),
			"Dbname":         "ERCOLE",
			"Environment":    "TST",
			"Finding":        "SQL statements consuming significant database time were found. These statements offer a good opportunity for performance improvement.",
			"Hostname":       "test-db",
			"Location":       "Germany",
			"Recommendation": "SQL Tuning",
			"_id":            utils.Str2oid("5e8c234b24f648a08585bd43"),
		},
		map[string]interface{}{
			"Action":         "Look at the \"Top SQL Statements\" finding for SQL statements consuming significant I/O on this segment. For example, the SELECT statement with SQL_ID \"4ztz048yfq32s\" is responsible for 100% of \"User I/O\" and \"Cluster\" waits for this segment.",
			"Benefit":        68.24,
			"CreatedAt":      utils.P("2020-04-07T08:52:59.872+02:00"),
			"Dbname":         "ERCOLE",
			"Environment":    "TST",
			"Finding":        "Individual database segments responsible for significant \"User I/O\" and \"Cluster\" waits were found.",
			"Hostname":       "test-db",
			"Location":       "Germany",
			"Recommendation": "Segment Tuning",
			"_id":            utils.Str2oid("5e8c234b24f648a08585bd43"),
		},
	}

	as.EXPECT().
		SearchAddms("", "Benefit", true, -1, -1, "", "", utils.MAX_TIME).
		Return(expectedRes, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchAddms)
	req, err := http.NewRequest("GET", "/addms", nil)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestSearchSegmentAdvisors_JSONPaged(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	expectedRes := map[string]interface{}{
		"Content": []interface{}{
			map[string]interface{}{
				"CreatedAt":      utils.P("2020-04-07T08:52:59.82+02:00"),
				"Dbname":         "4wcqjn-ecf040bdfab7695ab332aef7401f185c",
				"Environment":    "SVIL",
				"Hostname":       "publicitate-36d06ca83eafa454423d2097f4965517",
				"Location":       "Germany",
				"PartitionName":  "",
				"Reclaimable":    "\u003c1",
				"Recommendation": "3d7e603f515ed171fc99bdb908f38fb2",
				"SegmentName":    "nascar1-f9b3703bf8b3cc7ae070cd28e7fed7b3",
				"SegmentOwner":   "Brittany-424f6a749eef846fa40a1ad1ee3d3674",
				"SegmentType":    "TABLE",
				"_id":            utils.Str2oid("5e8c234b24f648a08585bd32"),
			},
			map[string]interface{}{
				"CreatedAt":      utils.P("2020-04-07T08:52:59.872+02:00"),
				"Dbname":         "ERCOLE",
				"Environment":    "TST",
				"Hostname":       "test-db",
				"Location":       "Germany",
				"PartitionName":  "iyyiuyyoy",
				"Reclaimable":    "\u003c1",
				"Recommendation": "32b36a77e7481343ef175483c086859e",
				"SegmentName":    "pasta-973e4d1f937da4d9bc1b092f934ab0ec",
				"SegmentOwner":   "Brittany-424f6a749eef846fa40a1ad1ee3d3674",
				"SegmentType":    "TABLE",
				"_id":            utils.Str2oid("5e8c234b24f648a08585bd43"),
			},
		},
		"Metadata": map[string]interface{}{
			"Empty":         false,
			"First":         true,
			"Last":          true,
			"Number":        0,
			"Size":          20,
			"TotalElements": 25,
			"TotalPages":    1,
		},
	}

	resFromService := []map[string]interface{}{
		expectedRes,
	}

	as.EXPECT().
		SearchSegmentAdvisors("foobar", "Reclaimable", true, 2, 3, "Italy", "TST", utils.P("2020-06-10T11:54:59Z")).
		Return(resFromService, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchSegmentAdvisors)
	req, err := http.NewRequest("GET", "/segment-advisors?search=foobar&sort-by=Reclaimable&sort-desc=true&page=2&size=3&location=Italy&environment=TST&older-than=2020-06-10T11%3A54%3A59Z", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(expectedRes), rr.Body.String())
}

func TestSearchSegmentAdvisors_JSONUnpaged(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	expectedRes := []map[string]interface{}{
		map[string]interface{}{
			"CreatedAt":      utils.P("2020-04-07T08:52:59.82+02:00"),
			"Dbname":         "4wcqjn-ecf040bdfab7695ab332aef7401f185c",
			"Environment":    "SVIL",
			"Hostname":       "publicitate-36d06ca83eafa454423d2097f4965517",
			"Location":       "Germany",
			"PartitionName":  "",
			"Reclaimable":    "\u003c1",
			"Recommendation": "3d7e603f515ed171fc99bdb908f38fb2",
			"SegmentName":    "nascar1-f9b3703bf8b3cc7ae070cd28e7fed7b3",
			"SegmentOwner":   "Brittany-424f6a749eef846fa40a1ad1ee3d3674",
			"SegmentType":    "TABLE",
			"_id":            utils.Str2oid("5e8c234b24f648a08585bd32"),
		},
		map[string]interface{}{
			"CreatedAt":      utils.P("2020-04-07T08:52:59.872+02:00"),
			"Dbname":         "ERCOLE",
			"Environment":    "TST",
			"Hostname":       "test-db",
			"Location":       "Germany",
			"PartitionName":  "iyyiuyyoy",
			"Reclaimable":    "\u003c1",
			"Recommendation": "32b36a77e7481343ef175483c086859e",
			"SegmentName":    "pasta-973e4d1f937da4d9bc1b092f934ab0ec",
			"SegmentOwner":   "Brittany-424f6a749eef846fa40a1ad1ee3d3674",
			"SegmentType":    "TABLE",
			"_id":            utils.Str2oid("5e8c234b24f648a08585bd43"),
		},
	}

	as.EXPECT().
		SearchSegmentAdvisors("", "", false, -1, -1, "", "", utils.MAX_TIME).
		Return(expectedRes, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchSegmentAdvisors)
	req, err := http.NewRequest("GET", "/segment-advisors", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(expectedRes), rr.Body.String())
}

func TestSearchSegmentAdvisors_JSONUnprocessableEntity1(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchSegmentAdvisors)
	req, err := http.NewRequest("GET", "/segment-advisors?sort-desc=asasdasd", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchSegmentAdvisors_JSONUnprocessableEntity2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchSegmentAdvisors)
	req, err := http.NewRequest("GET", "/segment-advisors?page=asasdasd", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchSegmentAdvisors_JSONUnprocessableEntity3(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchSegmentAdvisors)
	req, err := http.NewRequest("GET", "/segment-advisors?size=asasdasd", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchSegmentAdvisors_JSONUnprocessableEntity4(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchSegmentAdvisors)
	req, err := http.NewRequest("GET", "/segment-advisors?older-than=asasdasd", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchSegmentAdvisors_JSONInternalServerError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	as.EXPECT().
		SearchSegmentAdvisors("", "", false, -1, -1, "", "", utils.MAX_TIME).
		Return(nil, aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchSegmentAdvisors)
	req, err := http.NewRequest("GET", "/segment-advisors", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestSearchSegmentAdvisors_XLSXSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Log: utils.NewLogger("TEST"),
	}

	expectedRes := []map[string]interface{}{
		map[string]interface{}{
			"CreatedAt":      utils.P("2020-04-07T08:52:59.82+02:00"),
			"Dbname":         "4wcqjn-ecf040bdfab7695ab332aef7401f185c",
			"Environment":    "SVIL",
			"Hostname":       "publicitate-36d06ca83eafa454423d2097f4965517",
			"Location":       "Germany",
			"PartitionName":  "",
			"Reclaimable":    "\u003c1",
			"Recommendation": "3d7e603f515ed171fc99bdb908f38fb2",
			"SegmentName":    "nascar1-f9b3703bf8b3cc7ae070cd28e7fed7b3",
			"SegmentOwner":   "Brittany-424f6a749eef846fa40a1ad1ee3d3674",
			"SegmentType":    "TABLE",
			"_id":            utils.Str2oid("5e8c234b24f648a08585bd32"),
		},
		map[string]interface{}{
			"CreatedAt":      utils.P("2020-04-07T08:52:59.872+02:00"),
			"Dbname":         "ERCOLE",
			"Environment":    "TST",
			"Hostname":       "test-db",
			"Location":       "Germany",
			"PartitionName":  "iyyiuyyoy",
			"Reclaimable":    "\u003c1",
			"Recommendation": "32b36a77e7481343ef175483c086859e",
			"SegmentName":    "pasta-973e4d1f937da4d9bc1b092f934ab0ec",
			"SegmentOwner":   "Brittany-424f6a749eef846fa40a1ad1ee3d3674",
			"SegmentType":    "TABLE",
			"_id":            utils.Str2oid("5e8c234b24f648a08585bd43"),
		},
	}

	as.EXPECT().
		SearchSegmentAdvisors("foobar", "Reclaimable", true, -1, -1, "Italy", "TST", utils.P("2020-06-10T11:54:59Z")).
		Return(expectedRes, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchSegmentAdvisors)
	req, err := http.NewRequest("GET", "/segment-advisors?search=foobar&sort-by=Reclaimable&sort-desc=true&location=Italy&environment=TST&older-than=2020-06-10T11%3A54%3A59Z", nil)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	sp, err := xlsx.Open(rr.Body)
	require.NoError(t, err)
	sh := sp.SheetByName("Segment_Advisor")
	require.NotNil(t, sh)
	assert.Equal(t, "4wcqjn-ecf040bdfab7695ab332aef7401f185c", sh.Cell(0, 1).String())
	assert.Equal(t, "SVIL", sh.Cell(1, 1).String())
	assert.Equal(t, "publicitate-36d06ca83eafa454423d2097f4965517", sh.Cell(2, 1).String())
	assert.Equal(t, "", sh.Cell(3, 1).String())
	assert.Equal(t, "\u003c1", sh.Cell(4, 1).String())
	assert.Equal(t, "3d7e603f515ed171fc99bdb908f38fb2", sh.Cell(5, 1).String())
	assert.Equal(t, "nascar1-f9b3703bf8b3cc7ae070cd28e7fed7b3", sh.Cell(6, 1).String())
	assert.Equal(t, "Brittany-424f6a749eef846fa40a1ad1ee3d3674", sh.Cell(7, 1).String())
	assert.Equal(t, "TABLE", sh.Cell(8, 1).String())

	assert.Equal(t, "ERCOLE", sh.Cell(0, 2).String())
	assert.Equal(t, "TST", sh.Cell(1, 2).String())
	assert.Equal(t, "test-db", sh.Cell(2, 2).String())
	assert.Equal(t, "iyyiuyyoy", sh.Cell(3, 2).String())
	assert.Equal(t, "\u003c1", sh.Cell(4, 2).String())
	assert.Equal(t, "32b36a77e7481343ef175483c086859e", sh.Cell(5, 2).String())
	assert.Equal(t, "pasta-973e4d1f937da4d9bc1b092f934ab0ec", sh.Cell(6, 2).String())
	assert.Equal(t, "Brittany-424f6a749eef846fa40a1ad1ee3d3674", sh.Cell(7, 2).String())
	assert.Equal(t, "TABLE", sh.Cell(8, 2).String())
}

func TestSearchSegmentAdvisors_XLSXUnprocessableEntity1(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Log: utils.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchSegmentAdvisors)
	req, err := http.NewRequest("GET", "/segment-advisors?sort-desc=sadasddasasd", nil)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchSegmentAdvisors_XLSXUnprocessableEntity2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Log: utils.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchSegmentAdvisors)
	req, err := http.NewRequest("GET", "/segment-advisors?older-than=sadasddasasd", nil)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchSegmentAdvisors_XLSXInternalServerError1(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Log: utils.NewLogger("TEST"),
	}

	as.EXPECT().
		SearchSegmentAdvisors("", "", false, -1, -1, "", "", utils.MAX_TIME).
		Return(nil, aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchSegmentAdvisors)
	req, err := http.NewRequest("GET", "/segment-advisors", nil)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestSearchSegmentAdvisors_XLSXInternalServerError2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	expectedRes := []map[string]interface{}{
		map[string]interface{}{
			"CreatedAt":      utils.P("2020-04-07T08:52:59.82+02:00"),
			"Dbname":         "4wcqjn-ecf040bdfab7695ab332aef7401f185c",
			"Environment":    "SVIL",
			"Hostname":       "publicitate-36d06ca83eafa454423d2097f4965517",
			"Location":       "Germany",
			"PartitionName":  "",
			"Reclaimable":    "\u003c1",
			"Recommendation": "3d7e603f515ed171fc99bdb908f38fb2",
			"SegmentName":    "nascar1-f9b3703bf8b3cc7ae070cd28e7fed7b3",
			"SegmentOwner":   "Brittany-424f6a749eef846fa40a1ad1ee3d3674",
			"SegmentType":    "TABLE",
			"_id":            utils.Str2oid("5e8c234b24f648a08585bd32"),
		},
		map[string]interface{}{
			"CreatedAt":      utils.P("2020-04-07T08:52:59.872+02:00"),
			"Dbname":         "ERCOLE",
			"Environment":    "TST",
			"Hostname":       "test-db",
			"Location":       "Germany",
			"PartitionName":  "iyyiuyyoy",
			"Reclaimable":    "\u003c1",
			"Recommendation": "32b36a77e7481343ef175483c086859e",
			"SegmentName":    "pasta-973e4d1f937da4d9bc1b092f934ab0ec",
			"SegmentOwner":   "Brittany-424f6a749eef846fa40a1ad1ee3d3674",
			"SegmentType":    "TABLE",
			"_id":            utils.Str2oid("5e8c234b24f648a08585bd43"),
		},
	}

	as.EXPECT().
		SearchSegmentAdvisors("", "", false, -1, -1, "", "", utils.MAX_TIME).
		Return(expectedRes, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchSegmentAdvisors)
	req, err := http.NewRequest("GET", "/segment-advisors", nil)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestSearchPatchAdvisors_JSONPaged(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	expectedRes := map[string]interface{}{
		"Content": []interface{}{
			map[string]interface{}{
				"CreatedAt":   utils.P("2020-04-07T08:52:59.82+02:00"),
				"Date":        utils.P("2012-04-16T02:00:00+02:00"),
				"Dbname":      "4wcqjn-ecf040bdfab7695ab332aef7401f185c",
				"Dbver":       "11.2.0.3.0 Enterprise Edition",
				"Description": "PSU 11.2.0.3.2",
				"Environment": "SVIL",
				"Hostname":    "publicitate-36d06ca83eafa454423d2097f4965517",
				"Location":    "Germany",
				"Status":      "KO",
				"_id":         utils.Str2oid("5e8c234b24f648a08585bd32"),
			},
			map[string]interface{}{
				"CreatedAt":   utils.P("2020-04-07T08:52:59.872+02:00"),
				"Date":        utils.P("2012-04-16T02:00:00+02:00"),
				"Dbname":      "ERCOLE",
				"Dbver":       "12.2.0.1.0 Enterprise Edition",
				"Description": "PSU 11.2.0.3.2",
				"Environment": "TST",
				"Hostname":    "test-db",
				"Location":    "Germany",
				"Status":      "KO",
				"_id":         utils.Str2oid("5e8c234b24f648a08585bd43"),
			},
		},
		"Metadata": map[string]interface{}{
			"Empty":         false,
			"First":         true,
			"Last":          true,
			"Number":        0,
			"Size":          20,
			"TotalElements": 25,
			"TotalPages":    1,
		},
	}

	resFromService := []map[string]interface{}{
		expectedRes,
	}

	as.EXPECT().
		SearchPatchAdvisors("foobar", "Hostname", true, 2, 3, utils.P("2019-03-05T14:02:03Z"), "Italy", "TST", utils.P("2020-06-10T11:54:59Z"), "KO").
		Return(resFromService, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchPatchAdvisors)
	req, err := http.NewRequest("GET", "/patch-advisors?search=foobar&sort-by=Hostname&sort-desc=true&page=2&size=3&window-time=8&status=KO&location=Italy&environment=TST&older-than=2020-06-10T11%3A54%3A59Z", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(expectedRes), rr.Body.String())
}

func TestSearchPatchAdvisors_JSONUnpaged(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	expectedRes := []map[string]interface{}{
		map[string]interface{}{
			"CreatedAt":   utils.P("2020-04-07T08:52:59.82+02:00"),
			"Date":        utils.P("2012-04-16T02:00:00+02:00"),
			"Dbname":      "4wcqjn-ecf040bdfab7695ab332aef7401f185c",
			"Dbver":       "11.2.0.3.0 Enterprise Edition",
			"Description": "PSU 11.2.0.3.2",
			"Environment": "SVIL",
			"Hostname":    "publicitate-36d06ca83eafa454423d2097f4965517",
			"Location":    "Germany",
			"Status":      "KO",
			"_id":         utils.Str2oid("5e8c234b24f648a08585bd32"),
		},
		map[string]interface{}{
			"CreatedAt":   utils.P("2020-04-07T08:52:59.872+02:00"),
			"Date":        utils.P("2012-04-16T02:00:00+02:00"),
			"Dbname":      "ERCOLE",
			"Dbver":       "12.2.0.1.0 Enterprise Edition",
			"Description": "PSU 11.2.0.3.2",
			"Environment": "TST",
			"Hostname":    "test-db",
			"Location":    "Germany",
			"Status":      "KO",
			"_id":         utils.Str2oid("5e8c234b24f648a08585bd43"),
		},
	}

	as.EXPECT().
		SearchPatchAdvisors("", "", false, -1, -1, utils.P("2019-05-05T14:02:03Z"), "", "", utils.MAX_TIME, "").
		Return(expectedRes, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchPatchAdvisors)
	req, err := http.NewRequest("GET", "/patch-advisors", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(expectedRes), rr.Body.String())
}

func TestSearchPatchAdvisors_JSONUnprocessableEntity1(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchPatchAdvisors)
	req, err := http.NewRequest("GET", "/patch-advisors?sort-desc=sdasdasdasd", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchPatchAdvisors_JSONUnprocessableEntity2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchPatchAdvisors)
	req, err := http.NewRequest("GET", "/patch-advisors?page=sdasdasdasd", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchPatchAdvisors_JSONUnprocessableEntity3(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchPatchAdvisors)
	req, err := http.NewRequest("GET", "/patch-advisors?size=sdasdasdasd", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchPatchAdvisors_JSONUnprocessableEntity4(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchPatchAdvisors)
	req, err := http.NewRequest("GET", "/patch-advisors?window-time=sdasdasdasd", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchPatchAdvisors_JSONUnprocessableEntity5(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchPatchAdvisors)
	req, err := http.NewRequest("GET", "/patch-advisors?status=sdasdasdasd", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchPatchAdvisors_JSONUnprocessableEntity6(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchPatchAdvisors)
	req, err := http.NewRequest("GET", "/patch-advisors?older-than=sdasdasdasd", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchPatchAdvisors_JSONInternalServerError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	as.EXPECT().
		SearchPatchAdvisors("", "", false, -1, -1, utils.P("2019-05-05T14:02:03Z"), "", "", utils.MAX_TIME, "").
		Return(nil, aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchPatchAdvisors)
	req, err := http.NewRequest("GET", "/patch-advisors", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestSearchPatchAdvisors_XLSXSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Log: utils.NewLogger("TEST"),
	}

	expectedRes := []map[string]interface{}{
		map[string]interface{}{
			"CreatedAt":   utils.P("2020-04-07T08:52:59.82+02:00"),
			"Date":        utils.PDT("2012-04-16T02:00:00+02:00"),
			"Dbname":      "4wcqjn-ecf040bdfab7695ab332aef7401f185c",
			"Dbver":       "11.2.0.3.0 Enterprise Edition",
			"Description": "PSU 11.2.0.3.2",
			"Environment": "SVIL",
			"Hostname":    "publicitate-36d06ca83eafa454423d2097f4965517",
			"Location":    "Germany",
			"Status":      "KO",
			"_id":         utils.Str2oid("5e8c234b24f648a08585bd32"),
		},
		map[string]interface{}{
			"CreatedAt":   utils.P("2020-04-07T08:52:59.872+02:00"),
			"Date":        utils.PDT("2012-04-16T02:00:00+02:00"),
			"Dbname":      "ERCOLE",
			"Dbver":       "12.2.0.1.0 Enterprise Edition",
			"Description": "PSU 11.2.0.3.2",
			"Environment": "TST",
			"Hostname":    "test-db",
			"Location":    "Germany",
			"Status":      "KO",
			"_id":         utils.Str2oid("5e8c234b24f648a08585bd43"),
		},
	}

	as.EXPECT().
		SearchPatchAdvisors("foobar", "Hostname", true, -1, -1, utils.P("2019-03-05T14:02:03Z"), "Italy", "TST", utils.P("2020-06-10T11:54:59Z"), "KO").
		Return(expectedRes, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchPatchAdvisors)
	req, err := http.NewRequest("GET", "/patch-advisors?search=foobar&sort-by=Hostname&sort-desc=true&window-time=8&status=KO&location=Italy&environment=TST&older-than=2020-06-10T11%3A54%3A59Z", nil)
	require.NoError(t, err)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	sp, err := xlsx.Open(rr.Body)
	require.NoError(t, err)
	sh := sp.SheetByName("Patch_Advisor")
	require.NotNil(t, sh)
	assert.Equal(t, "PSU 11.2.0.3.2", sh.Cell(0, 1).String())
	assert.Equal(t, "publicitate-36d06ca83eafa454423d2097f4965517", sh.Cell(1, 1).String())
	assert.Equal(t, "4wcqjn-ecf040bdfab7695ab332aef7401f185c", sh.Cell(2, 1).String())
	assert.Equal(t, "11.2.0.3.0 Enterprise Edition", sh.Cell(3, 1).String())
	assert.Equal(t, utils.P("2012-04-16T00:00:00Z").String(), sh.Cell(4, 1).String())
	assert.Equal(t, "KO", sh.Cell(5, 1).String())

	assert.Equal(t, "PSU 11.2.0.3.2", sh.Cell(0, 2).String())
	assert.Equal(t, "test-db", sh.Cell(1, 2).String())
	assert.Equal(t, "ERCOLE", sh.Cell(2, 2).String())
	assert.Equal(t, "12.2.0.1.0 Enterprise Edition", sh.Cell(3, 2).String())
	assert.Equal(t, utils.P("2012-04-16T00:00:00Z").String(), sh.Cell(4, 2).String())
	assert.Equal(t, "KO", sh.Cell(5, 2).String())
}

func TestSearchPatchAdvisors_XLSXUnprocessableEntity1(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Log: utils.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchPatchAdvisors)
	req, err := http.NewRequest("GET", "/patch-advisors?sort-desc=dsasdasd", nil)
	require.NoError(t, err)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchPatchAdvisors_XLSXUnprocessableEntity2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Log: utils.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchPatchAdvisors)
	req, err := http.NewRequest("GET", "/patch-advisors?window-time=dsasdasd", nil)
	require.NoError(t, err)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchPatchAdvisors_XLSXUnprocessableEntity3(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Log: utils.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchPatchAdvisors)
	req, err := http.NewRequest("GET", "/patch-advisors?older-than=dsasdasd", nil)
	require.NoError(t, err)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchPatchAdvisors_XLSXUnprocessableEntity4(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Log: utils.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchPatchAdvisors)
	req, err := http.NewRequest("GET", "/patch-advisors?status=dsasdasd", nil)
	require.NoError(t, err)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchPatchAdvisors_XLSXInternalServerError1(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Log: utils.NewLogger("TEST"),
	}

	as.EXPECT().
		SearchPatchAdvisors("", "", false, -1, -1, utils.P("2019-05-05T14:02:03Z"), "", "", utils.MAX_TIME, "").
		Return(nil, aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchPatchAdvisors)
	req, err := http.NewRequest("GET", "/patch-advisors", nil)
	require.NoError(t, err)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}
func TestSearchPatchAdvisors_XLSXInternalServerError2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Log:     utils.NewLogger("TEST"),
	}

	expectedRes := []map[string]interface{}{
		map[string]interface{}{
			"CreatedAt":   utils.P("2020-04-07T08:52:59.82+02:00"),
			"Date":        utils.PDT("2012-04-16T02:00:00+02:00"),
			"Dbname":      "4wcqjn-ecf040bdfab7695ab332aef7401f185c",
			"Dbver":       "11.2.0.3.0 Enterprise Edition",
			"Description": "PSU 11.2.0.3.2",
			"Environment": "SVIL",
			"Hostname":    "publicitate-36d06ca83eafa454423d2097f4965517",
			"Location":    "Germany",
			"Status":      "KO",
			"_id":         utils.Str2oid("5e8c234b24f648a08585bd32"),
		},
		map[string]interface{}{
			"CreatedAt":   utils.P("2020-04-07T08:52:59.872+02:00"),
			"Date":        utils.PDT("2012-04-16T02:00:00+02:00"),
			"Dbname":      "ERCOLE",
			"Dbver":       "12.2.0.1.0 Enterprise Edition",
			"Description": "PSU 11.2.0.3.2",
			"Environment": "TST",
			"Hostname":    "test-db",
			"Location":    "Germany",
			"Status":      "KO",
			"_id":         utils.Str2oid("5e8c234b24f648a08585bd43"),
		},
	}

	as.EXPECT().
		SearchPatchAdvisors("foobar", "Hostname", true, -1, -1, utils.P("2019-03-05T14:02:03Z"), "Italy", "TST", utils.P("2020-06-10T11:54:59Z"), "KO").
		Return(expectedRes, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchPatchAdvisors)
	req, err := http.NewRequest("GET", "/patch-advisors?search=foobar&sort-by=Hostname&sort-desc=true&window-time=8&status=KO&location=Italy&environment=TST&older-than=2020-06-10T11%3A54%3A59Z", nil)
	require.NoError(t, err)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}