package config

import (
	"errors"
	"flag"
	"os"
	"strings"
)

type Config struct {
	WordlistFilename string
	Keyword          string
	Command          string
}

// NewConfig create Config instance
func NewConfig() Config {
	// default value
	config := Config{Keyword: "FUZZ"}

	//flag value
	flag.StringVar(&config.WordlistFilename, "wordlist", "", "wordlist used by fuzzer")
	flag.StringVar(&config.WordlistFilename, "w", "", "wordlist used by fuzzer")

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
