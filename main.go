package main

import (
	"github.com/asdine/storm"
	"github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
	"net/http"
)

var Config configParsed
var log *logrus.Logger
var storage *storm.DB
var BotMulticast = make(chan Message, 4)

func init() {
	var err error
	log = logrus.New()
	Config = GetConfig()
	log.SetLevel(Config.LogLevel)
	storage, err = storm.Open(Config.STPath)
	if err != nil {
		log.Fatalf("Can not open storage, error: %s", err)
	}
}

func main() {
	bot, err := tb.NewBot(tb.Settings{
		Token:  Config.TgToken,
		Poller: &tb.LongPoller{Timeout: Config.TgPollInterval},
	})
	if err != nil {
		log.Fatal("Failed to setup telegram bot, error: %s", err)
	}

	go StartBot(bot)
	go BotSender(bot, BotMulticast)

	http.HandleFunc(Config.GitHubHookPath, GitHubHookHandler)
	http.HandleFunc(Config.TravisHookPath, TravisHookHandler)
	http.HandleFunc("/health", HealthCheckHandler)
	if err := http.ListenAndServe(Config.WebListen, nil); err != nil {
		log.Fatalf("Web server error: %s", err)
	}
}
