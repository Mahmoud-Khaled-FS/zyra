package main

import (
	"os"

	"github.com/Mahmoud-Khaled-FS/zyra/internal/parser"
	"github.com/Mahmoud-Khaled-FS/zyra/internal/zyra"
)

func main() {
	// cmd.Execute()

	configPath := "./examples/zyra.config"
	bytesConfig, err := os.ReadFile(configPath)
	if err != nil {
		panic(err)
	}

	config, err := parser.ParseConfig(string(bytesConfig))
	if err != nil {
		panic(err)
	}
	// 1. interpolate AST
	// fmt.Println(config)

	path := "./examples/test.zyra"
	bytes, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	doc, err := parser.ParseDocument(string(bytes))
	if err != nil {
		panic(err)
	}
	// fmt.Println(doc)

	z := zyra.NewZyra(config)
	z.Process(doc)

	// lexer := parser.NewTokenizer(string(bytes))
	// tokens := lexer.Tokenize()
	// for _, value := range tokens {
	// 	value.Print()
	// }
	// parser := parser.NewParser(tokens)
	// zyraParsed := parser.ParseRequest()

	// req := httpclient.NewRequest(zyraParsed.Method, zyraParsed.URL)

	// req.AddQueries(zyraParsed.Query)
	// req.AddHeaders(zyraParsed.Headers)
	// req.AddBody(zyraParsed.Body)

	// // interpolator := interpolator.NewInterpolator(zyraConfig.Context)
	// // req.InterpolateRequest(interpolator)

	// req.Send()

	// resp, err := http.Get(request.URL)
	// if err != nil {
	// 	panic(err)
	// }
	// defer resp.Body.Close()

	// body, _ := io.ReadAll(resp.Body)

}
