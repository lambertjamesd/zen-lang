package parser

import (
	"fmt"
	"zen/source"
	"zen/tokenizer"
)

type ParseError struct {
	atToken *tokenizer.Token
	message string
}

func CreateError(atToken *tokenizer.Token, message string) (result ParseError) {
	return ParseError{
		atToken,
		message,
	}
}

func FormatError(parseError ParseError) (result string) {
	return fmt.Sprintf("%s\n%s", parseError.message, source.FormatLine(parseError.atToken.Source, parseError.atToken.At))
}
