package fuzz

// DisplayMode: interface use to determine field to display in cfuzz errput
type Filter interface {
	IsOk(result ExecResult) bool
}

// StdoutMaxFilter: Filter that accept only result with less characters than a specific number
type StdoutMaxFilter struct {
	Max int
}

// IsOk: return true if the lenght of stdout output is smaller or equal than max
func (maxFilter StdoutMaxFilter) IsOk(result ExecResult) bool {
	return len(result.Stdout) <= maxFilter.Max
}

// StdoutMaxFilter: Filter that accept only result with more characters than a specific number
type StdoutMinFilter struct {
	Min int
}

// IsOk: return true if the lenght of stdout output is greater or equal than min
func (filter StdoutMinFilter) IsOk(result ExecResult) bool {
	return len(result.Stdout) >= filter.Min
}

type StdoutEqFilter struct {
	Eq int
}

func (filter StdoutEqFilter) IsOk(result ExecResult) bool {
	return len(result.Stdout) == filter.Eq
}

// StderrMaxFilter: Filter that accept only result with less characters than a specific number
type StderrMaxFilter struct {
	Max int
}

// IsOk: return true if the lenght of stderr errput is smaller or equal than max
func (maxFilter StderrMaxFilter) IsOk(result ExecResult) bool {
	return len(result.Stderr) <= maxFilter.Max
}

// StderrMaxFilter: Filter that accept only result with more characters than a specific number
type StderrMinFilter struct {
	Min int
}

// IsOk: return true if the lenght of stderr errput is greater or equal than min
func (filter StderrMinFilter) IsOk(result ExecResult) bool {
	return len(result.Stderr) >= filter.Min
}

type StderrEqFilter struct {
	Eq int
}

func (filter StderrEqFilter) IsOk(result ExecResult) bool {
	return len(result.Stderr) == filter.Eq
}

type TimeMaxFilter struct {
	Max int
}

func (maxFilter TimeMaxFilter) IsOk(result ExecResult) bool {
	return int(result.Time.Seconds()) <= maxFilter.Max
}

type TimeMinFilter struct {
	Min int
}

func (filter TimeMinFilter) IsOk(result ExecResult) bool {
	return int(result.Time.Seconds()) >= filter.Min
}

type TimeEqFilter struct {
	Eq int
}

func (filter TimeEqFilter) IsOk(result ExecResult) bool {
	return int(result.Time.Seconds()) == filter.Eq
}

// CodeSuccessFilter: filter wether result regarding the exit code
type CodeSuccessFilter struct {
	Zero bool
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
