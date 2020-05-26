// Copyright (c) 2019 Sorint.lab S.p.A.
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
	"net/http"
	"time"

	"github.com/amreo/ercole-services/api-service/database"
	"github.com/amreo/ercole-services/utils"
	"github.com/golang/gddo/httputil"
	"github.com/gorilla/mux"
	"github.com/plandem/xlsx"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SearchHosts search hosts data using the filters in the request
func (ctrl *APIController) SearchHosts(w http.ResponseWriter, r *http.Request) {
	choiche := httputil.NegotiateContentType(r, []string{"application/json", "application/vnd.oracle.lms+vnd.openxmlformats-officedocument.spreadsheetml.sheet", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", "application/vnd.ercole.mongohostdata+json"}, "application/json")

	switch choiche {
	case "application/json":
		ctrl.SearchHostsJSON(w, r)
	case "application/vnd.oracle.lms+vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		ctrl.SearchHostsLMS(w, r)
	case "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		ctrl.SearchHostsXLSX(w, r)
	}
}

// SearchHostsJSON search hosts data using the filters in the request returning it in JSON
func (ctrl *APIController) SearchHostsJSON(w http.ResponseWriter, r *http.Request) {
	var mode string
	var search string
	var sortBy string
	var sortDesc bool
	var pageNumber int
	var pageSize int
	var location string
	var environment string
	var olderThan time.Time
	var searchHostsFilters database.SearchHostsFilters
	var err utils.AdvancedErrorInterface
	//parse the query params
	mode = r.URL.Query().Get("mode")
	if mode == "" {
		mode = "full"
	} else if mode != "full" && mode != "summary" && mode != "lms" && mode != "mhd" {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, utils.NewAdvancedErrorPtr(errors.New("Invalid mode value"), http.StatusText(http.StatusUnprocessableEntity)))
		return
	}

	search = r.URL.Query().Get("search")

	searchHostsFilters, err = ctrl.GetSearchHostFilters(r)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

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

	//get the data
	hosts, err := ctrl.Service.SearchHosts(mode, search, searchHostsFilters, sortBy, sortDesc, pageNumber, pageSize, location, environment, olderThan)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	if pageNumber == -1 || pageSize == -1 {
		//Write the data
		utils.WriteJSONResponse(w, http.StatusOK, hosts)
	} else {
		//Write the data
		utils.WriteJSONResponse(w, http.StatusOK, hosts[0])
	}
}

// SearchHostsLMS search hosts data using the filters in the request returning it in LMS+XLSX
func (ctrl *APIController) SearchHostsLMS(w http.ResponseWriter, r *http.Request) {
	var search string
	var sortBy string
	var sortDesc bool
	var location string
	var environment string
	var olderThan time.Time
	var searchHostsFilters database.SearchHostsFilters

	var aerr utils.AdvancedErrorInterface
	//parse the query params
	search = r.URL.Query().Get("search")

	searchHostsFilters, aerr = ctrl.GetSearchHostFilters(r)
	if aerr != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, aerr)
		return
	}

	sortBy = r.URL.Query().Get("sort-by")
	if sortDesc, aerr = utils.Str2bool(r.URL.Query().Get("sort-desc"), false); aerr != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, aerr)
		return
	}

	location = r.URL.Query().Get("location")
	environment = r.URL.Query().Get("environment")

	if olderThan, aerr = utils.Str2time(r.URL.Query().Get("older-than"), utils.MAX_TIME); aerr != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, aerr)
		return
	}

	//get the data
	hosts, aerr := ctrl.Service.SearchHosts("lms", search, searchHostsFilters, sortBy, sortDesc, -1, -1, location, environment, olderThan)
	if aerr != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, aerr)
		return
	}

	//Open the sheet
	sheets, err := xlsx.Open(ctrl.Config.ResourceFilePath + "/templates/template_lms.xlsm")
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, utils.NewAdvancedErrorPtr(err, "READ_TEMPLATE"))
		return
	}

	sheet := sheets.SheetByName("Database_&_EBS")

	i := 0
	//Add the data to the sheet
	for _, val := range hosts {
		sheet.Cell(0, i+3).SetText(val["PhysicalServerName"])
		sheet.Cell(1, i+3).SetText(val["VirtualServerName"])
		sheet.Cell(2, i+3).SetText(val["VirtualizationTechnology"])
		sheet.Cell(3, i+3).SetText(val["DBInstanceName"])
		sheet.Cell(4, i+3).SetText(val["PluggableDatabaseName"])
		sheet.Cell(5, i+3).SetText(val["ConnectString"])
		sheet.Cell(7, i+3).SetText(val["ProductVersion"])
		sheet.Cell(8, i+3).SetText(val["ProductEdition"])
		sheet.Cell(9, i+3).SetText(val["Environment"])
		sheet.Cell(10, i+3).SetText(val["Features"])
		sheet.Cell(11, i+3).SetText(val["RacNodeNames"])
		sheet.Cell(12, i+3).SetText(val["ProcessorModel"])
		sheet.Cell(13, i+3).SetInt(int(val["Processors"].(float64)))
		sheet.Cell(14, i+3).SetInt(int(val["CoresPerProcessor"].(float64)))
		sheet.Cell(15, i+3).SetInt(int(val["PhysicalCores"].(float64)))
		sheet.Cell(16, i+3).SetInt(int(val["ThreadsPerCore"].(int32)))
		sheet.Cell(17, i+3).SetText(val["ProcessorSpeed"])
		sheet.Cell(18, i+3).SetText(val["ServerPurchaseDate"])
		sheet.Cell(19, i+3).SetText(val["OperatingSystem"])
		sheet.Cell(20, i+3).SetText(val["Notes"])
		i++
	}

	//Write it to the response
	utils.WriteXLSXResponse(w, sheets)
}

// SearchHostsXLSX search hosts data using the filters in the request returning it in XLSX
func (ctrl *APIController) SearchHostsXLSX(w http.ResponseWriter, r *http.Request) {
	var search string
	var sortBy string
	var sortDesc bool
	var location string
	var environment string
	var olderThan time.Time
	var searchHostsFilters database.SearchHostsFilters

	var aerr utils.AdvancedErrorInterface
	//parse the query params
	search = r.URL.Query().Get("search")

	searchHostsFilters, aerr = ctrl.GetSearchHostFilters(r)
	if aerr != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, aerr)
		return
	}

	sortBy = r.URL.Query().Get("sort-by")
	if sortDesc, aerr = utils.Str2bool(r.URL.Query().Get("sort-desc"), false); aerr != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, aerr)
		return
	}

	location = r.URL.Query().Get("location")
	environment = r.URL.Query().Get("environment")

	if olderThan, aerr = utils.Str2time(r.URL.Query().Get("older-than"), utils.MAX_TIME); aerr != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, aerr)
		return
	}

	//get the data
	hosts, aerr := ctrl.Service.SearchHosts("summary", search, searchHostsFilters, sortBy, sortDesc, -1, -1, location, environment, olderThan)
	if aerr != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, aerr)
		return
	}

	//Open the sheet
	sheets, err := xlsx.Open(ctrl.Config.ResourceFilePath + "/templates/template_hosts.xlsx")
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, utils.NewAdvancedErrorPtr(err, "READ_TEMPLATE"))
		return
	}

	sheet := sheets.SheetByName("Hosts")

	//Add the data to the sheet
	for i, val := range hosts {
		sheet.Cell(0, i+1).SetText(val["Hostname"])
		sheet.Cell(1, i+1).SetText(val["Environment"])
		sheet.Cell(2, i+1).SetText(val["HostType"])
		if val["Cluster"] != nil && val["PhysicalHost"] != nil {
			sheet.Cell(3, i+1).SetText(val["Cluster"])
			sheet.Cell(4, i+1).SetText(val["PhysicalHost"])
		}
		sheet.Cell(5, i+1).SetText(val["Version"])
		sheet.Cell(6, i+1).SetText(val["CreatedAt"].(primitive.DateTime).Time().UTC().String())
		sheet.Cell(7, i+1).SetText(val["Databases"])
		sheet.Cell(8, i+1).SetText(val["OS"])
		sheet.Cell(9, i+1).SetText(val["Kernel"])
		sheet.Cell(10, i+1).SetBool(val["OracleCluster"].(bool))
		sheet.Cell(11, i+1).SetBool(val["SunCluster"].(bool))
		sheet.Cell(12, i+1).SetBool(val["VeritasCluster"].(bool))
		sheet.Cell(13, i+1).SetBool(val["Virtual"].(bool))
		sheet.Cell(14, i+1).SetText(val["Type"])
		sheet.Cell(15, i+1).SetInt(int(val["CPUThreads"].(float64)))
		sheet.Cell(16, i+1).SetInt(int(val["CPUCores"].(float64)))
		sheet.Cell(17, i+1).SetInt(int(val["Socket"].(float64)))
		sheet.Cell(18, i+1).SetInt(int(val["MemTotal"].(float64)))
		sheet.Cell(19, i+1).SetInt(int(val["SwapTotal"].(float64)))
		sheet.Cell(20, i+1).SetText(val["CPUModel"])
	}

	//Write it to the response
	utils.WriteXLSXResponse(w, sheets)
}

// GetSearchHostFilters return the host search filters in the request
func (ctrl *APIController) GetSearchHostFilters(r *http.Request) (database.SearchHostsFilters, utils.AdvancedErrorInterface) {
	var aerr utils.AdvancedErrorInterface

	filters := database.SearchHostsFilters{}

	filters.Hostname = r.URL.Query().Get("hostname")
	filters.Database = r.URL.Query().Get("database")
	filters.Asset = r.URL.Query().Get("asset")
	filters.HardwareAbstractionTechnology = r.URL.Query().Get("hardware-abstraction-technology")
	if r.URL.Query().Get("cluster") == "NULL" {
		filters.Cluster = nil
	} else {
		filters.Cluster = new(string)
		*filters.Cluster = r.URL.Query().Get("cluster")
	}
	filters.PhysicalHost = r.URL.Query().Get("physical-host")
	filters.OperatingSystem = r.URL.Query().Get("operating-system")
	filters.Kernel = r.URL.Query().Get("kernel")
	if filters.LTEMemoryTotal, aerr = utils.Str2float32(r.URL.Query().Get("memory-total-lte"), -1); aerr != nil {
		return database.SearchHostsFilters{}, aerr
	}
	if filters.GTEMemoryTotal, aerr = utils.Str2float32(r.URL.Query().Get("memory-total-gte"), -1); aerr != nil {
		return database.SearchHostsFilters{}, aerr
	}
	if filters.LTESwapTotal, aerr = utils.Str2float32(r.URL.Query().Get("swap-total-lte"), -1); aerr != nil {
		return database.SearchHostsFilters{}, aerr
	}
	if filters.GTESwapTotal, aerr = utils.Str2float32(r.URL.Query().Get("swap-total-gte"), -1); aerr != nil {
		return database.SearchHostsFilters{}, aerr
	}
	if r.URL.Query().Get("is-member-of-cluster") == "" {
		filters.IsMemberOfCluster = nil
	} else {
		filters.IsMemberOfCluster = new(bool)
		if *filters.IsMemberOfCluster, aerr = utils.Str2bool(r.URL.Query().Get("is-member-of-cluster"), false); aerr != nil {
			return database.SearchHostsFilters{}, aerr
		}
	}
	filters.CPUModel = r.URL.Query().Get("cpu-model")
	if filters.LTECPUCores, aerr = utils.Str2int(r.URL.Query().Get("cpu-cores-lte"), -1); aerr != nil {
		return database.SearchHostsFilters{}, aerr
	}
	if filters.GTECPUCores, aerr = utils.Str2int(r.URL.Query().Get("cpu-cores-gte"), -1); aerr != nil {
		return database.SearchHostsFilters{}, aerr
	}
	if filters.LTECPUThreads, aerr = utils.Str2int(r.URL.Query().Get("cpu-threads-lte"), -1); aerr != nil {
		return database.SearchHostsFilters{}, aerr
	}
	if filters.GTECPUThreads, aerr = utils.Str2int(r.URL.Query().Get("cpu-threads-gte"), -1); aerr != nil {
		return database.SearchHostsFilters{}, aerr
	}
	return filters, nil
}

// GetHost return all'informations about the host requested in the id path variable
func (ctrl *APIController) GetHost(w http.ResponseWriter, r *http.Request) {
	choiche := httputil.NegotiateContentType(r, []string{"application/json", "application/vnd.ercole.mongohostdata+json"}, "application/json")

	switch choiche {
	case "application/json":
		ctrl.GetHostJSON(w, r)
	case "application/vnd.ercole.mongohostdata+json":
		ctrl.GetHostMongoJSON(w, r)
	}
}

// GetHostJSON return all'informations about the host requested in the id path variable
func (ctrl *APIController) GetHostJSON(w http.ResponseWriter, r *http.Request) {
	var olderThan time.Time
	var err utils.AdvancedErrorInterface

	hostname := mux.Vars(r)["hostname"]

	if olderThan, err = utils.Str2time(r.URL.Query().Get("older-than"), utils.MAX_TIME); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	//get the data
	host, err := ctrl.Service.GetHost(hostname, olderThan, false)
	if err == utils.AerrHostNotFound {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusNotFound, err)
		return
	} else if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	//Write the data
	utils.WriteJSONResponse(w, http.StatusOK, host)
}

// GetHostMongoJSON return all'informations about the host requested in the id path variable
func (ctrl *APIController) GetHostMongoJSON(w http.ResponseWriter, r *http.Request) {
	var olderThan time.Time
	var aerr utils.AdvancedErrorInterface

	hostname := mux.Vars(r)["hostname"]

	if olderThan, aerr = utils.Str2time(r.URL.Query().Get("older-than"), utils.MAX_TIME); aerr != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, aerr)
		return
	}

	//get the data
	host, aerr := ctrl.Service.GetHost(hostname, olderThan, true)
	if aerr == utils.AerrHostNotFound {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusNotFound, aerr)
		return
	} else if aerr != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, aerr)
		return
	}

	//Write the response
	utils.WriteExtJSONResponse(ctrl.Log, w, http.StatusOK, host)
}

// ListLocations list locations using the filters in the request
func (ctrl *APIController) ListLocations(w http.ResponseWriter, r *http.Request) {
	var location string
	var environment string
	var olderThan time.Time

	var err utils.AdvancedErrorInterface
	//parse the query params
	location = r.URL.Query().Get("location")
	environment = r.URL.Query().Get("environment")

	if olderThan, err = utils.Str2time(r.URL.Query().Get("older-than"), utils.MAX_TIME); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	//get the data
	locations, err := ctrl.Service.ListLocations(location, environment, olderThan)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	//Write the data
	utils.WriteJSONResponse(w, http.StatusOK, locations)
}

// ListEnvironments list the environments using the filters in the request
func (ctrl *APIController) ListEnvironments(w http.ResponseWriter, r *http.Request) {
	var location string
	var environment string
	var olderThan time.Time

	var err utils.AdvancedErrorInterface
	//parse the query params
	location = r.URL.Query().Get("location")
	environment = r.URL.Query().Get("environment")

	if olderThan, err = utils.Str2time(r.URL.Query().Get("older-than"), utils.MAX_TIME); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	//get the data
	environments, err := ctrl.Service.ListEnvironments(location, environment, olderThan)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	//Write the data
	utils.WriteJSONResponse(w, http.StatusOK, environments)
}

// ArchiveHost archive the specified host in the request
func (ctrl *APIController) ArchiveHost(w http.ResponseWriter, r *http.Request) {
	if ctrl.Config.APIService.ReadOnly {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusForbidden, utils.NewAdvancedErrorPtr(errors.New("The API is disabled because the service is put in read-only mode"), "FORBIDDEN_REQUEST"))
		return
	}

	//Get the id from the path variable
	hostname := mux.Vars(r)["hostname"]

	//set the value
	aerr := ctrl.Service.ArchiveHost(hostname)
	if aerr == utils.AerrHostNotFound {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusNotFound, aerr)
	} else if aerr != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, aerr)
		return
	}

	//Write the data
	utils.WriteJSONResponse(w, http.StatusOK, nil)
}
