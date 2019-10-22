package boundschecking

import (
	"sort"
	"strings"
	"zen/zmath"
)

const UNUSED = ^uint32(0)

type CheckResult struct {
	IsTrue bool
}

type productGroupEntry struct {
	index     uint32
	nodeArray *NormalizedNodeArray
}

type equationColumnInfo struct {
	sumGroup *SumGroup
	isZero   bool
}

type KnownConstraints struct {
	equationColumns        []equationColumnInfo
	productGroupRows       map[uint32]productGroupEntry
	equationTransformation *zmath.Matrixi64
	sumSpaceBoundingVolume ConvexNDVolume
}

func NewKnownConstraints() *KnownConstraints {
	var result = &KnownConstraints{
		make([]equationColumnInfo, 0),
		make(map[uint32]productGroupEntry, 0),
		zmath.NewMatrixi64(1, 1),
		ConvexNDVolume{},
	}

	result.equationTransformation.InitialzeIdentityi64()

	return result
}

func (constraints *KnownConstraints) negateColumnVector(columnVector *zmath.Matrixi64) *zmath.Matrixi64 {
	var result = columnVector.Scalei64(zmath.Ri64Fromi64(-1))
	result.SetEntryi64(0, 0, zmath.SubRi64(result.GetEntryi64(0, 0), zmath.Ri64_1()))
	return result
}

func (constraints *KnownConstraints) checkColumnVector(columnVector *zmath.Matrixi64) (result CheckResult, err error) {
	transformedVector, err := constraints.equationTransformation.Muli64(columnVector)

	if err != nil {
		return CheckResult{
			false,
		}, err
	}

	var isTrue = true

	for index := uint32(0); isTrue && index < transformedVector.Rows; index = index + 1 {
		entryValue := transformedVector.GetEntryi64(index, 0).Numerator

		var sumGroup *SumGroup = nil

		if index != 0 {
			sumGroup = constraints.equationColumns[index-1].sumGroup
		}

		if index == 0 {
			if entryValue < 0 {
				isTrue = false
			}
		} else if sumGroup == nil && entryValue != 0 {
			isTrue = false
		} else if entryValue < 0 && !constraints.equationColumns[index-1].isZero {
			isTrue = false
		}

		if !isTrue && constraints.sumSpaceBoundingVolume.GetMaybeSumGroupIndex(sumGroup) == -1 {
			return CheckResult{
				false,
			}, nil
		}
	}

	if !isTrue {
		sumGroups, values := constraints.extractVolumeValues(transformedVector)

		if !constraints.sumSpaceBoundingVolume.IsBounded(sumGroups, values) {
			return CheckResult{
				false,
			}, nil
		}
	}

	return CheckResult{
		true,
	}, nil
}

func (constraints *KnownConstraints) CheckSumGroup(equation *SumGroup) (result CheckResult, err error) {
	for _, productGroup := range equation.ProductGroups {
		_, ok := constraints.productGroupRows[productGroup.Values.uniqueID]
		if !ok {
			return CheckResult{
				false,
			}, nil
		}
	}

	columnVector := constraints.extractColumnVector(equation)
	return constraints.checkColumnVector(columnVector)
}

func (constraints *KnownConstraints) InsertSumGroup(equation *SumGroup) (isValid bool, err error) {
	columnVector := constraints.extractColumnVector(equation)

	var contradictionCheckVector = constraints.negateColumnVector(columnVector)

	contradictionCheck, err := constraints.checkColumnVector(contradictionCheckVector)

	if err != nil || contradictionCheck.IsTrue {
		return false, err
	}

	transformedVector, err := constraints.equationTransformation.Muli64(columnVector)

	if err != nil {
		return false, err
	}

	negativeCount := 0
	negativeIndex := uint32(0)
	positiveCount := 0
	positiveIndex := uint32(0)
	blankIndex := UNUSED

	for index := uint32(1); index < transformedVector.Rows; index = index + 1 {
		entryValue := transformedVector.GetEntryi64(index, 0).Numerator
		if blankIndex == UNUSED && constraints.equationColumns[index-1].sumGroup == nil && entryValue != 0 {
			blankIndex = index
		}

		if entryValue < 0 {
			negativeIndex = index
			negativeCount = negativeCount + 1
		} else if entryValue > 0 {
			positiveIndex = index
			positiveCount = positiveCount + 1
		}
	}

	if blankIndex != UNUSED {
		constraints.equationColumns[blankIndex-1] = equationColumnInfo{equation, false}
		constraints.rowReduceVector(transformedVector, blankIndex)
		return true, nil
	} else if negativeCount == 0 {
		return true, nil
	} else if negativeCount == 1 && positiveCount == 0 {
		constraints.equationColumns[negativeIndex-1].isZero = true
		return true, nil
	} else if positiveCount == 1 {
		constraints.equationColumns[positiveIndex-1] = equationColumnInfo{equation, false}
		constraints.rowReduceVector(transformedVector, positiveIndex)
		return true, nil
	} else {
		return constraints.insertSumGroupIntoNDimension(equation, transformedVector)
	}
}

func (constraints *KnownConstraints) extractVolumeValues(columnVector *zmath.Matrixi64) (sumGroups []*SumGroup, values []zmath.RationalNumberi64) {
	sumGroups = nil
	values = nil

	for row := uint32(0); row < columnVector.Rows; row = row + 1 {
		var vectorValue = columnVector.GetEntryi64(0, row)
		if !vectorValue.IsZero() {
			if row == 0 {
				sumGroups = append(sumGroups, nil)
			} else {
				sumGroups = append(sumGroups, constraints.equationColumns[row-1].sumGroup)
			}
			values = append(values, vectorValue)
		}
	}

	return sumGroups, values
}

func (constraints *KnownConstraints) insertSumGroupIntoNDimension(equation *SumGroup, columnVector *zmath.Matrixi64) (isValid bool, err error) {
	sumGroups, values := constraints.extractVolumeValues(columnVector)

	// TODO possibly pick replacement equation instead of always defaulting to new

	err = constraints.sumSpaceBoundingVolume.Extrude(sumGroups, values)

	if err != nil {
		return false, err
	}

	return true, nil
}

func (constraints *KnownConstraints) rowReduceVector(vector *zmath.Matrixi64, pivotIndex uint32) {
	pivotValue := vector.GetEntryi64(pivotIndex, 0)

	for index := uint32(0); index < vector.Rows; index = index + 1 {
		if pivotIndex != index {
			scalarValue := zmath.DivRi64(
				vector.GetEntryi64(index, 0),
				zmath.Muli64(pivotValue, -1),
			)

			constraints.equationTransformation.AddRowToRow(pivotIndex, index, scalarValue)
		}
	}

	constraints.equationTransformation.ScaleRowi64(pivotIndex, zmath.InvRi64(pivotValue))
}

func (constraints *KnownConstraints) extractColumnVector(equation *SumGroup) *zmath.Matrixi64 {
	for _, productGroup := range equation.ProductGroups {
		constraints.ensureProductGroup(productGroup.Values)
	}

	var result = zmath.NewMatrixi64(constraints.equationTransformation.Cols, 1)
	result.InitialzeIdentityi64()

	result.SetEntryi64(0, 0, zmath.Ri64Fromi64(equation.ConstantOffset))

	for _, productGroup := range equation.ProductGroups {
		var index = constraints.productGroupRows[productGroup.Values.uniqueID]
		result.SetEntryi64(index.index, 0, productGroup.ConstantScalar)
	}

	return result
}

func (constraints *KnownConstraints) ensureProductGroup(productGroup *NormalizedNodeArray) uint32 {
	var productGroupID = productGroup.uniqueID
	var result, ok = constraints.productGroupRows[productGroupID]

	if !ok {
		constraints.equationColumns = append(constraints.equationColumns, equationColumnInfo{nil, false})
		constraints.productGroupRows[productGroupID] = productGroupEntry{
			constraints.equationTransformation.Rows,
			productGroup,
		}
		constraints.equationTransformation.Resize(constraints.equationTransformation.Rows+1, constraints.equationTransformation.Cols+1)
	}

	return result.index
}

func (from *KnownConstraints) Copy() *KnownConstraints {
	var equationColumns = make([]equationColumnInfo, len(from.equationColumns))
	var productGroupRows = make(map[uint32]productGroupEntry)

	copy(equationColumns, from.equationColumns)

	for k, v := range from.productGroupRows {
		productGroupRows[k] = v
	}

	return &KnownConstraints{
		equationColumns,
		productGroupRows,
		from.equationTransformation.Copy(),
		from.sumSpaceBoundingVolume.Copy(),
	}
}

func (from *KnownConstraints) ToString() string {
	var result strings.Builder

	result.WriteString("Columns:\n")

	for _, sumGroup := range from.equationColumns {
		if sumGroup.sumGroup == nil {
			result.WriteString("nil\n")
		} else {
			result.WriteString(ToString(sumGroup.sumGroup))
			if sumGroup.isZero {
				result.WriteString(" == 0")
			}
			result.WriteString("\n")
		}
	}

	var rows []productGroupEntry = nil

	for _, entry := range from.productGroupRows {
		rows = append(rows, entry)
	}

	sort.Slice(rows, func(a, b int) bool {
		return rows[a].index < rows[b].index
	})

	result.WriteString("\nColumn Vector Product Types:\n")

	for _, entry := range rows {
		result.WriteString(ToString(entry.nodeArray))
		result.WriteString("\n")
	}

	result.WriteString("\nTransform\n")

	result.WriteString(from.equationTransformation.String())

	return result.String()
}
