package typechecker

import (
	"fmt"
	"zen/parser"
	"zen/tokenizer"
)

type TypeChecker struct {
	scopes        []*VariableScope
	errors        []parser.ParseError
	typeStack     []parser.TypeNode
	functionStack []*parser.FunctionInformation
}

type VariableReference struct {
	Type parser.TypeNode
}

type TypeReference struct {
	Type parser.TypeNode
}

type VariableScope struct {
	parentScope *VariableScope
	variableMap map[string]*VariableReference
	typeMap     map[string]*TypeReference
}

func (typeChecker *TypeChecker) pushFunctionInfo(fnInfo *parser.FunctionInformation) {
	typeChecker.functionStack = append(typeChecker.functionStack, fnInfo)
}

func (typeChecker *TypeChecker) popFunctionInfo() {
	typeChecker.functionStack = typeChecker.functionStack[:len(typeChecker.functionStack)-1]
}

func (typeChecker *TypeChecker) peekFunctionInfo() *parser.FunctionInformation {
	if len(typeChecker.functionStack) == 0 {
		return nil
	} else {
		return typeChecker.functionStack[len(typeChecker.functionStack)-1]
	}
}

func (scope *VariableScope) initializeDefaultTypes() {
	scope.typeMap["i32"] = &TypeReference{parser.NewIntegerType(32, true)}

	scope.typeMap["bool"] = &TypeReference{parser.NewBooleanType()}
}

func (typeChecker *TypeChecker) createScope() *VariableScope {
	var topScope *VariableScope = nil

	if len(typeChecker.scopes) != 0 {
		topScope = typeChecker.scopes[len(typeChecker.scopes)-1]
	}

	var result = &VariableScope{
		topScope,
		make(map[string]*VariableReference),
		make(map[string]*TypeReference),
	}

	typeChecker.scopes = append(typeChecker.scopes, result)

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

func (typeChecker *TypeChecker) peekScope() *VariableScope {
	if len(typeChecker.scopes) == 0 {
		return nil
	} else {
		return typeChecker.scopes[len(typeChecker.scopes)-1]
	}
}

func (typeChecker *TypeChecker) findType(name string) *TypeReference {
	for index := len(typeChecker.scopes) - 1; index >= 0; index = index - 1 {
		resultCheck, ok := typeChecker.scopes[index].typeMap[name]

		if ok {
			return resultCheck
		}
	}

	return nil
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
	typeChecker.errors = append(typeChecker.errors, parser.CreateError(at, message))
}

func (typeChecker *TypeChecker) pushType(typeValue parser.TypeNode) {
	typeChecker.typeStack = append(typeChecker.typeStack, typeValue)
}

func (typeChecker *TypeChecker) acceptSubType(statement parser.ParseNode) parser.TypeNode {
	lenBefore := len(typeChecker.typeStack)
	statement.Accept(typeChecker)

	if lenBefore < len(typeChecker.typeStack) {
		result := typeChecker.typeStack[lenBefore]
		typeChecker.typeStack = typeChecker.typeStack[:lenBefore]
		return result
	} else {
		return &parser.VoidType{}
	}
}

func CreateTypeChecker() *TypeChecker {
	return &TypeChecker{}
}

func (typeChecker *TypeChecker) VisitVoidExpression(expr *parser.VoidExpression) {
	typeChecker.pushType(&parser.VoidType{})
}

func (typeChecker *TypeChecker) VisitIdentifier(id *parser.Identifier) {
	variableReference := typeChecker.findVariable(id.Token.Value)

	if variableReference != nil {
		id.Type = variableReference.Type
		typeChecker.pushType(variableReference.Type)
	} else {
		typeChecker.reportError(id.Token.At, "Variable '"+id.Token.Value+"' is not defined")
		typeChecker.pushType(&parser.UndefinedType{})
	}
}

func (typeChecker *TypeChecker) VisitNumber(number *parser.Number) {
	var result = parser.NewIntegerType(32, true)

	number.Type = result
	typeChecker.pushType(result)
}

func (typeChecker *TypeChecker) VisitUnaryExpression(exp *parser.UnaryExpression) {
	if exp.Operator.TokenType == tokenizer.MinusToken {
		var subType = typeChecker.acceptSubType(exp.Expr)
		var subEnumType = subType.GetNodeType()

		if subEnumType == parser.IntegerNodeType {
			exp.Type = subType
			typeChecker.pushType(subType)
		} else {
			if subEnumType != parser.UndefinedNodeType {
				typeChecker.reportError(exp.Operator.At, "Could not apply operator '-' to type ")
			}
			typeChecker.pushType(&parser.UndefinedType{})
		}
	} else {
		typeChecker.reportError(exp.Operator.At, "Unknown operator '"+exp.Operator.Value+"'")
		typeChecker.pushType(&parser.UndefinedType{})
	}
}

func (typeChecker *TypeChecker) VisitBinaryExpression(exp *parser.BinaryExpression) {
	var leftType = typeChecker.acceptSubType(exp.Left)
	var rightType = typeChecker.acceptSubType(exp.Right)

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
			typeChecker.pushType(&parser.UndefinedType{})
		}
	} else if exp.Operator.TokenType == tokenizer.EqualToken {
		if leftType.GetNodeType() == rightType.GetNodeType() {
			exp.Type = &parser.BooleanType{}
			typeChecker.pushType(exp.Type)
		} else {
			typeChecker.reportError(exp.Operator.At, "Cannot compare incompatible types")
			typeChecker.pushType(&parser.UndefinedType{})
		}
	} else if exp.Operator.TokenType == tokenizer.LTEqToken ||
		exp.Operator.TokenType == tokenizer.LTToken ||
		exp.Operator.TokenType == tokenizer.GTEqToken ||
		exp.Operator.TokenType == tokenizer.GTToken {
		if leftType.GetNodeType() != rightType.GetNodeType() || leftType.GetNodeType() != parser.IntegerNodeType {
			if leftType.GetNodeType() != parser.UndefinedNodeType && rightType.GetNodeType() != parser.UndefinedNodeType {
				typeChecker.reportError(exp.Operator.At, "Operator '"+exp.Operator.Value+"' cannot be applied to given types")
			}
		}
		exp.Type = &parser.BooleanType{}
		typeChecker.pushType(exp.Type)
	} else if exp.Operator.TokenType == tokenizer.BooleanAndToken ||
		exp.Operator.TokenType == tokenizer.BooleanOrToken {
		if leftType.GetNodeType() != parser.BooleanNodeType {
			typeChecker.reportError(exp.Operator.At, "Operator '"+exp.Operator.Value+"' left hand side must evalualte to boolean")
		}
		if rightType.GetNodeType() != parser.BooleanNodeType {
			typeChecker.reportError(exp.Operator.At, "Operator '"+exp.Operator.Value+"' left hand side must evalualte to boolean")
		}
		exp.Type = &parser.BooleanType{}
		typeChecker.pushType(exp.Type)
	} else {
		typeChecker.reportError(exp.Operator.At, "Operator '"+exp.Operator.Value+"' not supported")
		typeChecker.pushType(&parser.UndefinedType{})
	}
}

func (typeChecker *TypeChecker) VisitIf(ifStatement *parser.IfStatement) {
	_, ok := typeChecker.acceptSubType(ifStatement.Expresssion).(*parser.BooleanType)

	if !ok {
		typeChecker.reportError(ifStatement.Expresssion.Begin(), "If expression must evaluate to boolean")
	}

	typeChecker.acceptSubType(ifStatement.Body)

	if ifStatement.ElseBody != nil {
		typeChecker.acceptSubType(ifStatement.ElseBody)
	}
}

func (typeChecker *TypeChecker) VisitBody(body *parser.Body) {
	for _, statement := range body.Statements {
		typeChecker.acceptSubType(statement)
	}
}

func (typeChecker *TypeChecker) VisitReturn(ret *parser.ReturnStatement) {
	var forFunction = typeChecker.peekFunctionInfo()
	if forFunction == nil {
		typeChecker.reportError(ret.Begin(), "Return statements must be inside a function")
	} else if len(forFunction.ReturnType.Entries) != len(ret.ExpressionList) {
		typeChecker.reportError(ret.Begin(), fmt.Sprintf(
			"Expected %d return values got %d",
			len(forFunction.ReturnType.Entries),
			len(ret.ExpressionList),
		))
	}

	ret.ForFunction = forFunction

	for index, expression := range ret.ExpressionList {
		var returnType = typeChecker.acceptSubType(expression)

		if forFunction != nil && !forFunction.ReturnType.Entries[index].Type.CanAssignFrom(returnType) {
			typeChecker.reportError(expression.Begin(), "Return type incomatible with function signature")
		}
	}
}

func (typeChecker *TypeChecker) VisitNamedType(namedType *parser.NamedType) {
	var typeResult = typeChecker.findType(namedType.Token.Value)

	if typeResult == nil {
		typeChecker.reportError(namedType.Begin(), fmt.Sprintf("Could not find type %s", namedType.Token.Value))
		typeChecker.pushType(&parser.UndefinedType{})
	} else {
		namedType.Type = typeResult.Type
		typeChecker.pushType(typeResult.Type)
	}
}

func (typeChecker *TypeChecker) VisitStructureType(structure *parser.StructureType) {
	var subEntries []*parser.StructureNamedEntryType
	var topScope = typeChecker.peekScope()

	for _, entry := range structure.Entries {
		var subType = &parser.StructureNamedEntryType{
			entry.Name.Value,
			typeChecker.acceptSubType(entry.TypeExp),
		}
		entry.Type = subType
		subEntries = append(subEntries, subType)
		topScope.variableMap[entry.Name.Value] = &VariableReference{subType.Type}
	}

	var result = parser.NewStructureTypeType(subEntries)

	structure.Type = result
	typeChecker.pushType(result)
}

func (typeChecker *TypeChecker) VisitFunctionType(fn *parser.FunctionType) {
	var inputType = typeChecker.acceptSubType(fn.Input)
	var outputType = typeChecker.acceptSubType(fn.Output)

	inputAsStructure, ok := inputType.(*parser.StructureTypeType)

	if !ok {
		typeChecker.reportError(fn.Input.Begin(), "Function input type must be a structure")
	}

	outputAsStructure, ok := outputType.(*parser.StructureTypeType)

	if !ok {
		typeChecker.reportError(fn.Output.Begin(), "Function output type must be a structure")
	}

	var result = parser.NewFunctionTypeType(inputAsStructure, outputAsStructure)

	fn.Type = result
	typeChecker.pushType(result)
}

func (typeChecker *TypeChecker) VisitWhereType(where *parser.WhereType) {
	var whereScope = typeChecker.createScope()

	var contrainedType = typeChecker.acceptSubType(where.TypeExp)
	contrainedType.SetWhereExpression(where.WhereExp)
	typeChecker.pushType(contrainedType)

	var selfReference = &VariableReference{contrainedType}
	whereScope.variableMap["self"] = selfReference

	_, ok := typeChecker.acceptSubType(where.WhereExp).(*parser.BooleanType)

	if !ok {
		typeChecker.reportError(where.Begin(), "Where expression must evaluate to a boolean")
	}

	typeChecker.popScope()
}

func (typeChecker *TypeChecker) VisitTypeDef(typeDef *parser.TypeDefinition) {
	typeChecker.createScope()
	var defType = typeChecker.acceptSubType(typeDef.TypeExp)
	typeChecker.popScope()
	typeDef.Type = defType
	var topScope = typeChecker.peekScope()
	topScope.typeMap[typeDef.Name.Value] = &TypeReference{defType}
}

func (typeChecker *TypeChecker) VisitFunction(fn *parser.Function) {
	var topScope = typeChecker.createScope()

	var fnType = typeChecker.acceptSubType(fn.TypeExp)
	asFunctionType, ok := fnType.(*parser.FunctionTypeType)

	if !ok {
		typeChecker.reportError(fn.TypeExp.Begin(), "Function must be a function type")
	}

	var returnType *parser.StructureTypeType

	if !ok {
		returnType = parser.NewStructureTypeType(nil)
	} else {
		returnType = asFunctionType.Output

		for _, subType := range asFunctionType.Input.Entries {
			topScope.variableMap[subType.Name] = &VariableReference{subType.Type}
		}

		for _, subType := range asFunctionType.Output.Entries {
			topScope.variableMap[subType.Name] = &VariableReference{subType.Type}
		}
	}

	typeChecker.pushFunctionInfo(&parser.FunctionInformation{
		returnType,
	})

	typeChecker.acceptSubType(fn.Body)

	typeChecker.popScope()
	typeChecker.popFunctionInfo()

	fn.Type = asFunctionType
	typeChecker.pushType(asFunctionType)
}

func (typeChecker *TypeChecker) VisitFnDef(fnDef *parser.FunctionDefinition) {
	var fnType = typeChecker.acceptSubType(fnDef.Function)
	var topScope = typeChecker.peekScope()
	topScope.variableMap[fnDef.Name.Value] = &VariableReference{fnType}
}

func (typeChecker *TypeChecker) VisitFile(fileDef *parser.FileDefinition) {
	typeChecker.createScope().initializeDefaultTypes()
	for _, definition := range fileDef.Definitions {
		typeChecker.acceptSubType(definition)
	}
	typeChecker.popScope()
}

func CheckTypes(parseNode parser.ParseNode) []parser.ParseError {
	var checker = CreateTypeChecker()
	parseNode.Accept(checker)
	return checker.errors
}
