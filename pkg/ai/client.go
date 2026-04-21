package ai

import (
	"errors"
	"os"

	anthropic "github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

const model = "claude-haiku-4-5-20251001"

// NewClient creates an Anthropic client from ANTHROPIC_API_KEY.
// Returns an error if the key is not set.
func NewClient() (*anthropic.Client, error) {
	key := os.Getenv("ANTHROPIC_API_KEY")
	if key == "" {
		return nil, errors.New("ANTHROPIC_API_KEY is not set")
	}
	c := anthropic.NewClient(option.WithAPIKey(key))
	return &c, nil
}
