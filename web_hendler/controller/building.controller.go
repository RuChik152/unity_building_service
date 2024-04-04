package controller

import (
	"net/http"
	"web_hendler/service"

	"github.com/gorilla/mux"
)

func BuildingController(r *mux.Router) {
	buildRouter := r.PathPrefix("/building").Subrouter()

	// buildRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprintf(w, "TEST BUILDING ENDPOINT")
	// })

	buildRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		service.Manager()
	})

}
