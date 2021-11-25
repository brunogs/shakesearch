package search

import (
	"pulley.com/shakesearch/src/search"
	"testing"
)

func TestFilterBooksByTitle(t *testing.T) {
	searcher := search.BookSearcher{}
	err := searcher.Load("../../completeworks.txt")

	result := searcher.FindContainsTitles("Tragedy")
	if len(result) < 6 || err != nil {
		t.Fatalf(`The searcher was not able to find by title %v  %d`, err, len(result))
	}
}

func TestFilterBooksByTitleIgnoringCase(t *testing.T) {
	searcher := search.BookSearcher{}
	err := searcher.Load("../../completeworks.txt")

	result := searcher.FindContainsTitles("TRAGEDY")
	if len(result) < 6 || err != nil {
		t.Fatalf(`The searcher was not able to find by title %v  %d`, err, len(result))
	}
	result = searcher.FindContainsTitles("tragedy")
	if len(result) < 6 || err != nil {
		t.Fatalf(`The searcher was not able to find by title %v  %d`, err, len(result))
	}
}

func TestFilterBooksByChapterName(t *testing.T) {
	searcher := search.BookSearcher{}
	err := searcher.Load("../../completeworks.txt")

	result := searcher.FindContainsChapterName("a room in the palace")
	if len(result) < 1 || err != nil {
		t.Fatalf(`The searcher was not able to find by chapter name %v  %d`, err, len(result))
	}
}

func TestFilterBooksByChapterContent(t *testing.T) {
	searcher := search.BookSearcher{}
	err := searcher.Load("../../completeworks.txt")

	result := searcher.FindContainsChapterContent("scene")
	if len(result) < 1 || err != nil {
		t.Fatalf(`The searcher was not able to find by chapter content %v  %d`, err, len(result))
	}
}

func TestFilterChapterContentInBeginOfChapter(t *testing.T) {
	searcher := search.BookSearcher{}
	err := searcher.Load("../../completeworks.txt")

	result := searcher.FindContainsChapterContent("music to hear")
	if len(result) < 1 || err != nil {
		t.Fatalf(`The searcher was not able to find by chapter content %v  %d`, err, len(result))
	}

	result = searcher.FindContainsChapterContent("How CareFul")
	if len(result) < 1 || err != nil {
		t.Fatalf(`The searcher was not able to find by chapter content %v  %d`, err, len(result))
	}
}

func TestFilterChapterContentInEndOfChapter(t *testing.T) {
	searcher := search.BookSearcher{}
	err := searcher.Load("../../completeworks.txt")

	result := searcher.FindContainsChapterContent("up peerless")
	if len(result) < 1 || err != nil {
		t.Fatalf(`The searcher was not able to find by chapter content %v  %d`, err, len(result))
	}
}

func TestFilterChapterContentWithSpecialCharacter(t *testing.T) {
	searcher := search.BookSearcher{}
	err := searcher.Load("../../completeworks.txt")

	result := searcher.FindContainsChapterContent("music sadly?")
	if len(result) < 1 || err != nil {
		t.Fatalf(`The searcher was not able to find by chapter content %v  %d`, err, len(result))
	}
	result = searcher.FindContainsChapterContent("To be, Or Not to BE")
	if len(result) < 1 || err != nil {
		t.Fatalf(`The searcher was not able to find by chapter content %v  %d`, err, len(result))
	}
}
