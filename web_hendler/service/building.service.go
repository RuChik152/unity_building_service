package service

import (
	"log"
	"os/exec"
)

func Manager() {
	if !STATUS_BUILDING {
		gitArgsFetch := []string{
			"fetch",
			"origin",
			PROJECT_FOLDER,
		}
		runGit := exec.Command("..\\git\\git_runner.exe", gitArgsFetch...)
		gitOutput, err := runGit.CombinedOutput()
		if err != nil {
			log.Println("Ошибка загрузки: ", string(gitOutput), "\n", err)
		} else {
			log.Println(string(gitOutput))
		}
	} else {
		//логика для перезапуска сборки
	}
}
