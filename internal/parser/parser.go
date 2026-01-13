package parser

import "strings"

type parser struct {
	lines  []Line
	pos    int
	doc    *Document
	config *Config
}

type Line struct {
	Text string
	Num  int
}

func (p *parser) current() Line {
	return p.lines[p.pos]
}

func collectLines(lines []Line) string {
	var b strings.Builder
	for i, l := range lines {
		if i > 0 {
			b.WriteByte('\n')
		}
		b.WriteString(l.Text)
	}
	return b.String()
}

func splitLines(src string) []Line {
	raw := strings.Split(src, "\n")
	lines := make([]Line, 0, len(raw))
	for i, l := range raw {
		lines = append(lines, Line{
			Text: strings.TrimRight(l, "\r"),
			Num:  i + 1,
		})
	}
	return lines
}
