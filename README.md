# Calendar Bot 📅

A Telegram bot that generates a link to a calendar event from a message. It uses OpenAI's GPT-4 to parse the title and time of the meeting and generate the link. The usecase is when a meeting is scheduled in a chat, with a message "Meeting at 3pm tomorrow", the message can be forwarded to the bot. The bot will then generate a link to a calendar event with the title "Meeting" and the time "3pm tomorrow".

The bot is currently hosted on Fly and can be found with the nickname [@KalenteriLinkBot](https://t.me/kalenterilinkbot).

## Usage

Send the bot a message with the title and time of the event. The bot will then generate a link to a calendar event with the title and time.

## Local development

1. Install go from [here](https://golang.org/dl/)
1. Pull the project from github
1. Run the following command in the project directory
1. Make a .env file in the root and fill in the variables defined in [.env.example](.env.example)
   1. Get the Telegram bot token from [here](https://core.telegram.org/bots/tutorial#obtain-your-bot-token)
   1. Get the OpenAI API key from [here](https://platform.openai.com/docs/quickstart/account-setup)
1. Run the project with `go run .`

## License

Licensed under the [MIT License](LICENSE).
