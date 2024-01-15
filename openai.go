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
		Todays date is %s Extract the date and time and name of the meeting from the following message. 
		
		Return a JSON string { date: the date and time in ISO format eg 2024-11-21T09:00:00Z, title: the title of the meeting}.
		
		Here is the message:
		%s
	`

	prompt := fmt.Sprintf(prompt_template, time.Now().Format("2006-01-02"), message)

	completion, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: "ft:gpt-3.5-turbo-0613:personal::8h3P6v9R",
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
