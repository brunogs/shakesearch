package book

import (
	"regexp"
	"strings"
)

type Novel struct {
	content               string
	chapterNumericPattern *regexp.Regexp
	chapterRomanPattern   *regexp.Regexp
}

func (n *Novel) match() bool {
	return n.chapterNumericPattern.MatchString(n.content) || n.chapterRomanPattern.MatchString(n.content)
}

func (n *Novel) parseChapter() []Chapter {
	chapters := []Chapter{}
	var chaptersArray []string

	chaptersArray = n.chapterNumericPattern.FindAllString(n.content, -1)
	if chaptersArray == nil || len(chaptersArray) == 0 {
		chaptersArray = n.chapterRomanPattern.FindAllString(n.content, -1)
	}

	for i := 0; i < len(chaptersArray); i++ {
		start := strings.Index(n.content, chaptersArray[i])
		var end int
		if i == len(chaptersArray)-1 {
			end = len(n.content)
		} else {
			end = strings.Index(n.content, chaptersArray[i+1])
		}
		chapters = append(chapters, Chapter{
			Name:    chaptersArray[i],
			Content: n.content[start:end],
		})
	}

	return chapters
}
