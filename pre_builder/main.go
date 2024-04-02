package main

import (
	"fmt"
	"log"
	"os"
)

func init() {

	OptionsData = make(map[interface{}]interface{})
	OptionsData["copy"] = "copy ## Производит копирование настроек для целевой платформы для XRGeneralSettings"
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
		if len(os.Args) < 5 {
			log.Println("Ошибка. Не переданы все необходимые аргументы.")
			os.Exit(1)
		}

		statusCode := createGlobalConstant(os.Args[2], os.Args[3], os.Args[4])
		os.Exit(statusCode)
	}
}
