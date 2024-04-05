package service

import (
	"log"
	"os/exec"
	"runtime"
	"web_hendler/proc"
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
			runCopyGeneralSettings()
			runCreateGlobalConstant()
			if CHECK_LIST.pre_build == 0 {

				var pid = 17916
				listPid := proc.GetListChildProcces(pid)
				//runBuild("Android", "PICO")
				log.Println(listPid)
			} else {
				log.Println("ERROR PREBUILD PROCCES")
				return
			}
		} else {
			log.Println("ERROR GIT PROCCESS")
			return
		}
	} else {
		log.Println("PID>>>", PID_PROCCES_BUILDING)
		//логика для перезапуска сборки
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

func runCopyGeneralSettings() {
	copySettingsArgs := []string{
		"copy",
		"G:\\project\\BeliVR\\web-hook-server\\assets\\PicoXRGeneralSettings.asset",
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

	cmdRunBuild := exec.Command("..\\builder\\builder.exe", createArgs...)

	err := cmdRunBuild.Start()

	// log.Println("PID: ", PID_PROCCES_BUILDING)
	if err != nil {
		log.Println("ОШИБКА ЗАПУСКА СБОРКИ. ", "ERR: ", err)
		CHECK_LIST.building = cmdRunBuild.ProcessState.ExitCode()
	}
	PID_PROCCES_BUILDING = cmdRunBuild.Process.Pid
	log.Printf("PID запущенного процесса: %d", cmdRunBuild.Process.Pid)

	err = cmdRunBuild.Wait()
	if err != nil {
		runBuilderOutput, err := cmdRunBuild.Output()
		if err != nil {
			log.Println("ОШИБКА СБОРКИ: ", string(runBuilderOutput), "ERR: ", err)
		} else {
			log.Println(string(runBuilderOutput), "Status code: ", cmdRunBuild.ProcessState.ExitCode())
			CHECK_LIST.building = cmdRunBuild.ProcessState.ExitCode()
		}
	}
}

func getOS() string {
	return runtime.GOOS
}

// func destroyed(pid *int) {
// 	sistem := getOS()
// 	switch sistem {
// 	case "windows":
// 		os.FindProcess(*pid)
// 		process.Processes()
// 	}
// }
