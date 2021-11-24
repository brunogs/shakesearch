package search

import (
	"fmt"
	"pulley.com/shakesearch/src/book"
	"pulley.com/shakesearch/src/resources"
	"regexp"
	"strings"
	"sync"
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
	for _, b := range s.Books {
		chapterChannel := s.matchChapterContent(query, b)
		var chapters []book.Chapter
		sentencesSet := make(map[string]struct{})
		for c := range chapterChannel {
			cleanSentence := strings.Trim(c.Content, "\\s")
			if _, sentenceUsed := sentencesSet[cleanSentence]; !sentenceUsed {
				chapters = append(chapters, c)
				sentencesSet[cleanSentence] = struct{}{}
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

func (s *BookSearcher) matchChapterContent(query string, b book.Book) <-chan book.Chapter {
	chapterChannel := make(chan book.Chapter)
	pattern, _ := regexp.Compile(fmt.Sprintf("(?im)([^.]* %s\\b[^.]*)\\.", query))
	var wg sync.WaitGroup
	go func() {
		for _, c := range b.Chapters {
			wg.Add(1)
			go s.matchContent(pattern, c, chapterChannel, &wg)
		}
		wg.Wait()
		defer close(chapterChannel)
	}()
	return chapterChannel
}

func (s *BookSearcher) matchContent(pattern *regexp.Regexp, c book.Chapter, chapters chan book.Chapter, wg *sync.WaitGroup) {
	if pattern.MatchString(c.Content) {
		sentences := pattern.FindAllString(c.Content, -1)
		for _, s := range sentences {
			cleanSentence := strings.Trim(s, "\\s")
			if len(cleanSentence) > 10 {
				chapters <- book.Chapter{Name: c.Name, Content: s}
			}
		}
	}
	defer wg.Done()
}
