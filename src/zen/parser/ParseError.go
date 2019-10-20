package parser

import (
	"fmt"
	"strings"
	"zen/source"
	"zen/tokenizer"
)

type ParseError struct {
	At      tokenizer.SourceLocation
	message string
}

func CreateErrorWithMultipleLocations(at tokenizer.SourceLocation, mainMessage string, otherErrors []ParseError) (result ParseError) {
	var messageResult strings.Builder

	messageResult.WriteString(mainMessage)

	for _, parseError := range otherErrors {
		messageResult.WriteString("\n")
		if len(parseError.message) != 0 {
			messageResult.WriteString(parseError.message)
		}
		messageResult.WriteString(source.FormatLine(parseError.At.Source, parseError.At.At))
	}

	return ParseError{at, messageResult.String()}
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
