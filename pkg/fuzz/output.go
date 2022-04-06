package fuzz

import (
	"fmt"
	"os"
	"text/tabwriter"
)

// Banner: Print the banner as it is trendy for this kind of tool. thanks to: https://patorjk.com/software/taag
func Banner() {
	banner := `
                 _/_/                                
    _/_/_/    _/      _/    _/  _/_/_/_/  _/_/_/_/   
 _/        _/_/_/_/  _/    _/      _/        _/      
_/          _/      _/    _/    _/        _/         
 _/_/_/    _/        _/_/_/  _/_/_/_/  _/_/_/_/									
`

	fmt.Println(banner)
	fmt.Println()
}

// PrintConfig: print configuration of cfuzz running
func PrintConfig(cfg Config) {
	// filters
	allFilters := ""
	fmt.Println("tototo", cfg.Filters)
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
	PrintLine("command fuzzed:", cfg.Command)
	PrintLine("wordlist:", cfg.Wordlists.String())
	if allDisplayModes != "" {
		allDisplayModes = allDisplayModes[:len(allDisplayModes)-2] //delete last comma
		PrintLine("columns:", allDisplayModes)
	}
	if allFilters != "" {
		allFilters = allFilters[:len(allFilters)-2] //delete last comma
		PrintLine("filters:", allFilters)
	}
	if cfg.Hide {
		fmt.Println("Only displays filter that do not pass the filter")
	}
	fmt.Println()
	fmt.Println(line)
	fmt.Println()
}

// Nice printing of a line containing 2 or more elements
func PrintLine(value string, element ...string) {
	w := new(tabwriter.Writer)

	// minwidth, tabwidth, padding, padchar, flags
	w.Init(os.Stdout, 40, 8, 0, '\t', 0)

	defer w.Flush()
	line := value
	for i := 0; i < len(element); i++ {
		line += "\t" + element[i]
	}

	fmt.Fprintf(w, "%s\n", line)

}
