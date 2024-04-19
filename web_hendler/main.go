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
	service.PICO_APP_ID, _ = os.LookupEnv("PICO_APP_ID")
	service.PICO_APP_SECRET, _ = os.LookupEnv("PICO_APP_SECRET")
	service.OCULUS_APP_ID, _ = os.LookupEnv("OCULUS_APP_ID")
	service.OCULUS_APP_SECRET, _ = os.LookupEnv("OCULUS_APP_SECRET")
	service.NAME_KEYSTORE, _ = os.LookupEnv("KEYSTORE_NAME")

}

func main() {
	server := &http.Server{
		Addr:    ":8080",
		Handler: appRouter(),
	}

	log.Printf("<<SERVER START>>\n http://localhost%s", server.Addr)
	log.Fatal(server.ListenAndServe())
}
