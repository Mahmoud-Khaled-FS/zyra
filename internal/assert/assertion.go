package assert

import (
	"fmt"
	"strconv"
	"strings"
)

type Operator string

const (
	OpEq        Operator = "=="
	OpNe        Operator = "!="
	OpGt        Operator = ">"
	OpGte       Operator = ">="
	OpLt        Operator = "<"
	OpLte       Operator = "<="
	OpIs        Operator = "is"
	OpHas       Operator = "has"
	OpLengthEq  Operator = "length=="
	OpLengthGte Operator = "length>="
	OpLengthLte Operator = "length<="
)

type Value struct {
	Raw  any
	Type string
}

type PathSegment struct {
	Key   *string
	Index *int
}

type Assertion struct {
	// Path: dot-separated access to the response object, e.g. res.body.data.id
	Path []PathSegment

	// Operator: "==", ">=", "<=", "is", "has", etc.
	Fn string

	// Value: optional, e.g. 200, "application/json", 3, 55
	// For type checks like `is int`, this can be a string: "int", "json", "object"
	Args []Value

	Line int
}

func ParseAssertionLine(line string, lineNum int) (*Assertion, error) {
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

	args := make([]Value, 0, len(tokens)-2)
	for _, t := range tokens[2:] {
		args = append(args, parseValue(t))
	}

	return &Assertion{
		Path: path,
		Fn:   fn,
		Args: args,
		Line: lineNum,
	}, nil
}

func parsePath(path string) []PathSegment {
	var segments []PathSegment
	var buf strings.Builder
	inBracket := false
	inQuotes := false

	flushKey := func() {
		if buf.Len() == 0 {
			return
		}
		s := buf.String()
		buf.Reset()
		segments = append(segments, PathSegment{Key: &s})
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
				segments = append(segments, PathSegment{Index: &idx})
			} else {
				// string key ["foo"]
				val = strings.Trim(val, `"`)
				segments = append(segments, PathSegment{Key: &val})
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

func parseValue(v string) Value {
	v = strings.TrimSpace(v)

	// string literal
	if strings.HasPrefix(v, `"`) && strings.HasSuffix(v, `"`) {
		return Value{Raw: strings.Trim(v, `"`), Type: "string"}
	}

	// int
	if i, err := strconv.Atoi(v); err == nil {
		return Value{Raw: i, Type: "int"}
	}

	// float
	if f, err := strconv.ParseFloat(v, 64); err == nil {
		return Value{Raw: f, Type: "float"}
	}

	// bool
	if v == "true" || v == "false" {
		return Value{Raw: v == "true", Type: "bool"}
	}

	if strings.HasPrefix(v, "body") || strings.HasPrefix(v, "status") || strings.HasPrefix(v, "headers") {
		return Value{Raw: parsePath(v), Type: "ID"}
	}
	// fallback (identifier like json, object, null)
	return Value{Raw: string(v), Type: "key"}
}

func (a *Assertion) GetPath() string {
	var path strings.Builder
	for _, p := range a.Path {
		path.WriteString(*p.Key)
	}
	return path.String()
}
