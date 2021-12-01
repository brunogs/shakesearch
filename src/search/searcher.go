package search

import (
	"fmt"
	"github.com/blevesearch/bleve/v2"
	"pulley.com/shakesearch/src/book"
	"pulley.com/shakesearch/src/resources"
	"reflect"
	"regexp"
	"strings"
	"sync"
)

const (
	minimumCharacterForSentence int    = 40
	regexSimpleTerm             string = "(?i)\\b%s\\b"
	regexChapterContent         string = "(?im)(([^.]*)?%s\\b[^.]*)\\."
)

type BookSearcher struct {
	Books        []book.Book
	bookIndex    *bleve.Index
	chapterIndex *bleve.Index
}

func (s *BookSearcher) Load(filename string) error {
	books, booksIndex, chapterIndex, err := book.Parse(filename)
	if err != nil {
		return fmt.Errorf("Load: %w", err)
	}
	s.Books = books
	s.bookIndex = booksIndex
	s.chapterIndex = chapterIndex
	return nil
}

func (s *BookSearcher) Search(query string) resources.QueryResponse {
	books := searchDocuments(query, "Title", s.bookIndex)
	quotes := searchDocuments(query, "Content", s.chapterIndex)

	return resources.QueryResponse{
		Books:  books,
		Quotes: quotes,
	}
}

func searchDocuments(query string, fieldToSearch string, index *bleve.Index) []book.Book {
	queryTerm := bleve.NewConjunctionQuery()
	for _, term := range strings.Fields(cleanQuery(query)) {
		q := bleve.NewMatchQuery(term)
		q.SetField(fieldToSearch)
		if reflect.DeepEqual(fieldToSearch, "Content") {
			q.Analyzer = "enWithStopWords"
		}
		queryTerm.AddQuery(q)
	}

	phraseQuery := bleve.NewPhraseQuery(strings.Fields(cleanQuery(query)), fieldToSearch)

	finalQuery := bleve.NewDisjunctionQuery(queryTerm, phraseQuery)
	search := bleve.NewSearchRequest(finalQuery)
	search.Fields = []string{"Title"}
	search.Highlight = bleve.NewHighlightWithStyle("html")
	searchResults, _ := (*index).Search(search)

	results := []book.Book{}
	for _, hit := range searchResults.Hits {
		for k, v := range hit.Fragments {
			c := book.Chapter{Name: k, Content: v[0]}
			results = append(results, book.Book{
				Title:    hit.Fields["Title"].(string),
				Chapters: []book.Chapter{c},
			})
		}
	}
	return results
}

func cleanQuery(query string) string {
	onlyAlphaNumeric := regexp.MustCompile("(?i)[^A-Z0-9\\s]+").ReplaceAllString(query, "")
	return strings.ToLower(onlyAlphaNumeric)
}

func (s *BookSearcher) FindBook(title string) []book.Book {
	queryById := bleve.NewDocIDQuery([]string{title})
	search := bleve.NewSearchRequest(queryById)
	search.Fields = []string{"*"}
	searchResults, _ := (*s.bookIndex).Search(search)

	results := []book.Book{}
	var chapters []book.Chapter

	for _, hit := range searchResults.Hits {
		title := hit.Fields["Title"].(string)
		contents := hit.Fields["Chapters.Content"]

		switch reflect.TypeOf(contents).Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(contents)
			for i := 0; i < s.Len(); i++ {
				chapters = append(chapters, book.Chapter{Content: s.Index(i).Interface().(string)})
			}
		}
		results = append(results, book.Book{Title: title, Chapters: chapters})
	}
	return results
}

func (s *BookSearcher) FindContainsTitles(query string) []book.Book {
	results := []book.Book{}
	pattern, err := regexp.Compile(fmt.Sprintf(regexSimpleTerm, query))
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
	pattern, err := regexp.Compile(fmt.Sprintf(regexSimpleTerm, query))
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
	pattern, _ := regexp.Compile(fmt.Sprintf(regexChapterContent, query))
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
			cleanSentence := strings.TrimSpace(s)
			if len(cleanSentence) > minimumCharacterForSentence {
				chapters <- book.Chapter{Name: c.Name, Content: s}
			}
		}
	}
	defer wg.Done()
}
