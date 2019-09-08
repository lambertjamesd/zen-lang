package boundschecking

import (
	"zen/zmath"
)

type NormalizedNode interface {
	GetHashCode() int32
	IsEqual(other NormalizedNode) bool
}

type VariableReference struct {
	name    string
	valueId int32
}

func (variableReference *VariableReference) GetHashCode() int32 {
	return variableReference.valueId
}

type ValueReference interface {
}

type ProductGroup struct {
	Values         []ValueReference
	uniqueID	   uint32
	ConstantScalar zmath.RationalNumberi64
}

type SumGroup struct {
	ProductGroups  []*ProductGroup
	ConstantOffset int64
}

type AndGroup struct {
	SumGroups []*SumGroup
}

type OrGroup struct {
	AndGroups []*AndGroup
}

type NormalizedEquation struct {
	Equation *OrGroup
}
