package service

import (
	"log"
	"os/exec"
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
				log.Println("TEST BRANCH")
				log.Println(CHECK_LIST)
			} else {
				log.Println("ERROR PREBUILD PROCCES")
				return
			}
		} else {
			log.Println("ERROR GIT PROCCESS")
			return
		}
	} else {
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
