package service

import (
	"log"
	"os/exec"
	"runtime"
	"time"
	"web_hendler/proc"
	"web_hendler/uploader"
)

func Manager() {

	CHECK_LIST.git = 1
	CHECK_LIST.pre_build = 1
	CHECK_LIST.building = 1
	CHECK_LIST.upload = 1

	if !STATUS_BUILDING {
		gitSubManager()
		if CHECK_LIST.git == 0 {
			runCopyKey()
			runCreateGlobalConstant()
			if CHECK_LIST.pre_build == 0 {

				//runBuild("Android", "PICO")
				// listPid := proc.GetListChildProcces(PROCCES_BUILDING.Process.Pid)
				// log.Println(listPid)

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
			} else {
				log.Println("ERROR PREBUILD PROCCES")
				return
			}
		} else {
			log.Println("ERROR GIT PROCCESS")
			return
		}
	} else {
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
	gitArgsFetch := []string{
		"fetch",
		"origin",
		PROJECT_FOLDER,
	}
	runGitFetch := exec.Command("..\\git\\git_runner.exe", gitArgsFetch...)
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
		PROJECT_FOLDER,
		"main",
	}
	runGitReset := exec.Command("..\\git\\git_runner.exe", gitArgsReset...)
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
		PROJECT_FOLDER,
		"main",
	}
	runGitPull := exec.Command("..\\git\\git_runner.exe", gitArgsPull...)
	gitOutputPull, err := runGitPull.CombinedOutput()
	if err != nil {
		log.Println("Error PULL: ", string(gitOutputPull), "\n", err)
	} else {
		log.Println(string(gitOutputPull), "Status code: ", runGitFetch.ProcessState.ExitCode())
	}
	CHECK_LIST.git = runGitFetch.ProcessState.ExitCode()
}

func runCopyKey() {
	copyKeyStore := []string{
		"copy",
		"C:\\Unity\\karga.keystore",
		"C:\\Unity\\clone_3\\karga.keystore",
	}

	runCopyKeyStore := exec.Command("..\\pre_builder\\pre_builder.exe", copyKeyStore...)
	copyOutputKey, err := runCopyKeyStore.CombinedOutput()
	if err != nil {
		log.Println("ERROR COPY KEYSTORE FILE: ", string(copyOutputKey), "\n", err)
	} else {
		log.Println(string(copyOutputKey), "Status code: ", runCopyKeyStore.ProcessState.ExitCode())
	}
	CHECK_LIST.pre_build = runCopyKeyStore.ProcessState.ExitCode()
}

func runCopyGeneralSettings(device string) {

	var pathToSettings string

	switch device {
	case "PICO":
		pathToSettings = "G:\\project\\BeliVR\\web-hook-server\\assets\\PicoXRGeneralSettings.asset"
	case "OCULUS":
		pathToSettings = "G:\\project\\BeliVR\\web-hook-server\\assets\\OculusXRGeneralSettings.asset"
	}

	copySettingsArgs := []string{
		"copy",
		pathToSettings,
		"C:\\Unity\\clone_3\\Assets\\XR\\XRGeneralSettings.asset",
	}
	runCopyGeneralSettings := exec.Command("..\\pre_builder\\pre_builder.exe", copySettingsArgs...)
	copyOutputSettings, err := runCopyGeneralSettings.CombinedOutput()
	if err != nil {
		log.Println("ERROR COPY GENERAL SETTINGS: ", string(copyOutputSettings), "\n", err)
	} else {
		log.Println(string(copyOutputSettings), "Status code: ", runCopyGeneralSettings.ProcessState.ExitCode())
	}
	CHECK_LIST.pre_build = runCopyGeneralSettings.ProcessState.ExitCode()
}

func runCreateGlobalConstant() {

	creatGlobalConstantArgs := []string{
		"create",
		PROJECT_FOLDER,
		"main",
		"C:\\Unity\\clone_3\\Assets\\Scripts\\GlobalConstants.cs",
	}
	runCreateGlobalConstant := exec.Command("..\\pre_builder\\pre_builder.exe", creatGlobalConstantArgs...)
	creatGlobalConstantOutput, err := runCreateGlobalConstant.CombinedOutput()
	if err != nil {
		log.Println("ERROR CREAT CONSTANT: ", string(creatGlobalConstantOutput), "\n", err)
	} else {
		log.Println(string(creatGlobalConstantOutput), "Status code: ", runCreateGlobalConstant.ProcessState.ExitCode())
	}
	CHECK_LIST.pre_build = runCreateGlobalConstant.ProcessState.ExitCode()
}

func runBuild(platform string, device string) {

	STATUS_BUILDING = true
	createArgs := []string{
		PATH_TO_EDITOR,
		PROJECT_FOLDER,
		platform,
		"Karga_test_VR",
		DEST_ANDROID_BUILD_FOLDER,
		PATH_TO_LOGS,
		"G:\\project\\BeliVR\\web-hook-server\\config.json",
		"karga.keystore",
		device,
	}

	PROCCES_BUILDING = exec.Command("..\\builder\\builder.exe", createArgs...)

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
			STATUS_BUILDING = false
			log.Println("ОШИБКА СБОРКИ: ", string(runBuilderOutput), "ERR: ", err)
			return
		} else {
			log.Println(string(runBuilderOutput), "Status code: ", PROCCES_BUILDING.ProcessState.ExitCode())
			STATUS_BUILDING = false
			CHECK_LIST.building = PROCCES_BUILDING.ProcessState.ExitCode()
			return
		}
	}
}

func getOS() string {
	return runtime.GOOS
}
