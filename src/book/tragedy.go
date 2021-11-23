package book

import (
	"regexp"
	"strings"
)

const (
	SCENE        = "SCENE"
	INTRODUCTION = "Introduction"
)

type Tragedy struct {
	content              string
	chapterScenesPattern *regexp.Regexp
}

func (t *Tragedy) match() bool {
	return regexp.MustCompile(SCENE).MatchString(t.content)
}

func (t *Tragedy) parseChapter() []Chapter {
	chapters := []Chapter{}
	chapters = append(chapters, Chapter{
		Name:    INTRODUCTION,
		Content: t.content[0:strings.Index(t.content, SCENE)],
	})
	contentAfterIntroduction := t.content[strings.Index(t.content, SCENE):]
	chaptersArray := t.chapterScenesPattern.FindAllString(contentAfterIntroduction, -1)
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
	return chapters
}
