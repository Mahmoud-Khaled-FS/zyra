package parser

import "strings"

type Config struct {
	Context map[string]string
}

func ParseConfig(src string) (*Config, error) {
	lines := splitLines(src)

	p := &parser{
		lines: lines,
		config: &Config{
			Context: make(map[string]string),
		},
	}

	if err := p.parseConfig(); err != nil {
		return nil, err
	}

	return p.config, nil
}

func (p *parser) parseConfig() error {
	for p.pos < len(p.lines) {
		line := strings.TrimSpace(p.current().Text)

		switch {
		case line == "":
			p.pos++

		case isSection(line):
			if err := p.parseConfigSection(); err != nil {
				return err
			}

		default:
			return p.error("unexpected content")
		}
	}
	return nil
}

func (p *parser) parseConfigSection() error {
	section := strings.ToLower(strings.Trim(p.current().Text, "[]"))
	p.pos++

	switch section {
	case "context":
		return p.parseKeyValueSection(p.config.Context)

	default:
		return p.error("unknown section: " + section)
	}
}
