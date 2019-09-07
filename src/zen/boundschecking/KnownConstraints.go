package boundschecking

import zmath

type KnownConstraints struct {
	equationColumns []*SumGroup
	productGroupRows map[[]ValueReference, uint32]
	equationTransformation *Matrixi64
}

func NewKnownConstraints() *KnownConstraints {
	var result = &KnownConstraints{
		make([]NormalizedEquation, 0),
		make([][]ValueReference, 0),
		NewMatrixi64(1, 1),
	}

	result.equationTransformation.InitialzeIdentityi64()

	return result
}

func (constraints *KnownConstraints) InsertSumGroup(equation *SumGroup) (isValid bool, err error) {
	columnVector := constraints.extractColumnVector(equation)
	transformedVector := constraints.equationTransformation.Muli64(columnVector)

	if transformedVector.GetEntryi64(0, 0).Numerator < 0 {
		return false, errors.New("There is a contradiction in the rules")
	}

	hasNegativeValues := false
	blankIndex := -1

	for index := 1; index < transformedVector.Rows; index = index + 1 {
		entryValue := transformedVector.GetEntryi64(index, 0).Numerator
		if constraints.equationColumns[index - 1] == nil && entryValue != 0 {
			blankIndex = index
			break
		} else if entryValue < 0 {
			hasNegativeValues = true
		}
	}

	if blankIndex != -1 {
		atIndexValue := transformedVector.GetEntryi64(blankIndex, 0),

		for index := 0; index < transformedVector.Rows; index = index + 1 {
			if blankIndex != index {
				scalarValue := DivRi64(
					transformedVector.GetEntryi64(index, 0),
					Muli64(AbsRi64(atIndexValue), -1),
				)

				constraints.equationTransformation.AddRowToRow(blankIndex, index, scalarValue)
			}
		}

		constraints.equationTransformation.ScaleRowi64(blankIndex, InvRi64(atIndexValue))

		return true, nil
	} else if !hasNegativeValues {
		return true, nil
	} else {
		return false, errors.New("Not implemented yet")
	}
}

func (constraints *KnownConstraints) extractColumnVector(equation *SumGroup) *Matrixi64 {
	for _, productGroup := range equation.ProductGroups {
		constraints.ensureProductGroup(productGroup.Values)
	}

	var result = NewMatrixi64(constraints.equationTransformation.Cols, 1)

	result.SetEntryi64(0, 0, Ri64Fromi64(equation.ConstantOffset))
	
	for _, productGroup := range equation.ProductGroups {
		var index = constraints.productGroupRows[productGroup.Values]
		result.SetEntryi64(index, 0, Ri64Fromi64(productGroup.ConstantScalar))
	}

	return result
}

func (constraints *KnownConstraints) ensureProductGroup(productGroup []ValueReference) uint32 {
	var result = constraints.productGroupRows[productGroup]

	if result = 0 {
		result.equationColumns = append(result.equationColumns, nil)
		result.productGroupRows[productGroup] = constraints.equationTransformation.Rows
		result.equationTransformation.Resize(result.equationTransformation.Rows + 1, result.equationTransformation.Cols + 1)
	}

	return result
}

func (from *KnownConstraints) Copy() *KnownConstraints {
	var equationColumns = make([]*SumGroups, len(from.equationColumns))
	var productGroupRows = make(map[[]ValueReference, uint32])

	copy(equationColumns, from.equationColumns)

	for k,v := range from.productGroupRows {
		productGroupRows[k] = v
	}

	return &KnownConstraints{
		equationColumns,
		productGroupRows,
		from.equationTransformation.Copy(),
	};
}