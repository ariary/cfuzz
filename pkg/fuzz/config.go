package fuzz

import (
	"errors"
	"log"
	"os"
	"strconv"
	"strings"
)

// Wordlists is a slice of wordlist file paths with a String() method for formatting.
type Wordlists []string

// String returns a comma-separated string representation of the wordlists.
func (w Wordlists) String() string {
	return strings.Join(w, ",")
}

// Config holds all runtime configuration for a cfuzz run.
type Config struct {
	Wordlists     Wordlists
	Keyword       string
	Command       string
	RoutineDelay  int64
	Shell         string
	Timeout       int64
	Input         string
	StdinFuzzing  bool
	Multiple      bool
	StdinWordlist bool
	Threads       int
	AIFilter      string
	OnlyWord      bool
	// AIFilterFn is an optional callback that returns true if the result should be shown.
	// Used to wire in AI filtering without importing pkg/ai from pkg/fuzz.
	// When nil, no AI filtering is applied.
	AIFilterFn    func(word, stdout, stderr, code string) bool
	DisplayModes  []DisplayMode
	FullDisplay   bool
	HideBanner    bool
	Hide          bool
	Filters       []Filter
	ResultLogger  *log.Logger
}

// DefaultConfig returns a Config with sensible defaults and a stdout logger.
func DefaultConfig() Config {
	return Config{
		Keyword:      "FUZZ",
		Shell:        "/bin/bash",
		Timeout:      30,
		Threads:      50,
		ResultLogger: log.New(os.Stdout, "", 0),
	}
}

// CheckConfig validates that all required fields are present and consistent.
func (c *Config) CheckConfig() error {
	if len(c.Wordlists) == 0 && !c.StdinWordlist {
		return errors.New("no wordlist provided: use -w/--wordlist or --stdin-wordlist")
	}
	if len(c.Wordlists) != 0 && c.StdinWordlist {
		return errors.New("-w/--wordlist cannot be combined with --stdin-wordlist")
	}
	if c.Keyword == "" {
		return errors.New("fuzzing keyword cannot be empty")
	}
	if c.Command == "" {
		return errors.New("no command provided: set CFUZZ_CMD or pass command as argument")
	}
	if c.Multiple && c.StdinWordlist {
		return errors.New("--spider cannot be combined with --stdin-wordlist")
	}
	if c.Multiple && len(c.Wordlists) < 2 {
		return errors.New("--spider requires at least 2 wordlists")
	}
	if !c.Multiple && len(c.Wordlists) > 1 {
		return errors.New("multiple wordlists provided without --spider")
	}
	if c.FullDisplay && len(c.DisplayModes) > 0 {
		return errors.New("--full-output cannot be combined with other display modes: " + c.DisplayModes[0].Name())
	}
	return checkKeywordsPresence(c)
}

func checkKeywordsPresence(c *Config) error {
	if c.StdinFuzzing {
		if c.Multiple {
			n := strings.Count(c.Input+c.Command, c.Keyword)
			if n != len(c.Wordlists) {
				return errors.New("keyword count (" + strconv.Itoa(n) + ") must match wordlist count (" + strconv.Itoa(len(c.Wordlists)) + ")")
			}
		} else if !strings.Contains(c.Input, c.Keyword) {
			return errors.New("keyword " + c.Keyword + " not found in stdin input: " + c.Input)
		}
	} else if c.Multiple {
		n := strings.Count(c.Command, c.Keyword)
		if n != len(c.Wordlists) {
			return errors.New("keyword count (" + strconv.Itoa(n) + ") must match wordlist count (" + strconv.Itoa(len(c.Wordlists)) + ")")
		}
	} else if !strings.Contains(c.Command, c.Keyword) {
		return errors.New("keyword " + c.Keyword + " not found in command: " + c.Command)
	}
	return nil
}

// BuildDisplayModes returns the display mode slice from parsed flag values.
// When fullDisplay is true, returns nil (PrintExec handles full display separately).
// When no flags are set, defaults to stdout character count.
func BuildDisplayModes(stdout, stderr, showTime, code, fullDisplay bool) []DisplayMode {
	if fullDisplay {
		return nil
	}
	var modes []DisplayMode
	if stdout {
		modes = append(modes, StdoutDisplay{})
	}
	if stderr {
		modes = append(modes, StderrDisplay{})
	}
	if showTime {
		modes = append(modes, TimeDisplay{})
	}
	if code {
		modes = append(modes, CodeDisplay{})
	}
	if len(modes) == 0 {
		modes = []DisplayMode{StdoutDisplay{}}
	}
	return modes
}

// BuildFilters returns the filter slice from parsed flag values.
// Use -1 as sentinel for "not set" on all int parameters.
func BuildFilters(
	stdoutMin, stdoutMax, stdoutEq int,
	stdoutWords []string,
	stderrMin, stderrMax, stderrEq int,
	stderrWords []string,
	timeMin, timeMax, timeEq int,
	success, failure bool,
) []Filter {
	var filters []Filter

	if stdoutMin >= 0 {
		filters = append(filters, StdoutMinFilter{Min: stdoutMin})
	}
	if stdoutMax >= 0 {
		filters = append(filters, StdoutMaxFilter{Max: stdoutMax})
	}
	if stdoutEq >= 0 {
		filters = append(filters, StdoutEqFilter{Eq: stdoutEq})
	}
	for _, w := range stdoutWords {
		filters = append(filters, StdoutWordFilter{TargetWord: w})
	}

	if stderrMin >= 0 {
		filters = append(filters, StderrMinFilter{Min: stderrMin})
	}
	if stderrMax >= 0 {
		filters = append(filters, StderrMaxFilter{Max: stderrMax})
	}
	if stderrEq >= 0 {
		filters = append(filters, StderrEqFilter{Eq: stderrEq})
	}
	for _, w := range stderrWords {
		filters = append(filters, StderrWordFilter{TargetWord: w})
	}

	if timeMin >= 0 {
		filters = append(filters, TimeMinFilter{Min: timeMin})
	}
	if timeMax >= 0 {
		filters = append(filters, TimeMaxFilter{Max: timeMax})
	}
	if timeEq >= 0 {
		filters = append(filters, TimeEqFilter{Eq: timeEq})
	}

	if success {
		filters = append(filters, CodeSuccessFilter{Zero: true})
	}
	if failure {
		filters = append(filters, CodeSuccessFilter{Zero: false})
	}

	return filters
}
