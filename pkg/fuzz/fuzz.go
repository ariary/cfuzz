package fuzz

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"sync/atomic"
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

// countLines counts newlines in a file without loading it fully into memory.
func countLines(filename string) int {
	f, err := os.Open(filename)
	if err != nil {
		return 0
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	n := 0
	for scanner.Scan() {
		n++
	}
	return n
}

// printProgress writes a progress indicator to stderr using \r to overwrite in place.
// total == 0 means unknown (stdin mode).
func printProgress(done, total int64) {
	if total > 0 {
		pct := float64(done) / float64(total) * 100
		fmt.Fprintf(os.Stderr, "\r[%d/%d] %.1f%%  ", done, total, pct)
	} else {
		fmt.Fprintf(os.Stderr, "\r[%d done]  ", done)
	}
}

//cartesianProduct: take two different string slices and return the cartesian product of both
func cartesianProduct(list1 []string, list2 []string) (product [][]string) {
	product = make([][]string, len(list1)*len(list2))
	productIndex := 0
	for i := range list1 {
		for j := range list2 {
			product[productIndex] = append(product[productIndex], list1[i])
			product[productIndex] = append(product[productIndex], list2[j])
			productIndex++
		}
	}
	return product
}

//cartesianProductPlusPlus: Perform cartesian product between a slice of string slice and a string slice. Beware: complexity -> quadratic
func cartesianProductPlusPlus(list1 [][]string, list2 []string) (product [][]string) {
	product = make([][]string, len(list1)*len(list2))
	productIndex := 0
	for i := range list1 {
		for j := range list2 {
			product[productIndex] = append(product[productIndex], list1[i]...)
			product[productIndex] = append(product[productIndex], list2[j])
			productIndex++
		}
	}
	return product
}

// PerformFuzzing executes the fuzz run over the configured wordlist(s).
// Concurrency is limited to cfg.Threads goroutines via a semaphore channel.
// Progress is printed to stderr unless cfg.OnlyWord is true.
func PerformFuzzing(cfg Config) {
	sem := make(chan struct{}, cfg.Threads)
	showProgress := !cfg.OnlyWord

	if !cfg.Multiple {
		var scanner *bufio.Scanner
		var total int64

		if cfg.StdinWordlist {
			scanner = bufio.NewScanner(os.Stdin)
		} else {
			wordlist, err := os.Open(cfg.Wordlists[0])
			if err != nil {
				log.Fatal(err)
			}
			defer wordlist.Close()
			scanner = bufio.NewScanner(wordlist)
			total = int64(countLines(cfg.Wordlists[0]))
		}

		var wg sync.WaitGroup
		var doneCount atomic.Int64

		for scanner.Scan() {
			time.Sleep(time.Duration(cfg.RoutineDelay) * time.Millisecond)
			word := scanner.Text()
			sem <- struct{}{}
			wg.Add(1)
			go func(w string) {
				defer func() {
					<-sem
					n := doneCount.Add(1)
					if showProgress {
						printProgress(n, total)
					}
				}()
				Exec(cfg, &wg, []string{w})
			}(word)
		}
		wg.Wait()
		if showProgress {
			fmt.Fprintln(os.Stderr)
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	} else {
		var wordlists [][]string
		for _, wlPath := range cfg.Wordlists {
			wordlists = append(wordlists, getLines(wlPath))
		}

		substitutes := cartesianProduct(wordlists[0], wordlists[1])
		for i := 2; i < len(wordlists); i++ {
			substitutes = cartesianProductPlusPlus(substitutes, wordlists[i])
		}

		total := int64(len(substitutes))
		var wg sync.WaitGroup
		var doneCount atomic.Int64

		for _, subs := range substitutes {
			sem <- struct{}{}
			wg.Add(1)
			go func() {
				defer func() {
					<-sem
					n := doneCount.Add(1)
					if showProgress {
						printProgress(n, total)
					}
				}()
				Exec(cfg, &wg, subs)
			}()
		}
		wg.Wait()
		if showProgress {
			fmt.Fprintln(os.Stderr)
		}
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
	for _, sub := range substitutesStr {
		nCommand = strings.Replace(nCommand, cfg.Keyword, sub, mode)
		input = strings.Replace(input, cfg.Keyword, sub, mode)
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
	if cfg.FullDisplay {
		PrintFullExecOutput(cfg, result)
		return
	} else {

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
		PrintLine(cfg, result.Substitute, fields...)
	}
}
