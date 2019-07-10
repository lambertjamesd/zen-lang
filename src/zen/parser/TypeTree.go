package parser

type TypeNode interface {
	CanAssignFrom(other TypeNode) bool
}

type UndefinedType struct {
}

func (undefinedType *UndefinedType) CanAssignFrom(other TypeNode) bool {
	return false
}

type IntegerType struct {
	BitCount int
	IsSigned bool
}

func (integerType *IntegerType) CanAssignFrom(other TypeNode) bool {
	var asNumber, ok = other.(*IntegerType)

	return ok && asNumber.BitCount == integerType.BitCount && asNumber.IsSigned == integerType.IsSigned
}
