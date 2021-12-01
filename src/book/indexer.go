package book

import (
	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/custom"
	"github.com/blevesearch/bleve/v2/analysis/lang/en"
	"github.com/blevesearch/bleve/v2/analysis/token/lowercase"
	"github.com/blevesearch/bleve/v2/analysis/token/porter"
	"github.com/blevesearch/bleve/v2/analysis/tokenizer/unicode"
	"github.com/blevesearch/bleve/v2/mapping"
	"strings"
)

func IndexDocuments(books []*Book) (*bleve.Index, *bleve.Index) {
	bookMapping := buildBookIndexMapping()
	chapterMapping := buildChapterIndexMapping()
	booksIndex, _ := bleve.NewMemOnly(bookMapping)
	chapterIndex, _ := bleve.NewMemOnly(chapterMapping)

	bookBatch := booksIndex.NewBatch()
	chapterBatch := chapterIndex.NewBatch()
	for _, b := range books {
		bookBatch.Index(strings.TrimSpace(b.Title), Book{Title: b.Title})
		for _, c := range b.Chapters {
			c.Title = b.Title
			chapterBatch.Index(b.Title+c.Name, c)
			if chapterBatch.Size() > 500 {
				chapterIndex.Batch(chapterBatch)
				chapterBatch = chapterIndex.NewBatch()
			}
		}
	}
	booksIndex.Batch(bookBatch)
	return &booksIndex, &chapterIndex
}

func buildBookIndexMapping() mapping.IndexMapping {
	titleFieldMapping := bleve.NewTextFieldMapping()
	bookMapping := bleve.NewDocumentMapping()
	bookMapping.AddFieldMappingsAt("Title", titleFieldMapping)
	indexMapping := bleve.NewIndexMapping()
	indexMapping.AddDocumentMapping("book", bookMapping)
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
