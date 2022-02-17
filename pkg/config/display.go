package config

// Mode: type of field to display to cfuzz result
type Mode int64

const (
	Stdout Mode = iota //Output is the default hence
	Stderr
	Time
	ReturnCode
)

func (m Mode) String() string {
	switch m {
	case Stdout:
		return "stdout"
	case Stderr:
		return "stderr"
	case Time:
		return "time"
	case ReturnCode:
		return "code"
	}
	return "unknown"
}
