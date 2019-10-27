package boundschecking

import "zen/parser"

func (orGroup *OrGroup) MapNodes(
	state *NormalizerState,
	mapping map[int]int,
	exprMapping map[uint32]parser.Expression,
	intoExprMapping map[uint32]parser.Expression,
) *OrGroup {
	var result []*AndGroup = make([]*AndGroup, len(orGroup.AndGroups))

	for index, productGroup := range orGroup.AndGroups {
		result[index] = productGroup.MapNodes(state, mapping, exprMapping, intoExprMapping)
	}

	return state.CreateOrGroup(result)
}

func (andGroup *AndGroup) MapNodes(
	state *NormalizerState,
	mapping map[int]int,
	exprMapping map[uint32]parser.Expression,
	intoExprMapping map[uint32]parser.Expression,
) *AndGroup {
	var result []*SumGroup = make([]*SumGroup, len(andGroup.SumGroups))

	for index, sumGroup := range andGroup.SumGroups {
		var migratedSum = sumGroup.MapNodes(state, mapping)
		intoExprMapping[migratedSum.uniqueId] = exprMapping[sumGroup.uniqueId]
		result[index] = migratedSum
	}

	return state.CreateAndGroup(result)
}

func (sumGroup *SumGroup) MapNodes(state *NormalizerState, mapping map[int]int) *SumGroup {
	var result []*ProductGroup = make([]*ProductGroup, len(sumGroup.ProductGroups))

	for index, productGroup := range sumGroup.ProductGroups {
		result[index] = productGroup.MapNodes(state, mapping)
	}

	return state.CreateSumGroup(result, sumGroup.ConstantOffset)
}

func (productGroup *ProductGroup) MapNodes(state *NormalizerState, mapping map[int]int) *ProductGroup {
	var result []NormalizedNode = make([]NormalizedNode, len(productGroup.Values.Array))

	for index, node := range productGroup.Values.Array {
		asVar, ok := node.(*VariableReference)
		if ok {
			result[index] = asVar.MapNodes(state, mapping)
			continue
		}

		asProp, ok := node.(*PropertyReference)
		if ok {
			result[index] = asProp.MapNodes(state, mapping)
			continue
		}
	}

	return state.CreateProductGroup(result, productGroup.ConstantScalar)
}

func (variableReference *VariableReference) MapNodes(state *NormalizerState, mapping map[int]int) *VariableReference {
	return state.CreateVariableReference(variableReference.Name, mapping[variableReference.valueId])
}

func (propertyReference *PropertyReference) MapNodes(state *NormalizerState, mapping map[int]int) *PropertyReference {
	return state.CreatePropertyReference(propertyReference.Left, propertyReference.Right, mapping[propertyReference.valueId])
}
