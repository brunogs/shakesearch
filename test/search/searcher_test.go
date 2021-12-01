package search

import (
	"pulley.com/shakesearch/src/search"
	"testing"
)

var searcher search.BookSearcher
var err error

func init() {
	searcher = search.BookSearcher{}
	err = searcher.Load("../../completeworks.txt")
}

func TestFilterBooksByTitle(t *testing.T) {
	result := searcher.Search("Tragedy")
	if len(result.Books) < 6 || err != nil {
		t.Fatalf(`The searcher was not able to find by title %v  %d`, err, len(result.Books))
	}
}

func TestFilterBooksByTitleIgnoringCase(t *testing.T) {
	result := searcher.Search("TRAGEDY")
	if len(result.Books) < 6 || err != nil {
		t.Fatalf(`The searcher was not able to find by title %v  %d`, err, len(result.Books))
	}
	result = searcher.Search("tragedy")
	if len(result.Books) < 6 || err != nil {
		t.Fatalf(`The searcher was not able to find by title %v  %d`, err, len(result.Books))
	}
}

func TestFilterBooksByChapterName(t *testing.T) {
	result := searcher.Search("a room in the palace")
	if len(result.Quotes) < 1 || err != nil {
		t.Fatalf(`The searcher was not able to find by chapter name %v  %d`, err, len(result.Quotes))
	}
}

func TestFilterBooksByChapterContent(t *testing.T) {
	result := searcher.Search("scene")
	if len(result.Quotes) < 1 || err != nil {
		t.Fatalf(`The searcher was not able to find by chapter content %v  %d`, err, len(result.Quotes))
	}
}

func TestFilterChapterContentInBeginOfChapter(t *testing.T) {
	result := searcher.Search("music to hear")
	if len(result.Quotes) < 1 || err != nil {
		t.Fatalf(`The searcher was not able to find by chapter content %v  %d`, err, len(result.Quotes))
	}

	result = searcher.Search("How CareFul")
	if len(result.Quotes) < 1 || err != nil {
		t.Fatalf(`The searcher was not able to find by chapter content %v  %d`, err, len(result.Quotes))
	}
}

func TestFilterChapterContentInEndOfChapter(t *testing.T) {
	result := searcher.Search("up peerless")
	if len(result.Quotes) < 1 || err != nil {
		t.Fatalf(`The searcher was not able to find by chapter content %v  %d`, err, len(result.Quotes))
	}
}

func TestFilterChapterContentWithSpecialCharacter(t *testing.T) {
	result := searcher.Search("music sadly?")
	if len(result.Quotes) < 1 || err != nil {
		t.Fatalf(`The searcher was not able to find by chapter content %v  %d`, err, len(result.Quotes))
	}
	result = searcher.Search("To be, Or Not to BE")
	if len(result.Quotes) < 1 || err != nil {
		t.Fatalf(`The searcher was not able to find by chapter content %v  %d`, err, len(result.Quotes))
	}
}
