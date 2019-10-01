package typechecker

import (
	"zen/parser"
	"zen/tokenizer"
)

type TypeChecker struct {
	scopes    []*VariableScope
	errors    []parser.ParseError
	typeStack []parser.TypeNode
}

type VariableReference struct {
	Type parser.TypeNode
}

type VariableScope struct {
	parentScope *VariableScope
	variableMap map[string]*VariableReference
}

func (typeChecker *TypeChecker) createScope() *VariableScope {
	var topScope *VariableScope = nil

	if len(typeChecker.scopes) != 0 {
		topScope = typeChecker.scopes[len(typeChecker.scopes)-1]
	}

	var result = append(typeChecker.scopes, &VariableScope{
		topScope,
		make(map[string]*VariableReference),
	})

	typeChecker.scopes = result

	return result
}

func (typeChecker *TypeChecker) popScope() *VariableScope {
	if len(typeChecker.scopes) == 0 {
		return nil
	} else {
		var result = typeChecker.scopes[len(typeChecker.scopes)-1]
		typeChecker.scopes = typeChecker.scopes[:len(typeChecker.scopes)-1]
		return result
	}
}

func (typeChecker *TypeChecker) findVariable(name string) *VariableReference {
	for index := len(typeChecker.scopes) - 1; index >= 0; index = index - 1 {
		resultCheck, ok := typeChecker.scopes[index].variableMap[name]

		if ok {
			return resultCheck
		}
	}

	return nil
}

func (typeChecker *TypeChecker) reportError(at tokenizer.SourceLocation, message string) {
	typeChecker.errors = append(typeChecker.errors, parser.ParseError(at, message))
}

func (typeChecker *TypeChecker) pushType(typeValue parser.TypeNode) {
	typeChecker.typeStack = append(typeChecker.typeStack, typeValue)
}

func (typeChecker *TypeChecker) popType() parser.TypeNode {
	if len(typeChecker.typeStack) == 0 {
		return &parser.UndefinedType{}
	} else {
		var result = typeChecker.typeStack[len(typeChecker.typeStack)-1]
		typeChecker.typeStack = typeChecker.typeStack[:len(typeChecker.typeStack)-1]
		return result
	}
}

func CreateTypeChecker() *TypeChecker {
	return &TypeChecker{}
}

func (typeChecker *TypeChecker) VisitIdentifier(id *parser.Identifier) {
	variableReference, ok := typeChecker.findVariable(id.Token.Value)

	if ok {
		id.Type = variableReference.Type
		typeChecker.pushType(variableReference.Type)
	} else {
		typeChecker.reportError("Variable '"+id.Token.Value+"' is not defined", id.Token.At)
		typeChecker.pushType(&parser.UndefinedType{})
	}
}

func (typeChecker *TypeChecker) VisitNumber(number *parser.Number) {
	var result = &parser.IntegerType{
		32,
		true,
	}

	number.Type = result
	pushType(typeChecker, result)
}

func (typeChecker *TypeChecker) VisitUnaryExpression(exp *parser.UnaryExpression) {
	if exp.Operator.TokenType == tokenizer.MinusToken {
		exp.Expr.Accept(typeChecker)

		var subType = typeChecker.popType()
		var subEnumType = subType.GetNodeType()

		if subEnumType == parser.IntegerNodeType {
			exp.Type = subType
			typeChecker.pushType(subType)
		} else {
			if subEnumType != parser.UndefinedNodeType {
				typeChecker.reportError(exp.Operator.At, "Could not apply operator '-' to type ")
			}
			typeChecker.pushType(&UndefinedType{})
		}
	} else {
		typeChecker.reportError(exp.Operator.At, "Unknown operator '"+exp.Operator.Value+"'")
		typeChecker.pushType(&UndefinedType{})
	}
}

func (typeChecker *TypeChecker) VisitBinaryExpression(exp *parser.BinaryExpression) {
	exp.Left.Accept(typeChecker)
	var leftType = typeChecker.popType()
	exp.Right.Accept(typeChecker)
	var rightType = typeChecker.popType()

	if exp.Operator.TokenType == tokenizer.AddToken ||
		exp.Operator.TokenType == tokenizer.MinusToken ||
		exp.Operator.TokenType == tokenizer.MultiplyToken ||
		exp.Operator.TokenType == tokenizer.DivideToken {
		if leftType.GetNodeType() == rightType.GetNodeType() && leftType.GetNodeType() == parser.IntegerNodeType {
			exp.Type = leftType
			typeChecker.pushType(leftType)
		} else {
			if leftType.GetNodeType() != parser.UndefinedNodeType && rightType.GetNodeType() != parser.UndefinedNodeType {
				typeChecker.reportError(exp.Operator.At, "Operator '"+exp.Operator.Value+"' cannot be applied to given types")
			}
			typeChecker.pushType(&UndefinedType{})
		}
	} else if exp.Operator.TokenType == tokenizer.EqualToken {
		if leftType.GetNodeType() == rightType.GetNodeType() {
			exp.Type = &parser.BooleanType{}
			typeChecker.pushType(exp.Type)
		} else {
			typeChecker.reportError(exp.Operator.At, "Cannot compare incompatible types")
			typeChecker.pushType(&UndefinedType{})
		}
	} else if exp.Operator.TokenType == tokenizer.LTEqToken ||
		exp.Operator.TokenType == tokenizer.LTToken ||
		exp.Operator.TokenType == tokenizer.GTEqToken ||
		exp.Operator.TokenType == tokenizer.GTToken {
		if leftType.GetNodeType() == rightType.GetNodeType() && leftType.GetNodeType() != parser.IntegerNodeType {
			exp.Type = &BooleanType{}
			typeChecker.pushType(exp.Type)
		} else {
			if leftType.GetNodeType() != parser.UndefinedNodeType && rightType.GetNodeType() != parser.UndefinedNodeType {
				typeChecker.reportError(exp.Operator.At, "Operator '"+exp.Operator.Value+"' cannot be applied to given types")
			}
			typeChecker.pushType(&UndefinedType{})
		}
	} else {
		typeChecker.reportError(exp.Operator.At, "Operator '"+exp.Operator.Value+"' not supported")
		typeChecker.pushType(&UndefinedType{})
	}
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
	for _, definition := range fileDef.Definitions {
		definition.Accept(typeChecker)
	}
}
