package assert

import (
	"encoding/json"
	"fmt"
	"strings"

	httpclient "github.com/Mahmoud-Khaled-FS/zyra/internal/httpClient"
)

func Evaluate(resp *httpclient.ZyraResponse, a Assertion) error {
	if len(a.Path) == 0 {
		return fmt.Errorf("invalid path")
	}

	root := strings.ToLower(a.Path[0])
	switch root {
	case "res":
		return evaluateRes(resp, a)
	default:
		return fmt.Errorf("unknown root: %s", a.Path[0])
	}
}

func evaluateRes(resp *httpclient.ZyraResponse, a Assertion) error {
	if len(a.Path) < 2 {
		return fmt.Errorf("invalid path: %v", a.Path)
	}

	switch strings.ToLower(a.Path[1]) {
	case "status":
		return evaluateStatus(resp.Status, a)

	case "headers":
		return evaluateHeaders(resp.Headers, a)

	case "body":
		return evaluateBody(resp.Body, resp.BodyType, a, a.Path[2:])

	default:
		return fmt.Errorf("unknown response field: %s", a.Path[1])
	}
}

func evaluateStatus(status int, a Assertion) error {
	if a.Operator != OpEq {
		return fmt.Errorf("unsupported operator for status: %s", a.Operator)
	}

	expected, ok := a.Value.(int)
	if !ok {
		return fmt.Errorf("expected value for status must be int")
	}

	if status != expected {
		return fmt.Errorf("status %d != %d", status, expected)
	}
	return nil
}

// --- Headers ---
func evaluateHeaders(headers map[string]string, a Assertion) error {
	if a.Operator == OpHas {
		key, ok := a.Value.(string)
		if !ok {
			return fmt.Errorf("header key must be string")
		}
		if _, exists := headers[key]; !exists {
			return fmt.Errorf("header %s not found", key)
		}
		return nil
	}

	if a.Operator == OpEq {
		key := a.Path[len(a.Path)-1]
		expected, ok := a.Value.(string)
		if !ok {
			return fmt.Errorf("header value must be string")
		}
		actual, exists := headers[key]
		if !exists {
			return fmt.Errorf("header %s not found", key)
		}
		if actual != expected {
			return fmt.Errorf("header %s = %s, expected %s", key, actual, expected)
		}
		return nil
	}

	return fmt.Errorf("unsupported operator for headers: %s", a.Operator)
}

// --- Body ---
func evaluateBody(body any, bodyType httpclient.BodyType, a Assertion, path []string) error {
	// if a.Operator == OpIs && len(path) == 0 {
	// 	expectedType, ok := a.Value.(string)
	// 	if !ok {
	// 		return fmt.Errorf("expected type must be string")
	// 	}
	// }

	// Traverse nested JSON
	current := body
	for _, p := range path {
		m, ok := current.(map[string]any)
		if !ok {
			return fmt.Errorf("cannot traverse path %v, not object", path)
		}
		val, exists := m[p]
		if !exists {
			return fmt.Errorf("path %s does not exist", strings.Join(path, "."))
		}
		current = val
	}

	// Evaluate operator
	switch a.Operator {
	case OpEq, OpNe, OpGt, OpGte, OpLt, OpLte:
		return compare(current, a.Operator, a.Value)
	case OpLengthEq, OpLengthGte, OpLengthLte:
		return compareLength(current, a.Operator, a.Value)
	case OpIs:
		return checkType(current, a.Value.(string))
	default:
		return fmt.Errorf("unsupported operator: %s", a.Operator)
	}
}

func compare(actual any, op Operator, expected any) error {
	// Try numeric comparison first
	aNum, aOk := toFloat64(actual)
	eNum, eOk := toFloat64(expected)

	if aOk && eOk {
		switch op {
		case OpEq:
			if aNum != eNum {
				return fmt.Errorf("%v != %v", aNum, eNum)
			}
		case OpNe:
			if aNum == eNum {
				return fmt.Errorf("%v == %v", aNum, eNum)
			}
		case OpGt:
			if aNum <= eNum {
				return fmt.Errorf("%v <= %v", aNum, eNum)
			}
		case OpGte:
			if aNum < eNum {
				return fmt.Errorf("%v < %v", aNum, eNum)
			}
		case OpLt:
			if aNum >= eNum {
				return fmt.Errorf("%v >= %v", aNum, eNum)
			}
		case OpLte:
			if aNum > eNum {
				return fmt.Errorf("%v > %v", aNum, eNum)
			}
		}
		return nil
	}

	// Otherwise string comparison
	aStr := fmt.Sprintf("%v", actual)
	eStr := fmt.Sprintf("%v", expected)

	switch op {
	case OpEq:
		if aStr != eStr {
			return fmt.Errorf("%v != %v", aStr, eStr)
		}
	case OpNe:
		if aStr == eStr {
			return fmt.Errorf("%v == %v", aStr, eStr)
		}
	default:
		return fmt.Errorf("operator %s not supported for strings", op)
	}
	return nil
}

func compareLength(actual any, op Operator, expected any) error {
	var l int
	switch v := actual.(type) {
	case string:
		l = len(v)
	case []any:
		l = len(v)
	case map[string]any:
		l = len(v)
	default:
		return fmt.Errorf("cannot get length of %T", actual)
	}

	eNum, ok := expected.(int)
	if !ok {
		return fmt.Errorf("expected value must be int for length")
	}

	switch op {
	case OpLengthEq:
		if l != eNum {
			return fmt.Errorf("length %d != %d", l, eNum)
		}
	case OpLengthGte:
		if l < eNum {
			return fmt.Errorf("length %d < %d", l, eNum)
		}
	case OpLengthLte:
		if l > eNum {
			return fmt.Errorf("length %d > %d", l, eNum)
		}
	}
	return nil
}

func checkType(value any, expectedType string) error {
	switch expectedType {
	case "json":
		switch value.(type) {
		case map[string]any, []any:
			return nil
		default:
			return fmt.Errorf("value is not json")
		}

	case "object":
		if _, ok := value.(map[string]any); !ok {
			return fmt.Errorf("value is not object")
		}

	case "array":
		if _, ok := value.([]any); !ok {
			return fmt.Errorf("value is not array")
		}

	case "string":
		if _, ok := value.(string); !ok {
			return fmt.Errorf("value is not string")
		}

	case "int":
		switch value.(type) {
		case int, int64, int32, json.Number:
			// json.Number could be int or float → check
			if n, ok := value.(json.Number); ok {
				if strings.Contains(n.String(), ".") {
					return fmt.Errorf("value is not int")
				}
			}
		default:
			return fmt.Errorf("value is not int")
		}

	case "float":
		switch value.(type) {
		case float32, float64, json.Number:
			if n, ok := value.(json.Number); ok {
				if !strings.Contains(n.String(), ".") {
					return fmt.Errorf("value is not float")
				}
			}
		default:
			return fmt.Errorf("value is not float")
		}

	case "bool":
		if _, ok := value.(bool); !ok {
			return fmt.Errorf("value is not bool")
		}

	case "null":
		if value != nil {
			return fmt.Errorf("value is not null")
		}

	default:
		return fmt.Errorf("unknown type %s", expectedType)
	}

	return nil
}

func toFloat64(v any) (float64, bool) {
	switch val := v.(type) {
	case int:
		return float64(val), true
	case int64:
		return float64(val), true
	case float64:
		return val, true
	case json.Number:
		f, err := val.Float64()
		if err != nil {
			return 0, false
		}
		return f, true
	default:
		return 0, false
	}
}
