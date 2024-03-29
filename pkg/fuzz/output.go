package fuzz

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/ariary/go-utils/pkg/color"
)

var nameDraw = color.Teal(`
                 _/_/                                
    _/_/_/    _/      _/    _/  _/_/_/_/  _/_/_/_/   
 _/        _/_/_/_/  _/    _/      _/        _/      
_/          _/      _/    _/    _/        _/         
 _/_/_/    _/        _/_/_/  _/_/_/_/  _/_/_/_/			

 `)
var author = color.Yellow("By @ariary (" + color.Underlined("https://github.com/ariary") + ")")

// Banner: Print the banner as it is trendy for this kind of tool. thanks to: https://patorjk.com/software/taag
func Banner() {
	banner := nameDraw + author
	fmt.Println(banner)
	fmt.Println()
}

// PrintConfig: print configuration of cfuzz running
func PrintConfig(cfg Config) {
	// filters
	allFilters := ""
	for i := 0; i < len(cfg.Filters); i++ {
		allFilters += cfg.Filters[i].Name() + ", "
	}

	allDisplayModes := ""
	for i := 0; i < len(cfg.DisplayModes); i++ {
		allDisplayModes += cfg.DisplayModes[i].Name() + ", "
	}

	line := `[*] ----------------------~~~~~~~~~~~~~~~~~~~---------------------- [*]`
	fmt.Println(line)
	fmt.Println()
	PrintLine(cfg, "command fuzzed:", cfg.Command)
	if len(cfg.Wordlists) != 0 {
		PrintLine(cfg, "wordlist:", cfg.Wordlists.String())
	} else if cfg.StdinWordlist {
		PrintLine(cfg, "wordlist:", "from stdin")
	}

	if allDisplayModes != "" {
		allDisplayModes = allDisplayModes[:len(allDisplayModes)-2] //delete last comma
		PrintLine(cfg, "columns:", allDisplayModes)
	}
	if allFilters != "" {
		allFilters = allFilters[:len(allFilters)-2] //delete last comma
		PrintLine(cfg, "filters:", allFilters)
	}
	if cfg.Hide {
		fmt.Println("Only displays words that do not pass the filter")
	}
	fmt.Println()
	fmt.Println(line)
	fmt.Println()
}

//PrintLine: Nice printing of a line containing 2 or more elements
func PrintLine(cfg Config, value string, element ...string) {
	// string builder and tabwriter
	var strBuilder strings.Builder
	tabwriter := new(tabwriter.Writer)
	// // minwidth, tabwidth, padding, padchar, flags
	tabwriter.Init(&strBuilder, 40, 8, 0, '\t', 0)

	line := value
	for i := 0; i < len(element); i++ {
		line += "\t" + element[i]
	}
	fmt.Fprintf(tabwriter, "%s", line) //write into tab -> write into string builder

	tabwriter.Flush() // Flush before calling String()
	cfg.ResultLogger.Println(strBuilder.String())

}

//PrintFullExecOutput Nice printing of command execution
func PrintFullExecOutput(cfg Config, result ExecResult) {
	//TODO
	cfg.ResultLogger.Println()
	cfg.ResultLogger.Println(color.Bold(color.Cyan(result.Substitute)))
	if result.Stdout != "" {
		cfg.ResultLogger.Println(color.GreenForeground("STDOUT:"))
		cfg.ResultLogger.Println(result.Stdout)
	}
	if result.Stderr != "" {
		cfg.ResultLogger.Println(color.RedForeground("STDERR:"))
		cfg.ResultLogger.Println(result.Stderr)
	}

	cfg.ResultLogger.Println()
}
