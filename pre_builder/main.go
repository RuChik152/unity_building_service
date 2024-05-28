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
	OptionsData["copy"] = "copy ## Производит копирование файлов\nEXAMPLE: copy <\"C:\\path\\to\\source\\file\\file.txt\"> <\"C:\\path\\to\\dest\\file\\file.txt\">\n\n"
	OptionsData["create"] = "create ## Производит создание\\перезапись файла конфигурации для нужной платформы, необходимо указать путь до корневого каталога проекта, имя ветки и полный путь куда сохранять файл \nEXAMPLE: create <\"C:\\path\\to\\project\\dirrctory\"> <\"C:\\path\\to\\dest\\file\\file.txt\">\n\n"
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
	case "copy":
		if len(os.Args) < 4 {
			log.Println("Ошибка. Не переданы все необходимые аргументы.")
			os.Exit(1)
		}

		statusCode := copy(os.Args[2], os.Args[3])
		os.Exit(statusCode)
	case "create":
		if len(os.Args) < 4 {
			log.Println("Ошибка. Не переданы все необходимые аргументы.")
			os.Exit(1)
		}

		statusCode := createGlobalConstant(os.Args[2], os.Args[3], Platform(os.Args[4]))
		os.Exit(statusCode)
	}
}
