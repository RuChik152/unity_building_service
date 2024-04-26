package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"unicode"
)

func getCount() int {
	dir, _ := os.LookupEnv("PATH_PROJECT")
	if dir == "" {
		log.Println("Не получен путь до каталога проекта. Завершение работы")
		return 1
	}

	countCommit, err := exec.Command("git", "-C", dir, "rev-list", "--count", "main").Output()
	if err != nil {
		return 1
	}
	var numberVersion int
	convert(countCommit, &numberVersion)

	fmt.Printf("%d", numberVersion)
	return 0

}

func convert(data []byte, number *int) {

	verStr := string(data)

	clearStr := removeControlCharacters(verStr)
	log.Println(verStr)
	verInt, err := strconv.Atoi(clearStr)
	if err != nil {
		log.Printf("Ошибка преобразования строки в число, err: %s", err)
		return
	}

	*number = verInt + 40
}

func removeControlCharacters(input string) string {
	return strings.Map(func(r rune) rune {
		if checkRune(r) {
			return -1
		}
		return r
	}, input)
}

func checkRune(r rune) bool {
	return unicode.IsControl(r)
}
