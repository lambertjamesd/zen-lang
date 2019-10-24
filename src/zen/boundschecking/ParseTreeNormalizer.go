package boundschecking

import (
	"errors"
	"strconv"
	"zen/parser"
	"zen/tokenizer"
	"zen/zmath"
)

func (state *NormalizerState) normalizeExpression(expression parser.Expression) *NormalizedEquation {
	return nil
}

func (state *NormalizerState) normalizeBinaryExpressionToOrGroup(expression *parser.BinaryExpression) *OrGroup {
	if expression.Operator.TokenType == tokenizer.BooleanOrToken {
		return state.combineOrGroups(
			state.NormalizeToOrGroup(expression.Left),
			state.NormalizeToOrGroup(expression.Right),
		)
	} else if expression.Operator.TokenType == tokenizer.BooleanAndToken {
		return state.combineOrGroupsWithAnd(
			state.NormalizeToOrGroup(expression.Left),
			state.NormalizeToOrGroup(expression.Right),
		)
	} else if expression.Operator.TokenType == tokenizer.LTToken ||
		expression.Operator.TokenType == tokenizer.LTEqToken ||
		expression.Operator.TokenType == tokenizer.GTToken ||
		expression.Operator.TokenType == tokenizer.GTEqToken {
		leftSumGroup, _ := state.NormalizeToSumGroup(expression.Left)
		rightSumGroup, _ := state.NormalizeToSumGroup(expression.Right)

		if leftSumGroup == nil || rightSumGroup == nil {
			return &OrGroup{nil}
		}

		if expression.Operator.TokenType == tokenizer.LTToken ||
			expression.Operator.TokenType == tokenizer.LTEqToken {
			leftSumGroup = state.negateSumGroup(leftSumGroup)
		} else {
			rightSumGroup = state.negateSumGroup(rightSumGroup)
		}

		var offset = int64(0)

		if expression.Operator.TokenType == tokenizer.LTToken ||
			expression.Operator.TokenType == tokenizer.GTToken {
			offset = -1
		}

		var sumGroup = state.recordExpressionMapping(state.addSumGroups(leftSumGroup, rightSumGroup, offset), expression)
		return state.nodeCache.GetNodeSingleton(&OrGroup{
			[]*AndGroup{
				state.nodeCache.GetNodeSingleton(&AndGroup{
					[]*SumGroup{sumGroup},
					state.getNextUniqueId(),
				}).(*AndGroup),
			},
		}).(*OrGroup)
	} else if expression.Operator.TokenType == tokenizer.EqualToken ||
		expression.Operator.TokenType == tokenizer.NotEqualToken {
		leftSumGroup, _ := state.NormalizeToSumGroup(expression.Left)
		rightSumGroup, _ := state.NormalizeToSumGroup(expression.Right)

		if leftSumGroup == nil || rightSumGroup == nil {
			return &OrGroup{nil}
		}

		rightSumGroup = state.negateSumGroup(rightSumGroup)

		joined := state.addSumGroups(leftSumGroup, rightSumGroup, 0)

		if joined.IsZero() {
			return &OrGroup{nil}
		}

		var joinedNegate = state.negateSumGroup(joined)

		if joined.Compare(joinedNegate) > 0 {
			var tmp = joined
			joined = joinedNegate
			joinedNegate = tmp
		}

		result, _ := state.nodeCache.GetNodeSingleton(&OrGroup{
			[]*AndGroup{
				state.nodeCache.GetNodeSingleton(&AndGroup{
					[]*SumGroup{
						state.recordExpressionMapping(joined, expression),
						state.recordExpressionMapping(joinedNegate, expression),
					},
					state.getNextUniqueId(),
				}).(*AndGroup),
			},
		}).(*OrGroup)

		if expression.Operator.TokenType == tokenizer.NotEqualToken {
			result = state.NotOrGroup(result)
		}

		return result
	}

	return &OrGroup{nil}
}

func (state *NormalizerState) NormalizeToOrGroup(expression parser.Expression) *OrGroup {
	asBinaryExpression, ok := expression.(*parser.BinaryExpression)

	if ok {
		return state.normalizeBinaryExpressionToOrGroup(asBinaryExpression)
	}

	return &OrGroup{nil}
}

func (state *NormalizerState) normalizeUnaryExpressionToSumGroup(expression *parser.UnaryExpression) (result *SumGroup, err error) {
	expr, err := state.NormalizeToSumGroup(expression.Expr)

	if err != nil {
		return nil, err
	}

	if expression.Operator.TokenType == tokenizer.MinusToken {
		return state.negateSumGroup(expr), nil
	} else {
		return nil, errors.New("Could not apply operator " + expression.Operator.Value)
	}
}

func (state *NormalizerState) normalizeBinaryExpressionToSumGroup(expression *parser.BinaryExpression) (result *SumGroup, err error) {
	left, err := state.NormalizeToSumGroup(expression.Left)

	if err != nil {
		return nil, err
	}

	right, err := state.NormalizeToSumGroup(expression.Right)

	if err != nil {
		return nil, err
	}

	if expression.Operator.TokenType == tokenizer.AddToken {
		return state.addSumGroups(left, right, 0), nil
	} else if expression.Operator.TokenType == tokenizer.MinusToken {
		return state.addSumGroups(left, state.negateSumGroup(right), 0), nil
	} else if expression.Operator.TokenType == tokenizer.MultiplyToken {
		return state.multiplySumGroups(left, right), nil
	} else {
		return nil, errors.New("Could not combine operator " + expression.Operator.Value)
	}
}

func (state *NormalizerState) NormalizeToNode(expression parser.Expression) (result NormalizedNode, err error) {
	asIdentifier, ok := expression.(*parser.Identifier)

	if ok {
		return state.nodeCache.GetNodeSingleton(&VariableReference{
			asIdentifier.Token.Value,
			state.identfierSourceMapping[asIdentifier.Token.Value].UniqueId,
		}), nil
	}

	asProperty, ok := expression.(*parser.PropertyExpression)

	if ok {
		subNode, err := state.NormalizeToNode(asProperty.Left)

		if err != nil {
			return nil, err
		}

		return state.nodeCache.GetNodeSingleton(&PropertyReference{
			subNode,
			asProperty.Property.Value,
			0,
		}), nil
	}

	return nil, errors.New("Cannot convert to node")
}

func (state *NormalizerState) NormalizeToSumGroup(expression parser.Expression) (result *SumGroup, err error) {
	asBinaryExpression, ok := expression.(*parser.BinaryExpression)

	if ok {
		return state.normalizeBinaryExpressionToSumGroup(asBinaryExpression)
	}

	asUnaryExpression, ok := expression.(*parser.UnaryExpression)

	if ok {
		return state.normalizeUnaryExpressionToSumGroup(asUnaryExpression)
	}

	asNumber, ok := expression.(*parser.Number)

	if ok {
		parsedNumber, err := strconv.Atoi(asNumber.Token.Value)

		if err != nil {
			return nil, err
		}

		return state.nodeCache.GetNodeSingleton(&SumGroup{
			nil,
			int64(parsedNumber),
			state.getNextUniqueId(),
		}).(*SumGroup), nil
	}

	asNode, err := state.NormalizeToNode(expression)

	if err != nil {
		return nil, err
	}

	return state.sumGroupFromNode(asNode), nil
}

func (state *NormalizerState) sumGroupFromNode(node NormalizedNode) *SumGroup {
	var productGroup = state.nodeCache.GetNodeSingleton(&ProductGroup{
		state.nodeCache.GetNodeSingleton(&NormalizedNodeArray{
			[]NormalizedNode{node},
			state.getNextUniqueId(),
		}).(*NormalizedNodeArray),
		zmath.Ri64_1(),
	}).(*ProductGroup)

	var sumGroup = state.nodeCache.GetNodeSingleton(&SumGroup{
		[]*ProductGroup{productGroup},
		int64(0),
		state.getNextUniqueId(),
	}).(*SumGroup)

	return sumGroup
}
