package ai

import (
	"context"
	"fmt"
	"os"

	anthropic "github.com/anthropics/anthropic-sdk-go"
)

// GenerateWordlist calls the LLM to produce a wordlist for the given description
// and prints each line to stdout. count hints at desired size (0 = let LLM decide).
func GenerateWordlist(client *anthropic.Client, description string, count int) error {
	sizeHint := ""
	if count > 0 {
		sizeHint = fmt.Sprintf(" Generate approximately %d entries.", count)
	}

	prompt := fmt.Sprintf(`Generate a wordlist for: %s%s
Output only the words or phrases, one per line. No numbering, no explanation, no blank lines, no headers.`,
		description, sizeHint,
	)

	msg, err := client.Messages.New(context.Background(), anthropic.MessageNewParams{
		Model:     anthropic.Model(model),
		MaxTokens: 2048,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		},
	})
	if err != nil {
		return fmt.Errorf("wordlist generation failed: %w", err)
	}
	if len(msg.Content) == 0 {
		return fmt.Errorf("empty response from API")
	}
	fmt.Fprintln(os.Stdout, msg.Content[0].Text)
	return nil
}
