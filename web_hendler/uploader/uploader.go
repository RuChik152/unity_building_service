package uploader

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	"web_hendler/bot"
	"web_hendler/loger"
)

type UploaderList struct {
	APK string
	OBB string
}

func GetllistFile(query string, dirPath string, list *UploaderList) {
	listAll, _ := getListDir(query, dirPath)

	filterdAPKList := fileFilter("apk", listAll)
	list.APK, _ = checkOldFile(filterdAPKList, dirPath)

	filterdOBBList := fileFilter("obb", listAll)
	list.OBB, _ = checkOldFile(filterdOBBList, dirPath)

}

func getListDir(query string, dirPath string) ([]fs.DirEntry, error) {

	var fileList []fs.DirEntry

	dir, err := os.Open(dirPath)
	if err != nil {
		loger.LogPrint.Package("UPLOADER").Log(fmt.Sprint("Ошибка доступа к директории: ", err))
		return nil, err
	}
	defer dir.Close()

	files, err := dir.ReadDir(-1)
	if err != nil {
		loger.LogPrint.Package("UPLOADER").Log(fmt.Sprint("Ошибка чтения директории: ", err))
		return nil, err
	}

	for _, file := range files {
		if !file.IsDir() {
			if strings.Contains(file.Name(), query) && (strings.HasSuffix(file.Name(), "apk") || strings.HasSuffix(file.Name(), "obb")) {
				fileList = append(fileList, file)
			}
		}
	}

	return fileList, nil
}

func fileFilter(sufux string, list []fs.DirEntry) []fs.DirEntry {
	var listFile []fs.DirEntry

	for _, file := range list {
		if strings.HasSuffix(file.Name(), sufux) {
			listFile = append(listFile, file)
		}
	}

	return listFile

}

func checkOldFile(list []fs.DirEntry, desDirBuild string) (string, error) {

	var myfile fs.DirEntry
	var myTime time.Time

	if list == nil {
		return "", errors.New("путой список")
	}

	for _, file := range list {
		infoFile, err := file.Info()
		if err != nil {
			loger.LogPrint.Package("UPLOADER").Log(fmt.Sprint("ошибка чтения файла ", err))
			continue
		}
		modTime := infoFile.ModTime()
		if myfile == nil || modTime.After(myTime) {
			myfile = file
			myTime = modTime
		}
	}

	return filepath.Join(desDirBuild, myfile.Name()), nil
}

func UploderBuild(msg *bot.BuildResultMessage, device string, apk string, obb string, app_id string, app_secret string, chanel string) {

	pathMudule, _ := os.LookupEnv("PATH_UPLOADER_MOD")
	if pathMudule == "" {
		loger.LogPrint.Package("UPLOADER").Log("Не установлен путь к исполняемому файлу модуля для работы с GIT")
		return
	}

	runArgs := []string{
		device,
		app_id,
		app_secret,
		apk,
		obb,
		chanel,
	}

	cmd := exec.Command(pathMudule, runArgs...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		loger.LogPrint.Package("UPLOADER").Log(fmt.Sprint("Не успешная загрузка сборки для: ", device))
		switch device {
		case "PICO":
			msg.Device.SendInfo = device + " отправка: ⚠️ Не успешно. " + string(output)
		case "OCULUS":
			msg.Device.SendInfo = device + " отправка: ⚠️ Не успешно. " + string(output)
		}
		loger.LogPrint.Package("UPLOADER").Log(fmt.Sprint("Ошибка загрузки: ", string(output), "\n", err))
		return
	} else {
		log.Println("Успешная загрузка сборки для ", device)
		switch device {
		case "PICO":
			msg.Device.SendInfo = device + " отправка: ✅ Успешно."
		case "OCULUS":
			msg.Device.SendInfo = device + " отправка: ✅ Успешно."
		}
		loger.LogPrint.Package("UPLOADER").Log(fmt.Sprint(string(output)))
		return
	}
}

func UploaderDesktopBuild(msg *bot.BuildResultMessage, device string, app_id string, app_secret string) {
	pathMudule, _ := os.LookupEnv("PATH_UPLOADER_MOD")
	if pathMudule == "" {
		loger.LogPrint.Package("UPLOADER").Log("Не установлен путь к исполняемому файлу модуля для работы с GIT")
		return
	}

	runArgs := []string{
		device,
		app_id,
		app_secret,
	}

	cmd := exec.Command(pathMudule, runArgs...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		loger.LogPrint.Package("UPLOADER").Log(fmt.Sprint("Не успешная загрузка сборки для: ", device))

		msg.Device.SendInfo = device + " отправка: ⚠️ Не успешно. " + string(output)
		loger.LogPrint.Package("UPLOADER").Log(fmt.Sprint("Ошибка загрузки: ", string(output), "\n", err))
		return
	} else {
		log.Println("Успешная загрузка сборки для ", device)

		msg.Device.SendInfo = device + " отправка: ✅ Успешно."
		loger.LogPrint.Package("UPLOADER").Log(fmt.Sprint(string(output)))
		return
	}
}
