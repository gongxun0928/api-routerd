// SPDX-License-Identifier: Apache-2.0

package timedate

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func RouterGetTimeDate(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	property := vars["property"]

	switch r.Method {
	case "GET":

		timedate := new(TimeDate)
		if property != "" {
			timedate.Property = property
		}

		err := GetTimeDate(rw, property)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}

		break
	}
}

func RouterSetTimeDate(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PUT":

		timedate := new(TimeDate)
		err := json.NewDecoder(r.Body).Decode(&timedate)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Println(timedate)
		err = timedate.SetTimeDate()
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break
	}
}

func RegisterRouterTimeDate(router *mux.Router) {
	s := router.PathPrefix("/timedate").Subrouter().StrictSlash(false)
	s.HandleFunc("", RouterGetTimeDate)
	s.HandleFunc("/get/{property}", RouterGetTimeDate)
	s.HandleFunc("/set", RouterSetTimeDate)
}
