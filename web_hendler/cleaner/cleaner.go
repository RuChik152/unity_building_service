package cleaner

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

/*
Сканирует указаный каталог на указанную глубину, удаляет из него файлы коорые старше указанного времени
*/
func ScanOldFile(dir string, day int, depth int, platform string) {
	log.Printf("[WEB_HENDLER][CLEANER] Провожу сканирование и удаляю файлы которые старше %d дней.\n", day)
	list, err := GetListFile(dir, depth)
	if err != nil {
		log.Println(err)
	}

	for _, value := range list {
		DeleteOldFile(value, day, platform)
	}

}

/*
Определение возраста файла, если файл старше указанного возраста то файл удалеятьеся
Нужно указать путь и возвраст.
*/
func DeleteOldFile(path string, day int, platform string) {
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
		//fmt.Printf("[WEB_HENDLER][CLEANER] Удаляю папку или файл %s, старше %d дней \n", path, day)
		if !fileInfo.IsDir() {
			fmt.Printf("[WEB_HENDLER][CLEANER] Удаляю папку или файл %s, старше %d дней \n", path, day)
			if strings.HasSuffix(fileInfo.Name(), "apk") || strings.HasSuffix(fileInfo.Name(), "obb") {
				if status := cloud(fileInfo.Name(), path, platform); status == 0 {
					fmt.Printf("[WEB_HENDLER][CLEANER] Удаляю папку или файл %s, старше %d дней \n", path, day)
					os.Remove(path)
				}
			}
		}
	} else {
		fmt.Printf("[WEB_HENDLER][CLEANER] Папка или файл %s, младше %d дней \n", path, day)

	}
}

func cloud(name string, path string, platform string) int {
	path_to_cloud_mod, _ := os.LookupEnv("PATH_CLOUD_MOD")
	if path_to_cloud_mod != "" {
		log.Println("PATH_CLOUD_MOD:", path_to_cloud_mod)
	} else {
		panic("Ошибка!!! Не установлен путь к исполняемому файлу модуля отправки файлов в облако")
	}

	cmdArg := []string{
		"-name", name,
		"-path", path,
		"-platform", platform,
	}

	log.Println("Получены аргументы для отправки в облако: ", cmdArg)
	cmd := exec.Command(path_to_cloud_mod, cmdArg...)
	if err := cmd.Run(); err != nil {
		log.Println("Ошибка при запуске модуля загрузки данных в облако", err)
	}

	return cmd.ProcessState.ExitCode()

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
