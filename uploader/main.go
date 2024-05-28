package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	log.Println("MOD_UPLOADER_INIT")
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
		os.Exit(1)
	} else {
		log.Print("Success, .env file found")
	}

	//EXAMPLE PICO <app-id> <app-secret> <path-apk> <path-obb> <chanel>
	if len(os.Args) < 2 {
		log.Println("Arguments not Found")
		os.Exit(1)
	} else {
		PLATFORM = os.Args[1]

		switch PLATFORM {
		case "OCULUS":
			CLI_PATH, _ = os.LookupEnv("PATH_OCULUS_TOOLS")
			if CLI_PATH == "" {
				log.Println("Путь к CLI не найден. Процесс завершен.")
				os.Exit(1)
			}
		case "PICO":
			CLI_PATH, _ = os.LookupEnv("PATH_PICO_TOOLS")
			if CLI_PATH == "" {
				log.Println("Путь к CLI не найден. Процесс завершен.")
				os.Exit(1)
			}
		case "PC":
			CLI_PATH, _ = os.LookupEnv("PATH_STEAM_TOOLS")
			if CLI_PATH == "" {
				log.Println("Путь к CLI не найден. Процесс завершен.")
				os.Exit(1)
			}
		default:
			log.Println("Arguments not Found")
			os.Exit(1)
		}

		if PLATFORM == "PICO" || PLATFORM == "OCULUS" {
			// updateCLI := exec.Command(CLI_PATH, "self-update")
			// version, _ := exec.Command(CLI_PATH, "version").Output()
			// err := updateCLI.Run()
			// if err != nil {
			// 	log.Println("Ошибка проверки обновлений: ", err)
			// }
			// log.Println(string(version))

			if len(os.Args) < 7 {
				log.Println("ERROR: APP_ID or APP_SECRET not found")
				os.Exit(1)
			} else {
				APP_ID = os.Args[2]
				APP_SECRET = os.Args[3]
				PATH_APK = os.Args[4]
				PATH_OBB = os.Args[5]
				CHANEL = os.Args[6]
			}
		}

		if PLATFORM == "PC" {
			APP_ID = os.Args[2]
			APP_SECRET = os.Args[3]
		}

	}

}

func main() {
	var statusCode int
	switch PLATFORM {
	case "OCULUS":
		statusCode = uploaderOCULUS()
	case "PICO":
		statusCode = uploaderPICO()
	case "PC":
		statusCode = uploaderSteam()
	}

	os.Exit(statusCode)
}
