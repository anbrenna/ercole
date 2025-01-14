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

package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang/gddo/httputil"
	"github.com/gorilla/context"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/ercole-io/ercole/v2/utils/exutils"
)

// SearchOracleDatabaseAddms search addms data using the filters in the request
func (ctrl *APIController) SearchOracleDatabaseAddms(w http.ResponseWriter, r *http.Request) {
	choice := httputil.NegotiateContentType(r, []string{"application/json", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"}, "application/json")

	switch choice {
	case "application/json":
		ctrl.SearchOracleDatabaseAddmsJSON(w, r)
	case "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		ctrl.SearchOracleDatabaseAddmsXLSX(w, r)
	}
}

// SearchOracleDatabaseAddmsJSON search addms data using the filters in the request returning it in JSON format
func (ctrl *APIController) SearchOracleDatabaseAddmsJSON(w http.ResponseWriter, r *http.Request) {
	var search, sortBy, location, environment string

	var sortDesc bool

	var pageNumber, pageSize int

	var olderThan time.Time

	var err error
	//parse the query params
	search = r.URL.Query().Get("search")
	sortBy = r.URL.Query().Get("sort-by")

	if sortDesc, err = utils.Str2bool(r.URL.Query().Get("sort-desc"), false); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	if pageNumber, err = utils.Str2int(r.URL.Query().Get("page"), -1); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	if pageSize, err = utils.Str2int(r.URL.Query().Get("size"), -1); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	location = r.URL.Query().Get("location")
	if location == "" {
		user := context.Get(r, "user")
		locations, errLocation := ctrl.Service.ListLocations(user)

		if errLocation != nil {
			utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, errLocation)
			return
		}

		location = strings.Join(locations, ",")
	}

	environment = r.URL.Query().Get("environment")

	if olderThan, err = utils.Str2time(r.URL.Query().Get("older-than"), utils.MAX_TIME); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	//get the data
	addms, err := ctrl.Service.SearchOracleDatabaseAddms(search, sortBy, sortDesc, pageNumber, pageSize, location, environment, olderThan)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	if pageNumber == -1 || pageSize == -1 {
		//Write the data
		utils.WriteJSONResponse(w, http.StatusOK, addms)
	} else {
		//Write the data
		utils.WriteJSONResponse(w, http.StatusOK, addms[0])
	}
}

// SearchOracleDatabaseAddmsXLSX search addms data using the filters in the request returning it in XLSX format
func (ctrl *APIController) SearchOracleDatabaseAddmsXLSX(w http.ResponseWriter, r *http.Request) {
	var search, location, environment string

	var olderThan time.Time

	var err error

	search = r.URL.Query().Get("search")

	location = r.URL.Query().Get("location")
	if location == "" {
		user := context.Get(r, "user")
		locations, errLocation := ctrl.Service.ListLocations(user)

		if errLocation != nil {
			utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, errLocation)
			return
		}

		location = strings.Join(locations, ",")
	}

	environment = r.URL.Query().Get("environment")

	if olderThan, err = utils.Str2time(r.URL.Query().Get("older-than"), utils.MAX_TIME); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	//get the data
	addms, err := ctrl.Service.SearchOracleDatabaseAddms(search, "benefit", true, -1, -1, location, environment, olderThan)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	sheet := "Addm"
	headers := []string{
		"Performance Impact",
		"Hostname",
		"Database",
		"Finding",
		"Recommendation",
		"Action",
		"Environment",
	}

	file, err := exutils.NewXLSX(ctrl.Config, sheet, headers...)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError,
			utils.NewError(err, "Can't create new xlsx"))
		return
	}

	for i, val := range addms {
		file.SetCellValue(sheet, fmt.Sprintf("A%d", i+2), val["benefit"])        //Benefit column
		file.SetCellValue(sheet, fmt.Sprintf("B%d", i+2), val["hostname"])       //Hostname column
		file.SetCellValue(sheet, fmt.Sprintf("C%d", i+2), val["dbname"])         //Dbname column
		file.SetCellValue(sheet, fmt.Sprintf("D%d", i+2), val["finding"])        //Finding column
		file.SetCellValue(sheet, fmt.Sprintf("E%d", i+2), val["recommendation"]) //Recommendation column
		file.SetCellValue(sheet, fmt.Sprintf("F%d", i+2), val["action"])         //Action column
		file.SetCellValue(sheet, fmt.Sprintf("G%d", i+2), val["environment"])    //Environment column
	}

	utils.WriteXLSXResponse(w, file)
}

// SearchOracleDatabaseSegmentAdvisors search segment advisors data using the filters in the request
func (ctrl *APIController) SearchOracleDatabaseSegmentAdvisors(w http.ResponseWriter, r *http.Request) {
	choice := httputil.NegotiateContentType(r, []string{"application/json", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"}, "application/json")

	switch choice {
	case "application/json":
		ctrl.SearchOracleDatabaseSegmentAdvisorsJSON(w, r)
	case "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		ctrl.SearchOracleDatabaseSegmentAdvisorsXLSX(w, r)
	}
}

// SearchOracleDatabaseSegmentAdvisorsJSON search segment advisors data using the filters in the request returning it in JSON format
func (ctrl *APIController) SearchOracleDatabaseSegmentAdvisorsJSON(w http.ResponseWriter, r *http.Request) {
	var search, sortBy, location, environment string

	var sortDesc bool

	var olderThan time.Time

	var err error

	search = r.URL.Query().Get("search")
	sortBy = r.URL.Query().Get("sort-by")

	if sortDesc, err = utils.Str2bool(r.URL.Query().Get("sort-desc"), false); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	location = r.URL.Query().Get("location")
	if location == "" {
		user := context.Get(r, "user")
		locations, errLocation := ctrl.Service.ListLocations(user)

		if errLocation != nil {
			utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, errLocation)
			return
		}

		location = strings.Join(locations, ",")
	}

	environment = r.URL.Query().Get("environment")

	if olderThan, err = utils.Str2time(r.URL.Query().Get("older-than"), utils.MAX_TIME); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	segmentAdvisors, err := ctrl.Service.SearchOracleDatabaseSegmentAdvisors(search, sortBy, sortDesc, location, environment, olderThan)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	resp := map[string]interface{}{
		"segmentAdvisors": segmentAdvisors,
	}
	utils.WriteJSONResponse(w, http.StatusOK, resp)
}

// SearchOracleDatabaseSegmentAdvisorsXLSX search segment advisors data using the filters in the request returning it in XLSX format
func (ctrl *APIController) SearchOracleDatabaseSegmentAdvisorsXLSX(w http.ResponseWriter, r *http.Request) {
	filter, err := dto.GetGlobalFilter(r)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	if filter.Location == "" {
		user := context.Get(r, "user")
		locations, errLocation := ctrl.Service.ListLocations(user)

		if errLocation != nil {
			utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, errLocation)
			return
		}

		filter.Location = strings.Join(locations, ",")
	}

	xlsx, err := ctrl.Service.SearchOracleDatabaseSegmentAdvisorsAsXLSX(*filter)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteXLSXResponse(w, xlsx)
}

// SearchOracleDatabasePatchAdvisors search patch advisors data using the filters in the request
func (ctrl *APIController) SearchOracleDatabasePatchAdvisors(w http.ResponseWriter, r *http.Request) {
	choice := httputil.NegotiateContentType(r, []string{"application/json", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"}, "application/json")

	switch choice {
	case "application/json":
		ctrl.SearchOracleDatabasePatchAdvisorsJSON(w, r)
	case "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		ctrl.SearchOracleDatabasePatchAdvisorsXLSX(w, r)
	}
}

// SearchOracleDatabasePatchAdvisorsJSON search patch advisors data using the filters in the request returning it in JSON format
func (ctrl *APIController) SearchOracleDatabasePatchAdvisorsJSON(w http.ResponseWriter, r *http.Request) {
	var search, sortBy, location, environment, status string

	var sortDesc bool

	var pageNumber, pageSize, windowTime int

	var olderThan time.Time

	var err error
	//parse the query params
	search = r.URL.Query().Get("search")
	sortBy = r.URL.Query().Get("sort-by")

	if sortDesc, err = utils.Str2bool(r.URL.Query().Get("sort-desc"), false); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	if pageNumber, err = utils.Str2int(r.URL.Query().Get("page"), -1); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	if pageSize, err = utils.Str2int(r.URL.Query().Get("size"), -1); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	if windowTime, err = utils.Str2int(r.URL.Query().Get("window-time"), 6); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	status = r.URL.Query().Get("status")
	if status != "" && status != "OK" && status != "KO" {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, utils.NewError(errors.New("invalid status"), "Invalid  status"))
		return
	}

	location = r.URL.Query().Get("location")
	if location == "" {
		user := context.Get(r, "user")
		locations, errLocation := ctrl.Service.ListLocations(user)

		if errLocation != nil {
			utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, errLocation)
			return
		}

		location = strings.Join(locations, ",")
	}

	environment = r.URL.Query().Get("environment")

	if olderThan, err = utils.Str2time(r.URL.Query().Get("older-than"), utils.MAX_TIME); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	//get the data
	patchAdvisorResponse, err := ctrl.Service.SearchOracleDatabasePatchAdvisors(search, sortBy, sortDesc, pageNumber, pageSize, ctrl.TimeNow().AddDate(0, -windowTime, 0), location, environment, olderThan, status)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, patchAdvisorResponse)
}

// SearchOracleDatabasePatchAdvisorsXLSX search patch advisors data using the filters in the request returning it in XLSX format
func (ctrl *APIController) SearchOracleDatabasePatchAdvisorsXLSX(w http.ResponseWriter, r *http.Request) {
	var windowTime int

	filter, err := dto.GetGlobalFilter(r)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	if filter.Location == "" {
		user := context.Get(r, "user")
		locations, errLocation := ctrl.Service.ListLocations(user)

		if errLocation != nil {
			utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, errLocation)
			return
		}

		filter.Location = strings.Join(locations, ",")
	}

	if windowTime, err = utils.Str2int(r.URL.Query().Get("window-time"), 6); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	xlsx, err := ctrl.Service.SearchOracleDatabasePatchAdvisorsAsXLSX(ctrl.TimeNow().AddDate(0, -windowTime, 0), *filter)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteXLSXResponse(w, xlsx)
}

// SearchOracleDatabases search databases data using the filters in the request
func (ctrl *APIController) SearchOracleDatabases(w http.ResponseWriter, r *http.Request) {
	choice := httputil.NegotiateContentType(r, []string{"application/json", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"}, "application/json")

	filter, err := dto.GetSearchOracleDatabasesFilter(r)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	if filter.Location == "" {
		user := context.Get(r, "user")
		locations, errLocation := ctrl.Service.ListLocations(user)

		if errLocation != nil {
			utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, errLocation)
			return
		}

		filter.Location = strings.Join(locations, ",")
	}

	switch choice {
	case "application/json":
		ctrl.SearchOracleDatabasesJSON(w, r, *filter)
	case "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		ctrl.SearchOracleDatabasesXLSX(w, r, *filter)
	}
}

// SearchOracleDatabasesJSON search databases data using the filters in the request returning it in JSON
func (ctrl *APIController) SearchOracleDatabasesJSON(w http.ResponseWriter, r *http.Request, filter dto.SearchOracleDatabasesFilter) {
	databases, err := ctrl.Service.SearchOracleDatabases(filter)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	if filter.PageNumber == -1 || filter.PageSize == -1 {
		utils.WriteJSONResponse(w, http.StatusOK, databases.Content)
	} else {
		utils.WriteJSONResponse(w, http.StatusOK, databases)
	}
}

// SearchOracleDatabasesXLSX search databases data using the filters in the request returning it in XLSX
func (ctrl *APIController) SearchOracleDatabasesXLSX(w http.ResponseWriter, r *http.Request, filter dto.SearchOracleDatabasesFilter) {
	file, err := ctrl.Service.SearchOracleDatabasesAsXLSX(filter)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteXLSXResponse(w, file)
}

// SearchOracleDatabaseUsedLicenses search licenses consumed by the hosts using the filters in the request
func (ctrl *APIController) SearchOracleDatabaseUsedLicenses(w http.ResponseWriter, r *http.Request) {
	var sortBy, location, environment string

	var sortDesc bool

	var pageNumber, pageSize int

	var olderThan time.Time

	var err error

	sortBy = r.URL.Query().Get("sort-by")

	if sortDesc, err = utils.Str2bool(r.URL.Query().Get("sort-desc"), false); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	if pageNumber, err = utils.Str2int(r.URL.Query().Get("page"), -1); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	if pageSize, err = utils.Str2int(r.URL.Query().Get("size"), -1); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	location = r.URL.Query().Get("location")
	environment = r.URL.Query().Get("environment")

	if olderThan, err = utils.Str2time(r.URL.Query().Get("older-than"), utils.MAX_TIME); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	response, err := ctrl.Service.SearchOracleDatabaseUsedLicenses("", sortBy, sortDesc, pageNumber, pageSize, location, environment, olderThan)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	if pageNumber == -1 || pageSize == -1 {
		utils.WriteJSONResponse(w, http.StatusOK, response.Content)
	} else {
		utils.WriteJSONResponse(w, http.StatusOK, response)
	}
}
