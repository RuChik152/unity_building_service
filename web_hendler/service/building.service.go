package service

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
	"web_hendler/bot"
	"web_hendler/proc"
	"web_hendler/uploader"
)

func Manager() {

	CHECK_LIST.git = 1
	CHECK_LIST.pre_build = 1
	CHECK_LIST.building = 1
	CHECK_LIST.upload = 1

	if !STATUS_BUILDING {

		bot.StandartMsg.Event = "allow"
		bot.StandartMsg.Message = "Запускаю сборку..."

		strStartMsgBot, err := json.Marshal(bot.StandartMsg)
		if err != nil {
			log.Println("Ошибка преобразования данных для запроса к боту")
		} else {
			bot.SendMessageBot(string(strStartMsgBot), "#pipline_build_start")
		}

		bot.ResultMsgBuild.Event = "build"
		bot.ResultMsgBuild.Info.DataVersion = "-- NO_DATA"
		bot.ResultMsgBuild.Info.OculusLogs = "-- NO_DATA"
		bot.ResultMsgBuild.Info.PicoLogs = "-- NO_DATA"

		gitSubManager()
		if CHECK_LIST.git == 0 {
			runCopyKey()
			runCreateGlobalConstant()
			if CHECK_LIST.pre_build == 0 {

				for platform, targetPlatform := range LIST_PLATFORM {
					if STATUS_RESET {
						break
					}
					switch platform {
					case "Android":
						for _, device := range targetPlatform {
							if STATUS_RESET {
								break
							}
							runCopyGeneralSettings(device)
							runBuild(platform, device)

							if STATUS_RESET {
								break
							}

							var pathListFile uploader.UploaderList

							uploader.GetllistFile(device, DEST_ANDROID_BUILD_FOLDER, &pathListFile)

							log.Println("FINALY: ", pathListFile)

							var app_id string
							var app_secret string

							switch device {
							case "PICO":
								app_id = PICO_APP_ID
								app_secret = PICO_APP_SECRET
								break
							case "OCULUS":
								app_id = OCULUS_APP_ID
								app_secret = OCULUS_APP_SECRET
								break
							}

							if pathListFile.APK != "" && pathListFile.OBB != "" {
								uploader.UploderBuild(device, pathListFile.APK, pathListFile.OBB, app_id, app_secret, "ALPHA")
							} else {
								log.Println("Не найден APK или OBB")
								log.Println("Получен путь к APK: ", pathListFile.APK)
								log.Println("Получен путь к OBB: ", pathListFile.OBB)
							}

						}
					}
				}

				if !STATUS_RESET {
					strMessage, err := json.Marshal(bot.ResultMsgBuild)
					if err != nil {
						log.Println("Ошибка преобразования данных для отправки боту")
					}

					bot.SendMessageBot(string(strMessage), "#pipline_check")
				}

			} else {
				log.Println("ERROR PREBUILD PROCCES")
				return
			}
		} else {
			log.Println("ERROR GIT PROCCESS")
			return
		}
	} else {

		bot.StandartMsg.Event = "allow"
		bot.StandartMsg.Message = "Сборка уже ведеться, выполняю перезапуск..."

		strStartMsgBot, err := json.Marshal(bot.StandartMsg)
		if err != nil {
			log.Println("Ошибка преобразования данных для запроса к боту")
		} else {
			bot.SendMessageBot(string(strStartMsgBot), "#pipline_build_restart")
		}

		STATUS_BUILDING = false
		STATUS_RESET = true
		proc.DestroyedBuilding(PROCCES_BUILDING.Process.Pid)

		time.Sleep(10 * time.Second)

		STATUS_RESET = false
		log.Println("PID>>>", PID_PROCCES_BUILDING)
		Manager()
	}
}

func gitSubManager() {

	pathMudule, _ := os.LookupEnv("PATH_GIT_MOD")
	if pathMudule == "" {
		log.Println("Не установлен путь к исполняемому файлу модуля для работы с GIT")
		return
	}

	gitArgsFetch := []string{
		"fetch",
		"origin",
	}
	runGitFetch := exec.Command(pathMudule, gitArgsFetch...)
	gitFetchOutput, err := runGitFetch.CombinedOutput()
	if err != nil {
		log.Println("Error FETCH: ", string(gitFetchOutput), "\n", err)
	} else {
		log.Println(string(gitFetchOutput), "Status code: ", runGitFetch.ProcessState.ExitCode())
	}
	CHECK_LIST.git = runGitFetch.ProcessState.ExitCode()

	gitArgsReset := []string{
		"reset",
		"origin",
	}
	runGitReset := exec.Command(pathMudule, gitArgsReset...)
	gitOutputReset, err := runGitReset.CombinedOutput()
	if err != nil {
		log.Println("Error RESET: ", string(gitOutputReset), "\n", err)
	} else {
		log.Println(string(gitOutputReset), "Status code: ", runGitFetch.ProcessState.ExitCode())
	}
	CHECK_LIST.git = runGitFetch.ProcessState.ExitCode()

	gitArgsPull := []string{
		"pull",
		"origin",
	}
	runGitPull := exec.Command(pathMudule, gitArgsPull...)
	gitOutputPull, err := runGitPull.CombinedOutput()
	if err != nil {
		log.Println("Error PULL: ", string(gitOutputPull), "\n", err)
	} else {
		log.Println(string(gitOutputPull), "Status code: ", runGitFetch.ProcessState.ExitCode())
	}
	CHECK_LIST.git = runGitFetch.ProcessState.ExitCode()
}

func runCopyKey() {

	pathMudule, _ := os.LookupEnv("PATH_PREBUILD_MOD")
	if pathMudule == "" {
		log.Println("Не установлен путь к исполняемому файлу модуля для работы с PREBUILD")
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
		log.Println("ERROR COPY KEYSTORE FILE: ", string(copyOutputKey), "\n", err)
	} else {
		log.Println(string(copyOutputKey), "Status code: ", runCopyKeyStore.ProcessState.ExitCode())
	}
	CHECK_LIST.pre_build = runCopyKeyStore.ProcessState.ExitCode()
}

func runCopyGeneralSettings(device string) {

	pathMudule, _ := os.LookupEnv("PATH_PREBUILD_MOD")
	if pathMudule == "" {
		log.Println("Не установлен путь к исполняемому файлу модуля для работы с PREBUILD")
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
		log.Println("ERROR COPY GENERAL SETTINGS: ", string(copyOutputSettings), "\n", err)
	} else {
		log.Println(string(copyOutputSettings), "Status code: ", runCopyGeneralSettings.ProcessState.ExitCode())
	}
	CHECK_LIST.pre_build = runCopyGeneralSettings.ProcessState.ExitCode()
}

func runCreateGlobalConstant() {

	pathMudule, _ := os.LookupEnv("PATH_PREBUILD_MOD")
	if pathMudule == "" {
		log.Println("Не установлен путь к исполняемому файлу модуля для работы с GIT")
		return
	}

	creatGlobalConstantArgs := []string{
		"create",
		PROJECT_FOLDER,
		filepath.Join(PROJECT_FOLDER, "\\Assets\\Scripts\\GlobalConstants.cs"),
	}
	runCreateGlobalConstant := exec.Command(pathMudule, creatGlobalConstantArgs...)
	creatGlobalConstantOutput, err := runCreateGlobalConstant.CombinedOutput()
	if err != nil {
		log.Println("ERROR CREAT CONSTANT: ", string(creatGlobalConstantOutput), "\n", err)
	} else {
		log.Println(string(creatGlobalConstantOutput), "Status code: ", runCreateGlobalConstant.ProcessState.ExitCode())
	}
	CHECK_LIST.pre_build = runCreateGlobalConstant.ProcessState.ExitCode()
}

func runBuild(platform string, device string) {

	path_to_logs := PATH_TO_LOGS + device + ".log"

	STATUS_BUILDING = true
	createArgs := []string{
		platform,
		"Karga_VR",
		DEST_ANDROID_BUILD_FOLDER,
		path_to_logs,
		"G:\\project\\BeliVR\\web-hook-server\\config.json",
		device,
	}

	log.Println("Получены аргументы запуска: ", createArgs)

	PROCCES_BUILDING = exec.Command(PATH_BUILDER_MOD, createArgs...)

	err := PROCCES_BUILDING.Start()

	if err != nil {
		log.Println("ОШИБКА ЗАПУСКА СБОРКИ. ", "ERR: ", err)
		STATUS_BUILDING = false
		CHECK_LIST.building = PROCCES_BUILDING.ProcessState.ExitCode()
		return
	}
	PID_PROCCES_BUILDING = PROCCES_BUILDING.Process.Pid
	log.Printf("PID запущенного процесса: %d", PROCCES_BUILDING.Process.Pid)

	err = PROCCES_BUILDING.Wait()
	if err != nil {
		runBuilderOutput, err := PROCCES_BUILDING.Output()
		if err != nil {
			switch device {
			case "PICO":
				bot.ResultMsgBuild.PicoMessage.Status = false
				bot.ResultMsgBuild.PicoMessage.Message = "⚠️ Не успешно"
				break
			case "OCULUS":
				bot.ResultMsgBuild.OculusMessage.Status = false
				bot.ResultMsgBuild.OculusMessage.Message = "⚠️ Не успешно"
				break
			}
			STATUS_BUILDING = false
			log.Println("ОШИБКА СБОРКИ: ", string(runBuilderOutput), "ERR: ", err)
			return
		} else {
			switch device {
			case "PICO":
				bot.ResultMsgBuild.PicoMessage.Status = false
				bot.ResultMsgBuild.PicoMessage.Message = "⚠️ Не успешно"
				break
			case "OCULUS":
				bot.ResultMsgBuild.OculusMessage.Status = false
				bot.ResultMsgBuild.OculusMessage.Message = "⚠️ Не успешно"
				break
			}
			log.Println(string(runBuilderOutput), "Status code: ", PROCCES_BUILDING.ProcessState.ExitCode())
			STATUS_BUILDING = false
			CHECK_LIST.building = PROCCES_BUILDING.ProcessState.ExitCode()
			return
		}

	} else {
		switch device {
		case "PICO":
			bot.ResultMsgBuild.PicoMessage.Status = true
			bot.ResultMsgBuild.PicoMessage.Message = "✅ Успешно"
			break
		case "OCULUS":
			bot.ResultMsgBuild.OculusMessage.Status = true
			bot.ResultMsgBuild.OculusMessage.Message = "✅ Успешно"
			break
		}
	}
}

func getOS() string {
	return runtime.GOOS
}
