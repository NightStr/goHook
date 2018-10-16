package main

import (
	"fmt"
	"github.com/Syfaro/telegram-bot-api"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Config struct {
	KEY  string
	HOST string
	PORT string
}

func telegramBotResponding(bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		if "/get_url" == update.Message.Text {
			messageText := os.Getenv("HOST_NAME") + strconv.FormatInt(update.Message.Chat.ID, 10) + "/"
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
	var bot, err = tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_KEY"))
	bot.Debug = true
	if err != nil {
		log.Panic(err)
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	go telegramBotResponding(bot, updates)

	router := mux.NewRouter()
	router.HandleFunc("/{chatId:\\d+}/", makeHandler(bot)).Methods("POST")

	if err != nil {
		log.Panic(err)
	}
	port := os.Getenv("PORT")

	http.Handle("/", router)
	http.ListenAndServe(":"+port, nil)
}
