package book

import "regexp"

type ChapterParser interface {
	match() bool
	parseChapter() []Chapter
}

func ParseChapters(content string) []Chapter {
	parsers := []ChapterParser{
		&Tragedy{
			content:              content,
			chapterScenesPattern: regexp.MustCompile("(?m)(SCENE.*)"),
		},
		&Novel{
			content:               content,
			chapterNumericPattern: regexp.MustCompile("(?m)^\\s+(\\d+)\\s+"),
			chapterRomanPattern:   regexp.MustCompile("(?m)^[(X)?(IX|IV|V?I{0,3})]{1,}[.]"),
		},
		&Poem{
			content:            content,
			chapterPoemPattern: regexp.MustCompile("(?m)\n+\\s+$"),
		},
	}
	for _, p := range parsers {
		if p.match() {
			return p.parseChapter()
		}
	}
	return nil
}
