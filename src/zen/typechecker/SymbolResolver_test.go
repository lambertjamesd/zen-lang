package typechecker

import (
	"testing"
	"zen/parser"
	"zen/source"
)

func TestResolveSymbols(t *testing.T) {
	var file, errors = parser.Parse(source.SourceFromString("10"))

	if len(errors) > 0 {
		t.Fatalf("Error parsing expression")
	}

	ResolveSymbols(file)
}
