package main

import (
	"log"
	"os"

	"tg-pc-control/logger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	if err := logger.Init("bot.log"); err != nil {
		log.Fatal(err)
	}

	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		logger.Log.Fatal("set BOT_TOKEN env var")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		logger.Log.Fatal(err)
	}

	logger.Log.Printf("Authorized: %s", bot.Self.UserName)

	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("📈 Status"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🕒 Uptime"),
		),
	)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		logger.Log.Printf("RAW: %+v", update)

		if update.Message == nil {
			continue
		}

		chatID := update.Message.Chat.ID
		text := update.Message.Text
		logger.Log.Printf("MSG %d: %s", chatID, text)

		switch text {
		case "/start":
			msg := tgbotapi.NewMessage(chatID, "Hello! I control your PC🖥")
			msg.ReplyMarkup = keyboard
			if _, err := bot.Send(msg); err != nil {
				logger.Log.Printf("SEND ERR: %v", err)
			}

		default:
			msg := tgbotapi.NewMessage(chatID, "Echo: "+text)
			if _, err := bot.Send(msg); err != nil {
				logger.Log.Printf("SEND ERR: %v", err)
			}
		}
	}
}
