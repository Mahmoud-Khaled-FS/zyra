package builtin

import (
	"encoding/json"
	"fmt"
)

func fnEq(actual any, args []any) error {
	if len(args) != 1 {
		return fmt.Errorf("eq expects 1 argument")
	}
	return compare(actual, args[0], "==")
}

func fnNe(actual any, args []any) error {
	return compare(actual, args[0], "!=")
}

func fnGt(actual any, args []any) error {
	return compare(actual, args[0], ">")
}

func fnGte(actual any, args []any) error {
	return compare(actual, args[0], ">=")
}

func fnLt(actual any, args []any) error {
	return compare(actual, args[0], "<")
}

func fnLte(actual any, args []any) error {
	return compare(actual, args[0], "<=")
}

func compare(actual any, expected any, op string) error {
	aNum, aOk := toFloat64(actual)
	eNum, eOk := toFloat64(expected)

	if aOk && eOk {
		switch op {
		case "==":
			if aNum != eNum {
				return fmt.Errorf("%v != %v", aNum, eNum)
			}
		case "!=":
			if aNum == eNum {
				return fmt.Errorf("%v == %v", aNum, eNum)
			}
		case ">":
			if aNum <= eNum {
				return fmt.Errorf("%v <= %v", aNum, eNum)
			}
		case ">=":
			if aNum < eNum {
				return fmt.Errorf("%v < %v", aNum, eNum)
			}
		case "<":
			if aNum >= eNum {
				return fmt.Errorf("%v >= %v", aNum, eNum)
			}
		case "<=":
			if aNum > eNum {
				return fmt.Errorf("%v > %v", aNum, eNum)
			}
		}
		return nil
	}

	aStr := fmt.Sprintf("%v", actual)
	eStr := fmt.Sprintf("%v", expected)

	if op == "==" && aStr != eStr {
		return fmt.Errorf("%s != %s", aStr, eStr)
	}
	if op == "!=" && aStr == eStr {
		return fmt.Errorf("%s == %s", aStr, eStr)
	}

	return nil
}

func fnIs(actual any, args []any) error {
	if len(args) != 1 {
		return fmt.Errorf("is expects 1 argument")
	}
	t, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("type must be string")
	}
	return checkType(actual, t)
}

func fnHas(actual any, args []any) error {
	if len(args) != 1 {
		return fmt.Errorf("has expects 1 argument")
	}

	key := fmt.Sprintf("%v", args[0])

	switch v := actual.(type) {
	case map[string]any:
		if _, ok := v[key]; !ok {
			return fmt.Errorf("missing key: %s", key)
		}
	case map[string]string:
		if _, ok := v[key]; !ok {
			return fmt.Errorf("missing key: %s", key)
		}
	case []any:
		for _, item := range v {
			if fmt.Sprintf("%v", item) == key {
				return nil
			}
		}
		return fmt.Errorf("value not found: %s", key)
	default:
		return fmt.Errorf("has not supported on %T", actual)
	}

	return nil
}

func fnLen(actual any, args []any) error {
	if len(args) != 1 {
		return fmt.Errorf("len expects 1 argument")
	}

	expected, ok := args[0].(int)
	if !ok {
		return fmt.Errorf("len argument must be int")
	}

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

	if l != expected {
		return fmt.Errorf("length %d != %d", l, expected)
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

// CheckType validates that `value` matches the expected type string.
// Supported types: "json", "object", "array", "string", "int", "float", "bool", "null".
func checkType(value any, expectedType string) error {
	switch expectedType {
	case "json":
		switch value.(type) {
		case map[string]any, []any:
			return nil
		default:
			return fmt.Errorf("value is not JSON (object or array), got %T", value)
		}

	case "object":
		if _, ok := value.(map[string]any); !ok {
			return fmt.Errorf("value is not object, got %T", value)
		}

	case "array":
		if _, ok := value.([]any); !ok {
			return fmt.Errorf("value is not array, got %T", value)
		}

	case "string":
		if _, ok := value.(string); !ok {
			return fmt.Errorf("value is not string, got %T", value)
		}

	case "int":
		switch v := value.(type) {
		case int, int32, int64:
			return nil
		case json.Number:
			if _, err := v.Int64(); err != nil {
				return fmt.Errorf("value is not int, got %v", v)
			}
			return nil
		default:
			return fmt.Errorf("value is not int, got %T", value)
		}

	case "float":
		switch v := value.(type) {
		case float32, float64:
			return nil
		case json.Number:
			if _, err := v.Float64(); err != nil {
				return fmt.Errorf("value is not float, got %v", v)
			}
			return nil
		default:
			return fmt.Errorf("value is not float, got %T", value)
		}

	case "bool":
		if _, ok := value.(bool); !ok {
			return fmt.Errorf("value is not bool, got %T", value)
		}

	case "null":
		if value != nil {
			return fmt.Errorf("value is not null, got %T", value)
		}

	default:
		return fmt.Errorf("unknown expected type: %s", expectedType)
	}

	return nil
}
