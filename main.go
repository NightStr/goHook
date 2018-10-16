package main

import (
	"fmt"
	"github.com/Syfaro/telegram-bot-api"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type Config struct {
	KEY  string
	HOST string
	PORT string
}

func loadConfing() Config {
	config := Config{}
	data, err := ioutil.ReadFile("./config.yaml")

	if err != nil {
		log.Panic(err)
	}

	err = yaml.Unmarshal(data, &config)

	if err != nil {
		log.Panic(err)
	}
	return config
}

func telegramBotResponding(bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel, config Config) {
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		if "/get_url" == update.Message.Text {
			messageText := config.HOST + strconv.FormatInt(update.Message.Chat.ID, 10) + "/"
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, messageText)
			bot.Send(msg)
		}
	}
}

func makeHandler(bot *tgbotapi.BotAPI) func(resp http.ResponseWriter, req *http.Request) {
	return func(resp http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		body, err := ioutil.ReadAll(req.Body)
		var message string

		if err != nil {
			message = err.Error()
		} else {
			message = string(body)
		}
		chatId, err := strconv.ParseInt(vars["chatId"], 10, 64)

		if err == nil {
			telegramMessage := tgbotapi.NewMessage(chatId, message)
			bot.Send(telegramMessage)
		}
		fmt.Fprint(resp, "Ok")
	}
}

func main() {
	config := loadConfing()
	var bot, err = tgbotapi.NewBotAPI(config.KEY)
	if err != nil {
		log.Panic(err)
	}
	router := mux.NewRouter()

	router.HandleFunc("/{chatId:\\d+}/", makeHandler(bot)).Methods("POST")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	if err != nil {
		log.Panic(err)
	}

	go telegramBotResponding(bot, updates, config)
	http.Handle("/", router)
	http.ListenAndServe(":"+config.PORT, nil)
}
