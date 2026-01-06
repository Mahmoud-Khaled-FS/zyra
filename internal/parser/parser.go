package parser

import (
	"fmt"
	"os"
	"strings"
)

type ParsedZyra struct {
	Doc     string
	Method  string
	URL     string
	Headers map[string]string
	Body    string
	Query   map[string]string
}

type Parser struct {
	tokens []Token
	pos    int
}

type ParseError struct {
	Line    int
	Column  int
	Message string
}

func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens: tokens,
		pos:    0,
	}
}

func (p *Parser) ParseRequest() ParsedZyra {
	req := ParsedZyra{
		Headers: make(map[string]string),
	}

	// optional doc comment
	if p.match(TokenDocComment) {
		req.Doc = strings.TrimRight(strings.TrimLeft(p.tokens[p.pos].Literal, "\"\"\"\n"), "\n\"\"\"")
		p.next()
	}

	// method
	methodTok := p.expect(TokenMethod)
	req.Method = methodTok.Literal

	// URL
	urlTok := p.expect(TokenIdentifier)
	req.URL = urlTok.Literal

	// parse optional sections
	for {
		if p.match(TokenEOF) {
			break
		}
		sectionTok := p.expect(TokenSection)
		switch sectionTok.Literal {
		case "headers":
			req.Headers = p.parseKeyValueBlock()
		case "query":
			req.Query = p.parseKeyValueBlock()
		case "body":
			req.Body = p.parseBodyBlock()
		default:
			p.unexpectedToken(sectionTok)
		}
	}

	return req
}

func (p *Parser) current() Token {
	if p.pos >= len(p.tokens) {
		return Token{Type: TokenEOF}
	}
	return p.tokens[p.pos]
}

func (p *Parser) next() Token {
	p.pos++
	return p.current()
}

func (p *Parser) expect(ttype TokenType) Token {
	tok := p.current()
	if tok.Type != ttype {
		p.unexpectedToken(tok)
	}
	p.next()
	return tok
}

func (p *Parser) match(ttype TokenType) bool {
	if p.current().Type == ttype {
		return true
	}
	return false
}

func (p *Parser) parseKeyValueBlock() map[string]string {
	kv := make(map[string]string)

	for {
		if !p.match(TokenIdentifier) {
			break
		}

		keyTok := p.expect(TokenIdentifier)
		p.expect(TokenAssign)
		valTok := p.expect(TokenIdentifier)

		kv[keyTok.Literal] = valTok.Literal
	}

	return kv
}

func (p *Parser) parseBodyBlock() string {
	var sb strings.Builder
	for p.match(TokenIdentifier) {
		t := p.current()
		sb.WriteString(t.Literal)
		p.next()
	}
	return sb.String()
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("Parse error at line %d, column %d: %s", e.Line, e.Column, e.Message)
}

func (p *Parser) unexpectedToken(tok Token) {
	fmt.Fprintln(os.Stderr, &ParseError{
		Line:    tok.Line,
		Column:  tok.Column,
		Message: fmt.Sprintf("Unexpected token: %q", tok.Literal),
	})
	os.Exit(1)
}
