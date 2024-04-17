package uploader

import (
	"errors"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
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
		log.Println("Ошибка доступа к директории: ", err)
		return nil, err
	}
	defer dir.Close()

	files, err := dir.ReadDir(-1)
	if err != nil {
		log.Println("Ошибка чтения директории: ", err)
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
		return "", errors.New("утой список")
	}

	for _, file := range list {
		infoFile, err := file.Info()
		if err != nil {
			log.Println("ошибка чтения файла ", err)
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

func UploderBuild(device string, apk string, obb string, app_id string, app_secret string, chanel string) {

	runArgs := []string{
		device,
		app_id,
		app_secret,
		apk,
		obb,
		chanel,
	}

	cmd := exec.Command("..\\uploader\\uploader.exe", runArgs...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("Ошибка загрузки: ", string(output), "\n", err)
		return
	} else {
		log.Println(string(output))
		return
	}
}
