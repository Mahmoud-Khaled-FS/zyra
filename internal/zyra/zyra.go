package zyra

import (
	"github.com/Mahmoud-Khaled-FS/zyra/internal/assert"
	httpclient "github.com/Mahmoud-Khaled-FS/zyra/internal/httpClient"
	"github.com/Mahmoud-Khaled-FS/zyra/internal/parser"
	"github.com/Mahmoud-Khaled-FS/zyra/internal/resolver"
)

type Zyra struct {
	Config *parser.Config
	NoTest bool
}

func NewZyra(config *parser.Config, noTest bool) *Zyra {
	if config == nil {
		config = &parser.Config{}
	}
	return &Zyra{
		Config: config,
		NoTest: noTest,
	}
}

func (z *Zyra) Process(zf ZyraFile) (ZyraResult, error) {
	ctx := resolver.NewContext()
	ctx.Merge(z.Config.Context)
	if len(zf.Doc.Vars) > 0 {
		ctx.Merge(zf.Doc.Vars)
	}

	doc, err := resolver.ResolveDocument(zf.Doc, ctx)
	if err != nil {
		return ZyraResult{}, err
	}

	// 2. build request
	req := httpclient.NewRequest(doc.Method, doc.Path)
	req.AddHeaders(doc.Headers)
	req.AddQueries(doc.Query)
	req.AddBody(doc.Body)
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
	for _, a := range doc.Assertions {
		err = assert.Evaluate(zr, a)
		if err != nil {
			result.Errors = append(result.Errors, err)
		}
	}
	return result, nil
}
