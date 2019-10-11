package parser

type TypeNodeType int

const (
	UndefinedNodeType TypeNodeType = iota
	VoidNodeType
	IntegerNodeType
	BooleanNodeType
	StructureNodeType
	FunctionNodeType
)

type TypeNode interface {
	CanAssignFrom(other TypeNode) bool
	GetNodeType() TypeNodeType
	GetWhereExpression() Expression
	SetWhereExpression(expression Expression)
	UniqueId() int
}

const (
	UndefinedId = -1
	VoidId      = -2
)

var nextTypeId int = 0

func getNextTypeId() int {
	nextTypeId = nextTypeId + 1
	return nextTypeId
}

type UndefinedType struct {
}

func (undefinedType *UndefinedType) CanAssignFrom(other TypeNode) bool {
	return false
}

func (undefinedType *UndefinedType) GetNodeType() TypeNodeType {
	return UndefinedNodeType
}

func (undefinedType *UndefinedType) GetWhereExpression() Expression {
	return &VoidExpression{}
}
func (undefinedType *UndefinedType) SetWhereExpression(expression Expression) {
	// noop
}

func (undefinedType *UndefinedType) UniqueId() int {
	return UndefinedId
}

type VoidType struct {
}

func (undefinedType *VoidType) CanAssignFrom(other TypeNode) bool {
	_, ok := other.(*VoidType)
	return ok
}

func (undefinedType *VoidType) GetNodeType() TypeNodeType {
	return VoidNodeType
}

func (undefinedType *VoidType) GetWhereExpression() Expression {
	return &VoidExpression{}
}
func (undefinedType *VoidType) SetWhereExpression(expression Expression) {
	// noop
}

func (undefinedType *VoidType) UniqueId() int {
	return VoidId
}

type IntegerType struct {
	BitCount        int
	IsSigned        bool
	unqiueId        int
	whereExpression Expression
}

func NewIntegerType(bitCount int, isSigned bool) *IntegerType {
	return &IntegerType{
		bitCount,
		isSigned,
		getNextTypeId(),
		&VoidExpression{},
	}
}

func (integerType *IntegerType) CanAssignFrom(other TypeNode) bool {
	var asNumber, ok = other.(*IntegerType)

	return ok && asNumber.BitCount == integerType.BitCount && asNumber.IsSigned == integerType.IsSigned
}

func (integerType *IntegerType) GetNodeType() TypeNodeType {
	return IntegerNodeType
}

func (integerType *IntegerType) GetWhereExpression() Expression {
	return integerType.whereExpression
}
func (integerType *IntegerType) SetWhereExpression(expression Expression) {
	integerType.whereExpression = expression
}

func (integerType *IntegerType) UniqueId() int {
	return integerType.unqiueId
}

type BooleanType struct {
	uniqueID        int
	whereExpression Expression
}

func NewBooleanType() *BooleanType {
	return &BooleanType{
		getNextTypeId(),
		&VoidExpression{},
	}
}

func (booleanType *BooleanType) CanAssignFrom(other TypeNode) bool {
	var _, ok = other.(*BooleanType)
	return ok
}

func (booleanType *BooleanType) GetNodeType() TypeNodeType {
	return BooleanNodeType
}

func (booleanType *BooleanType) GetWhereExpression() Expression {
	return booleanType.whereExpression
}

func (booleanType *BooleanType) SetWhereExpression(expression Expression) {
	booleanType.whereExpression = expression
}

func (booleanType *BooleanType) UniqueId() int {
	return booleanType.uniqueID
}

type StructureNamedEntryType struct {
	Name     string
	UniqueId int
	Type     TypeNode
}

type StructureTypeType struct {
	Entries         []*StructureNamedEntryType
	uniqueId        int
	whereExpression Expression
}

func NewStructureTypeType(entries []*StructureNamedEntryType) *StructureTypeType {
	return &StructureTypeType{
		entries,
		getNextTypeId(),
		&VoidExpression{},
	}
}

func (structureType *StructureTypeType) CanAssignFrom(other TypeNode) bool {
	var asStructure, ok = other.(*StructureTypeType)

	if ok {
		if len(structureType.Entries) != len(asStructure.Entries) {
			return false
		}

		for index, entry := range structureType.Entries {
			if !entry.Type.CanAssignFrom(asStructure.Entries[index].Type) {
				return false
			}
		}

		return true
	} else {
		return false
	}
}

func (structureType *StructureTypeType) GetNodeType() TypeNodeType {
	return StructureNodeType
}

func (structureType *StructureTypeType) GetWhereExpression() Expression {
	return structureType.whereExpression
}

func (structureType *StructureTypeType) SetWhereExpression(expression Expression) {
	structureType.whereExpression = expression
}

func (structureType *StructureTypeType) UniqueId() int {
	return structureType.uniqueId
}

type FunctionTypeType struct {
	Input           *StructureTypeType
	Output          *StructureTypeType
	uniqueId        int
	whereExpression Expression
}

func NewFunctionTypeType(input *StructureTypeType, output *StructureTypeType) *FunctionTypeType {
	return &FunctionTypeType{
		input,
		output,
		getNextTypeId(),
		&VoidExpression{},
	}
}

func (functionType *FunctionTypeType) CanAssignFrom(other TypeNode) bool {
	asFunction, ok := other.(*FunctionTypeType)

	if ok {
		return functionType.Input.CanAssignFrom(asFunction.Input) &&
			functionType.Output.CanAssignFrom(asFunction.Output)
	} else {
		return false
	}
}

func (functionType *FunctionTypeType) GetNodeType() TypeNodeType {
	return FunctionNodeType
}

func (functionType *FunctionTypeType) GetWhereExpression() Expression {
	return functionType.whereExpression
}

func (functionType *FunctionTypeType) SetWhereExpression(expression Expression) {
	functionType.whereExpression = expression
}

func (functionType *FunctionTypeType) UniqueId() int {
	return functionType.uniqueId
}
