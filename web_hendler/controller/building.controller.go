package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"web_hendler/bot"
	"web_hendler/service"

	"github.com/gorilla/mux"
)

func BuildingController(r *mux.Router) {
	buildRouter := r.PathPrefix("/building").Subrouter()

	buildRouter.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {

		if req.Method == http.MethodPost {
			if checkDataCommit(res, req) {
				res.WriteHeader(http.StatusOK)
				service.Manager()
			}

		}

	})

}

func checkDataCommit(res http.ResponseWriter, req *http.Request) bool {

	var payload map[string]interface{}

	decodeDataRequest := json.NewDecoder(req.Body)
	err := decodeDataRequest.Decode(&payload)
	if err != nil {
		log.Println("Ошибка декодирования данных запроса тела входящего запроса")
		http.Error(res, "Error decoding JSON", http.StatusBadRequest)
		return false
	}

	var coments string

	event := req.Header.Get("X-GitHub-Event")
	ref, okRef := payload["ref"].(string)
	commits, okCommits := payload["commits"].([]interface{})
	author, okAuthor := getValueByPath(payload, []string{"pusher", "name"})

	if okRef && okCommits && okAuthor {
		if ref == "refs/heads/main" && event == "push" {

			for _, value := range commits {
				if commit, ok := value.(map[string]interface{}); ok {
					if message, ok := commit["message"].(string); ok {
						coments += "\n" + message
					}
				}
			}

			bot.CommitMsg.Event = "commit"
			bot.CommitMsg.AUTHOR = fmt.Sprintf("Автор: %s", author)
			bot.CommitMsg.MESSAGE = fmt.Sprintf("Комментарий коммита: %s", coments)
			commitData, err := json.Marshal(bot.CommitMsg)
			if err != nil {
				log.Println("Ошибка декодирования: ", err)
			} else {
				bot.SendMessageBot(string(commitData), "#pipline_event")
			}

			return true
		}
		return false
	} else {
		return false
	}

}

func getValueByPath(data interface{}, path []string) (interface{}, bool) {
	if len(path) == 0 {
		return nil, false
	}

	switch d := data.(type) {
	case map[string]interface{}:
		{
			value, ok := d[path[0]]
			if !ok {
				return nil, false
			}

			if len(path) == 1 {
				return value, true
			}

			return getValueByPath(value, path[1:])
		}
	case []interface{}:
		{
			index, err := strconv.Atoi(path[0])
			if err != nil {
				return nil, false
			}

			if len(d) > index {
				value := d[index]
				if len(path) == 1 {
					return value, true
				}

				return getValueByPath(value, path[1:])
			}

			return nil, false
		}
	default:
		return nil, false
	}
}
