package constraintchecker

import (
	"testing"
	"zen/boundschecking"
	"zen/parser"
	"zen/test"
	"zen/typechecker"
)

func ParseType(t *testing.T, source string) parser.TypeExpression {
	typeExpression, ok := parser.ParseTypeTest(source)

	if !ok {
		t.Error("Could not parse '" + source + "'")
	}

	var typeCheckErrors = typechecker.CheckTypes(typeExpression)

	if len(typeCheckErrors) != 0 {
		for _, err := range typeCheckErrors {
			t.Logf("%s\n", parser.FormatError(err))
		}
		t.Error("Unexpected errors")
	}

	return typeExpression
}

func TestContradictionsChecking(t *testing.T) {
	var nodeState = boundschecking.NewNormalizerState()
	var typeConstraintDiffer = NewTypeConstraintDifferCache(nodeState)

	var typeWithContradiction = ParseType(t, "i32 where self > 0 && self < 0")
	typeInfo, err := typeConstraintDiffer.GetConstraintsForType(typeWithContradiction.GetType())
	if err != nil {
		t.Error(err)
	}
	test.Assert(t, len(typeInfo.contradictions) == 1, "Should have a contradiction")

	typeWithContradiction = ParseType(t, "[a: i32, b: i32, c: i32] where a < b && b < c && c < a")
	typeInfo, err = typeConstraintDiffer.GetConstraintsForType(typeWithContradiction.GetType())
	if err != nil {
		t.Error(err)
	}
	test.Assert(t, len(typeInfo.contradictions) == 1, "Should have a contradiction")
}

func TestTypeDiff(t *testing.T) {
	var nodeState = boundschecking.NewNormalizerState()
	var typeConstraintDiffer = NewTypeConstraintDifferCache(nodeState)

	var typeA = ParseType(t, "i32 where self >= 0 && self <= 10")
	var typeB = ParseType(t, "i32")

	typeDiff, err := typeConstraintDiffer.GetTypeDiff(typeB.GetType(), typeA.GetType())

	if err != nil {
		t.Error(err)
	}

	test.Assert(t, len(typeDiff.additionalConstrants.AndGroups) == 2, "Should require constraints to go one way")
}
