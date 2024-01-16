package main

import (
	"context"
	"fmt"
	"os"
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

func extractMeetingInfo(client *openai.Client, message string) (string, error) {
	prompt_template := `
		Todays time and date is %s Extract the date and time and name of the meeting from the following message.
		
		Return a JSON string { date: the date and time in ISO format eg 2024-11-21T09:00:00Z, title: the title of the meeting}.

		Examples if todays date is tuesday 16.1.2023:
		Input: Sovitaan maanantai 22.1. Klo 14-16 ekan miitin ajankohdaksi🤝🤝
		Output: { date: 2023-01-22T14:00:00Z, title: Eka miitti }

		Input: Sovitaan huomenna 9 aamulla miitti
		Output: { date: 2023-01-17T09:00:00Z, title: Miitti }

		Input: Sovitaan perjantaina ysistä eteenpäin lounas
		Output: { date: 2023-01-19T09:00:00Z, title: Lounas }

		Input: Otetaan pe puoleltapäivin tapaaminen
		Output: { date: 2023-01-20T12:00:00Z, title: Tapaaminen }

		Input: Valitaan ke klo 14-16 TUAS-talolla
		Output: { date: 2023-01-17T14:00:00Z, title: TUAS-talo }

		Input: Valitaan ens viikon tiistai kympiltä
		Output: { date: 2023-01-23T10:00:00Z, title: Miitti }
		
		Here is the message:p
		%s
	`
	current_date_time := time.Now().Format("2006-01-02 15:04:05")
	prompt := fmt.Sprintf(prompt_template, current_date_time, message)

	completion, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)

	return completion.Choices[0].Message.Content, err
}
