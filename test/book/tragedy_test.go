package book

import (
	"pulley.com/shakesearch/src/book"
	"strings"
	"testing"
)

func TestParseTragedyBook(t *testing.T) {
	books, _ := book.Parse("../../completeworks.txt")
	tragedyExample := books[1]

	if strings.Trim(tragedyExample.Title, "\r\n") != "ALL’S WELL THAT ENDS WELL" {
		t.Fatalf(`Expected book %s but receives (%s)`, "ALL’S WELL THAT ENDS WELL", tragedyExample.Title)
	}
}

func TestParseTragedyChapters(t *testing.T) {
	books, _ := book.Parse("../../completeworks.txt")
	tragedyExample := books[1]

	if tragedyExample.Chapters[0].Name != "Introduction" {
		t.Fatalf(`Expected book %s but receives (%s)`, "ALL’S WELL THAT ENDS WELL", tragedyExample.Title)
	}
	for _, c := range tragedyExample.Chapters[1:] {
		if strings.Index(c.Name, "SCENE") < 0 {
			t.Fatalf(`All the scenes must be a chapter. Received chapter %s`, c.Name)
		}
		if len(c.Content) == 0 {
			t.Fatalf(`All chapters must have content. Received chapter %s - Content = %s`, c.Name, c.Content)
		}
	}
}
