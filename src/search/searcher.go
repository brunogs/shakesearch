package search

import (
	"fmt"
	"pulley.com/shakesearch/src/book"
	"pulley.com/shakesearch/src/resources"
	"regexp"
	"strings"
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

func (s *BookSearcher) SearchSummaries(query string) resources.QueryResponse {
	books := s.FindContainsTitles(query)
	quotes := s.FindContainsChapterContent(query)
	return resources.QueryResponse{
		Books:  books,
		Quotes: quotes,
	}
}

func (s *BookSearcher) FindContainsTitles(query string) []book.Book {
	results := []book.Book{}
	pattern, err := regexp.Compile(fmt.Sprintf("(?i)\\b%s\\b", query))
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
	pattern, err := regexp.Compile(fmt.Sprintf("(?i)\\b%s\\b", query))
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
	pattern, err := regexp.Compile(fmt.Sprintf("(?im)([^.]* %s [^.]*)\\.", query))
	if err != nil {
		return results
	}
	for _, b := range s.Books {
		var chapters []book.Chapter
		sentencesSet := make(map[string]struct{})
		for _, c := range b.Chapters {
			if pattern.MatchString(c.Content) {
				sentences := pattern.FindAllString(c.Content, -1)
				for _, s := range sentences {
					cleanSentence := strings.Trim(s, "\\s")
					if _, sentenceUsed := sentencesSet[cleanSentence]; !sentenceUsed && len(cleanSentence) > 10 {
						chapters = append(chapters, book.Chapter{Name: c.Name, Content: s})
						sentencesSet[cleanSentence] = struct{}{}
					}
				}
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
