package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/joho/godotenv"
)

func init() {

	actionArgs := os.Args[1]

	switch actionArgs {
	case "--help":
		log.Println("Example arguments: <Path\\to\\Editor> <Path\\to\\project> <Android> <Karga> <Path\\to\\Build\\> <Path\\to\\log\\log.log> <Path\\to\\config\\json> <keaystore.keystore> <PICO>")
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

		pathToProject := os.Args[2]
		// Название платформы для сборки Android, Win64 и т.д.
		targetPlatform := os.Args[3]
		name := os.Args[4]
		pathDestBuild := os.Args[5]
		logsFile := os.Args[6]
		pathConfigJson := os.Args[7]
		keystoreName := os.Args[8]
		targetDevice := os.Args[9]

		Arguments = append(Arguments,
			"-projectPath", pathToProject,
			"-quit",
			"-batchmode",
			"-nographics",
			"-buildTarget", targetPlatform,
			"-customBuildTarget", targetPlatform,
			"-customBuildName", name,
			"-customBuildPath", pathDestBuild,
			"-executeMethod", "BuildScript.BuildAPK",
			"-logFile", logsFile,
			"-pathJsonConfig", pathConfigJson,
			"-nameKeyStore", keystoreName,
			"-targetBuildDevice", targetDevice,
		)
	}

}

func main() {
	// Example arguments: <Path\to\Editor> <Path\to\project> <Android> <Karga> <Path\to\Build\> <Path\to\log\log.log> <Path\to\config\json> <keaystore.keystore> <PICO>

	pathToEditor := os.Args[1]

	log.Println(pathToEditor)
	// for _, value := range Arguments {
	// 	log.Println(value)
	// }

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
