package main

import (
	"bytes"
	"log"
	"os"
	"os/exec"
)

func reset(url string) int {
	var statusCode int

	dir, _ := os.LookupEnv("PATH_PROJECT")
	if dir == "" {
		log.Println("Не получен путь до каталога проекта. Завершение работы")
		os.Exit(1)
	}

	var remoteBranch = url + "/main"

	cmd := exec.Command("git", "-C", dir, "reset", "--hard", remoteBranch)

	var output bytes.Buffer
	cmd.Stdout = &output
	err := cmd.Run()
	if err != nil {
		log.Printf("Ошибка выполнения RESET: %s", err)
		return 1
	}

	if statusCode = cmd.ProcessState.ExitCode(); statusCode != 0 {
		log.Printf("Команда RESET выполнена с ошибкой. Код ошибки: %d", statusCode)
		return statusCode
	} else {
		log.Printf("Команда RESET выполнена успешно. %d", statusCode)
		log.Println("Вывод выполнения данных при команде GIT RESET: \n" + output.String() + "\n")
		return statusCode
	}
}
