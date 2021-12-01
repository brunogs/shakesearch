package book

import (
	"strings"
	"testing"
)

func TestParsePoemBook(t *testing.T) {
	books := GetBooks()
	poemExample := books[41]

	if strings.Trim(poemExample.Title, "\r\n") != "THE PHOENIX AND THE TURTLE" {
		t.Fatalf(`Expected book %s but receives (%s)`, "THE PHOENIX AND THE TURTLE", poemExample.Title)
	}
}

func TestParsePoemChapters(t *testing.T) {
	books := GetBooks()
	poemExample := books[41]

	for _, c := range poemExample.Chapters {
		if len(c.Name) > 0 {
			t.Fatalf(`All the poems chapter must not have a title. Received chapter %s`, c.Name)
		}
		if len(c.Content) == 0 {
			t.Fatalf(`All chapters must have content. Received chapter %s - Content = %s`, c.Name, c.Content)
		}
	}
}
