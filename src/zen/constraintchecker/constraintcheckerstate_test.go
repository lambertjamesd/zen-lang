package constraintchecker

import (
	"testing"
	"zen/boundschecking"
	"zen/parser"
	"zen/test"
)

func stringToOrGroup(t *testing.T, nodeState *boundschecking.NormalizerState, source string) *boundschecking.OrGroup {
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

func TestMinSimulation(t *testing.T) {
	var nodeState = boundschecking.NewNormalizerState()
	var checkerState = NewConstraintCheckerState()
	var elseState = checkerState.Copy()

	nodeState.UseIdentifierMapping("a", 1)
	nodeState.UseIdentifierMapping("b", 2)
	nodeState.UseIdentifierMapping("result", 3)

	checkerState.addRules(stringToOrGroup(t, nodeState, "a < b"))
	checkerState.addRules(stringToOrGroup(t, nodeState, "a == result"))
	checkResult, err := checkerState.checkOrGroup(stringToOrGroup(t, nodeState, "result <= a && result <= b"))

	if err != nil {
		t.Error(err.Error())
	}

	test.Assert(t, checkResult, "A branch should to true")

	//elseState.addRules(nodeState.NotOrGroup(stringToOrGroup(t, nodeState, "a < b")))
	t.Log(elseState.knownConstraints[0].ToString())
	t.Log(boundschecking.ToString(stringToOrGroup(t, nodeState, "b <= a")))
	elseState.addRules(stringToOrGroup(t, nodeState, "b <= a"))
	t.Log(elseState.knownConstraints[0].ToString())
	elseState.addRules(stringToOrGroup(t, nodeState, "b == result"))
	t.Log(elseState.knownConstraints[0].ToString())
	checkResult, err = elseState.checkOrGroup(stringToOrGroup(t, nodeState, "result <= a && result <= b"))

	if err != nil {
		t.Error(err.Error())
	}

	t.Log(elseState.knownConstraints[0].ToString())

	test.Assert(t, checkResult, "B branch should to true")
}
