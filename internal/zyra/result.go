package zyra

import (
	"fmt"
	"strings"

	httpclient "github.com/Mahmoud-Khaled-FS/zyra/internal/httpClient"
	"github.com/Mahmoud-Khaled-FS/zyra/internal/utils"
)

type ZyraResult struct {
	Errors   []error
	File     string
	Response *httpclient.ZyraResponse
}

const (
	green  = "\033[32m"
	red    = "\033[31m"
	yellow = "\033[33m"
	reset  = "\033[0m"
	bold   = "\033[1m"
)

func BeautyLogger(results []ZyraResult) {
	for _, res := range results {
		fmt.Printf("%sFile:%s %s\n", bold, reset, res.File)

		if len(res.Errors) == 0 {
			fmt.Printf("  %s✔ PASSED%s\n", green, reset)
		} else {
			fmt.Printf("  %s✖ FAILED%s\n", red, reset)
			for i, err := range res.Errors {
				fmt.Printf("    %d) %s\n", i+1, err.Error())
			}
		}

		if res.Response != nil {
			fmt.Printf("  Response Status: %d\n", res.Response.Status)
			fmt.Printf("  Response Duration: %s\n", utils.PrettyDuration(res.Response.Duration))
		}

		fmt.Println(strings.Repeat("-", 40))
	}
}
