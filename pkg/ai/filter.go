package ai

import (
	"context"
	"fmt"
	"os"
	"strings"

	anthropic "github.com/anthropics/anthropic-sdk-go"
)

// IsInteresting sends the result to the LLM and asks if it matches the criterion.
// Returns true if the LLM responds YES, or if an API error occurs (fail-open).
func IsInteresting(client *anthropic.Client, criterion, word, stdout, stderr, code string) bool {
	prompt := fmt.Sprintf(`You are evaluating command execution results for a security researcher.
Criterion: %s

WORD: %s
STDOUT: %s
STDERR: %s
EXIT CODE: %s

Does this result match the criterion? Reply with exactly YES or NO.`,
		criterion, word,
		truncate(stdout, 500),
		truncate(stderr, 500),
		code,
	)

	msg, err := client.Messages.New(context.Background(), anthropic.MessageNewParams{
		Model:     anthropic.Model(model),
		MaxTokens: 16,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		},
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "ai-filter error: %v\n", err)
		return true // fail-open: show the result
	}
	if len(msg.Content) == 0 {
		return true
	}
	return strings.TrimSpace(msg.Content[0].Text) == "YES"
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "…"
}
