package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

// InitializeOpenAIClient initializes and returns an OpenAI client
func InitializeOpenAIClient() (*openai.Client, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable is not set")
	}

	client := openai.NewClient(apiKey)

	return client, nil
}

type MeetingInfo struct {
	Date  time.Time `json:"date"`
	Title string    `json:"title"`
}

func extractMeetingInfo(client *openai.Client, message string) string {
	current_date := time.Now().Format("2006-01-02")
	prompt := "Todays date is " + current_date + "Extract the date and time and name of the meeting from the following message. Return a JSON string { date: the date and time in ISO format eg 2024-11-21T09:00:00Z, title: the title of the meeting}. Here is the message:  " + message

	completion, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		fmt.Println("Error getting completion:", err)
		return "Error getting completion!"
	}

	fmt.Println("Completion:", completion.Choices[0].Message.Content)

	var meeting_info MeetingInfo

	parsing_err := json.Unmarshal([]byte(completion.Choices[0].Message.Content), &meeting_info)
	if parsing_err != nil {
		fmt.Println("Error parsing JSON:", parsing_err)
		return "Error parsing JSON!: " + completion.Choices[0].Message.Content
	}
	end_date := meeting_info.Date.Add(time.Hour * 1)
	formatted_start_date := meeting_info.Date.Format("20060102T150405") + "UTC+2"
	formatted_end_date := end_date.Format("20060102T150405") + "UTC+2"

	return "http://www.google.com/calendar/event?action=TEMPLATE&text=" + strings.ReplaceAll(meeting_info.Title, " ", "+") + "&dates=" + formatted_start_date + "/" + formatted_end_date + "UTC+02&details=&location=&trp=false"

}
