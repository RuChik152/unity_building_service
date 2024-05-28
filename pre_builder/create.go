package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type Platform string

func createGlobalConstant(dirrectory string, pathFile string, platform Platform) int {
	countCommit, _ := exec.Command("git", "-C", dirrectory, "rev-list", "--count", "main").Output()
	shortHashCommit, _ := exec.Command("git", "-C", dirrectory, "rev-parse", "--short", "main").Output()
	dataBuild := time.Now().Format("02_01_2006")

	log.Println(shortHashCommit)
	//hash := string(shortHashCommit)
	hash := removeControlCharacters(string(shortHashCommit))
	log.Println(hash)

	var numberVersion int
	convert(countCommit, &numberVersion)

	fileContent := fmt.Sprintf(`
	using System.Collections;
	using UnityEngine;

	public class GlobalConstants : MonoBehaviour
	{
		public const string MessageInTheSky = "%s, v0.%d";
		public const string VersionBandel = "%d";
		public const string ProjectVersion = "0.%d";
		public const string ShortHashCommit = "%s";
		public const string DataBuild = "%s";
		public const string Platform = "%s";

		void
		Start()
		{
			GL.BLD_V = MessageInTheSky;
			GL.PLTFRM = Platform;
		}
	}

	`, dataBuild, numberVersion, numberVersion, numberVersion, hash, dataBuild, platform)

	file, err := os.Create(pathFile)
	if err != nil {
		log.Printf("Ошибка создания файла. %s", err)
		return 1
	}
	defer file.Close()

	_, err = file.WriteString(fileContent)
	if err != nil {
		log.Printf("Ошибка записи в файл, err: %s", err)
		return 1
	} else {
		log.Printf("Запись успешна")
		return 0
	}
}

func convert(data []byte, number *int) {
	log.Println(data)
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
