package tokenizer

import (
	"unicode"
	"zen/source"
)

type TokenType uint

type tokenizerState func(next rune) (nextState tokenizerState, token TokenType)

const (
	NoToken          TokenType = 0
	IDToken          TokenType = 1
	ErrorToken       TokenType = 2
	WhitespaceToken  TokenType = 3
	EOFToken         TokenType = 4
	NumberToken      TokenType = 5
	OpenSqaureToken  TokenType = 6
	CloseSquareToken TokenType = 7
	OpenCurlyToken   TokenType = 8
	CloseCurlyToken  TokenType = 9
	OpenParenToken   TokenType = 10
	CloseParenToken  TokenType = 11
	ColonToken       TokenType = 12
	FatArrowToken    TokenType = 13
	AssignToken      TokenType = 14
	EqualToken       TokenType = 15
	AddToken         TokenType = 16
	MinusToken       TokenType = 17
	MultiplyToken    TokenType = 18
	DivideToken      TokenType = 19
	DotToken         TokenType = 20
	CommaToken       TokenType = 21
	SemicolonToken   TokenType = 22
	BitwiseOrToken   TokenType = 23
	BitwiseAndToken  TokenType = 24
	BooleanOrToken   TokenType = 25
	BooleanAndToken  TokenType = 26
	LTToken          TokenType = 27
	LTEqToken        TokenType = 28
	GTToken          TokenType = 29
	GTEqToken        TokenType = 30
	NotToken         TokenType = 31
	NotEqualToken    TokenType = 32
)

type Token struct {
	TokenType TokenType
	Value     string
	At        SourceLocation
}

func (token *Token) End() SourceLocation {
	return token.At.WithOffset(len(token.Value))
}

type TokenizeResult struct {
	Tokens []Token
}

func defaultState(next rune) (nextState tokenizerState, token TokenType) {
	return startState(next), NoToken
}

func startState(next rune) (nextState tokenizerState) {
	if unicode.IsSpace(next) {
		return whitespaceState
	} else if unicode.IsLetter(next) || next == '_' {
		return idState
	} else if unicode.IsNumber(next) {
		return integerState
	} else if next == '[' {
		return outputTokenState(OpenSqaureToken)
	} else if next == ']' {
		return outputTokenState(CloseSquareToken)
	} else if next == '{' {
		return outputTokenState(OpenCurlyToken)
	} else if next == '}' {
		return outputTokenState(CloseCurlyToken)
	} else if next == '(' {
		return outputTokenState(OpenParenToken)
	} else if next == ')' {
		return outputTokenState(CloseParenToken)
	} else if next == ':' {
		return outputTokenState(ColonToken)
	} else if next == ';' {
		return outputTokenState(SemicolonToken)
	} else if next == '+' {
		return outputTokenState(AddToken)
	} else if next == '-' {
		return outputTokenState(MinusToken)
	} else if next == '*' {
		return outputTokenState(MultiplyToken)
	} else if next == '/' {
		return outputTokenState(DivideToken)
	} else if next == '.' {
		return outputTokenState(DotToken)
	} else if next == ',' {
		return outputTokenState(CommaToken)
	} else if next == '=' {
		return equalState
	} else if next == '<' {
		return lessThanState
	} else if next == '>' {
		return greaterThanState
	} else if next == '&' {
		return andState
	} else if next == '|' {
		return orState
	} else if next == '!' {
		return notState
	}
	return errorState
}

func outputTokenState(tokenType TokenType) (result tokenizerState) {
	return func(next rune) (nextState tokenizerState, token TokenType) {
		return startState(next), tokenType
	}
}

func idState(next rune) (nextState tokenizerState, token TokenType) {
	if unicode.IsLetter(next) || unicode.IsDigit(next) || next == '_' {
		return idState, NoToken
	} else {
		return startState(next), IDToken
	}
}

func integerState(next rune) (nextState tokenizerState, token TokenType) {
	if unicode.IsDigit(next) {
		return integerState, NoToken
	} else if next == '.' {
		return fractionalState, NoToken
	} else {
		return startState(next), NumberToken
	}
}

func fractionalState(next rune) (nextState tokenizerState, token TokenType) {
	if unicode.IsDigit(next) {
		return fractionalState, NoToken
	} else {
		return startState(next), NumberToken
	}
}

func whitespaceState(next rune) (nextState tokenizerState, token TokenType) {
	if unicode.IsSpace(next) {
		return whitespaceState, NoToken
	} else {
		return startState(next), WhitespaceToken
	}
}

func equalState(next rune) (nextState tokenizerState, token TokenType) {
	if next == '>' {
		return outputTokenState(FatArrowToken), NoToken
	} else if next == '=' {
		return outputTokenState(EqualToken), NoToken
	} else {
		return startState(next), AssignToken
	}
}

func orState(next rune) (nextState tokenizerState, token TokenType) {
	if next == '|' {
		return outputTokenState(BooleanOrToken), NoToken
	} else {
		return startState(next), BitwiseOrToken
	}
}

func notState(next rune) (nextState tokenizerState, token TokenType) {
	if next == '=' {
		return outputTokenState(NotEqualToken), NoToken
	} else {
		return startState(next), NotToken
	}
}

func andState(next rune) (nextState tokenizerState, token TokenType) {
	if next == '&' {
		return outputTokenState(BooleanAndToken), NoToken
	} else {
		return startState(next), BitwiseAndToken
	}
}

func lessThanState(next rune) (nextState tokenizerState, token TokenType) {
	if next == '=' {
		return outputTokenState(LTEqToken), NoToken
	} else {
		return startState(next), LTToken
	}
}

func greaterThanState(next rune) (nextState tokenizerState, token TokenType) {
	if next == '=' {
		return outputTokenState(GTEqToken), NoToken
	} else {
		return startState(next), GTToken
	}
}

func errorState(next rune) (nextState tokenizerState, token TokenType) {
	if next == -1 {
		return errorState, ErrorToken
	} else if unicode.IsSpace(next) {
		return startState(next), ErrorToken
	} else {
		return errorState, NoToken
	}
}

func Tokenize(src *source.Source) (result TokenizeResult) {
	var tokens []Token
	var currentState tokenizerState = defaultState

	var textSource = source.GetSourceContent(src)

	var currentTokenStart int = 0

	for index, character := range textSource {
		var nextState, token = currentState(character)

		if token != NoToken {
			if token != WhitespaceToken {
				tokens = append(tokens, Token{
					token,
					textSource[currentTokenStart:index],
					SourceLocation{
						src,
						currentTokenStart,
					},
				})
			}

			currentTokenStart = index
		}

		currentState = nextState
	}

	var _, lastToken = currentState(-1)

	tokens = append(tokens, Token{
		lastToken,
		textSource[currentTokenStart:len(textSource)],
		SourceLocation{
			src,
			currentTokenStart,
		},
	}, Token{
		EOFToken,
		"",
		SourceLocation{
			src,
			len(textSource),
		},
	})

	return TokenizeResult{
		tokens,
	}
}
