package main

import (
	"fmt"
	"github.com/Syfaro/telegram-bot-api"
	"io/ioutil"
	"log"
	"net/http"
)


const CHAT_ID = 0
var bot, err = tgbotapi.NewBotAPI("")


func url_generator(bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		if "/get_url" == update.Message.Text {
			messageText := fmt.Sprintf("https://someurl.ru/%v/", update.Message.Chat.ID)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, messageText)
			bot.Send(msg)
		}
	}
}


func msgHandler(resp http.ResponseWriter, req *http.Request) {
	message, _ := ioutil.ReadAll(req.Body)
	telegramMessage := tgbotapi.NewMessage(CHAT_ID, string(message))
	bot.Send(telegramMessage)
}


func main() {
	// подключаемся к боту с помощью токена
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.GetUpdatesChan(u)

	go url_generator(bot, updates)
	http.HandleFunc("/", msgHandler)
	http.ListenAndServe(":8080", nil)
}
