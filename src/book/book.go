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
	Chapters []Chapter
}

type Chapter struct {
	Name    string
	Content string
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
	for _, content := range contentByBook[:len(contentByBook)-1] {
		books = append(books, parseBook(content))
	}
	return books, nil
}

func parseBook(bookContent string) Book {
	titleAndContent := strings.SplitAfterN(bookContent, "\n", 2)
	content := titleAndContent[1]
	chapters := ParseChapters(content)
	return Book{
		Title:    titleAndContent[0],
		Chapters: chapters,
	}
}
