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

	"github.com/gorilla/mux"
)

// SetupRoutesForThunderController setup the routes of the router using the handler in the controller as http handler
func SetupRoutesForThunderController(router *mux.Router, ctrl ThunderControllerInterface) {
	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Pong"))
	})

	router = router.NewRoute().Subrouter()

	setupProtectedRoutes(router, ctrl)
}

func setupProtectedRoutes(router *mux.Router, ctrl ThunderControllerInterface) {
	router.HandleFunc("/oracle-cloud/recommendations/{ids}", ctrl.GetOciRecommendations).Methods("GET")
	router.HandleFunc("/oracle-cloud/configurations", ctrl.GetOciProfiles).Methods("GET")
	router.HandleFunc("/oracle-cloud/configurations", ctrl.AddOciProfile).Methods("POST")
	router.HandleFunc("/oracle-cloud/configurations/{id}", ctrl.UpdateOciProfile).Methods("PUT")
	router.HandleFunc("/oracle-cloud/configurations/{id}", ctrl.DeleteOciProfile).Methods("DELETE")
}