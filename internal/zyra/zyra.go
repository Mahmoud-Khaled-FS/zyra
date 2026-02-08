package zyra

import (
	"fmt"

	"github.com/Mahmoud-Khaled-FS/zyra/internal/assert"
	httpclient "github.com/Mahmoud-Khaled-FS/zyra/internal/httpClient"
	"github.com/Mahmoud-Khaled-FS/zyra/internal/interpolator"
	"github.com/Mahmoud-Khaled-FS/zyra/internal/logger"
	"github.com/Mahmoud-Khaled-FS/zyra/internal/parser"
)

type Zyra struct {
	Interpolator interpolator.Interpolator
	Config       *parser.Config
	// Builder      runtime.RequestBuilder
}

func NewZyra(config *parser.Config) *Zyra {
	if config == nil {
		config = &parser.Config{}
	}
	return &Zyra{
		Interpolator: interpolator.Interpolator{
			Ctx: config.Context,
		},
		Config: config,
	}
}

func (z *Zyra) Process(doc *parser.Document) error {
	resolved, err := z.interpolateDocument(doc)
	if err != nil {
		return err
	}

	// 2. build request
	req := httpclient.NewRequest(resolved.Method, resolved.Path)
	req.AddHeaders(resolved.Headers)
	req.AddQueries(resolved.Query)
	req.AddBody(resolved.Body)
	zr, err := req.Run()

	if err != nil {
		panic(err)
	}

	for _, a := range doc.Assertions {
		err = assert.Evaluate(zr, a)
		if err != nil {
			printAssertionResult(a, doc.Lines[a.Line-1].Text, err)
		} else {
			printAssertionResult(a, doc.Lines[a.Line-1].Text, nil)
		}
	}
	return nil
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

func printAssertionResult(a *assert.Assertion, line string, err error) {

	if err != nil {
		logger.Failed("line %d | %s", a.Line, line)
		fmt.Printf("  → error: %v\n", err)
	} else {
		logger.Passed("line %d | %s", a.Line, line)
	}
}
