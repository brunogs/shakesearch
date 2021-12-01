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

type BookSearcher struct {
	Books        []book.Book
	BooksByTitle map[string]book.Book
	bookIndex    *bleve.Index
	chapterIndex *bleve.Index
}

func (s *BookSearcher) Load(filename string) error {
	books, booksIndex, chapterIndex, err := book.Parse(filename)
	if err != nil {
		return fmt.Errorf("Load: %w", err)
	}
	s.Books = books
	s.BooksByTitle = make(map[string]book.Book)
	for _, b := range s.Books {
		s.BooksByTitle[strings.TrimSpace(b.Title)] = b
	}
	s.bookIndex = booksIndex
	s.chapterIndex = chapterIndex
	return nil
}

func (s *BookSearcher) Search(query string) resources.QueryResponse {
	response := &resources.QueryResponse{}
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		books := searchDocuments(query, "Title", s.bookIndex)
		response.Books = books
		wg.Done()
	}()
	go func() {
		quotes := searchDocuments(query, "Content", s.chapterIndex)
		response.Quotes = quotes
		wg.Done()
	}()
	wg.Wait()
	return *response
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
	return []book.Book{s.BooksByTitle[title]}
}
