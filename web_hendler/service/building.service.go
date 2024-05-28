package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
	"web_hendler/bot"
	"web_hendler/cleaner"
	"web_hendler/db"
	"web_hendler/loger"
	"web_hendler/proc"
	"web_hendler/uploader"
)

func Manager() {

	CHECK_LIST.git = 1
	CHECK_LIST.pre_build = 1
	CHECK_LIST.building = 1
	CHECK_LIST.upload = 1

	if !STATUS_BUILDING {

		go bot.StartBuildMsgBot()
		go handleBuildProcess()

	} else {
		go handelRestartBuild()
	}
}

func handleBuildProcess() {
	bot.ResultBuildMessage.Event = "build"
	gitSubManager()
	if CHECK_LIST.git == 0 {

		countVersion, _ := GetCountCurrentVersion()
		bot.ResultBuildMessage.Info.Log = fmt.Sprintf("-- /version_%d", countVersion)

		db.Commit.ID = countVersion
		go db.InsertOneDbCommit(db.Commit, "commits")
		runCopyKey()
		if CHECK_LIST.pre_build == 0 {

			for platform, targetPlatform := range LIST_PLATFORM {
				if STATUS_RESET {
					break
				}
				switch platform {
				case "Android":
					if age, err := strconv.Atoi(cleaner.AGE_FILE); err != nil {
						loger.LogPrint.Package("SERVICE").Log("Ошибка преобразования значения возраста файлов")
					} else {
						go cleaner.ScanOldFile(DEST_ANDROID_BUILD_FOLDER, age, 1, platform)
					}

					for _, device := range targetPlatform {
						if STATUS_RESET {
							break
						}
						runCreateGlobalConstant(device)
						runCopyGeneralSettings(device)
						runBuild(platform, device)

						if STATUS_RESET {
							break
						}

						go func(dev string, msg bot.BuildResultMessage) {

							done := make(chan bool)
							messageBot := msg
							if CHECK_LIST.building == 0 {
								go handelUploadBuild(dev, done, &messageBot)
							} else {
								done <- true
							}
							go handelBotMessage(done, &messageBot)

						}(device, bot.ResultBuildMessage)

					}
				case "Win64":
					for _, device := range targetPlatform {
						if STATUS_RESET {
							break
						}

						runCreateGlobalConstant(device)
						cleaner.RemoveAllFileDir(DEST_WIN_BUILD_FOLDER)
						runPCBuild(platform, device)

						if STATUS_RESET {
							break
						}

						go func(dev string, msg bot.BuildResultMessage) {

							done := make(chan bool)
							messageBot := msg
							if CHECK_LIST.building == 0 {
								go handelUploadBuild(dev, done, &messageBot)
							} else {
								done <- true
							}
							go handelBotMessage(done, &messageBot)

						}(device, bot.ResultBuildMessage)
					}
				}
				STATUS_BUILDING = false
			}

		} else {
			loger.LogPrint.Package("SERVICE").Log("ERROR PREBUILD PROCCES")
			return
		}
	} else {
		loger.LogPrint.Package("SERVICE").Log("ERROR GIT PROCCESS")
		return
	}
}

func handelBotMessage(done chan bool, msg *bot.BuildResultMessage) {

	defer close(done)
	if <-done {
		loger.LogPrint.Package("SERVICE").Log(fmt.Sprint("Отправил сообщение боту: ", msg))
		if data, err := json.Marshal(msg); err != nil {
			log.Panic("ошибка преобразования handelBotMessage")
		} else {
			go bot.SendMsgBot(&data)
		}

	}

}

func handelRestartBuild() {
	// bot.StandartMsg.Event = "allow"
	// bot.StandartMsg.Message = "Сборка уже ведеться, выполняю перезапуск..."

	// strStartMsgBot, err := json.Marshal(bot.StandartMsg)
	// if err != nil {
	// 	log.Println("Ошибка преобразования данных для запроса к боту")
	// } else {
	// 	bot.SendMessageBot(string(strStartMsgBot), "#pipline_build_restart")
	// }

	bot.RestartBuildMsg()

	STATUS_BUILDING = false
	STATUS_RESET = true
	proc.DestroyedBuilding(PROCCES_BUILDING.Process.Pid)

	time.Sleep(10 * time.Second)

	STATUS_RESET = false
	loger.LogPrint.Package("SERVICE").Log(fmt.Sprint("PID>>>", PID_PROCCES_BUILDING))
	Manager()
}

func gitSubManager() {

	pathMudule, _ := os.LookupEnv("PATH_GIT_MOD")
	if pathMudule == "" {
		loger.LogPrint.Package("SERVICE").Log("Не установлен путь к исполняемому файлу модуля для работы с GIT")
		return
	}

	gitArgsFetch := []string{
		"fetch",
		"origin",
	}
	runGitFetch := exec.Command(pathMudule, gitArgsFetch...)
	gitFetchOutput, err := runGitFetch.CombinedOutput()
	if err != nil {
		loger.LogPrint.Package("SERVICE").Log(fmt.Sprint("Error FETCH: ", string(gitFetchOutput), err))

	} else {
		loger.LogPrint.Package("SERVICE").Log(fmt.Sprint(string(gitFetchOutput), "Status code: ", runGitFetch.ProcessState.ExitCode()))
	}
	CHECK_LIST.git = runGitFetch.ProcessState.ExitCode()

	gitArgsReset := []string{
		"reset",
		"origin",
	}
	runGitReset := exec.Command(pathMudule, gitArgsReset...)
	gitOutputReset, err := runGitReset.CombinedOutput()
	if err != nil {
		loger.LogPrint.Package("SERVICE").Log(fmt.Sprint("Error RESET: ", string(gitOutputReset), err))
	} else {
		loger.LogPrint.Package("SERVICE").Log(fmt.Sprint(string(gitOutputReset), "Status code: ", runGitFetch.ProcessState.ExitCode()))
	}
	CHECK_LIST.git = runGitFetch.ProcessState.ExitCode()

	gitArgsPull := []string{
		"pull",
		"origin",
	}
	runGitPull := exec.Command(pathMudule, gitArgsPull...)
	gitOutputPull, err := runGitPull.CombinedOutput()
	if err != nil {
		loger.LogPrint.Package("SERVICE").Log(fmt.Sprint("Error PULL: ", string(gitOutputPull), "\n", err))
	} else {
		loger.LogPrint.Package("SERVICE").Log(fmt.Sprint(string(gitOutputPull), "Status code: ", runGitFetch.ProcessState.ExitCode()))
	}
	CHECK_LIST.git = runGitFetch.ProcessState.ExitCode()
}

func runCopyKey() {

	pathMudule, _ := os.LookupEnv("PATH_PREBUILD_MOD")
	if pathMudule == "" {
		loger.LogPrint.Package("SERVICE").Log("Не установлен путь к исполняемому файлу модуля для работы с PREBUILD")
		return
	}

	name_keystore, _ := os.LookupEnv("KEYSTORE_NAME")
	if name_keystore == "" {
		panic("Не установлено имя хранилища ключей <name.keystore>")
	}

	path_store_keystore, _ := os.LookupEnv("PATH_STORE_KEYSTORE")
	if path_store_keystore == "" {
		panic("Не установлен путь до хранилища ключей <path\\to\\storage\\folder>")
	}

	copyKeyStore := []string{
		"copy",
		filepath.Join(path_store_keystore, name_keystore),
		filepath.Join(PROJECT_FOLDER, name_keystore),
	}

	runCopyKeyStore := exec.Command(pathMudule, copyKeyStore...)
	copyOutputKey, err := runCopyKeyStore.CombinedOutput()
	if err != nil {
		loger.LogPrint.Package("SERVICE").Log(fmt.Sprint("ERROR COPY KEYSTORE FILE: ", string(copyOutputKey), "\n", err))
	} else {
		loger.LogPrint.Package("SERVICE").Log(fmt.Sprint(string(copyOutputKey), "Status code: ", runCopyKeyStore.ProcessState.ExitCode()))
	}
	CHECK_LIST.pre_build = runCopyKeyStore.ProcessState.ExitCode()
}

func runCopyGeneralSettings(device string) {

	pathMudule, _ := os.LookupEnv("PATH_PREBUILD_MOD")
	if pathMudule == "" {
		loger.LogPrint.Package("SERVICE").Log("Не установлен путь к исполняемому файлу модуля для работы с PREBUILD")
		return
	}

	path_to_assets, _ := os.LookupEnv("PATH_ASSETS_FOLDER")
	if path_to_assets == "" {
		panic("Не передан путь к файлам конфиграции платформ")
	}

	var pathToSettings string

	switch device {
	case "PICO":
		pathToSettings = filepath.Join(path_to_assets, "PicoXRGeneralSettings.asset")
	case "OCULUS":
		pathToSettings = filepath.Join(path_to_assets, "OculusXRGeneralSettings.asset")
	}

	copySettingsArgs := []string{
		"copy",
		pathToSettings,
		filepath.Join(PROJECT_FOLDER, "\\Assets\\XR\\XRGeneralSettings.asset"),
	}
	runCopyGeneralSettings := exec.Command(pathMudule, copySettingsArgs...)
	copyOutputSettings, err := runCopyGeneralSettings.CombinedOutput()
	if err != nil {
		loger.LogPrint.Package("SERVICE").Log(fmt.Sprint("ERROR COPY GENERAL SETTINGS: ", string(copyOutputSettings), "\n", err))
	} else {
		loger.LogPrint.Package("SERVICE").Log(fmt.Sprint(string(copyOutputSettings), "Status code: ", runCopyGeneralSettings.ProcessState.ExitCode()))
	}
	CHECK_LIST.pre_build = runCopyGeneralSettings.ProcessState.ExitCode()
}

func runCreateGlobalConstant(device string) {

	pathMudule, _ := os.LookupEnv("PATH_PREBUILD_MOD")
	if pathMudule == "" {
		loger.LogPrint.Package("SERVICE").Log("Не установлен путь к исполняемому файлу модуля для работы с GIT")
		return
	}

	creatGlobalConstantArgs := []string{
		"create",
		PROJECT_FOLDER,
		filepath.Join(PROJECT_FOLDER, "\\Assets\\Scripts\\GlobalConstants.cs"),
		device,
	}
	runCreateGlobalConstant := exec.Command(pathMudule, creatGlobalConstantArgs...)
	creatGlobalConstantOutput, err := runCreateGlobalConstant.CombinedOutput()
	if err != nil {
		loger.LogPrint.Package("SERVICE").Log(fmt.Sprint("ERROR CREAT CONSTANT: ", string(creatGlobalConstantOutput), "\n", err))
	} else {
		loger.LogPrint.Package("SERVICE").Log(fmt.Sprint(string(creatGlobalConstantOutput), "Status code: ", runCreateGlobalConstant.ProcessState.ExitCode()))
	}
	CHECK_LIST.pre_build = runCreateGlobalConstant.ProcessState.ExitCode()
}

func GetCountCurrentVersion() (int, error) {
	pathMudule, _ := os.LookupEnv("PATH_GIT_MOD")
	if pathMudule == "" {
		return 0, errors.New("не установлен путь к исполняемому файлу модуля для работы с GIT")
	}

	args := []string{
		"count",
	}
	cmd := exec.Command(pathMudule, args...)
	output, err := cmd.Output()
	if err != nil {
		return 0, errors.New("ошибка полчения данных от модуля GIT")
	} else {

		if i, err := strconv.Atoi(string(output)); err != nil {
			return 0, errors.New("ошибка преобразования данных")
		} else {
			return i, nil
		}

	}
}

func handelUploadBuild(device string, done chan bool, msg *bot.BuildResultMessage) {

	var app_id string
	var app_secret string

	switch device {
	case "PICO":
		app_id = PICO_APP_ID
		app_secret = PICO_APP_SECRET

	case "OCULUS":
		app_id = OCULUS_APP_ID
		app_secret = OCULUS_APP_SECRET
	case "PC":
		app_id = STEAM_LOGIN
		app_secret = STEAM_PASS
	}

	if device == "PICO" || device == "OCULUS" {
		var pathListFile uploader.UploaderList

		uploader.GetllistFile(device, DEST_ANDROID_BUILD_FOLDER, &pathListFile)

		loger.LogPrint.Package("SERVICE").Log(fmt.Sprint("Получен список путей файлам для загрузки: ", pathListFile))

		if pathListFile.APK != "" && pathListFile.OBB != "" {
			uploader.UploderBuild(msg, device, pathListFile.APK, pathListFile.OBB, app_id, app_secret, "ALPHA")
		} else {
			loger.LogPrint.Package("SERVICE").Log("Не найден APK или OBB")
			loger.LogPrint.Package("SERVICE").Log(fmt.Sprint("Получен путь к APK: ", pathListFile.APK))
			loger.LogPrint.Package("SERVICE").Log(fmt.Sprint("Получен путь к OBB: ", pathListFile.OBB))

		}
		done <- true
	}

	if device == "PC" {
		// нужна логика для упаковки и загрзуки в облако.
		uploader.UploaderDesktopBuild(msg, device, app_id, app_secret)
		done <- true
	}

}
