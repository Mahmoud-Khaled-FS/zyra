package zyra

import (
	"fmt"
	"strings"

	"github.com/Mahmoud-Khaled-FS/zyra/internal/assert"
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
		printAssertionResult(a, "1", err)
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

func formatAssertion(a assert.Assertion) string {
	path := strings.Join(a.Path, ".")
	if a.Value == nil {
		return fmt.Sprintf("%s %s", path, a.Operator)
	}
	return fmt.Sprintf("%s %s %v", path, a.Operator, a.Value)
}

func printAssertionResult(a assert.Assertion, actual any, err error) {
	msg := formatAssertion(a)

	if err != nil {
		fmt.Printf("[FAILED] line %d | %s\n", a.Line, msg)
		fmt.Printf("  → got: %v\n", actual)
	} else {
		fmt.Printf("[PASSED] line %d | %s\n", a.Line, msg)
	}
}
