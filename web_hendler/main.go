package main

import (
	"log"
	"net/http"
	"os"
	"web_hendler/service"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	} else {
		log.Print("Success, .env file found")
	}

	service.PROJECT_FOLDER, _ = os.LookupEnv("PATH_PROJECT")
}

func main() {

	server := &http.Server{
		Addr:    ":8080",
		Handler: appRouter(),
	}

	log.Printf("<<SERVER START>>\n http://localhost%s", server.Addr)
	log.Fatal(server.ListenAndServe())
}
