package config

import (
	"errors"
	"flag"
	"fmt"
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
		return "return"
	}
	return "unknown"
}

type Config struct {
	WordlistFilename string
	Keyword          string
	Command          string
	FilterType       Filter
	RoutineDelay     int64
	Shell            string
	Timeout          int64
}

var usage = `Usage of cfuzz: cfuzz [flags values] [command] or cfuzz [flags values] [command] with CFUZZ_CMD environment variable set
Fuzz command line execution and filter results
  -w, --wordlist     wordlist used by fuzzer
  -f, --filter	     filter type to sort execution results (among 'output','time','return')
  -d, --delay        delay in ms between each thread launching. A thread execute the command. (default: 0)
  -k, --keyword      keyword use to determine which zone to fuzz (default: FUZZ)
  -s, --shell        shell to use for execution (default: /bin/bash)
  -t, --timeout      command execution timeout in s. After reaching it the command is killed. (default: 30)
  -h, --help         prints help information 
`

// NewConfig create Config instance
func NewConfig() Config {
	// default value
	config := Config{Keyword: "FUZZ", FilterType: Output} //FilterType is already by default Output but to keep it in mind

	//flag wordlist
	flag.StringVar(&config.WordlistFilename, "wordlist", "", "wordlist used by fuzzer")
	flag.StringVar(&config.WordlistFilename, "w", "", "wordlist used by fuzzer")

	//flag keyword
	flag.StringVar(&config.Keyword, "keyword", "FUZZ", "keyword use to determine which zone to fuzz")
	flag.StringVar(&config.Keyword, "k", "FUZZ", "keyword use to determine which zone to fuzz")

	//flag shell
	flag.StringVar(&config.Shell, "shell", "/bin/bash", "shell to use for execution")
	flag.StringVar(&config.Shell, "s", "/bin/bash", "shell to use for execution")

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

	//flag RoutineDelay
	flag.Int64Var(&config.RoutineDelay, "d", 0, "delay in ms between each thread launching. A thread execute the command. (default: 0)")
	flag.Int64Var(&config.RoutineDelay, "delay", 0, "delay in ms between each thread launching. A thread execute the command. (default: 0)")

	//flag timeout
	flag.Int64Var(&config.Timeout, "t", 30, "Command execution timeout in s. After reaching it the command is killed. (default: 30)")
	flag.Int64Var(&config.Timeout, "timeout", 30, "Command execution timeout in s. After reaching it the command is killed. (default: 30)")

	flag.Usage = func() { fmt.Print(usage) }
	flag.Parse()

	// command
	if cmdEnv := os.Getenv("CFUZZ_CMD"); cmdEnv != "" {
		config.Command = cmdEnv
	} else if flag.NArg() > 0 {
		cmdArg := strings.Join(flag.Args(), " ")
		config.Command = cmdArg
	}

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
		return errors.New("No command provided. Please indicate it using environment variable CFUZZ_CMD or cfuzz [flag:value] [command]")
	}

	// check field consistency
	if !strings.Contains(c.Command, c.Keyword) {
		return errors.New("Fuzzing keyword has not been found in command. keyword:" + c.Keyword + " command:" + c.Command)
	}

	return nil
}
