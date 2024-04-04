package main

import (
	"log"
	"os/exec"
)

func uploaderOCULUS() int {
	cliArguments := []string{
		"upload-quest-build",
		"--age-group", "TEENS_AND_ADULTS",
		"--app_id", APP_ID,
		"--app_secret", APP_SECRET,
		"--apk", PATH_APK,
		"--obb", PATH_OBB,
		"--channel", CHANEL,
	}

	cmd := exec.Command(CLI_PATH, cliArguments...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("Ошибка загрузки: ", string(output), "\n", err)
		return 1
	} else {
		log.Println(string(output))
		return 0
	}

}

func uploaderPICO() int {
	cliArguments := []string{
		"upload-build",
		"--app-id", APP_ID,
		"--app-secret", APP_SECRET,
		"--region", "noncn",
		"--device", "PICO 4",
		"--apk", PATH_APK,
		"--obb", PATH_OBB,
		"--channel", "4",
	}

	cmd := exec.Command(CLI_PATH, cliArguments...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("Ошибка загрузки: ", string(output), "\n", err)
		return 1
	} else {
		log.Println(string(output))
		return 0
	}
}
