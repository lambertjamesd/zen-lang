package parser

import (
	"fmt"
	"zen/source"
	"zen/tokenizer"
)

type ParseError struct {
	At      tokenizer.SourceLocation
	message string
}

func CreateError(at tokenizer.SourceLocation, message string) (result ParseError) {
	return ParseError{
		at,
		message,
	}
}

func FormatError(parseError ParseError) (result string) {
	return fmt.Sprintf("%s\n%s", parseError.message, source.FormatLine(parseError.At.Source, parseError.At.At))
}
