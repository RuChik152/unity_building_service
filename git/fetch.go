package main

import (
	"bytes"
	"log"
	"os/exec"
)

func fetch(url string, dirrectory string) int {

	var statusCode int

	cmd := exec.Command("git", "-C", dirrectory, "fetch", url, "master")
	var output bytes.Buffer
	cmd.Stdout = &output
	err := cmd.Run()
	if err != nil {
		log.Printf("Ошибка выполнения FETCH: %s", err)
		return 1
	}

	if statusCode = cmd.ProcessState.ExitCode(); statusCode != 0 {
		log.Printf("Команда FETCH выполнена с ошибкой. Код ошибки: %d", statusCode)
		return statusCode
	} else {
		log.Printf("Команда FETCH выполнена успешно. %d", statusCode)
		log.Println("Вывод выполнения данных при команде GIT FETCH: \n" + output.String() + "\n")
		return statusCode
	}
}
