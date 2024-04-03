package main

import (
	"io"
	"log"
	"os"
)

func copy(pathSourceFile string, pathDestFile string) int {
	log.Printf(":::Был получен путь %s к целефому файлу \n", pathDestFile)
	log.Printf(":::Был получен путь %s к исходному файлу \n", pathSourceFile)

	sourceFile, err := os.Open(pathSourceFile)
	if err != nil {
		log.Println("::::Ошибка открытия целевого файла. ERROR: ", err)
		return 1
	}
	defer sourceFile.Close()

	destFile, err := os.Create(pathDestFile)
	if err != nil {
		log.Printf("::::Ошибка создания целевого файла %s, ERROR: %s \n", pathDestFile, err)
		return 1
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		log.Printf("::Копирование не успешное. Ошибка: %s", err)
		return 1
	} else {
		log.Printf(":::Копирование данных из %s в %s успешное \n", pathSourceFile, pathDestFile)
		return 0
	}
}
