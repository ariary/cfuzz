package fuzz

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type wordlists []string

type Config struct {
	WordlistFilename string
	Wordlists        wordlists
	Keyword          string
	Command          string
	RoutineDelay     int64
	Shell            string
	Timeout          int64
	Input            string
	StdinFuzzing     bool
	Multiple         bool
	DisplayModes     []DisplayMode
	Hide             bool
	Filters          []Filter
}

var usage = `Usage of cfuzz: cfuzz [flags values] [command] or cfuzz [flags values] [command] with CFUZZ_CMD environment variable set
Fuzz command line execution and filter results

CONFIGURATION
  -w, --wordlist              wordlist used by fuzzer
  -d, --delay                 delay in ms between each thread launching. A thread executes the command. (default: 0)
  -k, --keyword               keyword used to determine which zone to fuzz (default: FUZZ)
  -s, --shell                 shell to use for execution (default: /bin/bash)
  -to, --timeout              command execution timeout in s. After reaching it the command is killed. (default: 30)
  -i, --input                 provide stdin
  -if, --stdin-fuzzing        fuzz sdtin instead of command line
  -m, --spider                fuzz multiple keyword places. You must provide as many wordlists as keywords. Provide them in order you wan them to be applied

DISPLAY
  -oc, --stdout               display stdout number of characters
  -ec, --stderr               display stderr number of characters
  -t, --time                  display execution time
  -c, --code                  display exit code

FILTER

  -H, --hide                  display only if result does not pass the filters

 STDOUT:
  -omin, --stdout-min         filter to only display if stdout characters number is lesser than n
  -omax, --stdout-max         filter to only display if stdout characters number is greater than n
  -oeq,  --stdout-equal       filter to only display if stdout characters number is equal to n
  -ow,   --stdout-word        filter to only display if stdout cointains specific word

 STDERR:
  -emin, --stderr-min         filter to only display if stderr characters number is lesser than n
  -emax, --stderr-max         filter to only display if stderr characters number is greater than n
  -eeq,  --stderr-equal       filter to only display if stderr characters number is equal to n

 TIME:
  -tmin, --time-min           filter to only display if exectuion time is shorter than n seconds
  -tmax, --time-max           filter to only display if exectuion time is longer than n seconds
  -teq,  --time-equal         filter to only display if exectuion time is shorter than n seconds

 CODE:
  --success                   filter to only display if execution return a zero exit code
  --failure                   filter to only display if execution return a non-zero exit code

  -h, --help                  prints help information 
`

func (i *wordlists) String() string {

	return strings.Join(*i, ",")
}

func (i *wordlists) Set(value string) error {
	*i = append(*i, value)
	return nil
}

// NewConfig create Config instance
func NewConfig() Config {
	// default value
	config := Config{Keyword: "FUZZ"}

	// flag wordlist
	// flag.StringVar(&config.WordlistFilename, "wordlist", "", "wordlist used by fuzzer")
	// flag.StringVar(&config.WordlistFilename, "w", "", "wordlist used by fuzzer")
	flag.Var(&config.Wordlists, "worldlist", "Some description for this param.")

	// flag keyword
	flag.StringVar(&config.Keyword, "keyword", "FUZZ", "keyword use to determine which zone to fuzz")
	flag.StringVar(&config.Keyword, "k", "FUZZ", "keyword use to determine which zone to fuzz")

	// flag shell
	flag.StringVar(&config.Shell, "shell", "/bin/bash", "shell to use for execution")
	flag.StringVar(&config.Shell, "s", "/bin/bash", "shell to use for execution")

	// flag RoutineDelay
	flag.Int64Var(&config.RoutineDelay, "d", 0, "delay in ms between each thread launching. A thread execute the command. (default: 0)")
	flag.Int64Var(&config.RoutineDelay, "delay", 0, "delay in ms between each thread launching. A thread execute the command. (default: 0)")

	//flag timeout
	flag.Int64Var(&config.Timeout, "to", 30, "Command execution timeout in s. After reaching it the command is killed. (default: 30)")
	flag.Int64Var(&config.Timeout, "timeout", 30, "Command execution timeout in s. After reaching it the command is killed. (default: 30)")

	// flag input
	flag.StringVar(&config.Input, "input", "", "fuzz stdin")
	flag.StringVar(&config.Input, "i", "", "fuzz stdin")

	// flag stdin-fuzzing
	flag.BoolVar(&config.StdinFuzzing, "stdin-fuzzing", false, "fuzz stdin")
	flag.BoolVar(&config.StdinFuzzing, "if", false, "fuzz stdin")

	// flag spider
	flag.BoolVar(&config.Multiple, "spider", false, "fuzz multiple keyword")
	flag.BoolVar(&config.Multiple, "m", false, "fuzz multiple keyword")

	// flag hide
	flag.BoolVar(&config.Hide, "H", false, "hide fields that pass the filter")
	flag.BoolVar(&config.Hide, "hide", false, "hide fields that pass the filter")

	// display mode
	var stdoutDisplay bool
	flag.BoolVar(&stdoutDisplay, "oc", false, "display command execution  number of characters in stdout.")
	flag.BoolVar(&stdoutDisplay, "stdout-characters", false, "display execution command number of characters in stdout.")

	var stderrDisplay bool
	flag.BoolVar(&stderrDisplay, "ec", false, "display command execution  number of characters in stderr.")
	flag.BoolVar(&stderrDisplay, "stderr-characters", false, "display execution command number of characters in stderr.")

	var timeDisplay bool
	flag.BoolVar(&timeDisplay, "t", false, "display command execution  time.")
	flag.BoolVar(&timeDisplay, "time", false, "display command execution time.")

	var codeDisplay bool
	flag.BoolVar(&codeDisplay, "c", false, "display command execution exit code.")
	flag.BoolVar(&codeDisplay, "code", false, "display command execution exit code.")

	// filters
	parseFilters(&config)

	flag.Usage = func() { fmt.Print(usage) }
	flag.Parse()

	// command
	if cmdEnv := os.Getenv("CFUZZ_CMD"); cmdEnv != "" {
		config.Command = cmdEnv
	} else if flag.NArg() > 0 {
		cmdArg := strings.Join(flag.Args(), " ")
		config.Command = cmdArg
	}

	// parse display mode
	config.DisplayModes = parseDisplayMode(stdoutDisplay, stderrDisplay, timeDisplay, codeDisplay)
	return config
}

//CheckConfig: assert that all required fields are present in config, and are adequate to cfuzz run
func (c *Config) CheckConfig() error {
	// check field
	// if c.WordlistFilename == "" {
	// 	return errors.New("No wordlist provided. Please indicate a wordlist to use for fuzzing (-w,--wordlist)")
	// }
	if len(c.Wordlists) == 0 {
		return errors.New("No wordlist provided. Please indicate a wordlist to use for fuzzing (-w,--wordlist)")
	}

	if c.Keyword == "" {
		return errors.New("Fuzzing Keyword can't be empty string")
	}
	if c.Command == "" {
		return errors.New("No command provided. Please indicate it using environment variable CFUZZ_CMD or cfuzz [flag:value] [command]")
	}

	// check field consistency
	err := checkKeywordsPresence(c)

	return err
}

//checkKeywordsPresence: check the consistency between flag and keyword presence (ie Keyword is present in stdin or command and if --spider check
//there are as many keyword than wordlist)
func checkKeywordsPresence(c *Config) error {
	if c.StdinFuzzing {
		if c.Multiple { //stdin + multiple
			keywordNum := strings.Count(c.Input+c.Command, c.Keyword)
			if keywordNum != len(c.Wordlists) {
				return errors.New("Please provide as many wordlists as keyword. keyword:" + c.Keyword + " input:" + c.Input + "  command:" + c.Command + "wordlist number:" + strconv.Itoa(len(c.Wordlists)))
			}
		} else if !strings.Contains(c.Input, c.Keyword) { //stdin simple
			return errors.New("Fuzzing keyword has not been found in stdin. keyword:" + c.Keyword + " input:" + c.Input)
		} else {
			return nil
		}
	} else if c.Multiple { // multiple w/o stdin
		keywordNum := strings.Count(c.Command, c.Keyword)
		if keywordNum != len(c.Wordlists) {
			return errors.New("Please provide as many wordlists as keyword. keyword:" + c.Keyword + "  command:" + c.Command + "wordlist number:" + strconv.Itoa(len(c.Wordlists)))
		}
	} else if !strings.Contains(c.Command, c.Keyword) { //simple w/o stdin
		return errors.New("Fuzzing keyword has not been found in command. keyword:" + c.Keyword + " command:" + c.Command)
	}
	return nil
}

//parseDisplayMode: Return array of display mode interface chosen with flags. If none, default is stdout characters display mode
func parseDisplayMode(stdout bool, stderr bool, time bool, code bool) (modes []DisplayMode) {
	if stdout {
		modes = append(modes, StdoutDisplay{})
	}
	if stderr {
		modes = append(modes, StderrDisplay{})
	}
	if time {
		modes = append(modes, TimeDisplay{})
	}
	if code {
		modes = append(modes, CodeDisplay{})
	}

	//default, if none
	if len(modes) == 0 {
		stdoutDisplay := StdoutDisplay{}
		modes = []DisplayMode{stdoutDisplay}
	}
	return modes
}

//parseFilters: parse all flags and determine the filters, add them in the config struct given in parameter
func parseFilters(config *Config) {
	// stdout filters
	flag.Func("omax", "filter to display only results with less than n characters", func(max string) error {
		n, err := strconv.Atoi(max)
		if err != nil {
			return err
		}
		filter := StdoutMaxFilter{Max: n}
		config.Filters = append(config.Filters, filter)
		return nil
	})

	flag.Func("omin", "filter to display only results with more than n characters", func(min string) error {
		n, err := strconv.Atoi(min)
		if err != nil {
			return err
		}
		filter := StdoutMinFilter{Min: n}
		config.Filters = append(config.Filters, filter)
		return nil
	})

	flag.Func("oeq", "filter to display only results with exactly n characters", func(eq string) error {
		n, err := strconv.Atoi(eq)
		if err != nil {
			return err
		}
		filter := StdoutEqFilter{Eq: n}
		config.Filters = append(config.Filters, filter)
		return nil
	})

	flag.Func("ow", "filter to display only results cointaing specific in stdout", func(word string) error {
		filter := StdoutWordFilter{TargetWord: word}
		config.Filters = append(config.Filters, filter)
		return nil
	})

	// stderr filters
	flag.Func("emax", "filter to display only results with less than n characters", func(max string) error {
		n, err := strconv.Atoi(max)
		if err != nil {
			return err
		}
		filter := StderrMaxFilter{Max: n}
		config.Filters = append(config.Filters, filter)
		return nil
	})

	flag.Func("emin", "filter to display only results with more than n characters", func(min string) error {
		n, err := strconv.Atoi(min)
		if err != nil {
			return err
		}
		filter := StderrMinFilter{Min: n}
		config.Filters = append(config.Filters, filter)
		return nil
	})

	flag.Func("eeq", "filter to display only results with exactly n characters", func(eq string) error {
		n, err := strconv.Atoi(eq)
		if err != nil {
			return err
		}
		filter := StderrEqFilter{Eq: n}
		config.Filters = append(config.Filters, filter)
		return nil
	})

	// time filters
	flag.Func("tmax", "filter to display only results with a time lesser than n seconds", func(max string) error {
		n, err := strconv.Atoi(max)
		if err != nil {
			return err
		}
		filter := TimeMaxFilter{Max: n}
		config.Filters = append(config.Filters, filter)
		return nil
	})

	flag.Func("tmin", "filter to display only results with a time greater than n seconds", func(min string) error {
		n, err := strconv.Atoi(min)
		if err != nil {
			return err
		}
		filter := TimeMinFilter{Min: n}
		config.Filters = append(config.Filters, filter)
		return nil
	})

	flag.Func("teq", "filter to  display only results with a time equal to n seconds", func(eq string) error {
		n, err := strconv.Atoi(eq)
		if err != nil {
			return err
		}
		filter := TimeEqFilter{Eq: n}
		config.Filters = append(config.Filters, filter)
		return nil
	})

	// code filters
	var success, failure bool
	flag.BoolVar(&success, "success", false, "filter to display only command with exit code 0.")
	flag.BoolVar(&failure, "failure", false, "filter to display only command with a non-zero exit .")
	if success {
		filter := CodeSuccessFilter{Zero: true}
		config.Filters = append(config.Filters, filter)
	}
	if failure {
		filter := CodeSuccessFilter{Zero: false}
		config.Filters = append(config.Filters, filter)
	}
}
