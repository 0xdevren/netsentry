// Package cisco provides parsers for Cisco IOS and IOS-XE device configurations.
package cisco

import (
	"bufio"
	"bytes"
	"strings"
)

// TokenType classifies a lexer token from a Cisco IOS configuration.
type TokenType int

const (
	// TokenLine is a plain configuration line.
	TokenLine TokenType = iota
	// TokenBlockStart is the beginning of an indented configuration block.
	TokenBlockStart
	// TokenBlockEnd marks the return to a lower indentation level.
	TokenBlockEnd
	// TokenComment is an exclamation-mark comment line.
	TokenComment
)

// Token is a single lexical unit from the IOS configuration.
type Token struct {
	// Type classifies this token.
	Type TokenType
	// Text is the raw line text with leading whitespace preserved.
	Text string
	// Depth is the indentation depth (number of leading spaces / 1 space unit).
	Depth int
}

// Lexer tokenises Cisco IOS configuration text into a flat token stream.
// IOS uses single-space indentation, so Depth reflects leading space count.
type Lexer struct{}

// NewLexer constructs a new IOS Lexer.
func NewLexer() *Lexer { return &Lexer{} }

// Tokenise scans data line by line and returns the token stream.
func (l *Lexer) Tokenise(data []byte) []Token {
	var tokens []Token
	scanner := bufio.NewScanner(bytes.NewReader(data))
	prevDepth := 0

	for scanner.Scan() {
		raw := scanner.Text()
		trimmed := strings.TrimLeft(raw, " ")
		depth := len(raw) - len(trimmed)

		// Skip comment lines.
		if strings.HasPrefix(trimmed, "!") {
			tokens = append(tokens, Token{Type: TokenComment, Text: raw, Depth: depth})
			continue
		}
		// Skip empty lines.
		if trimmed == "" {
			continue
		}

		tt := TokenLine
		if depth > prevDepth {
			tt = TokenBlockStart
		} else if depth < prevDepth {
			tt = TokenBlockEnd
		}
		tokens = append(tokens, Token{Type: tt, Text: trimmed, Depth: depth})
		prevDepth = depth
	}

	return tokens
}
