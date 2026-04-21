package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ariary/cfuzz/pkg/ai"
	"github.com/ariary/cfuzz/pkg/fuzz"
	"github.com/ariary/cfuzz/pkg/mcp"
	"github.com/spf13/cobra"
)

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Start cfuzz as an MCP server over stdio",
	Long: `Start cfuzz as a Model Context Protocol (MCP) server.

Communicates over stdio using JSON-RPC 2.0 and exposes a "fuzz" tool.

To register with Claude Desktop, add to ~/.claude/claude_desktop_config.json:
  {
    "mcpServers": {
      "cfuzz": { "command": "cfuzz", "args": ["mcp"] }
    }
  }`,
	SilenceUsage:  true,
	SilenceErrors: true,
	Run:           runMCP,
}

func init() {
	rootCmd.AddCommand(mcpCmd)
}

func runMCP(_ *cobra.Command, _ []string) {
	fmt.Fprintln(os.Stderr, "cfuzz MCP server started (stdio)")
	mcp.Serve(handleFuzzCall)
}

func handleFuzzCall(args map[string]any) (string, error) {
	cfg := fuzz.DefaultConfig()
	cfg.HideBanner = true
	cfg.OnlyWord = false

	command, ok := args["command"].(string)
	if !ok || command == "" {
		return "", fmt.Errorf("command is required")
	}
	cfg.Command = command

	rawList, ok := args["wordlist"].([]any)
	if !ok || len(rawList) == 0 {
		return "", fmt.Errorf("wordlist is required and must be a non-empty array")
	}

	tmp, err := os.CreateTemp("", "cfuzz-mcp-*.txt")
	if err != nil {
		return "", fmt.Errorf("failed to create temp wordlist: %w", err)
	}
	defer os.Remove(tmp.Name())
	for _, entry := range rawList {
		if s, ok := entry.(string); ok {
			fmt.Fprintln(tmp, s)
		}
	}
	tmp.Close()
	cfg.Wordlists = []string{tmp.Name()}

	if v, ok := args["threads"].(float64); ok {
		cfg.Threads = int(v)
	}
	if v, ok := args["timeout"].(float64); ok {
		cfg.Timeout = int64(v)
	}
	if v, ok := args["success_only"].(bool); ok && v {
		cfg.Filters = append(cfg.Filters, fuzz.CodeSuccessFilter{Zero: true})
	}
	if v, ok := args["stdout_word"].(string); ok && v != "" {
		cfg.Filters = append(cfg.Filters, fuzz.StdoutWordFilter{TargetWord: v})
	}
	if v, ok := args["ai_filter"].(string); ok && v != "" {
		client, err := ai.NewClient()
		if err != nil {
			return "", fmt.Errorf("ai_filter requires ANTHROPIC_API_KEY: %w", err)
		}
		criterion := v
		cfg.AIFilterFn = func(word, stdout, stderr, code string) bool {
			return ai.IsInteresting(client, criterion, word, stdout, stderr, code)
		}
	}

	var buf strings.Builder
	cfg.ResultLogger = log.New(&buf, "", 0)
	cfg.DisplayModes = fuzz.BuildDisplayModes(false, false, false, false, false)

	fuzz.PerformFuzzing(cfg)
	return buf.String(), nil
}
