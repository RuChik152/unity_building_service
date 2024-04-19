package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/joho/godotenv"
)

func init() {
	log.Println("MOD_BUILD_INIT")
	actionArgs := os.Args[1]

	switch actionArgs {
	case "--help":
		log.Println("Example arguments: <Android> <Name_APK> <Path\\to\\Build\\> <Path\\to\\log\\log.log> <Path\\to\\config\\json> <PICO>")
		os.Exit(0)
	default:
		if err := godotenv.Load(); err != nil {
			log.Print("No .env file found")
		} else {
			log.Print("Success, .env file found")
		}

		//Если перевенные среды не будут отпределяться C# скриптом то раскоментровать и проверить
		// KEYSTORE_PASS, _ := os.LookupEnv("KEYSTORE_PASS")
		// KEY_ALIAS_NAME, _ := os.LookupEnv("KEY_ALIAS_NAME")
		// KEY_ALIAS_PASS, _ := os.LookupEnv("KEY_ALIAS_PASS")

		os.Setenv("KEYSTORE_PASS", "qwertyui")
		os.Setenv("KEY_ALIAS_NAME", "karga")
		os.Setenv("KEY_ALIAS_PASS", "qwertyui")

		pathToProject, _ := os.LookupEnv("PATH_TO_FOLDER_FOR_PROJECT")
		if pathToProject == "" {
			log.Println("Не установлен путь в переменных среды к папке с проектом. Завершаю работу...")
			os.Exit(1)
		}
		// Название платформы для сборки Android, Win64 и т.д.
		// targetPlatform := os.Args[1] //Android
		// name := os.Args[2]           //Karga_VR
		// pathDestBuild := os.Args[3]  //C:\Unity\build\
		// logsFile := os.Args[4]       //C:\Unity\logs\test_runner_logs.log
		// pathConfigJson := os.Args[5] //G:\project\BeliVR\web-hook-server\config.json
		keystoreName, _ := os.LookupEnv("KEYSTORE_NAME")
		// targetDevice := os.Args[6] //karga.keystore
		//PICO

		log.Println("<<<ARGUMENTS>>", os.Args)

		Arguments = append(Arguments,
			"-projectPath", pathToProject,
			"-quit",
			"-batchmode",
			"-nographics",
			"-buildTarget", os.Args[1],
			"-customBuildTarget", os.Args[1],
			"-customBuildName", os.Args[2],
			"-customBuildPath", os.Args[3],
			"-executeMethod", "BuildScript.BuildAPK",
			"-logFile", os.Args[4],
			"-pathJsonConfig", os.Args[5],
			"-nameKeyStore", keystoreName,
			"-targetBuildDevice", os.Args[6],
		)
	}

}

func main() {
	// Example arguments: <Android> <Karga> <Path\to\Build\> <Path\to\log\log.log> <Path\to\config\json> <PICO>
	pathToEditor, _ := os.LookupEnv("PATH_UNITY_EDITOR")
	if pathToEditor == "" {
		log.Println("Не установлен путь в переменных средах к редактору Unity. Завершаю работу...")
		os.Exit(1)
	}

	cmd := exec.Command(pathToEditor, Arguments...)

	if err := cmd.Run(); err != nil {
		log.Println("::Ошибка запуска процесса Unity.exe", err)
	}

	pid := cmd.Process.Pid
	fmt.Println("PID, процесса: ", pid)

	exitCode := cmd.ProcessState.ExitCode()
	log.Printf("<<FINALY BUILDING>> Exit code: %d\n", exitCode)
	os.Exit(exitCode)
}
