package boundschecking

import (
	"testing"
	"zen/test"
)

func assertTrue(t *testing.T, nodeState *NormalizerState, constraints *KnownConstraints, toCheck string, expectedValue bool, assertMessage string) {
	checkResult, err := constraints.CheckSumGroup(nodeState.stringToSumGroup(t, toCheck))

	if err != nil {
		t.Error("Failed to do check " + err.Error() + " " + assertMessage)
		return
	}

	test.Assert(t, checkResult.IsTrue == expectedValue, assertMessage)
}

func TestSimpleChecks(t *testing.T) {
	var constraints = NewKnownConstraints()
	var nodeState = NewNormalizerState()

	nodeState.UseIdentifierMapping("a", 1)
	nodeState.UseIdentifierMapping("b", 2)
	nodeState.UseIdentifierMapping("c", 3)

	constraints.InsertSumGroup(nodeState.stringToSumGroup(t, "a"))

	assertTrue(t, nodeState, constraints, "a", true, "Trival check")

	constraints.InsertSumGroup(nodeState.stringToSumGroup(t, "b"))

	assertTrue(t, nodeState, constraints, "b", true, "Trival check b")
	assertTrue(t, nodeState, constraints, "a + b", true, "Trival check a + b")
	assertTrue(t, nodeState, constraints, "a + 1", true, "Trival check a + 1")
	assertTrue(t, nodeState, constraints, "a - 1", false, "Trival check a - 1")
}

func TransitiveChecks(t *testing.T) {
	var constraints = NewKnownConstraints()
	var nodeState = NewNormalizerState()

	nodeState.UseIdentifierMapping("a", 1)
	nodeState.UseIdentifierMapping("b", 2)
	nodeState.UseIdentifierMapping("c", 3)

	constraints.InsertSumGroup(nodeState.stringToSumGroup(t, "a - b"))
	constraints.InsertSumGroup(nodeState.stringToSumGroup(t, "b - c"))

	assertTrue(t, nodeState, constraints, "a - c", true, "Transitive checks")

	constraints = NewKnownConstraints()

	constraints.InsertSumGroup(nodeState.stringToSumGroup(t, "a + b"))
	constraints.InsertSumGroup(nodeState.stringToSumGroup(t, "-a"))
	assertTrue(t, nodeState, constraints, "b", true, "a + b >= 0 && a <= 0 => b >= 0")
}

func TestContradictions(t *testing.T) {
	var constraints = NewKnownConstraints()
	var nodeState = NewNormalizerState()

	nodeState.UseIdentifierMapping("a", 1)
	nodeState.UseIdentifierMapping("b", 2)
	nodeState.UseIdentifierMapping("c", 3)

	insertResult, err := constraints.InsertSumGroup(nodeState.stringToSumGroup(t, "a - 10")) // a >= 10
	test.Assert(t, insertResult, "The first insert should pass")

	insertResult, err = constraints.InsertSumGroup(nodeState.stringToSumGroup(t, "-a")) // a < 0
	test.Assert(t, !insertResult, "The contradiction insert should fail")
	test.Assert(t, err == nil, "The contradiction insert should not have failed with an error")
}
