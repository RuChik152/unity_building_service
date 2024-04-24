package cleaner

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

/*
Сканирует указаный каталог на указанную глубину, удаляет из него файлы коорые старше указанного времени
*/
func ScanOldFile(dir string, day int, depth int) {
	log.Printf("[WEB_HENDLER][CLEANER] Провожу сканирование и удаляю файлы которые старше %d дней.\n", day)
	list, err := GetListFile(dir, depth)
	if err != nil {
		log.Println(err)
	}

	for _, value := range list {
		DeleteOldFile(value, day)
	}

}

/*
Определение возраста файла, если файл старше указанного возраста то файл удалеятьеся
Нужно указать путь и возвраст.
*/
func DeleteOldFile(path string, day int) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("[WEB_HENDLER][CLEANER] Путь %s не существует\n", path)
		} else {
			fmt.Printf("[WEB_HENDLER][CLEANER] Ошибка при получении информации о пути %s: %s\n", path, err)
		}
		return
	}

	modeTime := fileInfo.ModTime()
	age := time.Since(modeTime)

	if age > (time.Duration(day) * 24 * time.Hour) {
		fmt.Printf("[WEB_HENDLER][CLEANER] Удаляю папку или файл %s, старше %d дней \n", path, day)
		if !fileInfo.IsDir() {
			os.Remove(path)
		}
	} else {
		fmt.Printf("[WEB_HENDLER][CLEANER] Папка или файл %s, младше %d дней \n", path, day)

	}
}

/*
Возвращает срез содержащий пути файлам\каталогам в указаной дирректории.
Принимает путь папки которую нужно просканирвоать.
Так же нужно указать глубину сканирвоания.
*/
func GetListFile(rootPath string, maxDepth int) ([]string, error) {
	var list []string

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("[WEB_HENDLER][CLEANER] Ошибка при поиске %v\n", err)
		}

		lenthPath := len(strings.Split(path, Slash()))
		lenthRoot := len(strings.Split(rootPath, Slash()))

		currentDepth := lenthPath - lenthRoot
		if currentDepth <= maxDepth-1 {
			if info.IsDir() {
				//fmt.Printf("Папка %s\n", path)
				list = append(list, path)
			} else {
				//fmt.Printf("Файл %s\n", path)
				list = append(list, path)
			}
			return nil
		}
		return filepath.SkipDir

	})
	if err != nil {
		fmt.Printf("[WEB_HENDLER][CLEANER] Ошибка при попытке обхода папок %v\n", err)
		return nil, errors.New("[WEB_HENDLER][CLEANER] Ошибка при попытке обхода папок")
	}
	return list, nil
}

func Slash() string {
	if runtime.GOOS == "windows" {
		return "\\"
	} else {
		return "//"
	}
}
