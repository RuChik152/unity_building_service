package uploader

import (
	"errors"
	"io/fs"
	"log"
	"os"
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
	list.APK, _ = checkOldFile(filterdAPKList)

	filterdOBBList := fileFilter("obb", listAll)
	list.OBB, _ = checkOldFile(filterdOBBList)

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

func checkOldFile(list []fs.DirEntry) (string, error) {

	var myfile fs.DirEntry
	var myTime time.Time

	if list == nil {
		return "", errors.New("путой список")
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

	return myfile.Name(), nil
}
