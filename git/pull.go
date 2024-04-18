package main

import (
	"bytes"
	"log"
	"os"
	"os/exec"
)

func pull(url string) int {
	var statusCode int

	dir, _ := os.LookupEnv("PATH_PROJECT")
	if dir == "" {
		log.Println("Не получен путь до каталога проекта. Завершение работы")
		os.Exit(1)
	}

	cmd := exec.Command("git", "-C", dir, "pull", url, "main")
	var output bytes.Buffer
	cmd.Stdout = &output
	err := cmd.Run()
	if err != nil {
		log.Printf("Ошибка выполнения PULL: %s", err)
		return 1
	}

	if statusCode = cmd.ProcessState.ExitCode(); statusCode != 0 {
		log.Printf("Команда PULL выполнена с ошибкой. Код ошибки: %d", statusCode)
		return statusCode
	} else {
		log.Printf("Команда PULL выполнена успешно. %d", statusCode)
		log.Println("Вывод выполнения данных при команде GIT PULL: \n" + output.String() + "\n")
		return statusCode
	}
}
