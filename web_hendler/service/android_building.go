package service

import (
	"fmt"
	"os/exec"
	"web_hendler/bot"
	"web_hendler/loger"
)

func runBuild(platform string, device string) {

	path_to_logs := PATH_TO_LOGS + device + ".log"

	STATUS_BUILDING = true
	createArgs := []string{
		platform,
		"Karga_VR",
		DEST_ANDROID_BUILD_FOLDER,
		path_to_logs,
		PATH_TO_CONFIG_JSON,
		device,
	}

	loger.LogPrint.Package("SERVICE").Log(fmt.Sprint("Получены аргументы запуска: ", createArgs))

	PROCCES_BUILDING = exec.Command(PATH_BUILDER_MOD, createArgs...)

	err := PROCCES_BUILDING.Start()

	if err != nil {
		loger.LogPrint.Package("SERVICE").Log(fmt.Sprint("ОШИБКА ЗАПУСКА СБОРКИ. ", "ERR: ", err))

		STATUS_BUILDING = false
		CHECK_LIST.building = PROCCES_BUILDING.ProcessState.ExitCode()
		return
	}
	PID_PROCCES_BUILDING = PROCCES_BUILDING.Process.Pid
	loger.LogPrint.Package("SERVICE").Log(fmt.Sprintf("PID запущенного процесса: %d", PROCCES_BUILDING.Process.Pid))
	err = PROCCES_BUILDING.Wait()
	if err != nil {
		runBuilderOutput, err := PROCCES_BUILDING.Output()
		if err != nil {
			switch device {
			case "PICO":
				bot.ResultBuildMessage.Device.BuildInfo = device + " сборка: ⚠️ Не успешно: " + fmt.Sprintf("%s", err)

			case "OCULUS":
				bot.ResultBuildMessage.Device.BuildInfo = device + " сборка: ⚠️ Не успешно: " + fmt.Sprintf("%s", err)

			}
			STATUS_BUILDING = false
			loger.LogPrint.Package("SERVICE").Log(fmt.Sprint("ОШИБКА СБОРКИ: ", string(runBuilderOutput), "ERR: ", err))
			return
		} else {
			switch device {
			case "PICO":
				bot.ResultBuildMessage.Device.BuildInfo = device + " сборка: ⚠️ Не успешно"

			case "OCULUS":
				bot.ResultBuildMessage.Device.BuildInfo = device + " сборка: ⚠️ Не успешно"
			}
			loger.LogPrint.Package("SERVICE").Log(fmt.Sprint(string(runBuilderOutput), "Status code: ", PROCCES_BUILDING.ProcessState.ExitCode()))
			STATUS_BUILDING = false
			CHECK_LIST.building = PROCCES_BUILDING.ProcessState.ExitCode()
			return
		}

	} else {
		switch device {
		case "PICO":
			bot.ResultBuildMessage.Device.BuildInfo = device + " сборка: ✅ Успешно"
			CHECK_LIST.building = PROCCES_BUILDING.ProcessState.ExitCode()
		case "OCULUS":
			bot.ResultBuildMessage.Device.BuildInfo = device + " сборка: ✅ Успешно"
			CHECK_LIST.building = PROCCES_BUILDING.ProcessState.ExitCode()
		}
	}
}
