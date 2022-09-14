package fuzz

import (
	"strings"
)

// Filter: interface used to determine which results will be displayed
type Filter interface {
	//IsOk: filter logic. Return true if the result pass the filter conditon
	IsOk(result ExecResult) bool
	//Name: return filter name
	Name() string
}

// STDOUT FILTERS

// StdoutMaxFilter: Filter that accept only result with less characters than a specific number
type StdoutMaxFilter struct {
	Max int
}

func (filter StdoutMaxFilter) Name() string {
	return "stdout characters max"
}

//IsOk: return true if the lenght of stdout output is smaller or equal than max
func (maxFilter StdoutMaxFilter) IsOk(result ExecResult) bool {
	return len(result.Stdout) <= maxFilter.Max
}

//StdoutMaxFilter: Filter that accept only result with more characters than a specific number
type StdoutMinFilter struct {
	Min int
}

func (filter StdoutMinFilter) Name() string {
	return "stdout characters min"
}

//IsOk: return true if the lenght of stdout output is greater or equal than min
func (filter StdoutMinFilter) IsOk(result ExecResult) bool {
	return len(result.Stdout) >= filter.Min
}

//StdoutEqFilter: filter struct that accept only result with exact amount of characters
type StdoutEqFilter struct {
	Eq int
}

func (filter StdoutEqFilter) Name() string {
	return "stdout characters equal"
}

//IsOK: return true if the number of characters is eqal
func (filter StdoutEqFilter) IsOk(result ExecResult) bool {
	return len(result.Stdout) == filter.Eq
}

//StdoutWordFilter: filter struct that accept only result containing specific word
type StdoutWordFilter struct {
	TargetWord string
}

//Name: return StdoutWordFilter name
func (filter StdoutWordFilter) Name() string {
	return "stdout word containing"
}

//IsOK: return true a specific word is found in result stdout. Note that this is equivalent to grep:
//the word could be surrounded by non-space characters
func (filter StdoutWordFilter) IsOk(result ExecResult) bool {
	return strings.Contains(result.Stdout, filter.TargetWord)
}

// STDERR FILTERS

//StderrMaxFilter: Filter that accept only result with less characters than a specific number
type StderrMaxFilter struct {
	Max int
}

func (filter StderrMaxFilter) Name() string {
	return "stderr characters max"
}

// IsOk: return true if the lenght of stderr errput is smaller or equal than max
func (maxFilter StderrMaxFilter) IsOk(result ExecResult) bool {
	return len(result.Stderr) <= maxFilter.Max
}

// StderrMaxFilter: Filter that accept only result with more characters than a specific number
type StderrMinFilter struct {
	Min int
}

func (filter StderrMinFilter) Name() string {
	return "stderr characters min"
}

// IsOk: return true if the lenght of stderr errput is greater or equal than min
func (filter StderrMinFilter) IsOk(result ExecResult) bool {
	return len(result.Stderr) >= filter.Min
}

type StderrEqFilter struct {
	Eq int
}

func (filter StderrEqFilter) Name() string {
	return "stderr characters equal"
}

func (filter StderrEqFilter) IsOk(result ExecResult) bool {
	return len(result.Stderr) == filter.Eq
}

//StderrWordFilter: filter struct that accept only result containing specific word for stderr
type StderrWordFilter struct {
	TargetWord string
}

//Name: return StdoutWordFilter name
func (filter StderrWordFilter) Name() string {
	return "stderr word containing"
}

//IsOK: return true a specific word is found in result stdout. Note that this is equivalent to grep:
//the word could be surrounded by non-space characters
func (filter StderrWordFilter) IsOk(result ExecResult) bool {
	return strings.Contains(result.Stderr, filter.TargetWord)
}

// TIME FILTERS

type TimeMaxFilter struct {
	Max int
}

func (maxFilter TimeMaxFilter) IsOk(result ExecResult) bool {
	return int(result.Time.Seconds()) <= maxFilter.Max
}

func (filter TimeMaxFilter) Name() string {
	return "time max"
}

type TimeMinFilter struct {
	Min int
}

func (filter TimeMinFilter) Name() string {
	return "time min"
}

func (filter TimeMinFilter) IsOk(result ExecResult) bool {
	return int(result.Time.Seconds()) >= filter.Min
}

type TimeEqFilter struct {
	Eq int
}

func (filter TimeEqFilter) Name() string {
	return "time equal"
}

func (filter TimeEqFilter) IsOk(result ExecResult) bool {
	return int(result.Time.Seconds()) == filter.Eq
}

// CODE FILTERS

// CodeSuccessFilter: filter wether result regarding the exit code
type CodeSuccessFilter struct {
	Zero bool
}

func (filter CodeSuccessFilter) Name() string {
	if filter.Zero {
		return "on success"
	} else {
		return "non-zero exit code"
	}
}

// IsOk: return true if the exit code is 0 and the filter is parametrize with Zero set at true. return true if the exit code is != 0 and the filter is parametrize with Zero set at false.
// return false otherwise
func (filter CodeSuccessFilter) IsOk(result ExecResult) bool {
	if filter.Zero { // --success
		return result.Code == "0"
	} else {
		return result.Code != "0"
	}
}
