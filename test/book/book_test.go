package book

import (
	"testing"
)

func TestSplitFileInBooks(t *testing.T) {
	books := GetBooks()
	if len(books) < 44 {
		t.Fatalf(` book.Parse books parsed %d`, len(books))
	}
}

func TestAllBooksMustBeTitleAndChapters(t *testing.T) {
	books := GetBooks()
	for index, book := range books {
		if len(book.Title) == 0 || len(book.Chapters) == 0 {
			t.Fatalf(`The book in index %d is without title (%s) or chapters (%v)`, index, book.Title, book.Chapters)
		}
	}
}

func TestAllBooksMustBeSomeChapters(t *testing.T) {
	books := GetBooks()
	for _, book := range books {
		if len(book.Chapters) <= 1 {
			t.Fatalf(`The book %s have only %d chapter`, book.Title, len(book.Chapters))
		}
	}
}
