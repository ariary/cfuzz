package fuzz

import (
	"strconv"
	"time"
)

// DisplayMode: interface use to determine field to display in cfuzz output
type DisplayMode interface {
	DisplayString(result ExecResult) string
}

// StdoutDisplay: display mode that displays field concerning stdout
type StdoutDisplay struct {
}

// DisplayString: return number of character in stdout
func (stdout StdoutDisplay) DisplayString(result ExecResult) string {
	nbChar := len(result.Stdout)
	return strconv.Itoa(nbChar)
}

// StderrDisplay: display mode that displays field concerning stderr
type StderrDisplay struct {
}

// DisplayString: return number of character in stderr
func (stderr StderrDisplay) DisplayString(result ExecResult) string {
	nbChar := len(result.Stderr)
	return strconv.Itoa(nbChar)
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

// CodeDisplay: display mode that displays field concerning exit code of command exectuion
type CodeDisplay struct {
}

// DisplayString: return exit code
func (code CodeDisplay) DisplayString(result ExecResult) string {
	exitCode := result.Code
	return exitCode
}
