package zyra

import (
	"fmt"

	httpclient "github.com/Mahmoud-Khaled-FS/zyra/internal/httpClient"
	"github.com/Mahmoud-Khaled-FS/zyra/internal/interpolator"
	"github.com/Mahmoud-Khaled-FS/zyra/internal/parser"
)

type Zyra struct {
	Interpolator interpolator.Interpolator
	Config       *parser.Config
	// Builder      runtime.RequestBuilder
}

func NewZyra(config *parser.Config) *Zyra {
	return &Zyra{
		Interpolator: interpolator.Interpolator{
			Ctx: config.Context,
		},
		Config: config,
	}
}

func (z *Zyra) Process(doc *parser.Document) (*httpclient.Request, error) {
	resolved, err := z.interpolateDocument(doc)
	if err != nil {
		return nil, err
	}

	fmt.Println(resolved)

	// 2. build request
	req := httpclient.NewRequest(resolved.Method, resolved.Path)
	req.AddHeaders(resolved.Headers)
	req.AddQueries(resolved.Query)
	req.AddBody(resolved.Body)
	return req, nil
}

func (z *Zyra) interpolateDocument(doc *parser.Document) (*parser.Document, error) {
	cp := doc.Clone()

	var err error

	cp.Path, err = z.Interpolator.Interpolate(doc.Path)
	if err != nil {
		return nil, err
	}

	for k, v := range doc.Headers {
		cp.Headers[k], err = z.Interpolator.Interpolate(v)
		if err != nil {
			return nil, err
		}
	}

	for k, v := range doc.Query {
		cp.Query[k], err = z.Interpolator.Interpolate(v)
		if err != nil {
			return nil, err
		}
	}

	cp.Body, err = z.Interpolator.Interpolate(doc.Body)
	if err != nil {
		return nil, err
	}

	return cp, nil
}
