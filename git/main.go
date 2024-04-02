package main

import (
	"fmt"
	"log"
	"os"
)

func init() {

	OptionsData = make(map[interface{}]interface{})
	OptionsData["fetch"] = "fetch ## Считать обновления с репозитория"
	OptionsData["reset"] = "reset ## Сбросить состояние локального репозиоря до состояния удаленного"
	OptionsData["pull"] = "pull  ## Получить все полседние изменения"
}

func main() {

	actionArg := os.Args[1]

	switch actionArg {
	case "--help":
		fmt.Println("<< Возможные аргументы >>")
		for _, value := range OptionsData {
			fmt.Println(value)
		}
		os.Exit(0)
	case "fetch":
		if len(os.Args) < 4 {
			log.Println("Ошибка. Не переданы все необходимые аргументы. Нужно передать название удаленного репозиотрия и путь к каталогу с проектом")
			os.Exit(1)
		}

		remoteRepoURL := os.Args[2]
		localRepoDeirrectory := os.Args[3]

		if remoteRepoURL == "" {
			log.Println("Ошибка. Не передан удаленный репозиторий")
			os.Exit(1)
		}

		if localRepoDeirrectory == "" {
			log.Println("Ошибка. Не передан путь к каталогку проекта")
			os.Exit(1)
		}

		statusCode := fetch(remoteRepoURL, localRepoDeirrectory)
		os.Exit(statusCode)
	case "reset":
		if len(os.Args) < 5 {
			log.Println("Ошибка. Не переданы все необходимые аргументы. Нужно передать название удаленного репозиотрия и путь к каталогу с проектом и имя удаленной ветки")
			os.Exit(1)
		}

		statusCode := reset(os.Args[2], os.Args[3], os.Args[4])
		os.Exit(statusCode)
	case "pull":
		if len(os.Args) < 5 {
			log.Println("Ошибка. Не переданы все необходимые аргументы. Нужно передать название удаленного репозиотрия и путь к каталогу с проектом и имя удаленной ветки")
			os.Exit(1)
		}

		statusCode := pull(os.Args[2], os.Args[3], os.Args[4])
		os.Exit(statusCode)
	}

}
