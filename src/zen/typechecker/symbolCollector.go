package typechecker

import "zen/parser"

type scopeTypeReferences struct {
	typeSymbols map[string]parser.TypeSymbolDefinition
	symbols     map[string]parser.SymbolDefinition
}

type symbolCollector struct {
	symbols           map[uint64]*scopeTypeReferences
	currentScope      *parser.Scope
	currentReferences *scopeTypeReferences
}

func CollectSymbols(fileDef *parser.FileDefinition) *symbolCollector {
	var result = createSymbolCollector()
	fileDef.Accept(result)
	return result
}

func createScopeTypeReferences() *scopeTypeReferences {
	return &scopeTypeReferences{
		make(map[string]parser.TypeSymbolDefinition),
		make(map[string]parser.SymbolDefinition),
	}
}

func createSymbolCollector() *symbolCollector {
	return &symbolCollector{
		make(map[uint64]*scopeTypeReferences),
		nil,
		nil,
	}
}

func startScope(symbolCollector *symbolCollector, scope *parser.Scope) {
	scope.ParentScope = symbolCollector.currentScope
	symbolCollector.currentScope = scope
	symbolCollector.currentReferences = createScopeTypeReferences()
	symbolCollector.symbols[scope.Id] = symbolCollector.currentReferences
}

func endScope(symbolCollector *symbolCollector) {
	symbolCollector.currentScope = symbolCollector.currentScope.ParentScope

	if symbolCollector.currentScope == nil {
		symbolCollector.currentReferences = nil
	} else {
		symbolCollector.currentReferences = symbolCollector.symbols[symbolCollector.currentScope.Id]
	}
}

func (symbolCollector *symbolCollector) VisitIdentifier(id *parser.Identifier) {

}

func (symbolCollector *symbolCollector) VisitNumber(number *parser.Number) {

}

func (symbolCollector *symbolCollector) VisitBinaryExpression(exp *parser.BinaryExpression) {
	exp.Left.Accept(symbolCollector)
	exp.Right.Accept(symbolCollector)
}

func (symbolCollector *symbolCollector) VisitIf(ifStatement *parser.IfStatement) {
	ifStatement.Expresssion.Accept(symbolCollector)
	ifStatement.Body.Accept(symbolCollector)
	ifStatement.ElseBody.Accept(symbolCollector)
}

func (symbolCollector *symbolCollector) VisitBody(body *parser.Body) {
	startScope(symbolCollector, body.Scope)

	for _, entry := range body.Statements {
		entry.Accept(symbolCollector)
	}

	endScope(symbolCollector)
}

func (symbolCollector *symbolCollector) VisitReturn(ret *parser.ReturnStatement) {
	ret.Expression.Accept(symbolCollector)
}

func (symbolCollector *symbolCollector) VisitNamedType(namedType *parser.NamedType) {

}

func (symbolCollector *symbolCollector) VisitStructureNamedEntry(structureEntry *parser.StructureNamedEntry) {
	structureEntry.TypeExp.Accept(symbolCollector)
}

func (symbolCollector *symbolCollector) VisitStructureType(structure *parser.StructureType) {
	for _, entry := range structure.Entries {
		symbolCollector.currentReferences.symbols[entry.Name.Value] = entry
		entry.TypeExp.Accept(symbolCollector)
	}
}

func (symbolCollector *symbolCollector) VisitFunction(fn *parser.Function) {
	startScope(symbolCollector, fn.Scope)
	fn.TypeExp.Accept(symbolCollector)
	fn.Body.Accept(symbolCollector)
	endScope(symbolCollector)
}

func (symbolCollector *symbolCollector) VisitFunctionType(fn *parser.FunctionType) {
	fn.Input.Accept(symbolCollector)
	fn.Output.Accept(symbolCollector)
}

func (symbolCollector *symbolCollector) VisitWhereType(where *parser.WhereType) {
	where.TypeExp.Accept(symbolCollector)
	where.WhereExp.Accept(symbolCollector)
}

func (symbolCollector *symbolCollector) VisitTypeDef(typeDef *parser.TypeDefinition) {
	startScope(symbolCollector, typeDef.Scope)
	symbolCollector.currentReferences.typeSymbols[typeDef.Name.Value] = typeDef
	typeDef.TypeExp.Accept(symbolCollector)
	endScope(symbolCollector)
}

func (symbolCollector *symbolCollector) VisitFnDef(fnDef *parser.FunctionDefinition) {
	symbolCollector.currentReferences.symbols[fnDef.Name.Value] = fnDef.Function
	fnDef.Function.Accept(symbolCollector)
}

func (symbolCollector *symbolCollector) VisitFile(fileDef *parser.FileDefinition) {
	startScope(symbolCollector, fileDef.Scope)
	for _, definition := range fileDef.Definitions {
		definition.Accept(symbolCollector)
	}
	endScope(symbolCollector)
}
