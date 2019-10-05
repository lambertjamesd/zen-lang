package typechecker

import (
	"testing"
	"zen/parser"
	"zen/source"
	"zen/test"
)

func TestNumberType(t *testing.T) {
	var expr, errors = parser.ParseExpression(source.SourceFromString("10"))

	if len(errors) > 0 {
		t.Fatalf("Error parsing expression")
	}

	var typeChecker = CreateTypeChecker()

	expr.Accept(typeChecker)

	test.Assert(t, expr.GetType().CanAssignFrom(parser.NewIntegerType(32, true)), "Number literals evaluate to integers")
}
