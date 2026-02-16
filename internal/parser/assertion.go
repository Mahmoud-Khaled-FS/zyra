package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Mahmoud-Khaled-FS/zyra/internal/model"
)

func ParseAssertionLine(line string, lineNum int) (*model.Assertion, error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil, fmt.Errorf("empty line")
	}

	tokens := tokenizeAssertion(line)
	if len(tokens) < 2 {
		return nil, fmt.Errorf("invalid assertion syntax: %s", line)
	}

	path := parsePath(tokens[0])
	fn := tokens[1]

	args := make([]model.Value, 0, len(tokens)-2)
	for _, t := range tokens[2:] {
		args = append(args, parseValue(t))
	}

	return &model.Assertion{
		Path: path,
		Fn:   fn,
		Args: args,
		Line: lineNum,
	}, nil
}

func parsePath(path string) []model.PathSegment {
	var segments []model.PathSegment
	var buf strings.Builder
	inBracket := false
	inQuotes := false

	flushKey := func() {
		if buf.Len() == 0 {
			return
		}
		s := buf.String()
		buf.Reset()
		segments = append(segments, model.PathSegment{Key: &s})
	}

	for i := 0; i < len(path); i++ {
		c := path[i]

		switch c {
		case '.':
			if !inBracket {
				flushKey()
				continue
			}
			buf.WriteByte(c)

		case '[':
			flushKey()
			inBracket = true

		case ']':
			val := buf.String()
			buf.Reset()
			inBracket = false

			// index [0]
			if idx, err := strconv.Atoi(val); err == nil {
				segments = append(segments, model.PathSegment{Index: &idx})
			} else {
				// string key ["foo"]
				val = strings.Trim(val, `"`)
				segments = append(segments, model.PathSegment{Key: &val})
			}

		case '"':
			inQuotes = !inQuotes

		default:
			buf.WriteByte(c)
		}
	}

	flushKey()
	return segments
}

func tokenizeAssertion(s string) []string {
	var tokens []string
	var buf strings.Builder
	inQuotes := false
	brackets := 0

	for i := 0; i < len(s); i++ {
		c := s[i]

		switch c {
		case '"':
			inQuotes = !inQuotes
			buf.WriteByte(c)

		case '[':
			brackets++
			buf.WriteByte(c)

		case ']':
			brackets--
			buf.WriteByte(c)

		case ' ', '\t':
			if inQuotes || brackets > 0 {
				buf.WriteByte(c)
			} else if buf.Len() > 0 {
				tokens = append(tokens, buf.String())
				buf.Reset()
			}

		default:
			buf.WriteByte(c)
		}
	}

	if buf.Len() > 0 {
		tokens = append(tokens, buf.String())
	}

	return tokens
}

func parseValue(v string) model.Value {
	v = strings.TrimSpace(v)

	// string literal
	if strings.HasPrefix(v, `"`) && strings.HasSuffix(v, `"`) {
		return model.Value{Raw: strings.Trim(v, `"`), Type: "string"}
	}

	// int
	if i, err := strconv.Atoi(v); err == nil {
		return model.Value{Raw: i, Type: "int"}
	}

	// float
	if f, err := strconv.ParseFloat(v, 64); err == nil {
		return model.Value{Raw: f, Type: "float"}
	}

	// bool
	if v == "true" || v == "false" {
		return model.Value{Raw: v == "true", Type: "bool"}
	}

	if strings.HasPrefix(v, "body") || strings.HasPrefix(v, "status") || strings.HasPrefix(v, "headers") {
		return model.Value{Raw: parsePath(v), Type: "ID"}
	}

	if strings.HasPrefix(v, "{{") && strings.HasSuffix(v, "}}") {
		return model.Value{Raw: v, Type: "template"}
	}
	// fallback (identifier like json, object, null)
	return model.Value{Raw: string(v), Type: "key"}
}
