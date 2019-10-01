package boundschecking

import (
	"testing"
	"zen/parser"
	"zen/test"
)

func (nodeState *NormalizerState) stringToSumGroup(t *testing.T, source string) *SumGroup {
	a, ok := parser.ParseTest(source)

	if !ok {
		t.Errorf("Failed to parse %s", source)
		return nil
	} else {
		result, _ := nodeState.normalizeToSumGroup(a)
		return result
	}
}

func TestCombiningGroups(t *testing.T) {
	var nodeState = NewNormalizerState()

	nodeState.identfierSourceMapping["a"] = 1
	nodeState.identfierSourceMapping["b"] = 2

	a, ok := parser.ParseTest("a")

	if !ok {
		t.Error("Failed to parse expression")
		return
	}

	b, ok := parser.ParseTest("b")

	if !ok {
		t.Error("Failed to parse expression")
		return
	}

	sumGroupA, err := nodeState.normalizeToSumGroup(a)

	if err != nil {
		t.Error("Error normalizing")
		return
	}

	sumGroupB, err := nodeState.normalizeToSumGroup(b)

	if err != nil {
		t.Error("Error normalizing")
		return
	}

	var addAB = nodeState.addSumGroups(sumGroupA, sumGroupB)

	addABExpr, ok := parser.ParseTest("a + b")

	if !ok {
		t.Error("Failed to parse expression")
		return
	}

	addAB2, err := nodeState.normalizeToSumGroup(addABExpr)

	if err != nil {
		t.Error("Error normalizing")
		return
	}

	t.Log(ToString(addAB))
	t.Log(ToString(addAB2))

	if addAB != addAB2 {
		t.Error("Normalized expressions should match")
	}
}

func TestNormalizeEquality(t *testing.T) {
	var nodeState = NewNormalizerState()

	nodeState.identfierSourceMapping["a"] = 1
	nodeState.identfierSourceMapping["b"] = 2
	nodeState.identfierSourceMapping["c"] = 3

	test.Assert(t, nodeState.stringToSumGroup(t, "a + b") != nodeState.stringToSumGroup(t, "a + c"), "Not equal should not equal")
	test.Assert(t, nodeState.stringToSumGroup(t, "a + b") == nodeState.stringToSumGroup(t, "b + a"), "Sum order doesn't matter")
	test.Assert(t, nodeState.stringToSumGroup(t, "a*b + a*c") == nodeState.stringToSumGroup(t, "a*(b + c)"), "Multiplies distribute")
	test.Assert(t, nodeState.stringToSumGroup(t, "a - b") == nodeState.stringToSumGroup(t, "-b + a"), "Negation doesn't matter")
	test.Assert(t, nodeState.stringToSumGroup(t, "a - b") == nodeState.stringToSumGroup(t, "a + -b"), "Minus is same as negate")
	test.Assert(t, nodeState.stringToSumGroup(t, "a - b") == nodeState.stringToSumGroup(t, "a + -1 * b"), "Minus is same as negative multiply")
}
