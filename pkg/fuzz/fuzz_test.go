package fuzz

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"sync"
	"testing"
)

func TestCartesianProduct_IncludesAllCombinations(t *testing.T) {
	list1 := []string{"a", "b"}
	list2 := []string{"1", "2", "3"}

	result := cartesianProduct(list1, list2)

	if len(result) != 6 {
		t.Fatalf("expected 6 combinations, got %d", len(result))
	}

	pairs := make([]string, len(result))
	for i, pair := range result {
		pairs[i] = strings.Join(pair, ":")
	}
	sort.Strings(pairs)

	want := []string{"a:1", "a:2", "a:3", "b:1", "b:2", "b:3"}
	for i, w := range want {
		if pairs[i] != w {
			t.Errorf("index %d: want %q got %q", i, w, pairs[i])
		}
	}
}

func TestCartesianProduct_FirstEntryOfList2Included(t *testing.T) {
	list1 := []string{"user"}
	list2 := []string{"firstpass", "secondpass"}

	result := cartesianProduct(list1, list2)

	if len(result) != 2 {
		t.Fatalf("expected 2 combinations, got %d", len(result))
	}
	if result[0][1] != "firstpass" {
		t.Errorf("expected first list2 entry to be included, got %q", result[0][1])
	}
}

type syncWriter struct {
	mu  sync.Mutex
	buf strings.Builder
}

func (w *syncWriter) Write(b []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.buf.Write(b)
}

func TestPerformFuzzing_ProcessesAllWords(t *testing.T) {
	tmp, err := os.CreateTemp("", "cfuzz-wl-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmp.Name())

	const wordCount = 20
	for i := range wordCount {
		fmt.Fprintf(tmp, "word%d\n", i)
	}
	tmp.Close()

	sw := &syncWriter{}
	cfg := DefaultConfig()
	cfg.HideBanner = true
	cfg.Threads = 3
	cfg.Command = "echo FUZZ"
	cfg.Wordlists = []string{tmp.Name()}
	cfg.DisplayModes = BuildDisplayModes(false, false, false, false, false)
	cfg.ResultLogger = log.New(sw, "", 0)

	PerformFuzzing(cfg)

	lines := strings.Split(strings.TrimSpace(sw.buf.String()), "\n")
	if len(lines) != wordCount {
		t.Errorf("expected %d result lines, got %d:\n%s", wordCount, len(lines), sw.buf.String())
	}
}
