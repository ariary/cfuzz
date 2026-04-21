package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/ariary/cfuzz/pkg/ai"
	"github.com/ariary/cfuzz/pkg/fuzz"
	"github.com/spf13/cobra"
)

func init() {
	f := rootCmd.Flags()

	// Configuration
	f.StringArrayP("wordlist", "w", nil, "wordlist file(s) for fuzzing (repeatable with --spider)")
	f.StringP("keyword", "k", "FUZZ", "keyword to replace in command")
	f.StringP("shell", "s", "/bin/bash", "shell to use for execution")
	f.Int64P("delay", "d", 0, "delay in ms between goroutine launches")
	f.Int64("timeout", 30, "command execution timeout in seconds")
	f.StringP("input", "i", "", "provide command stdin")
	f.Bool("stdin-fuzzing", false, "fuzz stdin instead of command line")
	f.BoolP("spider", "m", false, "fuzz multiple keyword positions (requires multiple -w)")
	f.Bool("stdin-wordlist", false, "read wordlist from cfuzz stdin")
	f.IntP("threads", "j", 50, "max concurrent workers")

	// Display
	f.Bool("no-banner", false, "hide banner")
	f.BoolP("only-word", "r", false, "print only matched words (no metadata columns)")
	f.BoolP("hide", "H", false, "show results that do NOT pass the filters")
	f.Bool("stdout-chars", false, "display stdout character count")
	f.Bool("stderr-chars", false, "display stderr character count")
	f.BoolP("time", "t", false, "display execution time")
	f.BoolP("code", "c", false, "display exit code")
	f.BoolP("full-output", "f", false, "display full command output")

	// Filters — stdout (-1 = not set)
	f.Int("stdout-min", -1, "show only if stdout chars >= n")
	f.Int("stdout-max", -1, "show only if stdout chars <= n")
	f.Int("stdout-eq", -1, "show only if stdout chars == n")
	f.StringArray("stdout-word", nil, "show only if stdout contains word (repeatable)")

	// Filters — stderr (-1 = not set)
	f.Int("stderr-min", -1, "show only if stderr chars >= n")
	f.Int("stderr-max", -1, "show only if stderr chars <= n")
	f.Int("stderr-eq", -1, "show only if stderr chars == n")
	f.StringArray("stderr-word", nil, "show only if stderr contains word (repeatable)")

	// Filters — time (-1 = not set)
	f.Int("time-min", -1, "show only if execution time >= n seconds")
	f.Int("time-max", -1, "show only if execution time <= n seconds")
	f.Int("time-eq", -1, "show only if execution time == n seconds")

	// Filters — exit code
	f.Bool("success", false, "show only commands with exit code 0")
	f.Bool("failure", false, "show only commands with non-zero exit code")

	// AI
	f.String("ai-filter", "", "AI natural language filter description (requires ANTHROPIC_API_KEY)")

	rootCmd.RunE = runFuzz
}

func runFuzz(cmd *cobra.Command, args []string) error {
	f := cmd.Flags()
	cfg := fuzz.DefaultConfig()

	cfg.Wordlists, _ = f.GetStringArray("wordlist")
	cfg.Keyword, _ = f.GetString("keyword")
	cfg.Shell, _ = f.GetString("shell")
	cfg.RoutineDelay, _ = f.GetInt64("delay")
	cfg.Timeout, _ = f.GetInt64("timeout")
	cfg.Input, _ = f.GetString("input")
	cfg.StdinFuzzing, _ = f.GetBool("stdin-fuzzing")
	cfg.Multiple, _ = f.GetBool("spider")
	cfg.StdinWordlist, _ = f.GetBool("stdin-wordlist")
	cfg.Threads, _ = f.GetInt("threads")
	cfg.HideBanner, _ = f.GetBool("no-banner")
	cfg.OnlyWord, _ = f.GetBool("only-word")
	cfg.Hide, _ = f.GetBool("hide")
	cfg.FullDisplay, _ = f.GetBool("full-output")
	cfg.AIFilter, _ = f.GetString("ai-filter")

	// Command from env or positional args
	if cmdEnv := os.Getenv("CFUZZ_CMD"); cmdEnv != "" {
		cfg.Command = cmdEnv
	} else if len(args) > 0 {
		cfg.Command = strings.Join(args, " ")
	}

	// Display modes (skipped when --only-word)
	if !cfg.OnlyWord {
		stdoutChars, _ := f.GetBool("stdout-chars")
		stderrChars, _ := f.GetBool("stderr-chars")
		showTime, _ := f.GetBool("time")
		showCode, _ := f.GetBool("code")
		cfg.DisplayModes = fuzz.BuildDisplayModes(stdoutChars, stderrChars, showTime, showCode, cfg.FullDisplay)
	}

	// Filters
	stdoutMin, _ := f.GetInt("stdout-min")
	stdoutMax, _ := f.GetInt("stdout-max")
	stdoutEq, _ := f.GetInt("stdout-eq")
	stdoutWords, _ := f.GetStringArray("stdout-word")
	stderrMin, _ := f.GetInt("stderr-min")
	stderrMax, _ := f.GetInt("stderr-max")
	stderrEq, _ := f.GetInt("stderr-eq")
	stderrWords, _ := f.GetStringArray("stderr-word")
	timeMin, _ := f.GetInt("time-min")
	timeMax, _ := f.GetInt("time-max")
	timeEq, _ := f.GetInt("time-eq")
	success, _ := f.GetBool("success")
	failure, _ := f.GetBool("failure")

	cfg.Filters = fuzz.BuildFilters(
		stdoutMin, stdoutMax, stdoutEq, stdoutWords,
		stderrMin, stderrMax, stderrEq, stderrWords,
		timeMin, timeMax, timeEq,
		success, failure,
	)

	if cfg.AIFilter != "" {
		aiClient, err := ai.NewClient()
		if err != nil {
			return fmt.Errorf("--ai-filter requires ANTHROPIC_API_KEY: %w", err)
		}
		criterion := cfg.AIFilter
		cfg.AIFilterFn = func(word, stdout, stderr, code string) bool {
			return ai.IsInteresting(aiClient, criterion, word, stdout, stderr, code)
		}
	}

	if !cfg.HideBanner {
		fuzz.Banner()
		fuzz.PrintConfig(cfg)
	}

	if err := cfg.CheckConfig(); err != nil {
		return err
	}

	fuzz.PerformFuzzing(cfg)
	return nil
}
