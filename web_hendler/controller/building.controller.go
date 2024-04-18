package controller

import (
	"net/http"
	"web_hendler/service"

	"github.com/gorilla/mux"
)

func BuildingController(r *mux.Router) {
	buildRouter := r.PathPrefix("/building").Subrouter()

	buildRouter.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {

		if req.Method == http.MethodPost {
			service.Manager()
		}

	})

}
