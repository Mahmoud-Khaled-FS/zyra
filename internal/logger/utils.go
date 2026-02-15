package logger

import (
	"fmt"
)

type RequestMeta struct {
	FilePath   string `json:"filePath"`
	Method     string `json:"method"`
	URL        string `json:"url"`
	Assertions int    `json:"assertions"`
	HasVars    bool   `json:"hasVars"`
	HasHeaders bool   `json:"hasHeaders"`
	HasBody    bool   `json:"hasBody"`
}

// ANSI background colors
const (
	bgReset   = "\033[0m"
	bgRed     = "\033[41m"
	bgGreen   = "\033[42m"
	bgYellow  = "\033[43m"
	bgBlue    = "\033[44m"
	bgCyan    = "\033[46m"
	bgMagenta = "\033[45m"
	bgWhite   = "\033[47m"
)

const methodWidth = 6

// MethodColor returns the colored string for a HTTP method
func MethodColor(method string) string {
	switch method {
	case "GET":
		return fmt.Sprintf("%s %s %s", bgBlue, padMethod(method), bgReset)
	case "POST":
		return fmt.Sprintf("%s %s %s", bgGreen, padMethod(method), bgReset)
	case "PUT":
		return fmt.Sprintf("%s %s %s", bgYellow, padMethod(method), bgReset)
	case "PATCH":
		return fmt.Sprintf("%s %s %s", bgCyan, padMethod(method), bgReset)
	case "DELETE":
		return fmt.Sprintf("%s %s %s", bgRed, padMethod(method), bgReset)
	case "OPTIONS":
		return fmt.Sprintf("%s %s %s", bgMagenta, padMethod(method), bgReset)
	default:
		return fmt.Sprintf("%s %s %s", bgWhite, padMethod(method), bgReset)
	}
}

func padMethod(method string) string {
	if len(method) > methodWidth {
		return method[:methodWidth]
	}
	return fmt.Sprintf("%-*s", methodWidth, method) // left-align with spaces
}

func PrintList(reqs []RequestMeta) {
	maxMethod := 0
	maxURL := 0
	for _, r := range reqs {
		if len(r.Method) > maxMethod {
			maxMethod = len(r.Method)
		}
		if len(r.URL) > maxURL {
			maxURL = len(r.URL)
		}
	}

	for _, r := range reqs {
		fmt.Printf(
			"%s  %-*s  %s\n", MethodColor(r.Method),
			maxURL, r.URL,
			r.FilePath,
		)
	}
}
