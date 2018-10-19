package main

import (
	"log"
	"os"
	"telegaBot/hookBot"
)

func main() {
	bot, err := hookBot.NewChatBot(
		os.Getenv("TELEGRAM_KEY"),
		os.Getenv("PORT"),
		os.Getenv("HOST_NAME"),
		os.Getenv("TELEGRAM_DEBUG") == "true",
	)
	if err != nil {
		log.Panic(err)
	}
	bot.Start()
}
