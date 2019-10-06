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
		result, _ := nodeState.NormalizeToSumGroup(a)

		if result == nil {
			t.Errorf("Failed to normalize %s", source)
		}

		return result
	}
}

func (nodeState *NormalizerState) stringToOrGroup(t *testing.T, source string) *OrGroup {
	a, ok := parser.ParseTest(source)

	if !ok {
		t.Errorf("Failed to parse %s", source)
		return nil
	} else {
		result := nodeState.NormalizeToOrGroup(a)

		if result == nil {
			t.Errorf("Failed to normalize %s", source)
		}

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

	sumGroupA, err := nodeState.NormalizeToSumGroup(a)

	if err != nil {
		t.Error("Error normalizing")
		return
	}

	sumGroupB, err := nodeState.NormalizeToSumGroup(b)

	if err != nil {
		t.Error("Error normalizing")
		return
	}

	var addAB = nodeState.addSumGroups(sumGroupA, sumGroupB, 0)

	addABExpr, ok := parser.ParseTest("a + b")

	if !ok {
		t.Error("Failed to parse expression")
		return
	}

	addAB2, err := nodeState.NormalizeToSumGroup(addABExpr)

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

func TestRemoveZeros(t *testing.T) {
	var nodeState = NewNormalizerState()

	nodeState.identfierSourceMapping["a"] = 1
	nodeState.identfierSourceMapping["b"] = 2
	nodeState.identfierSourceMapping["c"] = 3

	test.Assert(t, nodeState.stringToSumGroup(t, "a - a") == nodeState.stringToSumGroup(t, "0"), "a - a")
	test.Assert(t, nodeState.stringToSumGroup(t, "a*b - b*a") == nodeState.stringToSumGroup(t, "0"), "a*b - b*a")
	test.Assert(t, nodeState.stringToSumGroup(t, "a*0 - b*a") == nodeState.stringToSumGroup(t, "-a*b"), "a*0 - b*a")
}

func TestOrGroups(t *testing.T) {
	var nodeState = NewNormalizerState()

	nodeState.identfierSourceMapping["a"] = 1
	nodeState.identfierSourceMapping["b"] = 2
	nodeState.identfierSourceMapping["c"] = 3

	test.Assert(t, nodeState.stringToOrGroup(t, "a > b") == nodeState.stringToOrGroup(t, "b < a"), "a > b")
	test.Assert(t, nodeState.stringToOrGroup(t, "a > b") == nodeState.stringToOrGroup(t, "a - b > 0"), "a > b, a - b > 0")
	test.Assert(t, nodeState.stringToOrGroup(t, "a >= b") == nodeState.stringToOrGroup(t, "a > b - 1"), "a >= b")
	test.Assert(t, nodeState.stringToOrGroup(t, "a >= b") == nodeState.stringToOrGroup(t, "a + 1 > b"), "a >= b")
	test.Assert(t, nodeState.stringToOrGroup(t, "a == b") == nodeState.stringToOrGroup(t, "b == a"), "a == b")
	test.Assert(t, nodeState.stringToOrGroup(t, "a == b") == nodeState.stringToOrGroup(t, "a >= b && b >= a"), "a == b")
	test.Assert(t, nodeState.stringToOrGroup(t, "a != b") == nodeState.stringToOrGroup(t, "b != a"), "a != b")
	test.Assert(t, nodeState.stringToOrGroup(t, "a != b") == nodeState.stringToOrGroup(t, "a > b || a < b"), "a != b")
}
