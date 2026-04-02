package llmclientfunctions

import (
	"context"
	"strings"

	"github.com/Adebusy/receiptonboardingsvc/utilities"
	"github.com/joho/godotenv"
	"github.com/openai/openai-go"
)

func SuggestCompanyNames(clientName, fields string) []string {
	godotenv.Load()
	client := openai.NewClient()
	ctx := context.Background()
	resp, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: openai.ChatModelGPT4o,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(utilities.CompanyNamesPrompt(clientName, fields)),
		},
	})

	if err != nil {
		panic(err)
	}
	jsonData := (resp.Choices[0].Message.Content)
	clean := utilities.TrimString(jsonData)
	return strings.Split(clean, "\n")
}

func SuggestCompanyLogo(clientName, fields string) []string {
	godotenv.Load()
	message := utilities.LogoConceptPrompt(clientName, fields)
	client := openai.NewClient()
	ctx := context.Background()

	resp, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: openai.ChatModelGPT4o,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(message),
		},
	})

	if err != nil {
		panic(err)
	}
	jsonData := (resp.Choices[0].Message.Content)
	clean := utilities.TrimString(jsonData)
	return strings.Split(clean, "\n")
}

func SuggestCompanySignature(firstName, lastName string) []string {
	godotenv.Load()
	message := utilities.CompanySignaturePrompt(firstName + " " + lastName)
	client := openai.NewClient()
	ctx := context.Background()

	resp, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: openai.ChatModelGPT4o,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(message),
		},
	})

	if err != nil {
		panic(err)
	}
	jsonData := (resp.Choices[0].Message.Content)
	clean := utilities.TrimString(jsonData)
	return strings.Split(clean, "\n")
}

func SuggestCompanyReceiptTemplate(clientName, fields string) string {

	return ""
}
