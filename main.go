package main

import (
	"context"
	"io"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

var (
	bot *tgbotapi.BotAPI
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	b, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatal("")
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
		go handleMessage(ctx, update)
	}

	defer cancel()
}

func handleMessage(ctx context.Context, u tgbotapi.Update) {
	select {
	case <-ctx.Done():
		return
	default:
		log.Printf("receive message from %s\n", u.Message.From.UserName)

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
		bot.Send(photo)
	}
}
