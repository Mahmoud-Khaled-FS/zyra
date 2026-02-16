package zyra

import (
	"github.com/Mahmoud-Khaled-FS/zyra/internal/assert"
	httpclient "github.com/Mahmoud-Khaled-FS/zyra/internal/httpClient"
	"github.com/Mahmoud-Khaled-FS/zyra/internal/parser"
	"github.com/Mahmoud-Khaled-FS/zyra/internal/resolver"
	"github.com/Mahmoud-Khaled-FS/zyra/internal/utils"
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
	url, err := getRequestUrl(doc.Path, z.Config)
	if err != nil {
		return ZyraResult{}, err
	}
	req := httpclient.NewRequest(doc.Method, url)
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

	for _, a := range z.Config.Assertions {
		err = assert.Evaluate(zr, a)
		if err != nil {
			result.Errors = append(result.Errors, err)
		}
	}

	for _, a := range doc.Assertions {
		err = assert.Evaluate(zr, a)
		if err != nil {
			result.Errors = append(result.Errors, err)
		}
	}
	return result, nil
}

func getRequestUrl(path string, config *parser.Config) (string, error) {
	base, ok := config.Options["base_url"]
	if !ok {
		return path, nil
	}

	if utils.IsValidURL(path) {
		return path, nil
	}

	url, err := utils.JoinURL(base, path)
	if err != nil {
		return "", err
	}
	return url, nil
}
