package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v3"
)

func main() {
	env_err := godotenv.Load()
	if env_err != nil {
		log.Fatal("Error loading .env file")
	}

	preferences := tele.Settings{
		Token:  os.Getenv("BOT_TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(preferences)
	if err != nil {
		log.Fatal(err)
		return
	}

	openai_client, err := InitializeOpenAIClient()
	if err != nil {
		fmt.Println(err)
		return
	}

	b.Handle(tele.OnText, func(c tele.Context) error {
		message := c.Message().Text
		extraction_result := extractMeetingInfo(openai_client, message)
		return c.Send(extraction_result)
	})

	b.Start()
}
