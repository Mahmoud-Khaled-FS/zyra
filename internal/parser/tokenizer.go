package parser

import (
	"strings"
)

type TokenType string

const (
	// Special
	TokenEOF     TokenType = "EOF"
	TokenIllegal TokenType = "ILLEGAL"
	TokenNewline TokenType = "NEWLINE"

	// Core structure
	TokenDocComment TokenType = "DOC_COMMENT" // """ ... """
	TokenMethod     TokenType = "METHOD"      // GET, POST, etc
	TokenURL        TokenType = "URL"

	// Sections
	TokenSection TokenType = "SECTION" // [headers], [body], [expect]

	// Key-value
	TokenIdentifier TokenType = "IDENTIFIER" // Content-Type
	TokenAssign     TokenType = "ASSIGN"     // =
	TokenValue      TokenType = "VALUE"      // application/json

	// Body
	TokenRaw TokenType = "RAW" // raw body content

	// Separators
	TokenSeparator TokenType = "SEPARATOR" // ---
)

type Token struct {
	Type    TokenType
	Literal string

	Line   int
	Column int
}

type Tokenizer struct {
	src []rune

	pos     int
	readPos int
	ch      rune

	line   int
	column int
}

func NewTokenizer(src string) *Tokenizer {
	l := &Tokenizer{
		src:  []rune(src),
		line: 1,
	}
	l.ch = l.src[l.pos]
	return l
}

func (l *Tokenizer) Tokenize() []Token {
	tokens := []Token{}
	t := l.nextToken()
	for {
		tokens = append(tokens, t)
		if t.Type == TokenEOF {
			break
		}
		t = l.nextToken()
	}
	return tokens
}

func (l *Tokenizer) nextToken() Token {
	l.eatWhitespace()

	tok := Token{
		Line:   l.line,
		Column: l.column,
	}

	switch {
	case l.ch == 0:
		tok.Type = TokenEOF
		l.readChar()
	case l.ch == '"' && l.peekChar() == '"' && l.peekSecondChar() == '"':
		tok.Type = TokenDocComment
		tok.Literal = l.readDocComment()
	case l.ch == '[':
		tok.Type = TokenSection
		tok.Literal = l.readSection()
	case l.ch == '=':
		tok.Type = TokenAssign
		tok.Literal = "="
		l.readChar()
	default:
		literal := l.readIdentifier()
		if isHTTPMethod(literal) {
			tok.Type = TokenMethod
		} else {
			tok.Type = TokenIdentifier
		}
		tok.Literal = literal
	}
	return tok
}

func (l *Tokenizer) readChar() {
	if l.readPos >= len(l.src) {
		l.ch = 0
	} else {
		l.ch = l.src[l.readPos]
	}
	l.pos = l.readPos
	l.readPos++

	if l.ch == '\n' {
		l.line++
		l.column = 0
	} else {
		l.column++
	}
}

func (l *Tokenizer) peekChar() rune {
	if l.readPos >= len(l.src) {
		return 0
	}
	return l.src[l.readPos]
}

func (l *Tokenizer) peekSecondChar() rune {
	if l.readPos+1 >= len(l.src) {
		return 0
	}
	return l.src[l.readPos+1]
}

func (l *Tokenizer) eatWhitespace() {
	for isWhiteSpace(l.ch) {

		l.readChar()
	}
}

func (l *Tokenizer) readDocComment() string {
	start := l.pos
	// Read """
	l.readChar()
	l.readChar()
	l.readChar()

	for !(l.ch == '"' && l.peekChar() == '"' && l.peekSecondChar() == '"') {
		if l.ch == 0 {
			break
		}
		l.readChar()
	}

	content := string(l.src[start : l.pos+3])

	// Read """
	l.readChar()
	l.readChar()
	l.readChar()

	return content
}

func (l *Tokenizer) readSection() string {
	l.readChar() // skip '['

	start := l.pos

	for l.ch != ']' && l.ch != 0 {
		l.readChar()
	}

	section := string(l.src[start:l.pos])

	l.readChar() // skip ']'

	return section
}

func (l *Tokenizer) readIdentifier() string {
	start := l.pos
	for !isWhiteSpace(l.ch) {
		if l.ch == 0 {
			break
		}
		l.readChar()
	}
	return string(l.src[start:l.pos])
}

func isHTTPMethod(text string) bool {
	lowerText := strings.ToLower(text)
	switch lowerText {
	case "get", "post", "put", "patch", "delete":
		return true
	default:
		return false
	}
}

func isLetter(ch rune) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z')
}

// func isDigit(ch rune) bool {
// 	return '0' <= ch && ch <= '9'
// }

func isWhiteSpace(ch rune) bool {
	return ch == ' ' || ch == '\n' || ch == '\r'
}
