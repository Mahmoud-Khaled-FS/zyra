package model

import (
	"github.com/Mahmoud-Khaled-FS/zyra/internal/utils"
)

type Line struct {
	Text string
	Num  int
}

type Document struct {
	Lines      []Line
	DocComment string

	Method string
	Path   string

	Headers map[string]string
	Query   map[string]string
	Vars    map[string]string
	Body    string

	Assertions []*Assertion
}

type Value struct {
	Raw  any
	Type string
}

func (d *Document) Clone() *Document {
	cp := &Document{
		DocComment: d.DocComment,
		Method:     d.Method,
		Path:       d.Path,
		Body:       d.Body,
		Headers:    utils.CloneMap(d.Headers),
		Query:      utils.CloneMap(d.Query),
	}
	return cp
}
