package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func init() {

	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
		os.Exit(1)
	} else {
		log.Print("Success, .env file found")
		YANDEX_TOKEN, _ = os.LookupEnv("YANDEX_TOKEN")
		if YANDEX_TOKEN == "" {
			log.Println("Не установлен OAuth токен для доступа к API Yandex")
			os.Exit(1)
		}
	}

	initArgs()

}

func main() {
	dataMap, errGetUrl := getTargetURL(ARG["-name"], ARG["-platform"])
	log.Println(dataMap)
	log.Println(errGetUrl)

	status, errSendFile := sendBuildToCloud(dataMap["href"].(string), ARG["-path"])
	if errSendFile != nil {
		log.Println("Ошибка отправки данных в YandexCloud", errSendFile)
		os.Exit(1)
	} else {
		log.Println("Отправка данных в YandexCloud успешна: ", status)
		os.Exit(0)
	}

}

/*
Получение URL для загрузки файлов в Yandex Disk
*/
func getTargetURL(name string, platform string) (map[string]interface{}, error) {
	url := "https://cloud-api.yandex.net/v1/disk/resources/upload"
	switch platform {
	case "Android":
		url += "?path=BeliVR/Kagra_builds/Android/" + name
	case "Win64":
		url += "?path=BeliVR/Kagra_builds/PC/" + name
	}
	url += "&overwrite=true"
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания HTTP запроса: %s", err)
	}

	req.Header.Add("Authorization", "OAuth "+YANDEX_TOKEN)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка HTTP запроса %s", err)
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	errs := json.NewDecoder(resp.Body).Decode(&data)
	if errs != nil {
		return nil, fmt.Errorf("ошибка чтения ответа: %s", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("ошибка запроса, код ответа: %d. Ответ: %s", resp.StatusCode, data["message"].(string))
	}

	return data, nil

}

/*
Отправка файла по указанному пути полученного от Yandex API
*/
func sendBuildToCloud(url string, filePath string) (int, error) {

	client := &http.Client{}

	file, err := os.Open(filePath)
	if err != nil {
		return 0, errors.New(fmt.Sprint(err))
	}
	defer file.Close()

	buffer, err := io.ReadAll(file)
	if err != nil {
		return 0, errors.New(fmt.Sprint(err))
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(buffer))
	if err != nil {
		return 0, errors.New(fmt.Sprint(err))
	}

	res, _ := client.Do(req)

	return res.StatusCode, nil
}
