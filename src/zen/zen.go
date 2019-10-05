package main

import (
	"log"
	"zen/parser"
	"zen/source"
	"zen/typechecker"
)

func main() {
	source, err := source.SourceFromFile("../../test/Min.zen")

	if err != nil {
		log.Fatalf("Error loading source %s", err)
	}

	var parseResult, errors = parser.Parse(source)

	if len(errors) > 0 {
		for _, element := range errors {
			log.Println(parser.FormatError(element))
		}
		log.Fatalf("Failed with %d errors", len(errors))
	} else {
		var typecheckErrors = typechecker.CheckTypes(parseResult)

		if len(typecheckErrors) > 0 {
			for _, element := range typecheckErrors {
				log.Println(parser.FormatError(element))
			}
		}
	}

	if parseResult == nil {
		log.Fatalf("Failed with unkown error")
	}

	log.Print("Success")
}
