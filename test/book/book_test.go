package book

import (
	"pulley.com/shakesearch/src/book"
	"testing"
)

func TestSplitFileInBooks(t *testing.T) {
	books, err := book.Parse("../../completeworks.txt")
	if len(books) < 44 || err != nil {
		t.Fatalf(` book.Parse %v books parsed %d`, err, len(books))
	}
}

func TestAllBooksMustBeTitleAndChapters(t *testing.T) {
	books, _ := book.Parse("../../completeworks.txt")
	for index, book := range books {
		if len(book.Title) == 0 || len(book.Chapters) == 0 {
			t.Fatalf(`The book in index %d is without title (%s) or chapters (%v)`, index, book.Title, book.Chapters)
		}
	}
}

func TestAllBooksMustBeSomeChapters(t *testing.T) {
	books, _ := book.Parse("../../completeworks.txt")
	for _, book := range books {
		if len(book.Chapters) <= 1 {
			t.Fatalf(`The book %s have only %d chapter`, book.Title, len(book.Chapters))
		}
	}
}
