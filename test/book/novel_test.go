package book

import (
	"regexp"
	"strings"
	"testing"
)

func TestParseNovelNumericBook(t *testing.T) {
	books := GetBooks()
	novelExample := books[0]

	if strings.Trim(novelExample.Title, "\r\n") != "THE SONNETS" {
		t.Fatalf(`Expected book %s but receives (%s)`, "THE SONNETS", novelExample.Title)
	}
}

func TestParseNovelWithNumericChapters(t *testing.T) {
	books := GetBooks()
	novelExample := books[0]

	for _, c := range novelExample.Chapters {
		if !regexp.MustCompile("^\\s+(\\d+)\\s+").MatchString(c.Name) {
			t.Fatalf(`All the chapter must be a number. Received chapter %s`, c.Name)
		}
		if len(c.Content) == 0 {
			t.Fatalf(`All chapters must have content. Received chapter %s - Content = %s`, c.Name, c.Content)
		}
	}
}

func TestParseNovelRomanBook(t *testing.T) {
	books := GetBooks()
	tragedyExample := books[40]

	if strings.Trim(tragedyExample.Title, "\r\n") != "THE PASSIONATE PILGRIM" {
		t.Fatalf(`Expected book %s but receives (%s)`, "THE PASSIONATE PILGRIM", tragedyExample.Title)
	}
}

func TestParseNovelWithRomanChapters(t *testing.T) {
	books := GetBooks()
	tragedyExample := books[40]

	for _, c := range tragedyExample.Chapters {
		if !regexp.MustCompile("^[(X)?(IX|IV|V?I{0,3})]{1,}[.]").MatchString(c.Name) {
			t.Fatalf(`All the chapter must be a roman number. Received chapter %s`, c.Name)
		}
		if len(c.Content) == 0 {
			t.Fatalf(`All chapters must have content. Received chapter %s - Content = %s`, c.Name, c.Content)
		}
	}
}
