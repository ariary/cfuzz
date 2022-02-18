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
	line := `[*] ----------------------~~~~~~~~~~~~~~~~~~~---------------------- [*]`
	fmt.Println(line)
	fmt.Println()
	fmt.Println("command fuzzed:\t\t", cfg.Command)
	fmt.Println("wordlist:\t\t", cfg.WordlistFilename)
	//fmt.Println("filtyer type:\t\t", cfg.FilterType.String())
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
