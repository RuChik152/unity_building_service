package main

import (
	"log"
	"os"
	"os/exec"
)

func init() {
	//EXAMPLE PICO <app-id> <app-secret> <path-apk> <path-obb> <chanel>
	if len(os.Args) < 2 {
		log.Println("Arguments not Found")
		os.Exit(1)
	} else {
		PLATFORM = os.Args[1]

		switch PLATFORM {
		case "OCULUS":
			CLI_PATH = "..\\uploader\\cli\\ovr-platform-util.exe"
		case "PICO":
			CLI_PATH = "..\\uploader\\cli\\pico-cli.exe"
		default:
			log.Println("Arguments not Found")
			os.Exit(1)
		}

		updateCLI := exec.Command(CLI_PATH, "self-update")
		version, _ := exec.Command(CLI_PATH, "version").Output()
		err := updateCLI.Run()
		if err != nil {
			log.Println("Ошибка проверки обновлений: ", err)
		}
		log.Println(string(version))

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

}

func main() {
	var statusCode int
	switch PLATFORM {
	case "OCULUS":
		statusCode = uploaderOCULUS()
	case "PICO":
		statusCode = uploaderPICO()
	}
	os.Exit(statusCode)
}
