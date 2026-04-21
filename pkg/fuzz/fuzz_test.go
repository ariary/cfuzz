package fuzz

import (
	"sort"
	"strings"
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
