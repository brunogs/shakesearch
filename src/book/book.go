package book

import (
	"fmt"
	"github.com/blevesearch/bleve/v2"
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
	Title   string
}

func Parse(filename string) ([]*Book, *bleve.Index, *bleve.Index, error) {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("parse failure: %w", err)
	}
	suffixarray := suffixarray.New(dat)
	indexes := suffixarray.Lookup([]byte(FIRST), 1)

	titleSplitPattern := regexp.MustCompile("(TITLE:|FINIS)")
	fullContent := string(dat)[indexes[0]:]
	contentByBook := titleSplitPattern.Split(fullContent, -1)

	books := []*Book{}

	for _, content := range contentByBook[:len(contentByBook)-1] {
		newBook := parseBook(content)
		for _, c := range newBook.Chapters {
			c.Title = newBook.Title
		}
		books = append(books, newBook)
	}
	booksIndex, chapterIndex := IndexDocuments(books)
	return books, booksIndex, chapterIndex, nil
}

func parseBook(bookContent string) *Book {
	titleAndContent := strings.SplitAfterN(bookContent, "\n", 2)
	content := titleAndContent[1]
	chapters := ParseChapters(content)
	return &Book{
		Title:    titleAndContent[0],
		Chapters: chapters,
	}
}
