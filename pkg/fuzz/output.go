package fuzz

import (
	"fmt"
	"os"
	"text/tabwriter"
)

//BAnner: Print the banner as it is trendy for this kind of tool. thanks to: https://patorjk.com/software/taag
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
	allFilters := ""
	for i := 0; i < len(cfg.Filters); i++ {
		allFilters += cfg.Filters[i].Name() + ", "
	}
	line := `[*] ----------------------~~~~~~~~~~~~~~~~~~~---------------------- [*]`
	fmt.Println(line)
	fmt.Println()
	Printline("command fuzzed", cfg.Command)
	Printline("wordlist:", cfg.WordlistFilename)
	if allFilters != "" {
		allFilters = allFilters[:len(allFilters)-2] //delete last comma
		Printline("filters:", allFilters)
	}
	fmt.Println()
	fmt.Println(line)
	fmt.Println()
}

// Nice printing of a line containing 2 eement the value and the data
func Printline(value string, data string) {
	w := new(tabwriter.Writer)

	// minwidth, tabwidth, padding, padchar, flags
	w.Init(os.Stdout, 40, 8, 0, '\t', 0)

	defer w.Flush()

	fmt.Fprintf(w, "%s\t%s\n", value, data)

}
