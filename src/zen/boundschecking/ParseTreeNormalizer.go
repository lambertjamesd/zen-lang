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
		return state.addSumGroups(left, right), nil
	} else if expression.Operator.TokenType == tokenizer.MinusToken {
		return state.addSumGroups(left, state.negateSumGroup(right)), nil
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
