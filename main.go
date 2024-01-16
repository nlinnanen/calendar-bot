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
	env_err := godotenv.Load()
	if env_err != nil {
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

	openai_client, err := InitializeOpenAIClient()
	if err != nil {
		fmt.Println(err)
		return
	}

	b.Handle(tele.OnText, func(c tele.Context) error {
		message := c.Message().Text

		extracted_meeting_info, err := extractMeetingInfo(openai_client, message)
		if err != nil {
			return handleError("Error in querying OpenAI", err, c)
		}

		parsed_meeting_info, err := parseMeetingInfo(extracted_meeting_info)
		if err != nil {
			return handleError("Error in parsing JSON", err, c)
		}

		calendar_link, err := generateCalendarLink(parsed_meeting_info.Date, parsed_meeting_info.Title)
		if err != nil {
			return handleError("Error in generating calendar link", err, c)
		}

		reply_template := "Here is the meeting info extracted from your message:\n- Date: %s\n- Title: %s\n[Link to google calendar](%s)"

		reply := fmt.Sprintf(reply_template, parsed_meeting_info.Date.Format("02.01.2006 15:04:05"), parsed_meeting_info.Title, calendar_link)

		return c.Send(reply, &tele.SendOptions{
			ParseMode: tele.ModeMarkdown,
		})
	})

	log.Output(1, "Bot is up and running 🚀")
	b.Start()
}
