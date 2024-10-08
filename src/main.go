package main

import (
	"context"
	"io"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	bot *tgbotapi.BotAPI
)

func main() {
	b, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}
	bot = b

	b.Debug = false
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	log.Println("bot service is running...")

	updates := b.GetUpdatesChan(u)
	for update := range updates {
		if update.Message != nil {
			if update.Message.IsCommand() {
				go handleCommand(ctx, update)
			} else {
				go handleMessage(ctx, update)
			}
		}
	}

	defer cancel()
}

func handleCommand(ctx context.Context, u tgbotapi.Update) {
	select {
	case <-ctx.Done():
		return
	default:
		switch u.Message.Command() {
		case "start":
			githubBtn := tgbotapi.NewInlineKeyboardButtonURL("Github of Creator", "https://github.com/dnridwn")
			message := "Hello! I am an Image Generation Bot.\nI was created by Den Ridwan Saputra\n\nPlease write your prompt."

			msg := tgbotapi.NewMessage(u.Message.Chat.ID, message)
			msg.ReplyToMessageID = u.Message.MessageID
			msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{githubBtn})
			_, err := bot.Send(msg)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func handleMessage(ctx context.Context, u tgbotapi.Update) {
	select {
	case <-ctx.Done():
		return
	default:
		log.Printf("receive message from %s\n", u.Message.From.UserName)

		msg := tgbotapi.NewMessage(u.Message.Chat.ID, "Please wait. We're working on it...")
		msg.ReplyToMessageID = u.Message.MessageID
		bot.Send(msg)

		file, err := SendImageGenerationRequest(u.Message.Text)
		if err != nil {
			log.Print(err)
			return
		}

		fileByte, err := io.ReadAll(file)
		if err != nil {
			log.Print(err)
			return
		}

		uploadFile := tgbotapi.FileBytes{
			Name:  u.Message.Text,
			Bytes: fileByte,
		}

		photo := tgbotapi.NewPhoto(u.Message.Chat.ID, uploadFile)
		photo.ReplyToMessageID = u.Message.MessageID
		bot.Send(photo)
	}
}
