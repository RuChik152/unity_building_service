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
	"web_hendler/loger"
)

var AGE_FILE string

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
			loger.LogPrint.Package("CLEANER").Log(fmt.Sprintf("Путь %s не существует\n", path))
		} else {
			loger.LogPrint.Package("CLEANER").Log(fmt.Sprintf("Ошибка при получении информации о пути %s: %s\n", path, err))
		}
		return
	}

	modeTime := fileInfo.ModTime()
	age := time.Since(modeTime)

	if age > (time.Duration(day) * 24 * time.Hour) {
		if !fileInfo.IsDir() {
			if strings.HasSuffix(fileInfo.Name(), "apk") || strings.HasSuffix(fileInfo.Name(), "obb") {
				if status := cloud(fileInfo.Name(), path, platform); status == 0 {
					loger.LogPrint.Package("CLEANER").Log(fmt.Sprintf("Удаляю папку или файл %s, старше %d дней \n", path, day))
					os.Remove(path)
				}
			}
		}
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
	loger.LogPrint.Package("CLEANER").Log(fmt.Sprint("Получены аргументы для отправки в облако: ", cmdArg))

	cmd := exec.Command(path_to_cloud_mod, cmdArg...)
	if err := cmd.Run(); err != nil {
		loger.LogPrint.Package("CLEANER").Log(fmt.Sprint("Ошибка при запуске модуля загрузки данных в облако", err))
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
			loger.LogPrint.Package("CLEANER").Log(fmt.Sprintf("Ошибка при поиске %v\n", err))
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
		loger.LogPrint.Package("CLEANER").Log(fmt.Sprintf("Ошибка при попытке обхода папок %v\n", err))
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

func RemoveAllFileDir(dir string) error {
	loger.LogPrint.Package("CLEANER").Log(fmt.Sprintf("Произвожу удаление файлов из дирректрории <<%s>> под сборку для STEAM\n", dir))
	if files, err := filepath.Glob(filepath.Join(dir, "*")); err != nil {
		return err
	} else {
		for _, file := range files {
			if err := os.Remove(file); err != nil {
				return err
			}
		}
	}
	return nil
}
