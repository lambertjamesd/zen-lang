package boundschecking

import (
	"fmt"
	"testing"
	"zen/test"
	"zen/zmath"
)

func veriyEdgeFrom(t *testing.T, volume *ConvexNDVolume) {
	for _, face := range volume.faces {
		for _, edge := range face.edges {
			test.Assert(t, edge.from == face, "From face links back to face")
			test.Assert(t, edge.to != face, "Face should not link back to self")
			test.Assert(
				t,
				edge.basisIndices.ToInt32() == edge.from.basisIndices.ToInt32()&edge.to.basisIndices.ToInt32(),
				"Edge should have any common basis between faces",
			)
		}
	}
}

func verifyBasisOrthogonal(t *testing.T, volume *ConvexNDVolume) {
	for _, face := range volume.faces {
		face.basisIndices.ForEach(func(basisIndex uint32) bool {
			test.Assert(t, zmath.MatrixDot(face.normal, volume.basisVectors[basisIndex]).IsZero(), "Normal should be orthogonal to face basis")
			return true
		})
	}
}

func verifyEdges(t *testing.T, volume *ConvexNDVolume) {
	for faceIndex, face := range volume.faces {
		for otherFaceIndex := faceIndex + 1; otherFaceIndex < len(volume.faces); otherFaceIndex = otherFaceIndex + 1 {
			var edgeCheck = *face.basisIndices
			var otherFace = volume.faces[otherFaceIndex]
			edgeCheck.Intersection(otherFace.basisIndices)

			if edgeCheck.Size()+2 >= volume.getDimensionCount() {
				var toEdgeIndex = face.getEdgeIndex(otherFace)
				var fromEdgeIndex = otherFace.getEdgeIndex(face)
				test.Assert(t, toEdgeIndex != -1, fmt.Sprintf("face %d to %d missing edge", faceIndex, otherFaceIndex))
				test.Assert(t, fromEdgeIndex != -1, fmt.Sprintf("face %d to %d missing edge", otherFaceIndex, faceIndex))

				test.Assert(t,
					edgeCheck.ToInt32() == face.edges[toEdgeIndex].basisIndices.ToInt32(),
					fmt.Sprintf(
						"face %d edge %d has wrong basis expected %b got %b",
						faceIndex,
						toEdgeIndex,
						edgeCheck.ToInt32(),
						face.edges[toEdgeIndex].basisIndices.ToInt32(),
					),
				)

				test.Assert(t,
					edgeCheck.ToInt32() == otherFace.edges[fromEdgeIndex].basisIndices.ToInt32(),
					fmt.Sprintf(
						"face %d edge %d has wrong basis expected %b got %b",
						otherFaceIndex,
						fromEdgeIndex,
						edgeCheck.ToInt32(),
						otherFace.edges[fromEdgeIndex].basisIndices.ToInt32(),
					),
				)
			} else {
				var toEdgeIndex = face.getEdgeIndex(otherFace)
				var fromEdgeIndex = otherFace.getEdgeIndex(face)
				test.Assert(t, toEdgeIndex == -1, fmt.Sprintf("face %d to %d shouldn't have edge", faceIndex, otherFaceIndex))
				test.Assert(t, fromEdgeIndex == -1, fmt.Sprintf("face %d to %d shouldn't have edge", otherFaceIndex, faceIndex))
			}
		}
	}
}

func verifyVolumeIsCorrect(t *testing.T, volume *ConvexNDVolume) {
	veriyEdgeFrom(t, volume)
	verifyBasisOrthogonal(t, volume)
	verifyEdges(t, volume)
}

func TestAddDimension(t *testing.T) {
	var nodeState = NewNormalizerState()

	nodeState.UseIdentifierMapping("a", 1)
	nodeState.UseIdentifierMapping("b", 2)
	nodeState.UseIdentifierMapping("c", 3)

	var ndVolume ConvexNDVolume

	ndVolume.extendDimension(nodeState.stringToSumGroup(t, "a"))
	ndVolume.extendDimension(nodeState.stringToSumGroup(t, "b"))
	ndVolume.extendDimension(nodeState.stringToSumGroup(t, "c"))

	test.Assert(t, len(ndVolume.faces) == 3, "Correct face count")
	verifyVolumeIsCorrect(t, &ndVolume)

}
