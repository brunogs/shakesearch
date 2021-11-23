package book

import (
	"regexp"
)

type Poem struct {
	content            string
	chapterPoemPattern *regexp.Regexp
}

func (p *Poem) match() bool {
	return p.chapterPoemPattern.MatchString(p.content)
}

func (p *Poem) parseChapter() []Chapter {
	chapters := []Chapter{}
	verses := p.chapterPoemPattern.Split(p.content, -1)
	for _, verse := range verses {
		if len(verse) > 0 {
			chapters = append(chapters, Chapter{
				Name:    "",
				Content: verse,
			})
		}
	}
	return chapters
}
