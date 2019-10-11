package constraintchecker

import (
	"zen/boundschecking"
	"zen/parser"
)

type preAndPostConditions struct {
	preConditions  *boundschecking.AndGroup
	postConditions *boundschecking.AndGroup
}

type functionStackFrame struct {
	outputNames      []*boundschecking.VariableReference
	conditions       []preAndPostConditions
	currentCondition int
}

func newFunctionStackFrame(normalizerState *boundschecking.NormalizerState, fnType *parser.FunctionTypeType) *functionStackFrame {
	var result functionStackFrame

	for _, outputType := range fnType.Output.Entries {
		result.outputNames = append(result.outputNames, normalizerState.CreateVariableReference(outputType.Name, outputType.UniqueId))
	}

	var conditions = normalizerState.NormalizeToOrGroup(fnType.GetWhereExpression())

	for _, andGroup := range conditions.AndGroups {
		result.conditions = append(result.conditions, splitAndGroup(&result, normalizerState, andGroup))
	}

	return &result
}

func isPostConditionIdentifier(fnStackFrame *functionStackFrame, name string) bool {
	for _, outputName := range fnStackFrame.outputNames {
		if outputName.Name == name {
			return true
		}
	}

	return false
}

func isPostConditionGroup(fnStackFrame *functionStackFrame, andGroup *boundschecking.SumGroup) bool {
	for _, productGroup := range andGroup.ProductGroups {
		for _, node := range productGroup.Values.Array {
			asIdentifier, ok := node.(*boundschecking.VariableReference)

			if ok {
				if isPostConditionIdentifier(fnStackFrame, asIdentifier.Name) {
					return true
				}
			}
		}
	}

	return false
}

func splitAndGroup(fnStackFrame *functionStackFrame, normalizerState *boundschecking.NormalizerState, andGroup *boundschecking.AndGroup) preAndPostConditions {
	var pre []*boundschecking.SumGroup = nil
	var post []*boundschecking.SumGroup = nil

	for _, group := range andGroup.SumGroups {
		if isPostConditionGroup(fnStackFrame, group) {
			post = append(post, group)
		} else {
			pre = append(pre, group)
		}
	}

	return preAndPostConditions{
		normalizerState.CreateAndGroup(pre),
		normalizerState.CreateAndGroup(post),
	}
}
