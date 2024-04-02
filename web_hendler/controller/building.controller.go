package controller

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func BuildingController(r *mux.Router) {
	buildRouter := r.PathPrefix("/building").Subrouter()

	buildRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "TEST BUILDING ENDPOINT")
	})

}
