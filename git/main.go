package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {

	log.Println("MOD_PREBUILD_INIT")
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
		os.Exit(1)
	} else {
		log.Print("Success, .env file found")
	}

	OptionsData = make(map[interface{}]interface{})
	OptionsData["fetch"] = "fetch ## Считать обновления с репозитория. EXAMPLE: fetch origin"
	OptionsData["reset"] = "reset ## Сбросить состояние локального репозиоря до состояния удаленного. EXAMPLE: reset origin"
	OptionsData["pull"] = "pull  ## Получить все полседние изменения. EXAMPLE: pull origin"
	OptionsData["count"] = "count ## Получить количество коммитов репозитории"

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
	case "count":
		checkArgs(2)
		statusCode := getCount()
		os.Exit(statusCode)
	case "fetch":
		checkArgs(3)
		statusCode := fetch(os.Args[2])
		os.Exit(statusCode)
	case "reset":
		checkArgs(3)
		statusCode := reset(os.Args[2])
		os.Exit(statusCode)
	case "pull":
		checkArgs(3)
		statusCode := pull(os.Args[2])
		os.Exit(statusCode)
	}

}

func checkArgs(count int) {
	if len(os.Args) < count {
		log.Println("Ошибка. Не переданы все необходимые аргументы. Нужно передать название удаленного репозиотрия и путь к каталогу с проектом и имя удаленной ветки")
		os.Exit(1)
	}
}
