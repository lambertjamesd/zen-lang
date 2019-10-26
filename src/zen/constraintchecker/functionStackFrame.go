package constraintchecker

import (
	"zen/boundschecking"
	"zen/parser"
)

type preAndPostConditions struct {
	preConditions  *boundschecking.AndGroup
	postConditions *boundschecking.OrGroup
}

type functionStackFrame struct {
	outputNames       []*boundschecking.VariableReference
	conditions        []preAndPostConditions
	currentCondition  int
	expressionMapping map[uint32]parser.Expression
}

const NO_PRECONDITIONS = ^uint32(0)

func newFunctionStackFrame(normalizerState *boundschecking.NormalizerState, fnType *parser.FunctionTypeType) *functionStackFrame {
	var result functionStackFrame

	var preConditionMapping = make(map[uint32]*boundschecking.AndGroup)
	var postConditionMapping = make(map[uint32][]*boundschecking.AndGroup)

	for _, outputType := range fnType.Output.Entries {
		result.outputNames = append(result.outputNames, normalizerState.CreateVariableReference(outputType.Name, outputType.UniqueId))
	}

	result.expressionMapping = normalizerState.StartTrackingExpressionMapping()

	var conditions = normalizerState.NormalizeToOrGroup(fnType.GetWhereExpression())

	normalizerState.StopTrackingExpressionMapping()

	for _, andGroup := range conditions.AndGroups {
		preGroup, postGroup := splitAndGroup(&result, normalizerState, andGroup)
		var groupID uint32
		if preGroup == nil {
			groupID = NO_PRECONDITIONS
		} else {
			groupID = preGroup.GetUniqueId()
		}
		preConditionMapping[groupID] = preGroup
		if postGroup != nil {
			postConditionMapping[groupID] = append(postConditionMapping[groupID], postGroup)
		}
	}

	for id, preGroup := range preConditionMapping {
		postConditions := postConditionMapping[id]

		result.conditions = append(result.conditions, preAndPostConditions{
			preGroup,
			normalizerState.CreateOrGroup(postConditions),
		})
	}

	if len(result.conditions) == 0 {
		result.conditions = append(result.conditions, preAndPostConditions{nil, nil})
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

func splitAndGroup(fnStackFrame *functionStackFrame, normalizerState *boundschecking.NormalizerState, andGroup *boundschecking.AndGroup) (preGroup *boundschecking.AndGroup, postGroup *boundschecking.AndGroup) {
	var pre []*boundschecking.SumGroup = nil
	var post []*boundschecking.SumGroup = nil

	for _, group := range andGroup.SumGroups {
		if isPostConditionGroup(fnStackFrame, group) {
			post = append(post, group)
		} else {
			pre = append(pre, group)
		}
	}

	return normalizerState.CreateAndGroup(pre), normalizerState.CreateAndGroup(post)
}
