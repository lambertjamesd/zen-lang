package boundschecking

import (
	"errors"
	"fmt"
	"strings"
	"zen/datastructures"
	"zen/zmath"
)

type boundsEdge struct {
	basisIndices *datastructures.BitSet32
	from         *boundsFace
	to           *boundsFace
}

type boundsFace struct {
	normal       *zmath.Matrixi64
	basisIndices *datastructures.BitSet32
	edges        []boundsEdge
}

type ConvexNDVolume struct {
	basisVectors   []*zmath.Matrixi64
	faces          []*boundsFace
	axisToSumGroup []*SumGroup
}

func (face *boundsFace) updateEdge(oldConnection *boundsFace, newConnection *boundsFace) {
	for edgeIndex, _ := range face.edges {
		var edge = &face.edges[edgeIndex]
		if edge.to == oldConnection {
			edge.to = newConnection
			edge.basisIndices.Clear()
			edge.basisIndices.Union(face.basisIndices)
			edge.basisIndices.Intersection(newConnection.basisIndices)
			return
		}
	}

	var basisIndices = face.basisIndices.Copy()
	basisIndices.Intersection(newConnection.basisIndices)

	face.edges = append(face.edges, boundsEdge{
		basisIndices,
		face,
		newConnection,
	})
}

func (face *boundsFace) getEdgeIndex(to *boundsFace) int {
	for edgeIndex, edge := range face.edges {
		if edge.to == to {
			return edgeIndex
		}
	}

	return -1
}

func (volume *ConvexNDVolume) getDimensionCount() uint32 {
	return uint32(len(volume.axisToSumGroup))
}

func (volume *ConvexNDVolume) getAllBasisIndices(basisCount uint32) *datastructures.BitSet32 {
	var resultMapping datastructures.BitSet32

	for i := uint32(0); i < basisCount; i = i + 1 {
		resultMapping.AddToSet(i)
	}

	return &resultMapping
}

func (volume *ConvexNDVolume) getEdgesToNewDimension(nextEdgeIndex uint32) []boundsEdge {
	var edges []boundsEdge = nil

	for _, face := range volume.faces {
		edges = append(edges, boundsEdge{
			face.basisIndices.Copy(),
			nil, // assigned later
			face,
		})

		face.basisIndices.AddToSet(nextEdgeIndex)

		for _, edge := range face.edges {
			edge.basisIndices.AddToSet(nextEdgeIndex)
		}

		face.normal.Resize(face.normal.Rows+1, 1)
	}

	return edges
}

func (volume *ConvexNDVolume) extendDimension(sumGroup *SumGroup) {
	volume.axisToSumGroup = append(volume.axisToSumGroup, sumGroup)

	var newBasisVector = zmath.NewMatrixi64(uint32(len(volume.axisToSumGroup)), 1)

	for i := 0; i < len(volume.axisToSumGroup); i = i + 1 {
		if i == len(volume.axisToSumGroup)-1 {
			newBasisVector.SetEntryi64(uint32(i), 0, zmath.Ri64_1())
		} else {
			newBasisVector.SetEntryi64(uint32(i), 0, zmath.Ri64_0())
		}
	}

	var nextBasisIndex = uint32(len(volume.basisVectors))

	var result = &boundsFace{
		newBasisVector,
		volume.getAllBasisIndices(nextBasisIndex),
		volume.getEdgesToNewDimension(nextBasisIndex),
	}

	volume.basisVectors = append(volume.basisVectors, newBasisVector)

	for edgeIndex, _ := range result.edges {
		var edge = &result.edges[edgeIndex]
		edge.from = result
		edge.to.updateEdge(result, result)
	}

	volume.faces = append(volume.faces, result)
}

func (volume *ConvexNDVolume) getMaybeSumGroupIndex(sumGroup *SumGroup) int {
	for index, otherGroup := range volume.axisToSumGroup {
		if otherGroup == sumGroup {
			return index
		}
	}

	return -1
}

func (volume *ConvexNDVolume) ensureSumGroupExists(sumGroup *SumGroup) int {
	index := volume.getMaybeSumGroupIndex(sumGroup)

	if index == -1 {
		volume.extendDimension(sumGroup)

		return len(volume.axisToSumGroup) - 1
	} else {
		return index
	}
}

func (volume *ConvexNDVolume) extractVector(sumGroup []*SumGroup, value []zmath.RationalNumberi64) *zmath.Matrixi64 {
	var result = zmath.NewMatrixi64(uint32(len(volume.axisToSumGroup)), 1)

	for subGroupIndex, sumGroup := range sumGroup {
		var matrixIndex = volume.getMaybeSumGroupIndex(sumGroup)
		if matrixIndex != -1 {
			result.SetEntryi64(uint32(matrixIndex), 0, value[subGroupIndex])
		} else if value[subGroupIndex].Numerator < 0 {
			return nil
		}
	}

	return result
}

func (volume *ConvexNDVolume) faceFromBitSet(bitSet *datastructures.BitSet32) (result *boundsFace, err error) {
	var neededBasisCount = volume.getDimensionCount() - 1
	if bitSet.Size() < neededBasisCount {
		return nil, errors.New("Not enough basis vectors to create face")
	}

	var basisVectors []*zmath.Matrixi64 = nil
	bitSet.ForEach(func(value uint32) bool {
		basisVectors = append(basisVectors, volume.basisVectors[value])
		if uint32(len(basisVectors)) == neededBasisCount {
			return false
		} else {
			return true
		}
	})

	orthoVector, err := zmath.OrthogonalVector(basisVectors)

	if err != nil {
		return nil, err
	}

	if zmath.MatrixSum(orthoVector).Numerator < 0 {
		orthoVector = orthoVector.Scalei64(zmath.Ri64Fromi64(int64(-1)))
	}

	return &boundsFace{
		orthoVector,
		bitSet.Copy(),
		nil,
	}, nil
}

func (volume *ConvexNDVolume) Extrude(sumGroup []*SumGroup, value []zmath.RationalNumberi64) error {
	if volume.IsBounded(sumGroup, value) {
		return nil
	}

	for _, sumGroup := range sumGroup {
		volume.ensureSumGroupExists(sumGroup)
	}

	var extractedVector = volume.extractVector(sumGroup, value)

	var toRemove = make(map[uint32]bool)
	var alreadyAdded = make(map[uint32]*boundsFace)

	for _, face := range volume.faces {
		var dotResult = zmath.MatrixDot(face.normal, extractedVector)

		if dotResult.Numerator < 0 {
			toRemove[face.basisIndices.ToInt32()] = true
		} else {
			alreadyAdded[face.basisIndices.ToInt32()] = face
		}
	}

	var nextFaces []*boundsFace = nil
	var newFaces []*boundsFace = nil
	var nDimensions = volume.getDimensionCount()
	var newAxisIndex = uint32(len(volume.basisVectors))

	volume.basisVectors = append(volume.basisVectors, extractedVector)

	for _, face := range volume.faces {
		if toRemove[face.basisIndices.ToInt32()] {
			if nDimensions > 2 {
				for _, faceEdge := range face.edges {
					if !toRemove[faceEdge.to.basisIndices.ToInt32()] {
						faceEdge.basisIndices.ForEachSubSet(nDimensions-2, func(subSet datastructures.BitSet32) {
							var subSetCopy = subSet.Copy()
							subSetCopy.AddToSet(newAxisIndex)

							var newFace = alreadyAdded[subSetCopy.ToInt32()]

							if newFace == nil {
								newFace, _ = volume.faceFromBitSet(subSetCopy)
								nextFaces = append(nextFaces, newFace)
								newFaces = append(newFaces, newFace)
							}

							faceEdge.to.updateEdge(face, newFace)
							newFace.updateEdge(face, face)
						})
					}
				}
			}
		} else {
			nextFaces = append(nextFaces, face)
		}
	}

	for newFaceIndex, newFace := range newFaces {
		for otherFaceIndex := newFaceIndex + 1; otherFaceIndex < len(newFaces); otherFaceIndex = otherFaceIndex + 1 {
			var otherNewFace = newFaces[otherFaceIndex]
			var faceIntersection = *newFace.basisIndices
			faceIntersection.Intersection(otherNewFace.basisIndices)
			if faceIntersection.Size() >= nDimensions-2 {
				newFace.updateEdge(otherNewFace, otherNewFace)
				otherNewFace.updateEdge(newFace, newFace)
			}
		}
	}

	volume.faces = nextFaces

	return nil
}

func (volume *ConvexNDVolume) IsBounded(sumGroup []*SumGroup, value []zmath.RationalNumberi64) bool {
	var extractedVector = volume.extractVector(sumGroup, value)

	if extractedVector == nil {
		return false
	}

	for _, face := range volume.faces {
		var dotResult = zmath.MatrixDot(face.normal, extractedVector)

		if dotResult.Numerator < 0 {
			return false
		}
	}

	return true
}

func (volume *ConvexNDVolume) faceIndex(face *boundsFace) int {
	for result, otherFace := range volume.faces {
		if otherFace == face {
			return result
		}
	}

	return -1
}

func (volume *ConvexNDVolume) ToString() string {
	var stringBuilder strings.Builder

	stringBuilder.WriteString("Axis to Sum Group\n")

	for sumGroupIndex, sumGroup := range volume.axisToSumGroup {
		stringBuilder.WriteString(fmt.Sprintf("  %d: ", sumGroupIndex))
		sumGroup.ToString(&stringBuilder)
		stringBuilder.WriteString("\n")
	}

	stringBuilder.WriteString("\nBasis Vectors\n")

	for basisIndex, basisVector := range volume.basisVectors {
		stringBuilder.WriteString(fmt.Sprintf("%d:\n", basisIndex))
		basisVector.BuildString(&stringBuilder, "  ")
		stringBuilder.WriteString("\n")
	}

	stringBuilder.WriteString("\nFaces\n")

	for faceIndex, face := range volume.faces {
		stringBuilder.WriteString(fmt.Sprintf("  \n  Face Index: %d\n  Basis Vectors: %b\n", faceIndex, face.basisIndices.ToInt32()))
		stringBuilder.WriteString("  Normal\n")
		face.normal.BuildString(&stringBuilder, "    ")
		stringBuilder.WriteString("\n  Edges\n")

		for _, edge := range face.edges {
			stringBuilder.WriteString(fmt.Sprintf(
				"    Edge Basis: %b toFace: %d fromFace: %d\n",
				edge.basisIndices.ToInt32(),
				volume.faceIndex(edge.to),
				volume.faceIndex(edge.from),
			))
		}
	}

	return stringBuilder.String()
}
