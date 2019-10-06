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
			state.normalizeToOrGroup(expression.Left),
			state.normalizeToOrGroup(expression.Right),
		)
	} else if expression.Operator.TokenType == tokenizer.BooleanAndToken {
		return state.combineOrGroupsWithAnd(
			state.normalizeToOrGroup(expression.Left),
			state.normalizeToOrGroup(expression.Right),
		)
	} else if expression.Operator.TokenType == tokenizer.LTToken ||
		expression.Operator.TokenType == tokenizer.LTEqToken ||
		expression.Operator.TokenType == tokenizer.GTToken ||
		expression.Operator.TokenType == tokenizer.GTEqToken {
		leftSumGroup, _ := state.normalizeToSumGroup(expression.Left)
		rightSumGroup, _ := state.normalizeToSumGroup(expression.Right)

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

		var sumGroup = state.addSumGroups(leftSumGroup, rightSumGroup, offset)
		return state.nodeCache.GetNodeSingleton(&OrGroup{
			[]*AndGroup{
				state.nodeCache.GetNodeSingleton(&AndGroup{
					[]*SumGroup{sumGroup},
				}).(*AndGroup),
			},
		}).(*OrGroup)
	} else if expression.Operator.TokenType == tokenizer.EqualToken ||
		expression.Operator.TokenType == tokenizer.NotEqualToken {
		leftSumGroup, _ := state.normalizeToSumGroup(expression.Left)
		rightSumGroup, _ := state.normalizeToSumGroup(expression.Right)

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
						joined,
						joinedNegate,
					},
				}).(*AndGroup),
			},
		}).(*OrGroup)

		if expression.Operator.TokenType == tokenizer.NotEqualToken {
			result = state.notOrGroup(result)
		}

		return result
	}

	return &OrGroup{nil}
}

func (state *NormalizerState) normalizeToOrGroup(expression parser.Expression) *OrGroup {
	asBinaryExpression, ok := expression.(*parser.BinaryExpression)

	if ok {
		return state.normalizeBinaryExpressionToOrGroup(asBinaryExpression)
	}

	return &OrGroup{nil}
}

func (state *NormalizerState) normalizeUnaryExpressionToSumGroup(expression *parser.UnaryExpression) (result *SumGroup, err error) {
	expr, err := state.normalizeToSumGroup(expression.Expr)

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
	left, err := state.normalizeToSumGroup(expression.Left)

	if err != nil {
		return nil, err
	}

	right, err := state.normalizeToSumGroup(expression.Right)

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

func (state *NormalizerState) normalizeToSumGroup(expression parser.Expression) (result *SumGroup, err error) {
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
		}).(*SumGroup), nil
	}

	asIdentifier, ok := expression.(*parser.Identifier)

	if ok {
		var varRef = state.nodeCache.GetNodeSingleton(&VariableReference{
			asIdentifier.Token.Value,
			state.identfierSourceMapping[asIdentifier.Token.Value],
		}).(*VariableReference)

		var productGroup = state.nodeCache.GetNodeSingleton(&ProductGroup{
			state.nodeCache.GetNodeSingleton(&NormalizedNodeArray{
				[]NormalizedNode{varRef},
				state.getNextUniqueId(),
			}).(*NormalizedNodeArray),
			zmath.Ri64_1(),
		}).(*ProductGroup)

		var sumGroup = state.nodeCache.GetNodeSingleton(&SumGroup{
			[]*ProductGroup{productGroup},
			int64(0),
		}).(*SumGroup)

		return sumGroup, nil
	}

	return nil, errors.New("Could not convert to sum group")
}
