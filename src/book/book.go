package book

import (
	"fmt"
	"index/suffixarray"
	"io/ioutil"
	"regexp"
	"strings"
)

const FIRST = "THE SONNETS"

type Book struct {
	Title    string
	Chapters []string
}

func Parse(filename string) ([]Book, error) {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("parse failure: %w", err)
	}

	suffixarray := suffixarray.New(dat)
	indexes := suffixarray.Lookup([]byte(FIRST), 1)

	titleSplitPattern := regexp.MustCompile("(TITLE:|FINIS)")
	fullContent := string(dat)[indexes[0]:]
	contentByBook := titleSplitPattern.Split(fullContent, -1)

	books := []Book{}
	for _, content := range contentByBook[1:] {
		titleAndContent := strings.SplitAfterN(content, "\n", 2)
		books = append(books, Book{Title: titleAndContent[0], Chapters: []string{titleAndContent[1]}})
	}
	return books, nil
}
