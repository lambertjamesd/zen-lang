package boundschecking

import (
	"strconv"
	"strings"
	"zen/zmath"
)

func JoinHash(a int32, b int32) int32 {
	var result = int32(17)
	result = result*13 + a
	result = result*13 + b
	return result
}

type NormalizedNodeType = int32

const (
	VariableReferenceType NormalizedNodeType = iota
	ProductGroupType
	NormalizedNodeArrayType
	SumGroupType
	AndGroupType
	OrGroupType
)

type NormalizedNode interface {
	GetHashCode() int32
	Compare(other NormalizedNode) int32
	GetNormalizedType() NormalizedNodeType
	ToString(builder *strings.Builder)
}

type VariableReference struct {
	name    string
	valueId int32
}

func (variableReference *VariableReference) GetHashCode() int32 {
	return variableReference.valueId
}

func (variable *VariableReference) Compare(other NormalizedNode) int32 {
	otherAsVariable, ok := other.(*VariableReference)

	if ok {
		return variable.valueId - otherAsVariable.valueId
	} else {
		return variable.GetNormalizedType() - other.GetNormalizedType()
	}
}

func (variable *VariableReference) GetNormalizedType() NormalizedNodeType {
	return VariableReferenceType
}

func (variable *VariableReference) ToString(builder *strings.Builder) {
	builder.WriteString(variable.name)
	builder.WriteString("_")
	builder.WriteString(strconv.Itoa(int(variable.valueId)))
}

type NormalizedNodeArray struct {
	Array    []NormalizedNode
	uniqueID uint32
}

type ProductGroup struct {
	Values         *NormalizedNodeArray
	ConstantScalar zmath.RationalNumberi64
}

func (productGroup *ProductGroup) GetHashCode() int32 {
	var result = JoinHash(
		int32(productGroup.ConstantScalar.Numerator),
		int32(productGroup.ConstantScalar.Denominator),
	)

	return JoinHash(result, productGroup.Values.GetHashCode())
}

func (productGroup *ProductGroup) Compare(other NormalizedNode) int32 {
	otherAsGroup, ok := other.(*ProductGroup)

	if ok {
		if productGroup.Values != otherAsGroup.Values {
			var arrayCompareResult = productGroup.Values.Compare(otherAsGroup.Values)

			if arrayCompareResult != 0 {
				return arrayCompareResult
			}
		}

		return productGroup.ConstantScalar.Compare(otherAsGroup.ConstantScalar)
	} else {
		return productGroup.GetNormalizedType() - other.GetNormalizedType()
	}
}

func (productGroup *ProductGroup) GetNormalizedType() NormalizedNodeType {
	return ProductGroupType
}

func (productGroup *ProductGroup) ToString(builder *strings.Builder) {
	if !productGroup.ConstantScalar.IsOne() {
		builder.WriteString(productGroup.ConstantScalar.ToString())

		if len(productGroup.Values.Array) != 0 {
			builder.WriteString("*")
		}
	}

	productGroup.Values.ToString(builder)
}

func (productGroup *NormalizedNodeArray) GetHashCode() int32 {
	var result = int32(0)

	for _, value := range productGroup.Array {
		result = JoinHash(result, value.GetHashCode())
	}

	return result
}

func (productGroup *NormalizedNodeArray) Compare(other NormalizedNode) int32 {
	asArray, ok := other.(*NormalizedNodeArray)

	if ok {
		var lengthDiff = len(productGroup.Array) - len(asArray.Array)

		if lengthDiff != 0 {
			return int32(lengthDiff)
		}

		for index, node := range productGroup.Array {
			var compareResult = node.Compare(asArray.Array[index])

			if compareResult != 0 {
				return compareResult
			}
		}

		return 0
	} else {
		return productGroup.GetNormalizedType() - other.GetNormalizedType()
	}
}

func (productGroup *NormalizedNodeArray) GetNormalizedType() NormalizedNodeType {
	return NormalizedNodeArrayType
}

func (productGroup *NormalizedNodeArray) ToString(builder *strings.Builder) {
	for index, group := range productGroup.Array {
		if index != 0 {
			builder.WriteString("*")
		}

		group.ToString(builder)
	}
}

type SumGroup struct {
	ProductGroups  []*ProductGroup
	ConstantOffset int64
}

func (sumGroup *SumGroup) IsZero() bool {
	return len(sumGroup.ProductGroups) == 0 && sumGroup.ConstantOffset == int64(0)
}

func (sumGroup *SumGroup) GetHashCode() int32 {
	var result = int32(sumGroup.ConstantOffset)

	for _, node := range sumGroup.ProductGroups {
		result = JoinHash(result, node.GetHashCode())
	}

	return result
}

func (sumGroup *SumGroup) Compare(other NormalizedNode) int32 {
	otherAsSumGroup, ok := other.(*SumGroup)

	if ok {
		if sumGroup == otherAsSumGroup {
			return 0
		}

		if len(sumGroup.ProductGroups) != len(otherAsSumGroup.ProductGroups) {
			return (int32)(len(sumGroup.ProductGroups) - len(otherAsSumGroup.ProductGroups))
		}

		scalarCompare := sumGroup.ConstantOffset - otherAsSumGroup.ConstantOffset

		if scalarCompare != 0 {
			return int32(scalarCompare)
		}

		for index, node := range sumGroup.ProductGroups {
			var result = node.Compare(otherAsSumGroup.ProductGroups[index])

			if result != 0 {
				return result
			}
		}

		return 0
	} else {
		return sumGroup.GetNormalizedType() - other.GetNormalizedType()
	}
}

func (sumGroup *SumGroup) GetNormalizedType() NormalizedNodeType {
	return SumGroupType
}

func (sumGroup *SumGroup) ToString(builder *strings.Builder) {
	for index, group := range sumGroup.ProductGroups {
		if index != 0 {
			builder.WriteString(" + ")
		}

		group.ToString(builder)
	}

	if len(sumGroup.ProductGroups) > 0 && sumGroup.ConstantOffset != 0 {
		builder.WriteString(" + ")
	}

	if sumGroup.ConstantOffset != 0 || len(sumGroup.ProductGroups) == 0 {
		builder.WriteString(strconv.Itoa(int(sumGroup.ConstantOffset)))
	}
}

type AndGroup struct {
	SumGroups []*SumGroup
}

func (andGroup *AndGroup) GetHashCode() int32 {
	var result = int32(0)

	for _, node := range andGroup.SumGroups {
		result = JoinHash(result, node.GetHashCode())
	}

	return result
}

func (andGroup *AndGroup) Compare(other NormalizedNode) int32 {
	otherAsAndGroup, ok := other.(*AndGroup)

	if ok {
		if andGroup == otherAsAndGroup {
			return 0
		}

		if len(andGroup.SumGroups) != len(otherAsAndGroup.SumGroups) {
			return (int32)(len(andGroup.SumGroups) - len(otherAsAndGroup.SumGroups))
		}

		for index, node := range andGroup.SumGroups {
			var result = node.Compare(otherAsAndGroup.SumGroups[index])

			if result != 0 {
				return result
			}
		}

		return 0
	} else {
		return andGroup.GetNormalizedType() - other.GetNormalizedType()
	}
}

func (andGroup *AndGroup) GetNormalizedType() NormalizedNodeType {
	return AndGroupType
}

func (andGroup *AndGroup) ToString(builder *strings.Builder) {
	for index, group := range andGroup.SumGroups {
		if index != 0 {
			builder.WriteString(" && ")
		}

		group.ToString(builder)
	}
}

type OrGroup struct {
	AndGroups []*AndGroup
}

func (orGroup *OrGroup) GetHashCode() int32 {
	var result = int32(0)

	for _, node := range orGroup.AndGroups {
		result = JoinHash(result, node.GetHashCode())
	}

	return result
}

func (orGroup *OrGroup) Compare(other NormalizedNode) int32 {
	otherAsOrGroup, ok := other.(*OrGroup)

	if ok {
		if orGroup == otherAsOrGroup {
			return 0
		}

		if len(orGroup.AndGroups) != len(otherAsOrGroup.AndGroups) {
			return (int32)(len(orGroup.AndGroups) - len(otherAsOrGroup.AndGroups))
		}

		for index, node := range orGroup.AndGroups {
			var result = node.Compare(otherAsOrGroup.AndGroups[index])

			if result != 0 {
				return result
			}
		}

		return 0
	} else {
		return orGroup.GetNormalizedType() - other.GetNormalizedType()
	}
}

func (orGroup *OrGroup) GetNormalizedType() NormalizedNodeType {
	return OrGroupType
}

func (orGroup *OrGroup) ToString(builder *strings.Builder) {
	for index, group := range orGroup.AndGroups {
		if index != 0 {
			builder.WriteString(" || ")
		}

		group.ToString(builder)
	}
}

type NormalizedEquation struct {
	Equation *OrGroup
}

func ToString(normalizedForm NormalizedNode) string {
	var result strings.Builder
	normalizedForm.ToString(&result)
	return result.String()
}
