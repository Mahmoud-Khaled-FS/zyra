package parser

import (
	"strings"

	"github.com/Mahmoud-Khaled-FS/zyra/internal/assert"
	"github.com/Mahmoud-Khaled-FS/zyra/internal/utils"
)

type Document struct {
	Lines      []Line
	DocComment string

	Method string
	Path   string

	Headers map[string]string
	Query   map[string]string
	Vars    map[string]string
	Body    string

	Assertions []*assert.Assertion
}

func ParseDocument(src string) (*Document, error) {
	lines := splitLines(src)

	p := &parser{
		lines: lines,
		doc: &Document{
			Headers: make(map[string]string),
			Query:   make(map[string]string),
			Vars:    make(map[string]string),
			Lines:   lines,
		},
	}

	if err := p.parseDocument(); err != nil {
		return nil, err
	}

	return p.doc, nil
}

func (p *parser) parseDocument() error {
	for p.pos < len(p.lines) {
		line := strings.TrimSpace(p.current().Text)

		// Comment Line
		if strings.HasPrefix(line, "#") {
			continue
		}

		switch {
		case line == "":
			p.pos++

		case line == `"""`:
			if err := p.parseDocComment(); err != nil {
				return err
			}

		case isRequestLine(line):
			if err := p.parseRequestLine(); err != nil {
				return err
			}

		case isSection(line):
			if err := p.parseDocumentSection(); err != nil {
				return err
			}

		default:
			return p.error("unexpected content")
		}
	}
	return nil
}

func (p *parser) parseDocComment() error {
	p.pos++ // skip opening """

	start := p.pos
	for p.pos < len(p.lines) && strings.TrimSpace(p.current().Text) != `"""` {
		p.pos++
	}

	if p.pos >= len(p.lines) {
		return p.error("unterminated doc comment")
	}

	p.doc.DocComment = collectLines(p.lines[start:p.pos])
	p.pos++ // skip closing """
	return nil
}

func (p *parser) parseRequestLine() error {
	line := p.current().Text
	parts := strings.Fields(line)

	if len(parts) < 2 {
		return p.error("invalid request line")
	}

	p.doc.Method = parts[0]
	p.doc.Path = parts[1]

	p.pos++
	return nil
}

func (p *parser) parseDocumentSection() error {
	section := strings.ToLower(strings.Trim(p.current().Text, "[]"))
	p.pos++

	switch section {
	case "headers":
		return p.parseKeyValueSection(p.doc.Headers)

	case "query":
		return p.parseKeyValueSection(p.doc.Query)

	case "body":
		return p.parseBody()

	case "assert":
		return p.parseAssertSection()

	case "vars":
		return p.parseKeyValueSection(p.doc.Vars)

	default:
		return p.error("unknown section: " + section)
	}
}

func (p *parser) parseAssertSection() error {
	for p.pos < len(p.lines) {
		line := strings.TrimSpace(p.current().Text)

		switch {
		case line == "" || strings.HasPrefix(line, "#"):
			p.pos++

		case isSection(line):
			return nil

		default:
			if err := p.parseAssertionLine(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *parser) parseAssertionLine() error {
	line := strings.TrimSpace(p.current().Text)

	assertion, err := assert.ParseAssertionLine(line, p.current().Num)
	if err != nil {
		return p.error(err.Error())
	}

	p.doc.Assertions = append(p.doc.Assertions, assertion)

	p.pos++
	return nil
}

func isRequestLine(line string) bool {
	if line == "" {
		return false
	}
	method := strings.Fields(line)[0]
	switch strings.ToUpper(method) {
	case "GET", "POST", "PUT", "PATCH", "DELETE":
		return true
	default:
		return false
	}
}

func (d *Document) Clone() *Document {
	cp := &Document{
		DocComment: d.DocComment,
		Method:     d.Method,
		Path:       d.Path,
		Body:       d.Body,
		Headers:    utils.CloneMap(d.Headers),
		Query:      utils.CloneMap(d.Query),
	}
	return cp
}
