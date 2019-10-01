package parser

type TypeNodeType int

const (
	UndefinedNodeType TypeNodeType = iota
	IntegerNodeType
	BooleanNodeType
)

type TypeNode interface {
	CanAssignFrom(other TypeNode) bool
	GetNodeType() TypeNodeType
}

type UndefinedType struct {
}

func (undefinedType *UndefinedType) CanAssignFrom(other TypeNode) bool {
	return false
}

func (undefinedType *UndefinedType) GetNodeType() TypeNodeType {
	return UndefinedNodeType
}

type IntegerType struct {
	BitCount int
	IsSigned bool
}

func (integerType *IntegerType) CanAssignFrom(other TypeNode) bool {
	var asNumber, ok = other.(*IntegerType)

	return ok && asNumber.BitCount == integerType.BitCount && asNumber.IsSigned == integerType.IsSigned
}

func (integerType *IntegerType) GetNodeType() TypeNodeType {
	return IntegerNodeType
}

type BooleanType struct{}

func (booleanType *BooleanType) CanAssignFrom(other TypeNode) bool {
	var _, ok = other.(*BooleanType)
	return ok
}

func (booleanType *BooleanType) GetNodeType() TypeNodeType {
	return BooleanNodeType
}
