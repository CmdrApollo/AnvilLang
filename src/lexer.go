package main

import (
	"fmt"
	"strings"
	"unicode"
)

type Token struct {
	TokenType  string
	TokenValue string
	Line       int
}

func (t Token) Print() {
	// used for debug only
	if len(t.TokenValue) == 0 {
		fmt.Println(fmt.Sprintf("%d { %s }", t.Line, t.TokenType))
	} else {
		fmt.Println(fmt.Sprintf("%d { %s : %s }", t.Line, t.TokenType, t.TokenValue))
	}
}

// struct that takes a stream of text and turns it into an array of tokens
type Lexer struct {
	Stream string
	Tokens []Token
}

// function that does the lexing
func (l *Lexer) Lex() {
	// first, set up a map of symbols to tokens
	symbolToks := make(map[string]Token)
	// arithmetic operators
	symbolToks["+"] = Token{TokenType: "addsub_operator", TokenValue: "plus"}
	symbolToks["-"] = Token{TokenType: "addsub_operator", TokenValue: "minus"}
	symbolToks["*"] = Token{TokenType: "muldiv_operator", TokenValue: "star"}
	symbolToks["/"] = Token{TokenType: "muldiv_operator", TokenValue: "slash"}
	symbolToks["%"] = Token{TokenType: "muldiv_operator", TokenValue: "percent"}
	// comparison operators
	symbolToks["="] = Token{TokenType: "equal_token"}
	symbolToks["=="] = Token{TokenType: "equality_operator", TokenValue: "equal_equal"}
	symbolToks["!="] = Token{TokenType: "equality_operator", TokenValue: "not_equal"}
	symbolToks["<="] = Token{TokenType: "comparison_operator", TokenValue: "lesser_equal"}
	symbolToks[">="] = Token{TokenType: "comparison_operator", TokenValue: "greater_equal"}
	symbolToks["<"] = Token{TokenType: "comparison_operator", TokenValue: "lesser"}
	symbolToks[">"] = Token{TokenType: "comparison_operator", TokenValue: "greater"}
	symbolToks["&&"] = Token{TokenType: "and_operator"}
	symbolToks["||"] = Token{TokenType: "or_operator"}
	symbolToks["!"] = Token{TokenType: "not_operator"}
	// braces and brackets
	symbolToks["("] = Token{TokenType: "left_paren_token"}
	symbolToks[")"] = Token{TokenType: "right_paren_token"}
	symbolToks["["] = Token{TokenType: "left_square_token"}
	symbolToks["]"] = Token{TokenType: "right_square_token"}
	symbolToks["{"] = Token{TokenType: "left_curly_token"}
	symbolToks["}"] = Token{TokenType: "right_curly_token"}
	// other tokens
	symbolToks[":"] = Token{TokenType: "colon_token"}
	symbolToks[";"] = Token{TokenType: "semicolon_token"}
	symbolToks[","] = Token{TokenType: "comma_token"}
	symbolToks["."] = Token{TokenType: "period_token"}
	symbolToks["?"] = Token{TokenType: "question_token"}
	symbolToks["->"] = Token{TokenType: "arrow_token"}
	// bool literals
	symbolToks["true"] = Token{TokenType: "bool_literal", TokenValue: "true"}
	symbolToks["false"] = Token{TokenType: "bool_literal", TokenValue: "false"}
	// keywords
	symbolToks["fn"] = Token{TokenType: "fn_keyword"}
	symbolToks["st"] = Token{TokenType: "st_keyword"}
	symbolToks["new"] = Token{TokenType: "new_keyword"}
	symbolToks["var"] = Token{TokenType: "var_keyword"}
	symbolToks["while"] = Token{TokenType: "while_keyword"}
	symbolToks["loop"] = Token{TokenType: "loop_keyword"}
	symbolToks["if"] = Token{TokenType: "if_keyword"}
	symbolToks["else"] = Token{TokenType: "else_keyword"}
	symbolToks["return"] = Token{TokenType: "return_keyword"}
	symbolToks["pass"] = Token{TokenType: "pass_keyword"}
	symbolToks["continue"] = Token{TokenType: "continue_keyword"}
	symbolToks["break"] = Token{TokenType: "break_keyword"}
	symbolToks["import"] = Token{TokenType: "import_keyword"}

	line := 1

	// now for the actual reading of data
	i := 0
	for i < len(l.Stream) {
		// get the current character
		character := rune(l.Stream[i])

		// get the next character if it exists
		nextCharacter := byte(0)
		if i+1 < len(l.Stream) {
			nextCharacter = l.Stream[i+1]
		}

		// skip whitespace
		if unicode.IsSpace(character) {
			if string(l.Stream[i]) == "\n" {
				line++
			}
			i++
			continue
		} else if character == '#' {
			// skip comments
			for i < len(l.Stream) {
				if string(l.Stream[i]) == "\n" {
					line++
					break
				}
				i++
			}
			continue
		} else if character == '"' {
			// eat the '"'
			i++

			// strings
			stringLiteral := make([]byte, 0)

			for i < len(l.Stream) {
				nextCharacter = byte(l.Stream[i])

				if nextCharacter == '"' {
					break
				} else if nextCharacter == '\\' {
					i++
					if i >= len(l.Stream) {
						// incomplete escape, add the backslash
						stringLiteral = append(stringLiteral, '\\')
						break
					}
					escapeChar := byte(l.Stream[i])
					switch escapeChar {
					case '"':
						stringLiteral = append(stringLiteral, '"')
					case '\\':
						stringLiteral = append(stringLiteral, '\\')
					case 'n':
						stringLiteral = append(stringLiteral, '\n')
					case 't':
						stringLiteral = append(stringLiteral, '\t')
					case 'r':
						stringLiteral = append(stringLiteral, '\r')
					case 'e':
						stringLiteral = append(stringLiteral, '\033')
					default:
						// unknown escape, add backslash and char
						stringLiteral = append(stringLiteral, '\\')
						stringLiteral = append(stringLiteral, escapeChar)
					}
					i++
				} else {
					stringLiteral = append(stringLiteral, nextCharacter)
					i++
				}
			}

			// add the identifier as a string
			l.Tokens = append(l.Tokens, Token{TokenType: "string_literal", TokenValue: string(stringLiteral), Line: line})
			// do not handle the closing '"' as a character, therefore skip it by not continuing
		} else if unicode.IsLetter(character) || character == '_' {
			// tokenize identifiers
			identifier := make([]byte, 0)

			// loop until invalid character, adding as we go
			for i < len(l.Stream) {
				nextCharacter = byte(l.Stream[i])

				if !(unicode.IsLetter(rune(nextCharacter)) || unicode.IsDigit(rune(nextCharacter)) || nextCharacter == '_') {
					break
				}

				identifier = append(identifier, nextCharacter)
				i++
			}

			// add the identifier as a string
			identifierString := string(identifier)

			tok, ok := symbolToks[identifierString]

			tok.Line = line

			if ok {
				tok.Line = line
				l.Tokens = append(l.Tokens, tok)
			} else {
				l.Tokens = append(l.Tokens, Token{TokenType: "identifier", TokenValue: identifierString, Line: line})
			}

			// handle the last found character normally, hence "continue"
			continue
		} else if unicode.IsDigit(character) {
			// tokenize number literals
			number := make([]byte, 0)

			// loop until invalid character, adding as we go
			for i < len(l.Stream) {
				nextCharacter = byte(l.Stream[i])

				if !(unicode.IsDigit(rune(nextCharacter)) || nextCharacter == '.') {
					break
				}

				number = append(number, nextCharacter)
				i++
			}

			// parse the number string
			numberString := string(number)
			if strings.Count(numberString, ".") == 1 {
				l.Tokens = append(l.Tokens, Token{TokenType: "float_literal", TokenValue: numberString, Line: line})
			} else if strings.Count(numberString, ".") == 0 {
				l.Tokens = append(l.Tokens, Token{TokenType: "int_literal", TokenValue: numberString, Line: line})
			} else {
				fmt.Println("ERROR: Malformed float literal.")
			}

			// handle the last found character normally, hence "continue"
			continue
		} else {
			// tokenize symbols
			tok, ok := Token{}, false

			tok.Line = line

			// double character symbols
			if nextCharacter != byte(0) {
				tok, ok = symbolToks[string(character)+string(nextCharacter)]
			}

			// if it exists, add it to the list
			if ok {
				tok.Line = line
				l.Tokens = append(l.Tokens, tok)
				// advance because double character
				i++
			} else {
				// check single character symbols
				tok, ok = symbolToks[string(character)]

				if ok {
					tok.Line = line
					// it exists, add it to the list
					l.Tokens = append(l.Tokens, tok)
				} else {
					// throw an unknown character error
					fmt.Println(fmt.Sprintf("ERROR: Unknown character '%s'.", string(character)))
				}
			}
		}

		// advance to the next character
		i++
	}
}

func (l Lexer) PrintTokens() {
	for _, tok := range l.Tokens {
		tok.Print()
	}
}
