package main

import (
	"os"

	httpclient "github.com/Mahmoud-Khaled-FS/zyra/internal/httpClient"
	"github.com/Mahmoud-Khaled-FS/zyra/internal/parser"
)

func main() {
	path := "./examples/test.zyra"
	bytes, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	lexer := parser.NewTokenizer(string(bytes))
	tokens := lexer.Tokenize()
	// for index, value := range tokens {
	// 	fmt.Printf("%d) %s\n", index, value.Type)
	// }
	parser := parser.NewParser(tokens)
	zyraParsed := parser.ParseRequest()

	req := httpclient.NewRequest(zyraParsed.Method, zyraParsed.URL)

	req.AddQueries(zyraParsed.Query)
	req.AddHeaders(zyraParsed.Headers)
	req.AddBody(zyraParsed.Body)

	req.Send()

	// resp, err := http.Get(request.URL)
	// if err != nil {
	// 	panic(err)
	// }
	// defer resp.Body.Close()

	// body, _ := io.ReadAll(resp.Body)

}
