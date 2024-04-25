package main

import (
	"flag"
	"fmt"
)

var YANDEX_TOKEN string
var ARG map[string]string = make(map[string]string, 0)

func initArgs() {
	name := flag.String("name", "", "Имя файла")
	path := flag.String("path", "", "Локальный путь к файлу")
	platform := flag.String("platform", "", "Имя платформы Android, Win64...")
	// path := flag.String("path", "", "Описание для path")
	// re := flag.String("re", "", "Описание для re")

	flag.Parse()

	ARG["-name"] = *name
	ARG["-path"] = *path
	ARG["-platform"] = *platform
	// ARG["-path"] = *path
	// ARG["-re"] = *re

	fmt.Println(ARG)
	// log.Println(ARG["--name"])
}
