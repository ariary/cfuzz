package config

import (
	"errors"
	"flag"
	"os"
	"strings"
)

// Filter: type of filter apply to cfuzz result
type Filter int64

const (
	Output Filter = iota //Output is the default hence
	Time
	ReturnCode
)

func (f Filter) String() string {
	switch f {
	case Output:
		return "output"
	case Time:
		return "time"
	case ReturnCode:
		return "returnCode"
	}
	return "unknown"
}

type Config struct {
	WordlistFilename string
	Keyword          string
	Command          string
	FilterType       Filter
}

// NewConfig create Config instance
func NewConfig() Config {
	// default value
	config := Config{Keyword: "FUZZ", FilterType: Output} //FilterType is already by default Output but to keep it in mind

	//flag wordlist
	flag.StringVar(&config.WordlistFilename, "wordlist", "", "wordlist used by fuzzer")
	flag.StringVar(&config.WordlistFilename, "w", "", "wordlist used by fuzzer")

	//flag filter
	var filter string
	flag.StringVar(&filter, "f", "output", "filter type to sort execution results")
	flag.StringVar(&filter, "filter", "output", "filter type to sort execution results")
	switch filter {
	case Output.String():
		config.FilterType = Output
	case Time.String():
		config.FilterType = Time
	case ReturnCode.String():
		config.FilterType = ReturnCode
	} //default: if unreadable keep output
	flag.Parse()

	config.Command = os.Getenv("CFUZZ_CMD")

	return config
}

// CheckConfig: assert that all required fields are present in config, and are adequate to cfuzz run
func (c *Config) CheckConfig() error {
	// check field
	if c.WordlistFilename == "" {
		return errors.New("No wordlist provided. Please indicate a wordlist to use for fuzzing (-w,--wordlist)")
	}

	if c.Keyword == "" {
		return errors.New("Fuzzing Keyword can't be empty string")
	}
	if c.Command == "" {
		return errors.New("No command provided. Please indicate it using environment variable CFUZZ_CMD")
	}

	// check field consistency
	if !strings.Contains(c.Command, c.Keyword) {
		return errors.New("Fuzzing keyword has not been found in command. keyword:" + c.Keyword + " command:" + c.Command)
	}

	return nil
}
