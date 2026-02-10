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
	NoTest       bool
}

func NewZyra(config *parser.Config, noTest bool) *Zyra {
	if config == nil {
		config = &parser.Config{}
	}
	return &Zyra{
		Interpolator: interpolator.Interpolator{
			Ctx: config.Context,
		},
		Config: config,
		NoTest: noTest,
	}
}

func (z *Zyra) Process(zf ZyraFile) (ZyraResult, error) {
	resolved, err := z.interpolateDocument(zf.Doc)
	if err != nil {
		return ZyraResult{}, err
	}

	// 2. build request
	req := httpclient.NewRequest(resolved.Method, resolved.Path)
	req.AddHeaders(resolved.Headers)
	req.AddQueries(resolved.Query)
	req.AddBody(resolved.Body)
	zr, err := req.Run()

	if err != nil {
		return ZyraResult{}, err
	}

	result := ZyraResult{
		File:     zf.File,
		Response: zr,
	}

	if z.NoTest {
		return result, nil
	}
	for _, a := range zf.Doc.Assertions {
		err = assert.Evaluate(zr, a)
		if err != nil {
			result.Errors = append(result.Errors, err)
		}
		// if err != nil {
		// 	printAssertionResult(a, zf.Doc.Lines[a.Line-1].Text, err)
		// } else {
		// 	printAssertionResult(a, zf.Doc.Lines[a.Line-1].Text, nil)
		// }
	}
	return result, nil
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
