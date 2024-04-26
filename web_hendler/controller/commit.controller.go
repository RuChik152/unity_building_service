package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"web_hendler/db"

	"github.com/gorilla/mux"
)

func CommitController(r *mux.Router) {
	router := r.PathPrefix("/commit").Subrouter()

	router.HandleFunc("/info/{id}", func(res http.ResponseWriter, req *http.Request) {
		params := mux.Vars(req)
		commitID, err := strconv.Atoi(params["id"])
		if err != nil {
			log.Println("Ошибка преобразования")
			res.WriteHeader(500)
		}

		commit, err := db.GetCommitData(commitID, "commits")
		if err != nil {
			log.Println(err)
			res.WriteHeader(http.StatusNotFound)
		} else {
			res.Header().Set("Content-Type", "application/json")

			responseData, _ := json.Marshal(commit)
			res.Write(responseData)
		}

	})
}
