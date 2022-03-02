package fuzz

import (
	"bufio"
	"bytes"
	"context"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type ExecResult struct {
	Substitute string
	Stdout     string
	Stderr     string
	Time       time.Duration
	Code       string
	Error      error
	Timeout    bool
}

// func (er *ExecResult) timeTrack(start time.Time) {
// 	elapsed := time.Since(start)
// 	er.Time = elapsed
// }

//getLines: read a file a return a slice containing each lines
func getLines(filename string) (wordlist []string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		wordlist = append(wordlist, scanner.Text())
	}

	return wordlist
}

//cartesianProduct: take two different string slices and return the cartesian product of both
func cartesianProduct(list1 []string, list2 []string) (product [][]string) {
	product = make([][]string, len(list1)*(len(list2)-1))
	productIndex := 0
	for i := 0; i < len(list1); i++ { //for each item of first list
		for j := 1; j < len(list2); j++ { //couple it with other
			product[productIndex] = append(product[productIndex], list1[i])
			product[productIndex] = append(product[productIndex], list2[j])
			productIndex++
		}
	}
	return product
}

//cartesianProductPlusPlus: Perform cartesian product between a slice of string slice and a string slice. Beware: complexity -> quadratic
func cartesianProductPlusPlus(list1 [][]string, list2 []string) (product [][]string) {
	product = make([][]string, len(list1)*(len(list2)))
	productIndex := 0
	for i := 0; i < len(list1); i++ { //for each item of first list
		for j := 0; j < len(list2); j++ { //couple it with other
			product[productIndex] = append(product[productIndex], list1[i]...)
			product[productIndex] = append(product[productIndex], list2[j])
			productIndex++
		}
	}
	return product
}

//PerformFuzzing: Exec specific crafted command for each wordlist file line read
func PerformFuzzing(cfg Config) {
	// read wordlist
	if !cfg.Multiple { /////////KEEP THIS ITERATION IF SIMPLE (not multiple) => AVOID BROWSING THE WORDLIST TWICE
		wordlist, err := os.Open(cfg.Wordlists[0])
		if err != nil {
			log.Fatal(err)
		}
		defer wordlist.Close()

		var wg sync.WaitGroup

		scanner := bufio.NewScanner(wordlist) // Caveat: Scanner will error with lines longer than 65536 characters. cf https://stackoverflow.com/questions/8757389/reading-a-file-line-by-line-in-go
		for scanner.Scan() {
			time.Sleep(time.Duration(cfg.RoutineDelay) * time.Millisecond)
			wg.Add(1)
			substituteStr := scanner.Text()

			go Exec(cfg, &wg, []string{substituteStr})
		}

		wg.Wait()

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	} else { //multiple
		//construct lists of word containing in wordlist
		var wordlists [][]string
		for i := 0; i < len(cfg.Wordlists); i++ {
			wordlists = append(wordlists, getLines(cfg.Wordlists[i]))
		}

		//Browse list
		substitutes := cartesianProduct(wordlists[0], wordlists[1])
		for i := 2; i < len(wordlists); i++ {
			substitutes = cartesianProductPlusPlus(substitutes, wordlists[i])

		}

		var wg sync.WaitGroup

		for i := 0; i < len(substitutes); i++ {
			wg.Add(1)
			go Exec(cfg, &wg, substitutes[i])
		}

		wg.Wait()

	}

}

//Exec: exec the new command and send result to print function
// Thanks to https://medium.com/@vCabbage/go-timeout-commands-with-os-exec-commandcontext-ba0c861ed738 for execution timeout
func Exec(cfg Config, wg *sync.WaitGroup, substitutesStr []string) {
	defer wg.Done()
	var mode int
	if !cfg.Multiple {
		mode = -1 //< 0 ~ all replace
	} else {
		mode = 1 //replace only first occurrence
	}

	nCommand := cfg.Command
	input := cfg.Input
	for i := 0; i < len(substitutesStr); i++ {
		nCommand = strings.Replace(nCommand, cfg.Keyword, substitutesStr[i], mode)

		input = strings.Replace(input, cfg.Keyword, substitutesStr[i], mode)
	}

	// Create a new context and add a timeout to it
	timeout := time.Duration(cfg.Timeout) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel() // The cancel should be deferred so resources are cleaned up

	// Create the command with our context
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, cfg.Shell, "-c", nCommand)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if cfg.Input != "" {
		cmd.Stdin = strings.NewReader(input)
	}

	// run
	start := time.Now()
	err := cmd.Run()
	elapsed := time.Since(start)

	substituteStr := strings.Join(substitutesStr, ",")
	result := ExecResult{Substitute: substituteStr}
	if ctx.Err() == context.DeadlineExceeded {
		result.Timeout = true
		result.Time = timeout
	} else {
		result.Timeout = false
		result.Time = elapsed
	}

	result.Stdout = stdout.String()
	result.Stderr = stderr.String()

	if err != nil {
		result.Error = err
		result.Code = err.Error() //killed, 2, etc
	} else {
		result.Code = "0"
	}

	PrintExec(cfg, result)
}

// PrintExec: Print execution result according to configuration and filter
func PrintExec(cfg Config, result ExecResult) {
	// filter
	// if cfg.Hide { //hide field that pass filter and show others
	// 	for i := 0; i < len(cfg.Filters); i++ {
	// 		if cfg.Filters[i].IsOk(result) {
	// 			return //don't display it
	// 		}
	// 	}
	// } else {
	// 	for i := 0; i < len(cfg.Filters); i++ {
	// 		if !cfg.Filters[i].IsOk(result) {
	// 			return //don't display it
	// 		}
	// 	}
	// }
	for i := 0; i < len(cfg.Filters); i++ {
		if cfg.Filters[i].IsOk(result) == cfg.Hide {
			return //don't display it
		}
	}
	// display
	var fields []string
	for i := 0; i < len(cfg.DisplayModes); i++ {
		fields = append(fields, cfg.DisplayModes[i].DisplayString(result))
	}
	PrintLine(result.Substitute, fields...)

}
