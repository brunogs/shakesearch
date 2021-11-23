package search

import (
	"fmt"
	"pulley.com/shakesearch/src/book"
	"regexp"
)

type BookSearcher struct {
	Books []book.Book
}

func (s *BookSearcher) Load(filename string) error {
	books, err := book.Parse(filename)
	if err != nil {
		return fmt.Errorf("Load: %w", err)
	}
	s.Books = books
	return nil
}

func (s *BookSearcher) SearchSummaries(query string) []book.Book {
	results := []book.Book{}
	for _, b := range s.Books[0:2] {
		b.Chapters = b.Chapters[:1]
		b.Chapters[0].Content = b.Chapters[0].Content[0:250] + "..."
		results = append(results, b)
	}
	return results
}

func (s *BookSearcher) FindContainsTitles(query string) []book.Book {
	results := []book.Book{}
	pattern, err := regexp.Compile(fmt.Sprintf("(?i)%s", query))
	if err != nil {
		return results
	}
	for _, b := range s.Books {
		if pattern.MatchString(b.Title) {
			results = append(results, b)
		}
	}
	return results
}

func (s *BookSearcher) FindContainsChapterName(query string) []book.Book {
	results := []book.Book{}
	pattern, err := regexp.Compile(fmt.Sprintf("(?i)%s", query))
	if err != nil {
		return results
	}
	for _, b := range s.Books {
		var chapters []book.Chapter
		for _, c := range b.Chapters {
			if pattern.MatchString(c.Name) {
				chapters = append(chapters, c)
			}
		}
		if len(chapters) > 0 {
			bookCopy := b
			bookCopy.Chapters = chapters
			results = append(results, bookCopy)
		}
	}
	return results
}

func (s *BookSearcher) FindContainsChapterContent(query string) []book.Book {
	results := []book.Book{}
	pattern, err := regexp.Compile(fmt.Sprintf("(?i)%s", query))
	if err != nil {
		return results
	}
	for _, b := range s.Books {
		var chapters []book.Chapter
		for _, c := range b.Chapters {
			if pattern.MatchString(c.Content) {
				chapters = append(chapters, c)
			}
		}
		if len(chapters) > 0 {
			bookCopy := b
			bookCopy.Chapters = chapters
			results = append(results, bookCopy)
		}
	}
	return results
}
