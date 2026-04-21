package main

import (
	"fmt"
	"os"

	"github.com/ariary/cfuzz/pkg/ai"
	"github.com/spf13/cobra"
)

var wordlistCmd = &cobra.Command{
	Use:   "wordlist \"description\"",
	Short: "Generate a wordlist using AI",
	Long: `Generate a context-aware wordlist by describing what you need.
Output is printed to stdout, one entry per line, suitable for piping into cfuzz.

Requires ANTHROPIC_API_KEY.

Examples:
  cfuzz wordlist "default credentials for network switches"
  cfuzz wordlist "common web admin paths" -n 50 | cfuzz --stdin-wordlist "curl -s http://target/FUZZ"`,
	Args:          cobra.ExactArgs(1),
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE:          runWordlist,
}

func init() {
	wordlistCmd.Flags().IntP("count", "n", 0, "hint for number of entries to generate (0 = let AI decide)")
	rootCmd.AddCommand(wordlistCmd)
}

func runWordlist(cmd *cobra.Command, args []string) error {
	client, err := ai.NewClient()
	if err != nil {
		return fmt.Errorf("cfuzz wordlist requires ANTHROPIC_API_KEY: %w", err)
	}

	count, _ := cmd.Flags().GetInt("count")
	if err := ai.GenerateWordlist(client, args[0], count); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return nil
}
