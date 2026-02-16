package assert

import (
	"fmt"
	"strings"

	"github.com/Mahmoud-Khaled-FS/zyra/internal/assert/builtin"
	httpclient "github.com/Mahmoud-Khaled-FS/zyra/internal/httpClient"
	"github.com/Mahmoud-Khaled-FS/zyra/internal/model"
)

func Evaluate(resp *httpclient.ZyraResponse, a *model.Assertion) error {
	value, err := ResolvePath(resp, a.Path)
	if err != nil {
		return fmt.Errorf("line %d: %w", a.Line, err)
	}

	fn, ok := builtin.Get(a.Fn)
	if !ok {
		return fmt.Errorf("line %d: unknown function '%s'", a.Line, a.Fn)
	}

	var args []any = make([]any, len(a.Args))
	for i, a := range a.Args {
		if a.Type == "ID" {
			path, ok := a.Raw.([]model.PathSegment)
			if ok {
				args[i], err = ResolvePath(resp, path)
				if err != nil {
					args[i] = a.Raw
				}
			}
			continue
		}
		args[i] = a.Raw
	}

	if err := fn(value, args); err != nil {
		return fmt.Errorf("line %d: %w", a.Line, err)
	}

	return nil
}

func ResolvePath(resp *httpclient.ZyraResponse, path []model.PathSegment) (any, error) {
	if len(path) == 0 {
		return nil, fmt.Errorf("empty path")
	}

	seg := path[0]

	if seg.Key == nil {
		return nil, fmt.Errorf("root must be key")
	}

	switch strings.ToLower(*seg.Key) {

	case "status":
		return resp.Status, nil

	case "headers":
		return resolveHeaders(resp.Headers, path[1:])

	case "body":
		return resolveBody(resp.Body, path[1:])

	default:
		return nil, fmt.Errorf("unknown root: %s", *seg.Key)
	}
}

func resolveHeaders(headers map[string]string, path []model.PathSegment) (any, error) {
	if len(path) == 0 {
		return headers, nil
	}

	if path[0].Key == nil {
		return nil, fmt.Errorf("invalid header path")
	}

	key := *path[0].Key
	val, ok := headers[key]
	if !ok {
		return nil, fmt.Errorf("header not found: %s", key)
	}

	return val, nil
}

func resolveBody(body any, path []model.PathSegment) (any, error) {
	current := body

	for _, seg := range path {
		switch v := current.(type) {

		case map[string]any:
			if seg.Key == nil {
				return nil, fmt.Errorf("expected key, got index")
			}
			val, ok := v[*seg.Key]
			if !ok {
				return nil, fmt.Errorf("field not found: %s", *seg.Key)
			}
			current = val

		case []any:
			if seg.Index == nil {
				return nil, fmt.Errorf("expected index, got key")
			}
			if *seg.Index < 0 || *seg.Index >= len(v) {
				return nil, fmt.Errorf("index out of range: %d", *seg.Index)
			}
			current = v[*seg.Index]

		default:
			return nil, fmt.Errorf("cannot traverse %T", current)
		}
	}

	return current, nil
}
