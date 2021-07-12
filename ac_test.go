package ac_test

import (
	_ "embed"
	ac2 "github.com/kuhsinyv/ac"
	"strings"
	"testing"
)

//go:embed testdata/test_mix.txt
var chinesePatterns string

func TestMultiPatternSearch(t *testing.T) {
	ac := new(ac2.Automaton)
	if err := ac.Build(strings.Split(chinesePatterns, "\n")); err != nil {
		t.Fatal(err)
	}

	terms := ac.MultiPatternSearch([]rune(strings.Replace("刘德华刚去了江 姐家学习 Golang", " ", "", -1)))
	for _, term := range terms {
		t.Log(term.Pos, string(term.Word))
	}
}
