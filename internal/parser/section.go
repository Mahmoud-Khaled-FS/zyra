package parser

import "strings"

func (p *parser) parseKeyValueSection(dst map[string]string) error {
	for p.pos < len(p.lines) {
		line := strings.TrimSpace(p.current().Text)

		if line == "" || strings.HasPrefix(line, "#") {
			p.pos++
			continue
		}

		if isSection(line) {
			return nil
		}

		key, val, ok := strings.Cut(line, "=")
		if !ok {
			return p.error("expected key = value")
		}

		dst[strings.TrimSpace(key)] = strings.TrimSpace(val)
		p.pos++
	}
	return nil
}

func (p *parser) parseBody() error {
	start := p.pos

	for p.pos < len(p.lines) && !isSection(strings.TrimSpace(p.current().Text)) {
		p.pos++
	}

	p.doc.Body = collectLines(p.lines[start:p.pos])
	return nil
}

func isSection(line string) bool {
	return strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]")
}
