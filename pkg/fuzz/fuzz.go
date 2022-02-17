package fuzz

import (
	"bufio"
	"bytes"
	"context"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ariary/cfuzz/pkg/config"
	"github.com/ariary/cfuzz/pkg/output"
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

// PerformFuzzing: Exec specific crafted command for each wordlist file line read
func PerformFuzzing(cfg config.Config) {
	// read wordlist
	wordlist, err := os.Open(cfg.WordlistFilename)
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

		go Exec(cfg, &wg, substituteStr)
	}

	wg.Wait()

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

// Exec: exec the new command and lsend result to print function
// Thanks to https://medium.com/@vCabbage/go-timeout-commands-with-os-exec-commandcontext-ba0c861ed738 for execution timeout
func Exec(cfg config.Config, wg *sync.WaitGroup, substituteStr string) {
	defer wg.Done()
	nCommand := strings.Replace(cfg.Command, cfg.Keyword, substituteStr, 1) //> 0, all replace

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

	// run
	start := time.Now()
	err := cmd.Run()
	elapsed := time.Since(start)

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
func PrintExec(cfg config.Config, result ExecResult) {
	//word counter
	filteredData := strconv.Itoa(len(result.Stdout))
	output.Printline(result.Substitute, filteredData)
}
