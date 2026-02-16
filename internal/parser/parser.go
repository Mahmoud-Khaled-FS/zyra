package parser

import (
	"strings"

	"github.com/Mahmoud-Khaled-FS/zyra/internal/model"
)

type parser struct {
	lines  []model.Line
	pos    int
	doc    *model.Document
	config *Config
}

func (p *parser) current() model.Line {
	return p.lines[p.pos]
}

func collectLines(lines []model.Line) string {
	var b strings.Builder
	for i, l := range lines {
		if i > 0 {
			b.WriteByte('\n')
		}
		b.WriteString(l.Text)
	}
	return b.String()
}

func splitLines(src string) []model.Line {
	raw := strings.Split(src, "\n")
	lines := make([]model.Line, 0, len(raw))
	for i, l := range raw {
		lines = append(lines, model.Line{
			Text: strings.TrimRight(l, "\r"),
			Num:  i + 1,
		})
	}
	return lines
}
