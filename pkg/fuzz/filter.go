package fuzz

// OutputFilter: output filter
type OutputFilter int64

const (
	CharacterExact OutputFilter = iota
	CharacterMax
	CharacterMin
)

// StderrFilter: output filter
type StderrFilter int64

const (
	StderrCharacterExact OutputFilter = iota
	StderrCharacterMax
	StderrCharacterMin
)

// TimeFilter: type of field to display to cfuzz result
type TimeFilter int64

const (
	TimeExact OutputFilter = iota
	TimeMax
	TimeMin
)

// CodeFilter: type of field to display to cfuzz result
type CodeFilter int64

const (
	CodeExact = iota
	CodeFalure
)
