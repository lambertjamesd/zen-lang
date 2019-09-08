package boundschecking

import (
	"errors"
	"zen/parser"
)

func checkBounds() {

}

type NormalizerState struct {
	nodeCache NodeCache
	currentUniqueID uint32
}

func (state *NormalizerState) NormalizeExpression(expression parser.Expression) *NormalizedEquation {
	return nil
}

func (state *NormalizerState) normalizeToSumGroup(expression parser.Expression) (result *SumGroup, err error) {
	asBinaryExpression, ok := expression.(*parser.BinaryExpression)

	if ok {
		return nil, nil
	}

	asNumber, ok := expression.(*parser.Number)

	if ok {
		return nil, nil
	}

	asIdentifier, ok := expression.(*parser.Identifier)

	if ok {
		return nil, nil
	}

	return nil, errors.New("Could not convert to sum group")
}