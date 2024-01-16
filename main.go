package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v3"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Output(1, "Error loading .env file")
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

	if os.Getenv("PUBLIC_URL") != "" {
		log.Output(1, "Setting webhook")
		b.SetWebhook(&tele.Webhook{
			Listen: os.Getenv("PUBLIC_URL"),
		})
	}

	// Initialize server to handle root path for health check
	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "OK")
		})
		http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "OK")
		})
		http.ListenAndServe(":8080", nil)
	}()

	openAIClient, err := InitializeOpenAIClient()
	if err != nil {
		fmt.Println(err)
		return
	}

	b.Handle("/start", func(c tele.Context) error {
		reply := "Hello 👋 I'm a bot that can generate a Google Calendar link from a meeting description. Just send me a message with the meeting details and I'll do the rest 🤖"
		return c.Send(reply)
	})

	b.Handle(tele.OnText, func(c tele.Context) error {
		message := c.Message().Text

		extractedMeetingInfo, err := extractMeetingInfo(openAIClient, message)
		if err != nil {
			return handleError("Error in querying OpenAI", err, c)
		}

		parsedMeetingInfo, err := parseMeetingInfo(extractedMeetingInfo)
		if err != nil {
			return handleError("Error in parsing JSON", err, c)
		}

		calendarLink, err := generateCalendarLink(parsedMeetingInfo.Date, parsedMeetingInfo.Title)
		if err != nil {
			return handleError("Error in generating calendar link", err, c)
		}

		replyTemplate := "📅 Date: %s\n📣 Title: %s\n[🔗 Link to google calendar 🔗](%s)"

		reply := fmt.Sprintf(replyTemplate, parsedMeetingInfo.Date.Format("02.01.2006 15:04"), parsedMeetingInfo.Title, calendarLink)

		return c.Send(reply, &tele.SendOptions{
			ParseMode: tele.ModeMarkdown,
		})
	})

	log.Output(1, "Bot is up and running 🚀")
	b.Start()
}
