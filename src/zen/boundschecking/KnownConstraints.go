package boundschecking

import (
	"errors"
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

type KnownConstraints struct {
	equationColumns        []*SumGroup
	productGroupRows       map[uint32]productGroupEntry
	equationTransformation *zmath.Matrixi64
}

func NewKnownConstraints() *KnownConstraints {
	var result = &KnownConstraints{
		make([]*SumGroup, 0),
		make(map[uint32]productGroupEntry, 0),
		zmath.NewMatrixi64(1, 1),
	}

	result.equationTransformation.InitialzeIdentityi64()

	return result
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
	transformedVector, err := constraints.equationTransformation.Muli64(columnVector)

	if err != nil {
		return CheckResult{
			false,
		}, err
	}

	if transformedVector.GetEntryi64(0, 0).Numerator < 0 {
		return CheckResult{
			false,
		}, nil
	}

	for index := uint32(1); index < transformedVector.Rows; index = index + 1 {
		entryValue := transformedVector.GetEntryi64(index, 0).Numerator
		if constraints.equationColumns[index-1] == nil && entryValue != 0 {
			return CheckResult{
				false,
			}, nil
		} else if entryValue < 0 {
			return CheckResult{
				false,
			}, nil
		}
	}

	return CheckResult{
		true,
	}, nil
}

func (constraints *KnownConstraints) InsertSumGroup(equation *SumGroup) (isValid bool, err error) {
	columnVector := constraints.extractColumnVector(equation)
	transformedVector, err := constraints.equationTransformation.Muli64(columnVector)

	if err != nil {
		return false, err
	}

	if transformedVector.GetEntryi64(0, 0).Numerator < 0 {
		return false, nil
	}

	negativeCount := 0
	positiveCount := 0
	blankIndex := UNUSED

	for index := uint32(1); index < transformedVector.Rows; index = index + 1 {
		entryValue := transformedVector.GetEntryi64(index, 0).Numerator
		if constraints.equationColumns[index-1] == nil && entryValue != 0 {
			blankIndex = index
			break
		} else if entryValue < 0 {
			negativeCount = negativeCount + 1
		} else if entryValue > 0 {
			positiveCount = positiveCount + 1
		}
	}

	if blankIndex != UNUSED {
		constraints.equationColumns[blankIndex-1] = equation
		constraints.rowReduceVector(transformedVector, blankIndex)
		return true, nil
	} else if negativeCount == 0 {
		return true, nil
	} else if positiveCount == 1 {
		positiveIndex := uint32(0)

		for index := uint32(1); index < transformedVector.Rows; index = index + 1 {
			entryValue := transformedVector.GetEntryi64(index, 0).Numerator

			if entryValue > 0 {
				positiveIndex = index
			}
		}

		constraints.equationColumns[positiveIndex-1] = equation
		constraints.rowReduceVector(transformedVector, positiveIndex)

		return true, nil
	} else {
		return false, errors.New("Not implemented yet")
	}
}

func (constraints *KnownConstraints) rowReduceVector(vector *zmath.Matrixi64, pivotIndex uint32) {
	pivotValue := vector.GetEntryi64(pivotIndex, 0)

	for index := uint32(0); index < vector.Rows; index = index + 1 {
		if pivotIndex != index {
			scalarValue := zmath.DivRi64(
				vector.GetEntryi64(index, 0),
				zmath.Muli64(zmath.AbsRi64(pivotValue), -1),
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
		constraints.equationColumns = append(constraints.equationColumns, nil)
		constraints.productGroupRows[productGroupID] = productGroupEntry{
			constraints.equationTransformation.Rows,
			productGroup,
		}
		constraints.equationTransformation.Resize(constraints.equationTransformation.Rows+1, constraints.equationTransformation.Cols+1)
	}

	return result.index
}

func (from *KnownConstraints) Copy() *KnownConstraints {
	var equationColumns = make([]*SumGroup, len(from.equationColumns))
	var productGroupRows = make(map[uint32]productGroupEntry)

	copy(equationColumns, from.equationColumns)

	for k, v := range from.productGroupRows {
		productGroupRows[k] = v
	}

	return &KnownConstraints{
		equationColumns,
		productGroupRows,
		from.equationTransformation.Copy(),
	}
}

// type KnownConstraints struct {
// 	equationColumns        []*SumGroup
// 	productGroupRows       map[uint32]uint32
// 	equationTransformation *zmath.Matrixi64
// }

func (from *KnownConstraints) ToString() string {
	var result strings.Builder

	result.WriteString("Columns:\n")

	for _, sumGroup := range from.equationColumns {
		if sumGroup == nil {
			result.WriteString("nil\n")
		} else {
			result.WriteString(ToString(sumGroup))
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
