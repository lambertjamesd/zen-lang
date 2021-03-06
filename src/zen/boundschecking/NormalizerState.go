package boundschecking

import (
	"sort"
	"zen/parser"
	"zen/zmath"
)

type IdentifierSource struct {
	UniqueId int
}

type NormalizerState struct {
	nodeCache                *NodeCache
	currentUniqueID          uint32
	identfierSourceMapping   map[string]IdentifierSource
	currentExpressionMapping map[uint32]parser.Expression
}

func NewNormalizerState() *NormalizerState {
	return &NormalizerState{
		NewNodeCache(),
		uint32(0),
		make(map[string]IdentifierSource),
		nil,
	}
}

func (state *NormalizerState) recordExpressionMapping(sumGroup *SumGroup, expression parser.Expression) *SumGroup {
	if state.currentExpressionMapping != nil {
		state.currentExpressionMapping[sumGroup.GetUniqueId()] = expression
	}
	return sumGroup
}

func (state *NormalizerState) StartTrackingExpressionMapping() map[uint32]parser.Expression {
	state.currentExpressionMapping = make(map[uint32]parser.Expression)
	return state.currentExpressionMapping
}

func (state *NormalizerState) StopTrackingExpressionMapping() {
	state.currentExpressionMapping = nil
}

func (state *NormalizerState) UseIdentifierMapping(name string, id int) {
	state.identfierSourceMapping[name] = IdentifierSource{id}
}

func (state *NormalizerState) getNextUniqueId() uint32 {
	state.currentUniqueID = state.currentUniqueID + 1
	return state.currentUniqueID
}

func (state *NormalizerState) multiplyProductGroupByScalar(a *ProductGroup, b zmath.RationalNumberi64) *ProductGroup {
	if b.IsZero() {
		return nil
	} else if b.IsOne() {
		return a
	} else {
		return state.nodeCache.GetNodeSingleton(&ProductGroup{
			a.Values,
			zmath.MulRi64(a.ConstantScalar, b).SimplifyRi64(),
		}).(*ProductGroup)
	}
}

func (state *NormalizerState) multiplyNodeArrays(a *NormalizedNodeArray, b *NormalizedNodeArray) *NormalizedNodeArray {
	var result []NormalizedNode = nil

	var aIndex = 0
	var bIndex = 0

	for aIndex < len(a.Array) || bIndex < len(b.Array) {
		var compareResult = 0

		if aIndex == len(a.Array) {
			compareResult = 1
		} else if bIndex == len(b.Array) {
			compareResult = -1
		} else {
			compareResult = a.Array[aIndex].Compare(b.Array[bIndex])
		}

		if compareResult <= 0 {
			result = append(result, a.Array[aIndex])
			aIndex = aIndex + 1
		}

		if compareResult >= 0 {
			result = append(result, b.Array[bIndex])
			bIndex = bIndex + 1
		}
	}

	return state.nodeCache.GetNodeSingleton(&NormalizedNodeArray{
		result,
		state.getNextUniqueId(),
	}).(*NormalizedNodeArray)
}

func (state *NormalizerState) multiplyProdctGroups(a *ProductGroup, b *ProductGroup) *ProductGroup {
	return state.nodeCache.GetNodeSingleton(&ProductGroup{
		state.multiplyNodeArrays(a.Values, b.Values),
		zmath.MulRi64(a.ConstantScalar, b.ConstantScalar).SimplifyRi64(),
	}).(*ProductGroup)
}

func (state *NormalizerState) multiplySumGroupByProductGroup(sumGroup *SumGroup, productGroup *ProductGroup) []*ProductGroup {
	var result []*ProductGroup = nil

	var constScalarResult = state.multiplyProductGroupByScalar(productGroup, zmath.Ri64Fromi64(sumGroup.ConstantOffset))

	for _, node := range sumGroup.ProductGroups {
		var nodeMultiplyResult = state.multiplyProdctGroups(node, productGroup)

		if constScalarResult != nil {
			var compareResult = constScalarResult.Values.Compare(nodeMultiplyResult.Values)

			if compareResult < 0 {
				result = append(result, constScalarResult)
				constScalarResult = nil
				result = append(result, nodeMultiplyResult)
			} else if compareResult == 0 {
				result = append(result, state.nodeCache.GetNodeSingleton(&ProductGroup{
					nodeMultiplyResult.Values,
					zmath.MulRi64(constScalarResult.ConstantScalar, nodeMultiplyResult.ConstantScalar),
				}).(*ProductGroup))
				constScalarResult = nil
			} else {
				result = append(result, nodeMultiplyResult)
			}

			result = append(result)
		} else {
			result = append(result, nodeMultiplyResult)
		}
	}

	if constScalarResult != nil {
		result = append(result, constScalarResult)
	}

	return result
}

func (state *NormalizerState) addProductGroups(a []*ProductGroup, b []*ProductGroup) []*ProductGroup {
	if a == nil {
		return b
	} else if b == nil {
		return a
	}

	var result []*ProductGroup = nil

	var aIndex = 0
	var bIndex = 0

	for aIndex < len(a) || bIndex < len(b) {
		var compareResult = 0

		if aIndex == len(a) {
			compareResult = 1
		} else if bIndex == len(b) {
			compareResult = -1
		} else {
			compareResult = a[aIndex].Values.Compare(b[bIndex].Values)
		}

		if compareResult < 0 {
			result = append(result, a[aIndex])
			aIndex = aIndex + 1
		} else if compareResult > 0 {
			result = append(result, b[bIndex])
			bIndex = bIndex + 1
		} else {
			scalarResult := zmath.AddRi64(a[aIndex].ConstantScalar, b[bIndex].ConstantScalar)
			aIndex = aIndex + 1
			bIndex = bIndex + 1

			if !scalarResult.IsZero() {
				result = append(result, state.nodeCache.GetNodeSingleton(&ProductGroup{
					a[aIndex].Values,
					scalarResult,
				}).(*ProductGroup))
			}
		}
	}

	return result
}

func (state *NormalizerState) negateSumGroup(a *SumGroup) *SumGroup {
	var values []*ProductGroup = nil

	for _, group := range a.ProductGroups {
		values = append(values, state.nodeCache.GetNodeSingleton(&ProductGroup{
			group.Values,
			zmath.NegateRi64(group.ConstantScalar),
		}).(*ProductGroup))
	}

	return state.nodeCache.GetNodeSingleton(&SumGroup{
		values,
		-a.ConstantOffset,
		state.getNextUniqueId(),
	}).(*SumGroup)
}

func (state *NormalizerState) addSumGroups(a *SumGroup, b *SumGroup, extraOffset int64) *SumGroup {
	if a == nil {
		return b
	} else if b == nil {
		return a
	}

	var values []*ProductGroup = nil

	var aIndex = 0
	var bIndex = 0

	for aIndex < len(a.ProductGroups) && bIndex < len(b.ProductGroups) {
		if a.ProductGroups[aIndex].Values == b.ProductGroups[bIndex].Values {
			var scalar = zmath.AddRi64(a.ProductGroups[aIndex].ConstantScalar, b.ProductGroups[bIndex].ConstantScalar)

			if !scalar.IsZero() {
				values = append(values, state.nodeCache.GetNodeSingleton(&ProductGroup{
					a.ProductGroups[aIndex].Values,
					scalar,
				}).(*ProductGroup))
			}

			aIndex = aIndex + 1
			bIndex = bIndex + 1
		} else {
			var compareResult = a.ProductGroups[aIndex].Compare(b.ProductGroups[bIndex])

			if compareResult < 0 {
				values = append(values, a.ProductGroups[aIndex])
				aIndex = aIndex + 1
			} else {
				values = append(values, b.ProductGroups[bIndex])
				bIndex = bIndex + 1
			}
		}
	}
	for aIndex < len(a.ProductGroups) {
		values = append(values, a.ProductGroups[aIndex])
		aIndex = aIndex + 1
	}

	for bIndex < len(b.ProductGroups) {
		values = append(values, b.ProductGroups[bIndex])
		bIndex = bIndex + 1
	}

	var result = &SumGroup{
		values,
		a.ConstantOffset + b.ConstantOffset + extraOffset,
		state.getNextUniqueId(),
	}

	return state.nodeCache.GetNodeSingleton(result).(*SumGroup)
}

func (state *NormalizerState) multiplySumGroups(a *SumGroup, b *SumGroup) *SumGroup {
	var values []*ProductGroup = nil

	for _, productGroup := range b.ProductGroups {
		var productResult = state.multiplySumGroupByProductGroup(a, productGroup)
		values = state.addProductGroups(values, productResult)
	}

	if b.ConstantOffset == 1 {
		values = state.addProductGroups(values, a.ProductGroups)
	} else if b.ConstantOffset != 0 {
		var scaledA []*ProductGroup = nil
		var scalarAsRational = zmath.Ri64Fromi64(b.ConstantOffset)

		if !scalarAsRational.IsZero() {
			for _, productGroup := range a.ProductGroups {
				scaledA = append(scaledA, state.multiplyProductGroupByScalar(productGroup, scalarAsRational))
			}
		}

		values = state.addProductGroups(values, scaledA)
	}

	return state.nodeCache.GetNodeSingleton(&SumGroup{
		values,
		a.ConstantOffset * b.ConstantOffset,
		state.getNextUniqueId(),
	}).(*SumGroup)
}

func ZipSumGroups(a []*SumGroup, b []*SumGroup) []*SumGroup {
	if len(a) == 0 {
		return b
	} else if len(b) == 0 {
		return a
	}
	var aIndex = 0
	var bIndex = 0

	var result []*SumGroup = nil

	for aIndex < len(a) || bIndex < len(b) {
		var compareResult = 0

		if aIndex == len(a) {
			compareResult = 1
		} else if bIndex == len(b) {
			compareResult = -1
		} else {
			compareResult = a[aIndex].Compare(b[bIndex])
		}

		if compareResult == 0 {
			result = append(result, a[aIndex])
			aIndex = aIndex + 1
			bIndex = bIndex + 1
		} else if compareResult < 0 {
			result = append(result, a[aIndex])
			aIndex = aIndex + 1
		} else {
			result = append(result, b[bIndex])
			bIndex = bIndex + 1
		}
	}

	return result
}

func (state *NormalizerState) combineAndGroups(a *AndGroup, b *AndGroup) *AndGroup {
	return state.nodeCache.GetNodeSingleton(&AndGroup{
		ZipSumGroups(a.SumGroups, b.SumGroups),
		state.getNextUniqueId(),
	}).(*AndGroup)
}

func (state *NormalizerState) combineOrGroups(a *OrGroup, b *OrGroup) *OrGroup {
	if a == nil {
		return b
	} else if b == nil {
		return a
	}

	var aIndex = 0
	var bIndex = 0

	var result []*AndGroup = nil

	for aIndex < len(a.AndGroups) || bIndex < len(b.AndGroups) {
		var compareResult = 0

		if aIndex == len(a.AndGroups) {
			compareResult = 1
		} else if bIndex == len(b.AndGroups) {
			compareResult = -1
		} else {
			compareResult = a.AndGroups[aIndex].Compare(b.AndGroups[bIndex])
		}

		if compareResult == 0 {
			result = append(result, a.AndGroups[aIndex])
			aIndex = aIndex + 1
			bIndex = bIndex + 1
		} else if compareResult < 0 {
			result = append(result, a.AndGroups[aIndex])
			aIndex = aIndex + 1
		} else {
			result = append(result, b.AndGroups[bIndex])
			bIndex = bIndex + 1
		}
	}

	return state.nodeCache.GetNodeSingleton(&OrGroup{
		result,
	}).(*OrGroup)
}

func (state *NormalizerState) combineOrWithAndGroup(orGroup *OrGroup, andGroup *AndGroup) *OrGroup {
	var result []*AndGroup = nil

	for _, subGroup := range orGroup.AndGroups {
		result = append(result, state.combineAndGroups(subGroup, andGroup))
	}

	return &OrGroup{
		result,
	}
}

func (state *NormalizerState) combineOrGroupsWithAnd(a *OrGroup, b *OrGroup) *OrGroup {
	var result *OrGroup = nil

	for _, subGroup := range a.AndGroups {
		result = state.combineOrWithAndGroup(b, subGroup)
	}

	return state.nodeCache.GetNodeSingleton(result).(*OrGroup)
}

func (state *NormalizerState) notSumGroup(sumGroup *SumGroup) *SumGroup {
	var result = state.negateSumGroup(sumGroup)

	return state.nodeCache.GetNodeSingleton(&SumGroup{
		result.ProductGroups,
		result.ConstantOffset - 1,
		state.getNextUniqueId(),
	}).(*SumGroup)
}

func (state *NormalizerState) notAndGroup(andGroup *AndGroup) *OrGroup {
	var result []*AndGroup = nil

	for _, sumGroup := range andGroup.SumGroups {
		var notSumGroup = state.notSumGroup(sumGroup)
		result = append(result, state.nodeCache.GetNodeSingleton(&AndGroup{
			[]*SumGroup{notSumGroup},
			state.getNextUniqueId(),
		}).(*AndGroup))
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Compare(result[j]) < 0
	})

	return state.nodeCache.GetNodeSingleton(&OrGroup{
		result,
	}).(*OrGroup)
}

func (state *NormalizerState) NotOrGroup(orGroup *OrGroup) *OrGroup {
	var result *OrGroup = nil

	for _, andGroup := range orGroup.AndGroups {
		result = state.combineOrGroups(result, state.notAndGroup(andGroup))
	}

	return result
}

func (state *NormalizerState) CreateEquality(sumGroups *SumGroup, id NormalizedNode) []*SumGroup {
	var halfGroup = state.addSumGroups(sumGroups, state.negateSumGroup(state.sumGroupFromNode(id)), int64(0))
	return []*SumGroup{
		halfGroup,
		state.negateSumGroup(halfGroup),
	}
}

func (state *NormalizerState) CreateNormalizedNodeArray(nodes []NormalizedNode) *NormalizedNodeArray {
	return state.nodeCache.GetNodeSingleton(&NormalizedNodeArray{
		nodes,
		state.getNextUniqueId(),
	}).(*NormalizedNodeArray)
}

func (state *NormalizerState) CreateProductGroup(nodes []NormalizedNode, constantScalar zmath.RationalNumberi64) *ProductGroup {
	return state.nodeCache.GetNodeSingleton(&ProductGroup{
		state.CreateNormalizedNodeArray(nodes),
		constantScalar,
	}).(*ProductGroup)
}

func (state *NormalizerState) CreateSumGroup(productGroups []*ProductGroup, constantOffset int64) *SumGroup {
	return state.nodeCache.GetNodeSingleton(&SumGroup{
		productGroups,
		constantOffset,
		state.getNextUniqueId(),
	}).(*SumGroup)
}

func (state *NormalizerState) CreateAndGroup(sumGroups []*SumGroup) *AndGroup {
	if len(sumGroups) == 0 {
		return nil
	}

	return state.nodeCache.GetNodeSingleton(&AndGroup{
		sumGroups,
		state.getNextUniqueId(),
	}).(*AndGroup)
}

func (state *NormalizerState) CreateOrGroup(andGroups []*AndGroup) *OrGroup {
	if len(andGroups) == 0 {
		return nil
	}

	return state.nodeCache.GetNodeSingleton(&OrGroup{
		andGroups,
	}).(*OrGroup)
}

func (state *NormalizerState) CreateVariableReference(name string, at int) *VariableReference {
	return state.nodeCache.GetNodeSingleton(&VariableReference{
		name,
		at,
	}).(*VariableReference)
}

func (state *NormalizerState) CreatePropertyReference(left NormalizedNode, right string, at int) *PropertyReference {
	return state.nodeCache.GetNodeSingleton(&PropertyReference{
		left,
		right,
		at,
	}).(*PropertyReference)
}
