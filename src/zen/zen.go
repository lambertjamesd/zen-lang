package main

import (
	"log"
	"zen/constraintchecker"
	"zen/parser"
	"zen/source"
	"zen/typechecker"
)

func checkErrors(errors []parser.ParseError) bool {
	if len(errors) == 0 {
		return true
	} else {
		for _, element := range errors {
			log.Println(parser.FormatError(element))
		}

		return false
	}
}

func main() {
	source, err := source.SourceFromFile("../../test/Range.zen")

	if err != nil {
		log.Fatalf("Error loading source %s", err)
	}

	var parseResult, errors = parser.Parse(source)

	if checkErrors(errors) &&
		checkErrors(typechecker.CheckTypes(parseResult)) &&
		checkErrors(constraintchecker.CheckConstraints(parseResult)) {
		log.Print("Success")
	} else {
		log.Print("Fail")
	}
}
