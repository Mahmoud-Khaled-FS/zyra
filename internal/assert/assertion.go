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

type Assertion struct {
	// Path: dot-separated access to the response object, e.g. res.body.data.id
	Path []string

	// Operator: "==", ">=", "<=", "is", "has", etc.
	Operator Operator

	// Value: optional, e.g. 200, "application/json", 3, 55
	// For type checks like `is int`, this can be a string: "int", "json", "object"
	Value any

	Line int
}

func ParseAssertionLine(line string, lineNum int) (Assertion, error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return Assertion{}, fmt.Errorf("empty line")
	}

	// List of operators in order of decreasing length (important!)
	operators := []string{
		" length >= ",
		" length <= ",
		" length == ",
		">=",
		"<=",
		"==",
		"!=",
		">",
		"<",
		" is ",
		" has ",
	}

	for _, op := range operators {
		if strings.Contains(line, op) {
			parts := strings.SplitN(line, op, 2)
			if len(parts) != 2 {
				return Assertion{}, fmt.Errorf("invalid assertion syntax")
			}

			left := strings.TrimSpace(parts[0])
			right := strings.TrimSpace(parts[1])

			// Convert numeric values if possible
			var value any = right
			if num, err := strconv.Atoi(right); err == nil {
				value = num
			} else if strings.HasPrefix(right, `"`) && strings.HasSuffix(right, `"`) {
				value = strings.Trim(right, `"`)
			}

			// Normalize operator for length
			var operator Operator
			switch strings.TrimSpace(op) {
			case "length >=":
				operator = OpLengthGte
			case "length <=":
				operator = OpLengthLte
			case "length ==":
				operator = OpLengthEq
			case "==":
				operator = OpEq
			case "!=":
				operator = OpNe
			case ">=":
				operator = OpGte
			case "<=":
				operator = OpLte
			case ">":
				operator = OpGt
			case "<":
				operator = OpLt
			case "is":
				operator = OpIs
			case "has":
				operator = OpHas
			default:
				return Assertion{}, fmt.Errorf("unknown operator: %s", op)
			}

			return Assertion{
				Path:     parsePath(left),
				Operator: operator,
				Value:    value,
				Line:     lineNum,
			}, nil
		}
	}

	return Assertion{}, fmt.Errorf("no valid operator found in line")
}

func parsePath(path string) []string {
	return strings.Split(path, ".")
}
