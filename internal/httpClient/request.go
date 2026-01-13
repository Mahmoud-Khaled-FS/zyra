package httpclient

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Request struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    string
	Queries map[string]string
}

func NewRequest(method string, url string) *Request {
	return &Request{
		Method: method,
		URL:    url,
	}
}

func (r *Request) AddHeaders(headers map[string]string) {
	r.Headers = headers
}

func (r *Request) AddQueries(queries map[string]string) {
	r.Queries = queries
}

func (r *Request) AddBody(body string) {
	r.Body = body
}

func (r *Request) Run() {
	start := time.Now()

	var sb strings.Builder

	sb.WriteString(r.URL)

	if len(r.Queries) > 0 {
		sb.WriteString("?")
		pairs := make([]string, 0, len(r.Queries))
		for k, v := range r.Queries {
			pairs = append(pairs, k+"="+v)
		}

		sb.WriteString(strings.Join(pairs, "&"))
	}

	url := sb.String()

	httpReq, err := http.NewRequest(
		strings.ToUpper(r.Method),
		url,
		bytes.NewBufferString(r.Body),
	)

	if err != nil {
		panic(err)
	}

	for k, v := range r.Headers {
		httpReq.Header.Set(k, v)
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(httpReq)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println("Request>>>")
	fmt.Printf("Method: %s\n", r.Method)
	fmt.Printf("URL: %s\n", url)
	fmt.Println("<<<Response")
	fmt.Printf("Status: %d\n", resp.StatusCode)
	fmt.Printf("StatusText: %s\n", resp.Status)
	fmt.Printf("Headers Count: %d\n", len(resp.Header))
	fmt.Printf("Body: %s\n", string(body))
	fmt.Printf("Duration: %d\n", time.Since(start))
}
