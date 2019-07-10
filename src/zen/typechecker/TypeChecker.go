package typechecker

import "zen/parser"

type TypeChecker struct {
}

func returnType(typeChecker *TypeChecker, typeNode parser.TypeNode) {

}

func CreateTypeChecker() *TypeChecker {
	return &TypeChecker{}
}

func (typeChecker *TypeChecker) VisitIdentifier(id *parser.Identifier) {

}

func (typeChecker *TypeChecker) VisitNumber(number *parser.Number) {
	returnType(typeChecker, &parser.IntegerType{
		32,
		true,
	})
}

func (typeChecker *TypeChecker) VisitBinaryExpression(exp *parser.BinaryExpression) {

}

func (typeChecker *TypeChecker) VisitIf(ifStatement *parser.IfStatement) {

}

func (typeChecker *TypeChecker) VisitBody(body *parser.Body) {

}

func (typeChecker *TypeChecker) VisitReturn(ret *parser.ReturnStatement) {

}

func (typeChecker *TypeChecker) VisitNamedType(namedType *parser.NamedType) {

}

func (typeChecker *TypeChecker) VisitStructureType(structure *parser.StructureType) {

}

func (typeChecker *TypeChecker) VisitFunctionType(fn *parser.FunctionType) {

}

func (typeChecker *TypeChecker) VisitWhereType(where *parser.WhereType) {

}

func (typeChecker *TypeChecker) VisitTypeDef(typeDef *parser.TypeDefinition) {

}

func (typeChecker *TypeChecker) VisitFnDef(fnDef *parser.FunctionDefinition) {

}

func (typeChecker *TypeChecker) VisitFile(fileDef *parser.FileDefinition) {

}
