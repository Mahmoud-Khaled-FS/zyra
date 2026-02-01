package httpclient

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"
)

type BodyType int

const (
	BodyTypeObject BodyType = iota
	BodyTypeArray
	BodyTypeString
	BodyTypeInt
	BodyTypeFloat
	BodyTypeBool
	BodyTypeNull
	BodyTypeUnknown
)

type ZyraResponse struct {
	Status   int
	RawBody  []byte
	Body     any
	Headers  map[string]string
	BodyType BodyType
	Duration time.Duration
}

func NewResponse(resp *http.Response) (*ZyraResponse, error) {
	defer resp.Body.Close()
	zr := &ZyraResponse{
		Status: resp.StatusCode,
	}

	var headers map[string]string = make(map[string]string, len(resp.Header))

	for k, v := range resp.Header {
		if len(v) > 0 {
			headers[k] = v[0]
		}
	}
	zr.Headers = headers

	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	zr.RawBody = rawBody

	var data any
	decoder := json.NewDecoder(bytes.NewReader(rawBody))
	decoder.UseNumber()

	if err := decoder.Decode(&data); err != nil {
		// Not JSON → treat as string
		zr.BodyType = BodyTypeString
		zr.Body = string(rawBody)
		return zr, nil
	}

	switch v := data.(type) {
	case map[string]any:
		zr.BodyType = BodyTypeObject
	case []any:
		zr.BodyType = BodyTypeArray
	case string:
		zr.BodyType = BodyTypeString
	case json.Number:
		if isFloat(v) {
			zr.BodyType = BodyTypeFloat
		} else {
			zr.BodyType = BodyTypeInt
		}
	case bool:
		zr.BodyType = BodyTypeBool
	case nil:
		zr.BodyType = BodyTypeNull
	default:
		zr.BodyType = BodyTypeUnknown
	}

	zr.Body = data

	return zr, nil
}

func isFloat(n json.Number) bool {
	s := n.String()
	return strings.ContainsAny(s, ".eE")
}
