package boundschecking

import zmath

type KnownConstraints struct {
	equationColumns []SumGroup
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

func (constraints *KnownConstraints) ExtractColumnVector(equation *SumGroup) *Matrixi64 {
	for _, productGroup := range equation.ProductGroups {
		constraints.EnsureProductGroup(productGroup.Values)
	}

	var result = NewMatrixi64(constraints.equationTransformation.Cols, 1)

	result.SetEntryi64(0, 0, Ri64Fromi64(equation.ConstantOffset))
	
	for _, productGroup := range equation.ProductGroups {
		var index = constraints.productGroupRows[productGroup.Values]
		result.SetEntryi64(index, 0, Ri64Fromi64(productGroup.ConstantScalar))
	}

	return result
}

func (constraints *KnownConstraints) EnsureProductGroup(productGroup []ValueReference) uint32 {
	var result = constraints.productGroupRows[productGroup]

	if result = 0 {
		result.productGroupRows[productGroup] = constraints.equationTransformation.Rows
		result.equationTransformation.Resize(result.equationTransformation.Rows + 1, result.equationTransformation.Cols + 1)
	}

	return result
}

func (from *KnownConstraints) Copy() *KnownConstraints {
	return &KnownConstraints{
		from.equationColumns,
		from.productGroupRows,
		from.equationTransformation.Copy(),
	};
}