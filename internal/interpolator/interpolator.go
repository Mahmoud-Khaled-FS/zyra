package interpolator

import (
	"fmt"
	"strings"
)

// TODO (MAHMOUD) - Default Values {{TOKEN|guest}}
// TODO (MAHMOUD) - Nested resolution {{API_{{ENV}}}}
// TODO (MAHMOUD) - Escaping \{{NOT_A_VAR}}

type Context map[string]string

type Interpolator struct {
	Ctx Context
}

func (i *Interpolator) Interpolate(raw string) (string, error) {
	var out strings.Builder

	for {
		start := strings.Index(raw, "{{")
		if start == -1 {
			out.WriteString(raw)
			break
		}

		end := strings.Index(raw[start+2:], "}}")
		if end == -1 {
			return "", fmt.Errorf("unterminated template")
		}

		out.WriteString(raw[:start])

		key := raw[start+2 : start+2+end]
		val, ok := i.Ctx[key]
		if !ok {
			return "", fmt.Errorf("undefined variable: %s", key)
		}

		out.WriteString(val)
		raw = raw[start+2+end+2:]
	}

	return out.String(), nil
}
