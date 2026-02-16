package model

import (
	"strings"
)

type Operator string

const (
	OpEq        Operator = "=="
	OpNe        Operator = "!="
	OpGt        Operator = ">"
	OpGte       Operator = ">="
	OpLt        Operator = "<"
	OpLte       Operator = "<="
	OpIs        Operator = "is"
	OpHas       Operator = "has"
	OpLengthEq  Operator = "length=="
	OpLengthGte Operator = "length>="
	OpLengthLte Operator = "length<="
)

type PathSegment struct {
	Key   *string
	Index *int
}

type Assertion struct {
	// Path: dot-separated access to the response object, e.g. res.body.data.id
	Path []PathSegment

	// Operator: "==", ">=", "<=", "is", "has", etc.
	Fn string

	// Value: optional, e.g. 200, "application/json", 3, 55
	// For type checks like `is int`, this can be a string: "int", "json", "object"
	Args []Value

	Line int
}

func (a *Assertion) GetPath() string {
	var path strings.Builder
	for _, p := range a.Path {
		path.WriteString(*p.Key)
	}
	return path.String()
}
