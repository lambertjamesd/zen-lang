package constraintchecker

import (
	"zen/boundschecking"
)

type ConstraintChecker struct {
	checkerStateStack []*ConstraintCheckerState
	normalizerState   *boundschecking.NormalizerState
}

func NewConstrantChecker() *ConstraintChecker {
	return &ConstraintChecker{
		nil,
		boundschecking.NewNormalizerState(),
	}
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
