package main

import (
	"github.com/Mahmoud-Khaled-FS/zyra/cmd"
)

func main() {
	cmd.Execute()
	return

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
