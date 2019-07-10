package typechecker

import (
	"testing"
	"zen/parser"
	"zen/source"
)

func TestSymbolCollector(t *testing.T) {
	var file, errors = parser.Parse(source.SourceFromString(`
			func Testing [a: u32, b: u32] => [r: u32] {
				
			}
		`))

	if len(errors) > 0 {
		for _, err := range errors {
			t.Log(parser.FormatError(err))
		}
		t.Fatalf("Error parsing expression")
	}

	var symbols = CollectSymbols(file)

	t.Log(symbols)
}
