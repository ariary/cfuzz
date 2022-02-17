package output

import (
	"fmt"

	"github.com/ariary/cfuzz/pkg/config"
)

//Print the banner as it is trendy for this kind of tool. thanks to: https://patorjk.com/software/taag
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

func PrintConfig(cfg config.Config) {
	line := `[*] ----------------------~~~~~~~~~~~~~~~~~~~---------------------- [*]`
	fmt.Println(line)
	fmt.Println()
	fmt.Println("command fuzzed:\t\t", cfg.Command)
	fmt.Println("wordlist:\t\t", cfg.WordlistFilename)
	fmt.Println("filtyer type:\t\t", cfg.FilterType.String())
	fmt.Println()
	fmt.Println(line)
	fmt.Println()
}
