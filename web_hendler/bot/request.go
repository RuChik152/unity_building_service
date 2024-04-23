package bot

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type ResultMsgBuildForJSON struct {
	Event         string           `json:"event"`
	OculusMessage DeviceBotMessage `json:"oculus"`
	PicoMessage   DeviceBotMessage `json:"pico"`
	PCMessage     DeviceBotMessage `json:"pc"`
	Info          BuildInfo        `json:"info"`
}

type BotMessage struct {
	Event         string           `json:"event"`
	OculusMessage DeviceBotMessage `json:"oculus"`
	PicoMessage   DeviceBotMessage `json:"pico"`
	PCMessage     DeviceBotMessage `json:"pc"`
	Info          BuildInfo        `json:"info"`

	Mu sync.Mutex
}

func (bm *BotMessage) MarshalJSON() ([]byte, error) {
	tempStruct := ResultMsgBuildForJSON{
		Event:         bm.Event,
		OculusMessage: bm.OculusMessage,
		PicoMessage:   bm.PicoMessage,
		PCMessage:     bm.PCMessage,
		Info:          bm.Info,
	}

	return json.Marshal(tempStruct)
}

type DeviceBotMessage struct {
	Status    bool   `json:"status"`
	Message   string `json:"message"`
	SendBuild string `json:"sendStatus"`
}

type BuildInfo struct {
	DataVersion string `json:"version"`
	OculusLogs  string `json:"oculus"`
	PicoLogs    string `json:"pico"`
}

var ResultMsgBuild BotMessage

var StandartMsg StandartMsgBot

var CommitMsg CommitData

type MessageBot struct {
	Msg  string    `json:"msg"`
	Data time.Time `json:"data"`
	Tag  string    `json:"tag"`
}

type StandartMsgBot struct {
	Event   string `json:"event"`
	Message string `json:"message"`
}

type CommitData struct {
	Event   string `json:"event"`
	ID      string `json:"ID"`
	SHA     string `json:"sha"`
	AUTHOR  string `json:"author"`
	MESSAGE string `json:"message"`
}

func SendMessageBot(msg string, tg string) {
	bot_port, empty_bot_port := os.LookupEnv("BOT_PORT")
	bot_url, empty_bot_url := os.LookupEnv("BOT_URL")

	if empty_bot_port && empty_bot_url {

		var currentTime time.Time = time.Now()
		message := MessageBot{
			Msg:  msg,
			Data: currentTime,
			Tag:  tg,
		}

		payload, err := json.Marshal(message)
		if err != nil {
			log.Println("Ошибка преобразования тела звпроса для бота в JSON формат данных")
		}

		resp, err := http.Post(bot_url+":"+bot_port+"/echo_build", "application/json", bytes.NewBuffer(payload))
		if err != nil {
			log.Println("Ошибка отправки сообщения боту")
		} else {
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Println("::Ошибка чтения ответа: ", err)
			} else {
				log.Println("::Получен ответ от бота: ", string(body))
			}
		}

	} else {
		log.Println("Отсуствую данные для подключния к боту")
	}
}
