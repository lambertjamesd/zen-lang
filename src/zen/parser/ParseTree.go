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
	VisitIdentifier(id *Identifier)
	VisitNumber(number *Number)
	VisitBinaryExpression(exp *BinaryExpression)
	VisitFunction(function *Function)
	VisitIf(ifStatement *IfStatement)
	VisitBody(body *Body)

	VisitReturn(ret *ReturnStatement)

	VisitNamedType(namedType *NamedType)
	VisitStructureNamedEntry(structureEntry *StructureNamedEntry)
	VisitStructureType(structure *StructureType)
	VisitFunctionType(fn *FunctionType)
	VisitWhereType(where *WhereType)

	VisitTypeDef(typeDef *TypeDefinition)
	VisitFnDef(fnDef *FunctionDefinition)
	VisitFile(fileDef *FileDefinition)
}

type ParseNode interface {
}

type Statement interface {
	Accept(visitor Visitor)
}

type Expression interface {
	Statement
	GetType() TypeNode
}

type TypeExpression interface {
	Accept(visitor Visitor)
}

type TypeSymbolDefinition interface {
	Accept(visitor Visitor)
	GetType() TypeNode
}

type SymbolDefinition interface {
	Accept(visitor Visitor)
	GetType() TypeNode
}

type Definition interface {
	Accept(visitor Visitor)
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

type Body struct {
	Statements []Statement
	Type       TypeNode
	Scope      *Scope
}

func (node *Body) Accept(visitor Visitor) {
	visitor.VisitBody(node)
}

func (node *Body) GetType() TypeNode {
	return node.Type
}

type IfStatement struct {
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

type NamedType struct {
	Token *tokenizer.Token
}

func (node *NamedType) Accept(visitor Visitor) {
	visitor.VisitNamedType(node)
}

type StructureNamedEntry struct {
	Name    *tokenizer.Token
	TypeExp TypeExpression
	Type    TypeNode
}

func (node *StructureNamedEntry) Accept(visitor Visitor) {
	visitor.VisitStructureNamedEntry(node)
}

func (node *StructureNamedEntry) GetType() TypeNode {
	return node.Type
}

type StructureType struct {
	Entries []*StructureNamedEntry
}

func (node *StructureType) Accept(visitor Visitor) {
	visitor.VisitStructureType(node)
}

type FunctionType struct {
	Input  TypeExpression
	Output TypeExpression
}

func (node *FunctionType) Accept(visitor Visitor) {
	visitor.VisitFunctionType(node)
}

type WhereType struct {
	TypeExp  TypeExpression
	WhereExp Expression
}

func (node *WhereType) Accept(visitor Visitor) {
	visitor.VisitWhereType(node)
}

type AsNamingType struct {
	TypeExp TypeExpression
	Name    *tokenizer.Token
}

type TypeDefinition struct {
	Name    *tokenizer.Token
	TypeExp TypeExpression
	Scope   *Scope
	Type    TypeNode
}

func (node *TypeDefinition) Accept(visitor Visitor) {
	visitor.VisitTypeDef(node)
}

func (node *TypeDefinition) GetType() TypeNode {
	return node.Type
}

type Function struct {
	TypeExp TypeExpression
	Body    *Body
	Scope   *Scope
	Type    TypeNode
}

func (node *Function) Accept(visitor Visitor) {
	visitor.VisitFunction(node)
}

func (node *Function) GetType() TypeNode {
	return node.Type
}

type FunctionDefinition struct {
	Name     *tokenizer.Token
	Function *Function
}

func (node *FunctionDefinition) Accept(visitor Visitor) {
	visitor.VisitFnDef(node)
}

type FileDefinition struct {
	Definitions []Definition
	Scope       *Scope
}

func (node *FileDefinition) Accept(visitor Visitor) {
	visitor.VisitFile(node)
}

type ReturnStatement struct {
	Expression Expression
}

func (node *ReturnStatement) Accept(visitor Visitor) {
	visitor.VisitReturn(node)
}
