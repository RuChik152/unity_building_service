package main

import (
	"fmt"
	"net/http"
	"web_hendler/controller"

	"github.com/gorilla/mux"
)

func appRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World")
	})

	controller.BuildingController(router)

	return router
}
