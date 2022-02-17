package output

import "fmt"

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
}
