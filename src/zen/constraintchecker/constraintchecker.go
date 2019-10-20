package constraintchecker

import (
	"zen/boundschecking"
	"zen/parser"
	"zen/tokenizer"
)

type ConstraintChecker struct {
	checkerStateStack []*ConstraintCheckerState
	functionStack     []*functionStackFrame
	normalizerState   *boundschecking.NormalizerState
	errors            []parser.ParseError
}

func NewConstrantChecker() *ConstraintChecker {
	return &ConstraintChecker{
		nil,
		nil,
		boundschecking.NewNormalizerState(),
		nil,
	}
}

func (constraintChecker *ConstraintChecker) reportErrorMessage(at tokenizer.SourceLocation, message string) {
	constraintChecker.errors = append(constraintChecker.errors, parser.CreateError(at, message))
}

func (constraintChecker *ConstraintChecker) reportError(parseError parser.ParseError) {
	constraintChecker.errors = append(constraintChecker.errors, parseError)
}

func (constraintChecker *ConstraintChecker) createState() *ConstraintCheckerState {
	var result *ConstraintCheckerState = nil
	if len(constraintChecker.checkerStateStack) == 0 {
		result = NewConstraintCheckerState()
	} else {
		result = constraintChecker.checkerStateStack[len(constraintChecker.checkerStateStack)-1].Copy()
	}
	constraintChecker.checkerStateStack = append(constraintChecker.checkerStateStack, result)
	return result
}

func (constraintChecker *ConstraintChecker) popState() {
	constraintChecker.checkerStateStack = constraintChecker.checkerStateStack[:len(constraintChecker.checkerStateStack)-1]
}

func (constraintChecker *ConstraintChecker) peekState() *ConstraintCheckerState {
	if len(constraintChecker.checkerStateStack) == 0 {
		return nil
	} else {
		return constraintChecker.checkerStateStack[len(constraintChecker.checkerStateStack)-1]
	}
}

func (constraintChecker *ConstraintChecker) peekFunctionStack() *functionStackFrame {
	if len(constraintChecker.checkerStateStack) == 0 {
		return nil
	} else {
		return constraintChecker.functionStack[len(constraintChecker.functionStack)-1]
	}
}

func (constraintChecker *ConstraintChecker) VisitVoidExpression(id *parser.VoidExpression) {

}

func (constraintChecker *ConstraintChecker) VisitIdentifier(id *parser.Identifier) {

}

func (constraintChecker *ConstraintChecker) VisitNumber(number *parser.Number) {

}

func (constraintChecker *ConstraintChecker) VisitUnaryExpression(exp *parser.UnaryExpression) {

}

func (constraintChecker *ConstraintChecker) VisitBinaryExpression(exp *parser.BinaryExpression) {

}

func (constraintChecker *ConstraintChecker) VisitFunction(function *parser.Function) {
	for _, input := range function.Type.Input.Entries {
		constraintChecker.normalizerState.UseIdentifierMapping(input.Name, input.UniqueId)
	}

	for _, output := range function.Type.Output.Entries {
		constraintChecker.normalizerState.UseIdentifierMapping(output.Name, output.UniqueId)
	}

	var functionStackFrame = newFunctionStackFrame(constraintChecker.normalizerState, function.Type)
	constraintChecker.functionStack = append(constraintChecker.functionStack, functionStackFrame)

	for index, conditions := range functionStackFrame.conditions {
		state := constraintChecker.createState()
		functionStackFrame.currentCondition = index

		if conditions.preConditions != nil {
			state.addRules([]*boundschecking.AndGroup{conditions.preConditions})
		}

		function.Body.Accept(constraintChecker)

		constraintChecker.popState()
	}

	constraintChecker.functionStack = constraintChecker.functionStack[:len(constraintChecker.functionStack)-1]

}

func (constraintChecker *ConstraintChecker) VisitIf(ifStatement *parser.IfStatement) {
	ifStatement.Expresssion.Accept(constraintChecker)
	var expresssionRules = constraintChecker.normalizerState.NormalizeToOrGroup(ifStatement.Expresssion)

	var ifBodyState = constraintChecker.createState()
	_, err := ifBodyState.addRules(expresssionRules.AndGroups)
	if err != nil {
		constraintChecker.reportErrorMessage(ifStatement.Expresssion.Begin(), err.Error())
	}
	ifStatement.Body.Accept(constraintChecker)
	constraintChecker.popState()

	if ifStatement.ElseBody != nil {
		var elseBodyState = constraintChecker.createState()
		_, err = elseBodyState.addRules(constraintChecker.normalizerState.NotOrGroup(expresssionRules).AndGroups)
		if err != nil {
			constraintChecker.reportErrorMessage(ifStatement.Expresssion.Begin(), err.Error())
		}
		ifStatement.ElseBody.Accept(constraintChecker)
		constraintChecker.popState()
	}
}

func (constraintChecker *ConstraintChecker) VisitBody(body *parser.Body) {
	for _, statement := range body.Statements {
		statement.Accept(constraintChecker)
	}
}

func (constraintChecker *ConstraintChecker) VisitReturn(ret *parser.ReturnStatement) {
	for _, returnValue := range ret.ExpressionList {
		returnValue.Accept(constraintChecker)
	}

	var functionStack = constraintChecker.peekFunctionStack()
	postCondition := functionStack.conditions[functionStack.currentCondition].postConditions

	if postCondition != nil {
		var state = constraintChecker.peekState()

		for index, returnValue := range ret.ExpressionList {
			sumGroup, err := constraintChecker.normalizerState.NormalizeToSumGroup(returnValue)

			if err != nil {
				constraintChecker.reportErrorMessage(returnValue.Begin(), err.Error())
			} else if index < len(functionStack.outputNames) {
				var rules = constraintChecker.normalizerState.CreateEquality(sumGroup, functionStack.outputNames[index])
				_, err := state.addSumGroups(rules)

				if err != nil {
					constraintChecker.reportErrorMessage(returnValue.Begin(), "Could not append to known data")
				}

			} else if index == len(functionStack.outputNames) {
				constraintChecker.reportErrorMessage(returnValue.Begin(), "Too many return arguments")
			}
		}

		result, err := state.checkOrGroup(postCondition)

		if err != nil {
			constraintChecker.reportErrorMessage(ret.Begin(), err.Error())
		} else if len(result) > 0 {
			constraintChecker.reportError(parser.CreateErrorWithMultipleLocations(
				ret.Begin(),
				"Could not verify post conditions",
				constraintChecker.formatErrorWithConstraints("With precondition at\n", result),
			))
		}
	}
}

func (constraintChecker *ConstraintChecker) formatErrorWithConstraints(lineMessage string, conditions []*boundschecking.SumGroup) []parser.ParseError {
	var sourceErrors []parser.ParseError = nil
	var topFrame = constraintChecker.peekFunctionStack()
	var alreadyFormatted = make(map[uint32]bool)

	if topFrame != nil {
		for _, condition := range conditions {
			var id = condition.GetUniqueId()

			_, ok := alreadyFormatted[id]

			if !ok {
				alreadyFormatted[id] = true

				expression, ok := topFrame.expressionMapping[id]

				if ok {
					sourceErrors = append(sourceErrors, parser.CreateError(expression.Begin(), lineMessage))
				}
			}
		}
	}

	return sourceErrors
}

func (constraintChecker *ConstraintChecker) VisitNamedType(namedType *parser.NamedType) {

}

func (constraintChecker *ConstraintChecker) VisitStructureType(structure *parser.StructureType) {

}

func (constraintChecker *ConstraintChecker) VisitFunctionType(fn *parser.FunctionType) {

}

func (constraintChecker *ConstraintChecker) VisitWhereType(where *parser.WhereType) {

}

func (constraintChecker *ConstraintChecker) VisitTypeDef(typeDef *parser.TypeDefinition) {
}

func (constraintChecker *ConstraintChecker) VisitFnDef(fnDef *parser.FunctionDefinition) {
	fnDef.Function.Accept(constraintChecker)
}

func (constraintChecker *ConstraintChecker) VisitFile(fileDef *parser.FileDefinition) {
	for _, definition := range fileDef.Definitions {
		definition.Accept(constraintChecker)
	}
}

func CheckConstraints(parseNode parser.ParseNode) []parser.ParseError {
	var checker = NewConstrantChecker()
	parseNode.Accept(checker)
	return checker.errors
}
