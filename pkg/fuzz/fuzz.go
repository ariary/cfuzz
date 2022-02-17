package fuzz

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/ariary/cfuzz/pkg/config"
)

// PerformFuzzing: Exec specific crafted command for each wordlist file line read
func PerformFuzzing(cfg config.Config) {
	// read wordlist
	wordlist, err := os.Open(cfg.WordlistFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer wordlist.Close()

	scanner := bufio.NewScanner(wordlist) // Caveat: Scanner will error with lines longer than 65536 characters. cf https://stackoverflow.com/questions/8757389/reading-a-file-line-by-line-in-go
	for scanner.Scan() {

		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
