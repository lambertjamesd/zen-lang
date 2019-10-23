package parser

import (
	"zen/tokenizer"
)

var currentScopeID uint64 = 0

type Scope struct {
	Id          uint64
	ParentScope *Scope
}

func CreateScope() *Scope {
	currentScopeID = currentScopeID + 1
	return &Scope{
		currentScopeID,
		nil,
	}
}

type Visitor interface {
	VisitVoidExpression(id *VoidExpression)
	VisitIdentifier(id *Identifier)
	VisitNumber(number *Number)
	VisitUnaryExpression(exp *UnaryExpression)
	VisitPropertyExpression(exp *PropertyExpression)
	VisitBinaryExpression(exp *BinaryExpression)
	VisitStructureExpression(exp *StructureExpression)
	VisitFunction(function *Function)
	VisitIf(ifStatement *IfStatement)
	VisitBody(body *Body)

	VisitReturn(ret *ReturnStatement)

	VisitNamedType(namedType *NamedType)
	VisitStructureType(structure *StructureType)
	VisitFunctionType(fn *FunctionType)
	VisitWhereType(where *WhereType)

	VisitTypeDef(typeDef *TypeDefinition)
	VisitFnDef(fnDef *FunctionDefinition)
	VisitFile(fileDef *FileDefinition)
}

type ParseNode interface {
	Accept(visitor Visitor)
	Begin() tokenizer.SourceLocation
	End() tokenizer.SourceLocation
}

type Statement interface {
	ParseNode
}

type Expression interface {
	Statement
	GetType() TypeNode
}

type TypeExpression interface {
	ParseNode
}

type TypeSymbolDefinition interface {
	ParseNode
	GetType() TypeNode
}

type SymbolDefinition interface {
	ParseNode
	GetType() TypeNode
}

type Definition interface {
	ParseNode
}

type VoidExpression struct {
	At tokenizer.SourceLocation
}

func (node *VoidExpression) Accept(visitor Visitor) {
	visitor.VisitVoidExpression(node)
}

func (node *VoidExpression) GetType() TypeNode {
	return &VoidType{}
}

func (node *VoidExpression) Begin() tokenizer.SourceLocation {
	return node.At
}

func (node *VoidExpression) End() tokenizer.SourceLocation {
	return node.At
}

type Identifier struct {
	Token *tokenizer.Token
	Type  TypeNode
}

func (node *Identifier) Accept(visitor Visitor) {
	visitor.VisitIdentifier(node)
}

func (node *Identifier) GetType() TypeNode {
	return node.Type
}

func (node *Identifier) Begin() tokenizer.SourceLocation {
	return node.Token.At
}

func (node *Identifier) End() tokenizer.SourceLocation {
	return node.Token.End()
}

type Number struct {
	Token *tokenizer.Token
	Type  TypeNode
}

func (node *Number) Accept(visitor Visitor) {
	visitor.VisitNumber(node)
}

func (node *Number) GetType() TypeNode {
	return node.Type
}

func (node *Number) Begin() tokenizer.SourceLocation {
	return node.Token.At
}

func (node *Number) End() tokenizer.SourceLocation {
	return node.Token.End()
}

type PropertyExpression struct {
	Left     Expression
	Property *tokenizer.Token
	Type     TypeNode
}

func (node *PropertyExpression) Accept(visitor Visitor) {
	visitor.VisitPropertyExpression(node)
}

func (node *PropertyExpression) GetType() TypeNode {
	return node.Type
}

func (node *PropertyExpression) Begin() tokenizer.SourceLocation {
	return node.Left.Begin()
}

func (node *PropertyExpression) End() tokenizer.SourceLocation {
	return node.Property.End()
}

type BinaryExpression struct {
	Left     Expression
	Operator *tokenizer.Token
	Right    Expression
	Type     TypeNode
}

func (node *BinaryExpression) Accept(visitor Visitor) {
	visitor.VisitBinaryExpression(node)
}

func (node *BinaryExpression) GetType() TypeNode {
	return node.Type
}

func (node *BinaryExpression) Begin() tokenizer.SourceLocation {
	return node.Left.Begin()
}

func (node *BinaryExpression) End() tokenizer.SourceLocation {
	return node.Right.End()
}

type UnaryExpression struct {
	Expr     Expression
	Operator *tokenizer.Token
	Type     TypeNode
}

func (node *UnaryExpression) Accept(visitor Visitor) {
	visitor.VisitUnaryExpression(node)
}

func (node *UnaryExpression) GetType() TypeNode {
	return node.Type
}

func (node *UnaryExpression) Begin() tokenizer.SourceLocation {
	return node.Operator.At
}

func (node *UnaryExpression) End() tokenizer.SourceLocation {
	return node.Expr.End()
}

type StructureExpressionEntry struct {
	Name *tokenizer.Token
	Expr Expression
}

type StructureExpression struct {
	openBracket  *tokenizer.Token
	Entries      []StructureExpressionEntry
	closeBracket *tokenizer.Token
	Type         *StructureTypeType
}

func (node *StructureExpression) Accept(visitor Visitor) {
	visitor.VisitStructureExpression(node)
}

func (node *StructureExpression) GetType() TypeNode {
	return node.Type
}

func (node *StructureExpression) Begin() tokenizer.SourceLocation {
	return node.openBracket.At
}

func (node *StructureExpression) End() tokenizer.SourceLocation {
	return node.closeBracket.End()
}

type Body struct {
	open       *tokenizer.Token
	Statements []Statement
	Type       TypeNode
	Scope      *Scope
	end        *tokenizer.Token
}

func (node *Body) Accept(visitor Visitor) {
	visitor.VisitBody(node)
}

func (node *Body) GetType() TypeNode {
	return node.Type
}

func (node *Body) Begin() tokenizer.SourceLocation {
	return node.open.At
}

func (node *Body) End() tokenizer.SourceLocation {
	return node.end.End()
}

type IfStatement struct {
	ifKeyword   *tokenizer.Token
	Expresssion Expression
	Body        *Body
	ElseBody    Expression
	Type        TypeNode
}

func (node *IfStatement) GetType() TypeNode {
	return node.Type
}

func (node *IfStatement) Accept(visitor Visitor) {
	visitor.VisitIf(node)
}

func (node *IfStatement) Begin() tokenizer.SourceLocation {
	return node.ifKeyword.At
}

func (node *IfStatement) End() tokenizer.SourceLocation {
	if node.ElseBody != nil {
		return node.ElseBody.End()
	} else {
		return node.Body.End()
	}
}

type NamedType struct {
	Token *tokenizer.Token
	Type  TypeNode
}

func (node *NamedType) Accept(visitor Visitor) {
	visitor.VisitNamedType(node)
}

func (node *NamedType) Begin() tokenizer.SourceLocation {
	return node.Token.At
}

func (node *NamedType) End() tokenizer.SourceLocation {
	return node.Token.End()
}

type StructureNamedEntry struct {
	Name    *tokenizer.Token
	TypeExp TypeExpression
	Type    *StructureNamedEntryType
}

func (node *StructureNamedEntry) Begin() tokenizer.SourceLocation {
	if node.Name != nil {
		return node.Name.At
	} else {
		return node.TypeExp.Begin()
	}
}

func (node *StructureNamedEntry) End() tokenizer.SourceLocation {
	return node.TypeExp.End()
}

type StructureType struct {
	open    *tokenizer.Token
	Entries []*StructureNamedEntry
	Type    *StructureTypeType
	close   *tokenizer.Token
}

func (node *StructureType) Accept(visitor Visitor) {
	visitor.VisitStructureType(node)
}

func (node *StructureType) Begin() tokenizer.SourceLocation {
	return node.open.At
}

func (node *StructureType) End() tokenizer.SourceLocation {
	return node.close.End()
}

type FunctionType struct {
	Input  TypeExpression
	Output TypeExpression
	Type   *FunctionTypeType
}

func (node *FunctionType) Accept(visitor Visitor) {
	visitor.VisitFunctionType(node)
}

func (node *FunctionType) Begin() tokenizer.SourceLocation {
	return node.Input.Begin()
}

func (node *FunctionType) End() tokenizer.SourceLocation {
	return node.Output.End()
}

type WhereType struct {
	whereKeyword *tokenizer.Token
	TypeExp      TypeExpression
	WhereExp     Expression
}

func (node *WhereType) Accept(visitor Visitor) {
	visitor.VisitWhereType(node)
}

func (node *WhereType) Begin() tokenizer.SourceLocation {
	return node.whereKeyword.At
}

func (node *WhereType) End() tokenizer.SourceLocation {
	return node.WhereExp.End()
}

type AsNamingType struct {
	TypeExp TypeExpression
	Name    *tokenizer.Token
}

type TypeDefinition struct {
	typeKeyword *tokenizer.Token
	Name        *tokenizer.Token
	TypeExp     TypeExpression
	Scope       *Scope
	Type        TypeNode
}

func (node *TypeDefinition) Accept(visitor Visitor) {
	visitor.VisitTypeDef(node)
}

func (node *TypeDefinition) GetType() TypeNode {
	return node.Type
}

func (node *TypeDefinition) Begin() tokenizer.SourceLocation {
	return node.typeKeyword.At
}

func (node *TypeDefinition) End() tokenizer.SourceLocation {
	return node.TypeExp.End()
}

type Function struct {
	TypeExp TypeExpression
	Body    *Body
	Scope   *Scope
	Type    *FunctionTypeType
}

func (node *Function) Accept(visitor Visitor) {
	visitor.VisitFunction(node)
}

func (node *Function) GetType() TypeNode {
	return node.Type
}

func (node *Function) Begin() tokenizer.SourceLocation {
	return node.TypeExp.Begin()
}

func (node *Function) End() tokenizer.SourceLocation {
	return node.Body.End()
}

type FunctionDefinition struct {
	Name     *tokenizer.Token
	Function *Function
}

func (node *FunctionDefinition) Accept(visitor Visitor) {
	visitor.VisitFnDef(node)
}

func (node *FunctionDefinition) Begin() tokenizer.SourceLocation {
	return node.Name.At
}

func (node *FunctionDefinition) End() tokenizer.SourceLocation {
	return node.Function.End()
}

type FileDefinition struct {
	start       tokenizer.SourceLocation
	Definitions []Definition
	Scope       *Scope
	end         tokenizer.SourceLocation
}

func (node *FileDefinition) Accept(visitor Visitor) {
	visitor.VisitFile(node)
}

func (node *FileDefinition) Begin() tokenizer.SourceLocation {
	return node.start
}

func (node *FileDefinition) End() tokenizer.SourceLocation {
	return node.end
}

type FunctionInformation struct {
	ReturnType *StructureTypeType
}

type ReturnStatement struct {
	returnKeyword  *tokenizer.Token
	ExpressionList []Expression
	ForFunction    *FunctionInformation
}

func (node *ReturnStatement) Accept(visitor Visitor) {
	visitor.VisitReturn(node)
}

func (node *ReturnStatement) Begin() tokenizer.SourceLocation {
	return node.returnKeyword.At
}

func (node *ReturnStatement) End() tokenizer.SourceLocation {
	if len(node.ExpressionList) > 0 {
		return node.ExpressionList[len(node.ExpressionList)-1].End()
	} else {
		return node.returnKeyword.End()
	}
}
