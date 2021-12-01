package book

import (
	"fmt"
	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/custom"
	"github.com/blevesearch/bleve/v2/analysis/lang/en"
	"github.com/blevesearch/bleve/v2/analysis/token/lowercase"
	"github.com/blevesearch/bleve/v2/analysis/token/porter"
	"github.com/blevesearch/bleve/v2/analysis/tokenizer/unicode"
	"github.com/blevesearch/bleve/v2/mapping"
	"index/suffixarray"
	"io/ioutil"
	"regexp"
	"strings"
)

const FIRST = "THE SONNETS"

type Book struct {
	Title    string
	Chapters []Chapter
}

type Chapter struct {
	Name    string
	Content string
	Title   string
}

func Parse(filename string) ([]Book, *bleve.Index, *bleve.Index, error) {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("parse failure: %w", err)
	}
	suffixarray := suffixarray.New(dat)
	indexes := suffixarray.Lookup([]byte(FIRST), 1)

	titleSplitPattern := regexp.MustCompile("(TITLE:|FINIS)")
	fullContent := string(dat)[indexes[0]:]
	contentByBook := titleSplitPattern.Split(fullContent, -1)

	books := []Book{}
	bookMapping := buildBookIndexMapping()
	chapterMapping := buildChapterIndexMapping()
	booksIndex, _ := bleve.NewMemOnly(bookMapping)
	chapterIndex, _ := bleve.NewMemOnly(chapterMapping)

	for _, content := range contentByBook[:len(contentByBook)-1] {
		newBook := parseBook(content)
		booksIndex.Index(strings.TrimSpace(newBook.Title), newBook)
		for _, c := range newBook.Chapters {
			c.Title = newBook.Title
			chapterIndex.Index(newBook.Title+c.Name, c)
		}
		books = append(books, newBook)
	}
	return books, &booksIndex, &chapterIndex, nil
}

func parseBook(bookContent string) Book {
	titleAndContent := strings.SplitAfterN(bookContent, "\n", 2)
	content := titleAndContent[1]
	chapters := ParseChapters(content)
	return Book{
		Title:    titleAndContent[0],
		Chapters: chapters,
	}
}

func buildBookIndexMapping() mapping.IndexMapping {
	titleFieldMapping := bleve.NewTextFieldMapping()
	titleFieldMapping.IncludeTermVectors = true
	titleFieldMapping.Analyzer = "enWithStopWords"

	bookMapping := bleve.NewDocumentMapping()
	bookMapping.AddFieldMappingsAt("Title", titleFieldMapping)

	indexMapping := bleve.NewIndexMapping()
	indexMapping.AddDocumentMapping("book", bookMapping)
	indexMapping.AddCustomAnalyzer("enWithStopWords",
		map[string]interface{}{
			"type":      custom.Name,
			"tokenizer": unicode.Name,
			"token_filters": []string{
				en.PossessiveName,
				lowercase.Name,
				porter.Name,
			},
		})
	indexMapping.DefaultMapping = bookMapping

	return indexMapping
}

func buildChapterIndexMapping() mapping.IndexMapping {
	englishTextFieldMapping := bleve.NewTextFieldMapping()
	englishTextFieldMapping.Analyzer = en.AnalyzerName
	englishTextFieldMapping.IncludeTermVectors = true

	contentFieldMapping := bleve.NewTextFieldMapping()
	contentFieldMapping.IncludeTermVectors = true
	contentFieldMapping.Analyzer = "enWithStopWords"

	chapterMapping := bleve.NewDocumentMapping()
	chapterMapping.AddFieldMappingsAt("Title", englishTextFieldMapping)
	chapterMapping.AddFieldMappingsAt("Content", contentFieldMapping)

	indexMapping := bleve.NewIndexMapping()
	indexMapping.AddDocumentMapping("bookChapter", chapterMapping)
	indexMapping.AddCustomAnalyzer("enWithStopWords",
		map[string]interface{}{
			"type":      custom.Name,
			"tokenizer": unicode.Name,
			"token_filters": []string{
				en.PossessiveName,
				lowercase.Name,
				porter.Name,
			},
		})
	indexMapping.DefaultMapping = chapterMapping
	return indexMapping
}
