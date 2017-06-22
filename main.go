package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"vederko-bot/amr"

	"gopkg.in/telegram-bot-api.v4"
)

func getAmrOwner(plate string) string {
	return amr.Names[plate]
}

// Hack for heroku
func mainHandler(resp http.ResponseWriter, _ *http.Request) {
	resp.Write([]byte("Hi there! I'm Vederko Telegram bot!"))
}

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_API_KEY"))
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	port := os.Getenv("PORT")
	if port != "" {
		http.HandleFunc("/", mainHandler)
		go http.ListenAndServe(":"+port, nil)
	}

	updates := bot.ListenForWebhook("/" + bot.Token)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		owner := getAmrOwner(update.Message.Text)
		if owner == "" {
			owner = "Владелец не найден"
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, owner)
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}
}
