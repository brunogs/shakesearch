package book

import (
	"fmt"
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
}

func Parse(filename string) ([]Book, error) {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("parse failure: %w", err)
	}
	suffixarray := suffixarray.New(dat)
	indexes := suffixarray.Lookup([]byte(FIRST), 1)

	titleSplitPattern := regexp.MustCompile("(TITLE:|FINIS)")
	fullContent := string(dat)[indexes[0]:]
	contentByBook := titleSplitPattern.Split(fullContent, -1)

	books := []Book{}
	for _, content := range contentByBook[:len(contentByBook)-1] {
		books = append(books, parseBook(content))
	}
	return books, nil
}

func parseBook(bookContent string) Book {
	titleAndContent := strings.SplitAfterN(bookContent, "\n", 2)
	chapterNumericPattern := regexp.MustCompile("(?m)^\\s+(\\d+)\\s+")
	chapterRomanPattern := regexp.MustCompile("(?m)^[(X)?(IX|IV|V?I{0,3})]{1,}[.]")
	chapterPoemPattern := regexp.MustCompile("(?m)\n+\\s+$")
	content := titleAndContent[1]
	chapters := []Chapter{}
	if chapterNumericPattern.MatchString(content) {
		chaptersArray := chapterNumericPattern.FindAllString(content, -1)
		for i := 0; i < len(chaptersArray); i++ {
			start := strings.Index(titleAndContent[1], chaptersArray[i])
			var end int
			if i == len(chaptersArray)-1 {
				end = len(content)
			} else {
				end = strings.Index(titleAndContent[1], chaptersArray[i+1])
			}
			chapters = append(chapters, Chapter{
				Name:    chaptersArray[i],
				Content: content[start:end],
			})
		}
	} else if regexp.MustCompile("SCENE").MatchString(content) {
		chapters = append(chapters, Chapter{
			Name:    "Introduction",
			Content: content[0:strings.Index(content, "SCENE")],
		})
		chapterScenesPattern := regexp.MustCompile("(?m)(SCENE.*)")
		contentAfterIntroduction := content[strings.Index(content, "SCENE"):]
		chaptersArray := chapterScenesPattern.FindAllString(contentAfterIntroduction, -1)
		for i := 0; i < len(chaptersArray); i++ {
			start := strings.Index(contentAfterIntroduction, chaptersArray[i])
			var end int
			if i == len(chaptersArray)-1 {
				end = len(contentAfterIntroduction[start:])
			} else {
				end = strings.Index(contentAfterIntroduction[start:], chaptersArray[i+1])
			}
			chapters = append(chapters, Chapter{
				Name:    chaptersArray[i],
				Content: contentAfterIntroduction[start:(start + end)],
			})
		}
	} else if chapterRomanPattern.MatchString(content) {
		chaptersArray := chapterRomanPattern.FindAllString(content, -1)
		for i := 0; i < len(chaptersArray); i++ {
			start := strings.Index(titleAndContent[1], chaptersArray[i])
			var end int
			if i == len(chaptersArray)-1 {
				end = len(content)
			} else {
				end = strings.Index(titleAndContent[1], chaptersArray[i+1])
			}
			chapters = append(chapters, Chapter{
				Name:    chaptersArray[i],
				Content: content[start:end],
			})
		}
	} else {
		verses := chapterPoemPattern.Split(content, -1)
		for _, verse := range verses {
			if len(verse) > 0 {
				chapters = append(chapters, Chapter{
					Name:    "",
					Content: verse,
				})
			}
		}
	}
	return Book{
		Title:    titleAndContent[0],
		Chapters: chapters,
	}
}
