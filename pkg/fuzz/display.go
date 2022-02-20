package fuzz

import (
	"strconv"
	"time"
)

// DisplayMode: interface used to determine field to display in cfuzz output
type DisplayMode interface {
	DisplayString(result ExecResult) string
	Name() string
}

// StdoutDisplay: display mode that displays field concerning stdout
type StdoutDisplay struct {
}

// DisplayString: return number of character in stdout
func (stdout StdoutDisplay) DisplayString(result ExecResult) string {
	nbChar := len(result.Stdout)
	return strconv.Itoa(nbChar)
}

func (stdout StdoutDisplay) Name() string {
	return "stdout characters number"
}

// StderrDisplay: display mode that displays field concerning stderr
type StderrDisplay struct {
}

// DisplayString: return number of character in stderr
func (stderr StderrDisplay) DisplayString(result ExecResult) string {
	nbChar := len(result.Stderr)
	return strconv.Itoa(nbChar)
}

func (stderr StderrDisplay) Name() string {
	return "stderr characters number"
}

// TimeDisplay: display mode that displays field concerning stderr
type TimeDisplay struct {
	timeout      time.Duration
	reachTimeout bool
}

// DisplayString: return command execution time in second
func (time TimeDisplay) DisplayString(result ExecResult) string {
	var timeExecution string
	if time.reachTimeout {
		timeExecution = time.timeout.String()
	} else {
		timeExecution = result.Time.String()
	}

	return timeExecution
}

func (time TimeDisplay) Name() string {
	return "execution time"
}

// CodeDisplay: display mode that displays field concerning exit code of command exectuion
type CodeDisplay struct {
}

// DisplayString: return exit code
func (code CodeDisplay) DisplayString(result ExecResult) string {
	exitCode := result.Code
	return exitCode
}

func (code CodeDisplay) Name() string {
	return "exit code"
}
