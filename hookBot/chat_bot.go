package hookBot

import (
	"github.com/Syfaro/telegram-bot-api"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

type HookBot struct {
	serverKey    string
	bot          *tgbotapi.BotAPI
	updates      tgbotapi.UpdatesChannel
	updateConfig tgbotapi.UpdateConfig
	debug        bool
	commands     map[string]func(message *tgbotapi.Message) (string, error)
	middleware   []func(string) string

	httpPort string
	hostName string
}

func (bot *HookBot) AddCommand(commandName string, handler func(message *tgbotapi.Message) (string, error)) {
	bot.commands[commandName] = handler
}

func (bot *HookBot) postHandler(resp http.ResponseWriter, req *http.Request) {
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
		go bot.sendMessage(chatId, message)
	}
}

func (bot HookBot) telegramUpdate() {
	var msg tgbotapi.MessageConfig
	for update := range bot.updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		message := update.Message
		handler := bot.commands[message.Text]
		if handler != nil {
			response, err := handler(message)
			if err != nil {
				if bot.debug {
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
				} else {
					log.Panic(err)
				}
			} else {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, response)
			}
			bot.bot.Send(msg)
		}
	}
}

func (bot *HookBot) sendMessage(chatId int64, message string) (tgbotapi.Message, error) {
	for _, mw := range bot.middleware {
		message = mw(message)
	}
	time.Sleep(2 * time.Second)
	msg := tgbotapi.NewMessage(chatId, message)
	return bot.bot.Send(msg)
}

func (bot HookBot) getUrl(message *tgbotapi.Message) (string, error) {
	msg := bot.hostName + strconv.FormatInt(message.Chat.ID, 10) + "/"
	return msg, nil
}

func (bot *HookBot) AddMiddleware(mw func(string) string) {
	bot.middleware = append(bot.middleware, mw)
}

func (bot *HookBot) Start() {
	go bot.telegramUpdate()
	http.ListenAndServe(":"+bot.httpPort, nil)
}

func NewChatBot(key string, httpPort string, hostName string, debug bool) (*HookBot, error) {
	var bot, err = tgbotapi.NewBotAPI(key)
	bot.Debug = debug
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates, err := bot.GetUpdatesChan(updateConfig)

	chatBot := HookBot{
		bot:          bot,
		debug:        debug,
		updates:      updates,
		updateConfig: updateConfig,
		httpPort:     httpPort,
		hostName:     hostName,
		commands:     make(map[string]func(message *tgbotapi.Message) (string, error)),
		middleware:   []func(string) string{},
	}
	chatBot.AddCommand("/get_url", chatBot.getUrl)

	router := mux.NewRouter()
	router.HandleFunc("/{chatId:\\d+}/", chatBot.postHandler).Methods("POST")

	if err != nil {
		log.Panic(err)
	}

	http.Handle("/", router)
	if err != nil {
		return nil, err
	} else {
		return &chatBot, nil
	}
}
