package resolver

import (
	"fmt"

	"github.com/Mahmoud-Khaled-FS/zyra/internal/model"
)

func ResolveDocument(doc *model.Document, ctx *Context) (*model.Document, error) {
	cp := doc.Clone()

	var err error

	cp.Path, err = interpolate(doc.Path, ctx)
	if err != nil {
		return nil, err
	}

	for k, v := range doc.Headers {
		cp.Headers[k], err = interpolate(v, ctx)
		if err != nil {
			return nil, err
		}
	}

	for k, v := range doc.Query {
		cp.Query[k], err = interpolate(v, ctx)
		if err != nil {
			return nil, err
		}
	}

	cp.Body, err = interpolate(doc.Body, ctx)
	if err != nil {
		return nil, err
	}

	for _, a := range cp.Assertions {
		for i, arg := range a.Args {
			if arg.Type == "template" {
				v, ok := arg.Raw.(string)
				if !ok {
					return nil, fmt.Errorf("Can not parse %v", arg.Raw)
				}
				raw, err := interpolate(v, ctx)
				if err != nil {
					return nil, err
				}
				a.Args[i] = model.Value{Raw: raw, Type: "string"}
			}
		}
	}

	return cp, nil
}
