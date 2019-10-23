package typechecker

import "zen/parser"

type SymbolResolver struct {
	symbolCollector *symbolCollector
	currentScope    *parser.Scope
}

func enterScope(symbolCollector *symbolCollector, scope *parser.Scope) {
	symbolCollector.currentScope = scope
}

func exitScope(symbolCollector *symbolCollector) {

}

func createSymbolResolver(symbolCollector *symbolCollector) *SymbolResolver {
	return &SymbolResolver{
		symbolCollector,
		nil,
	}
}

func (symbolResolver *SymbolResolver) VisitVoidExpression(id *parser.VoidExpression) {

}

func (symbolResolver *SymbolResolver) VisitIdentifier(id *parser.Identifier) {

}

func (symbolResolver *SymbolResolver) VisitNumber(number *parser.Number) {

}

func (symbolResolver *SymbolResolver) VisitUnaryExpression(exp *parser.UnaryExpression) {
	exp.Expr.Accept(symbolResolver)
}

func (symbolResolver *SymbolResolver) VisitPropertyExpression(exp *parser.PropertyExpression) {

}

func (symbolResolver *SymbolResolver) VisitBinaryExpression(exp *parser.BinaryExpression) {
	exp.Left.Accept(symbolResolver)
	exp.Right.Accept(symbolResolver)
}

func (symbolResolver *SymbolResolver) VisitStructureExpression(exp *parser.StructureExpression) {
	for _, entry := range exp.Entries {
		entry.Expr.Accept(symbolResolver)
	}
}

func (symbolResolver *SymbolResolver) VisitIf(ifStatement *parser.IfStatement) {
	ifStatement.Expresssion.Accept(symbolResolver)
	ifStatement.Body.Accept(symbolResolver)
	ifStatement.ElseBody.Accept(symbolResolver)
}

func (symbolResolver *SymbolResolver) VisitBody(body *parser.Body) {
	for _, entry := range body.Statements {
		entry.Accept(symbolResolver)
	}
}

func (symbolResolver *SymbolResolver) VisitReturn(ret *parser.ReturnStatement) {
	for _, expr := range ret.ExpressionList {
		expr.Accept(symbolResolver)
	}
}

func (symbolResolver *SymbolResolver) VisitNamedType(namedType *parser.NamedType) {

}

func (symbolResolver *SymbolResolver) VisitStructureType(structure *parser.StructureType) {
	for _, entry := range structure.Entries {
		entry.TypeExp.Accept(symbolResolver)
	}
}

func (symbolResolver *SymbolResolver) VisitFunctionType(fn *parser.FunctionType) {
	fn.Input.Accept(symbolResolver)
	fn.Output.Accept(symbolResolver)
}

func (symbolResolver *SymbolResolver) VisitWhereType(where *parser.WhereType) {
	where.TypeExp.Accept(symbolResolver)
	where.WhereExp.Accept(symbolResolver)
}

func (symbolResolver *SymbolResolver) VisitTypeDef(typeDef *parser.TypeDefinition) {
	typeDef.TypeExp.Accept(symbolResolver)
}

func (symbolResolver *SymbolResolver) VisitFunction(function *parser.Function) {
	function.TypeExp.Accept(symbolResolver)
	function.Body.Accept(symbolResolver)
}

func (symbolResolver *SymbolResolver) VisitFnDef(fnDef *parser.FunctionDefinition) {
	fnDef.Function.Accept(symbolResolver)
}

func (symbolResolver *SymbolResolver) VisitFile(fileDef *parser.FileDefinition) {
	for _, entry := range fileDef.Definitions {
		entry.Accept(symbolResolver)
	}
}

func ResolveSymbols(fileDef *parser.FileDefinition) {
	var symbols = CollectSymbols(fileDef)
	var resolver = createSymbolResolver(symbols)
	fileDef.Accept(resolver)
}
