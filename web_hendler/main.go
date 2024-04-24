package main

import (
	"log"
	"net/http"
	"os"
	"path"
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
	if service.PROJECT_FOLDER != "" {
		log.Println("PATH_PROJECT:", service.PROJECT_FOLDER)
	} else {
		log.Println("Ошибка!!! Не установлен путь к каталогу с проектом")
		os.Exit(1)
	}

	service.PICO_APP_ID, _ = os.LookupEnv("PICO_APP_ID")
	if service.PICO_APP_ID != "" {
		log.Println("PICO_APP_ID:", service.PICO_APP_ID)
	} else {
		log.Println("Ошибка!!! Не установлен APP_ID приложения для PICO CLI")
		os.Exit(1)
	}

	service.PICO_APP_SECRET, _ = os.LookupEnv("PICO_APP_SECRET")
	if service.PICO_APP_SECRET != "" {
		log.Println("PICO_APP_SECRET:", "*****************")
	} else {
		log.Println("Ошибка!!! Не установлен SECRED_ID приложения для PICO CLI")
		os.Exit(1)
	}

	service.OCULUS_APP_ID, _ = os.LookupEnv("OCULUS_APP_ID")
	if service.OCULUS_APP_ID != "" {
		log.Println("OCULUS_APP_ID:", service.OCULUS_APP_ID)
	} else {
		log.Println("Ошибка!!! Не установлен APP_ID приложения для OCULUS CLI")
		os.Exit(1)
	}

	service.OCULUS_APP_SECRET, _ = os.LookupEnv("OCULUS_APP_SECRET")
	if service.OCULUS_APP_SECRET != "" {
		log.Println("OCULUS_APP_SECRET:", "*******************")
	} else {
		log.Println("Ошибка!!! Не установлен SECRED_ID приложения для OCULUS CLI")
		os.Exit(1)
	}

	service.NAME_KEYSTORE, _ = os.LookupEnv("KEYSTORE_NAME")
	if service.NAME_KEYSTORE != "" {
		log.Println("KEYSTORE_NAME:", service.NAME_KEYSTORE)
	} else {
		log.Println("Ошибка!!! Не установлен путь с имененм хранилища ключей")
		os.Exit(1)
	}

	service.PATH_TO_LOGS, _ = os.LookupEnv("PATH_TO_LOGS")
	if service.PATH_TO_LOGS != "" {
		log.Println("PATH_TO_LOGS:", service.PATH_TO_LOGS)
	} else {
		log.Println("Ошибка!!! Не установлен путь к файлу логов")
		os.Exit(1)
	}

	service.DEST_ANDROID_BUILD_FOLDER, _ = os.LookupEnv("PATH_TO_ANDROID_BUILD_FOLDER")
	if service.DEST_ANDROID_BUILD_FOLDER != "" {
		service.DEST_ANDROID_BUILD_FOLDER = path.Join(service.DEST_ANDROID_BUILD_FOLDER + "\\")
		log.Println("PATH_TO_ANDROID_BUILD_FOLDER:", service.DEST_ANDROID_BUILD_FOLDER)

	} else {
		log.Println("Ошибка!!! Не установлен путь к каталогу для сборок под Android")
		os.Exit(1)
	}

	service.PATH_TO_EDITOR, _ = os.LookupEnv("PATH_TO_EDITOR")
	if service.PATH_TO_EDITOR != "" {
		log.Println("PATH_TO_EDITOR:", service.PATH_TO_EDITOR)
	} else {
		log.Println("Ошибка!!! Не установлен путь к исполняемому файлу Unity")
		os.Exit(1)
	}

	service.PATH_BUILDER_MOD, _ = os.LookupEnv("PATH_BUILDER_MOD")
	if service.PATH_BUILDER_MOD != "" {
		log.Println("PATH_BUILDER_MOD:", service.PATH_BUILDER_MOD)
	} else {
		log.Println("Ошибка!!! Не установлен путь к исполняемому файлу модуля для работы со сборкой")
		os.Exit(1)
	}

	service.PATH_CLOUD_MOD, _ = os.LookupEnv("PATH_CLOUD_MOD")
	if service.PATH_CLOUD_MOD != "" {
		log.Println("PATH_BUILDER_MOD:", service.PATH_CLOUD_MOD)
	} else {
		log.Println("Ошибка!!! Не установлен путь к исполняемому файлу модуля отправки файлов в облако")
		os.Exit(1)
	}

	service.PATH_TO_CONFIG_JSON, _ = os.LookupEnv("PATH_TO_CONFIG_JSON")
	if service.PATH_TO_CONFIG_JSON != "" {
		log.Println("PATH_TO_CONFIG_JSON:", service.PATH_TO_CONFIG_JSON)
	} else {
		log.Println("Ошибка!!! Не установлен путь к JSON конфигу")
		os.Exit(1)
	}

}

func main() {

	server := &http.Server{
		Addr:    ":8080",
		Handler: appRouter(),
	}

	log.Printf("<<SERVER START>>\n http://localhost%s", server.Addr)
	log.Fatal(server.ListenAndServe())
}
