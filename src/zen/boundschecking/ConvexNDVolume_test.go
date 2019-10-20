package boundschecking

import "testing"

func TestAddDimension(t *testing.T) {
	var nodeState = NewNormalizerState()

	nodeState.UseIdentifierMapping("a", 1)
	nodeState.UseIdentifierMapping("b", 2)
	nodeState.UseIdentifierMapping("c", 3)

	var ndVolume ConvexNDVolume

	ndVolume.extendDimension(nodeState.stringToSumGroup(t, "a"))
	ndVolume.extendDimension(nodeState.stringToSumGroup(t, "b"))
	ndVolume.extendDimension(nodeState.stringToSumGroup(t, "c"))

	t.Error(ndVolume.ToString())
}
