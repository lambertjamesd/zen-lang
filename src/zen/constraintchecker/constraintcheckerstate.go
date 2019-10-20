package constraintchecker

import (
	"zen/boundschecking"
)

type ConstraintCheckerState struct {
	knownConstraints []*boundschecking.KnownConstraints
}

func NewConstraintCheckerState() *ConstraintCheckerState {
	return &ConstraintCheckerState{
		[]*boundschecking.KnownConstraints{boundschecking.NewKnownConstraints()},
	}
}

func (state *ConstraintCheckerState) Copy() *ConstraintCheckerState {
	var constraintCopy []*boundschecking.KnownConstraints = nil

	for _, bounds := range state.knownConstraints {
		constraintCopy = append(constraintCopy, bounds.Copy())
	}

	return &ConstraintCheckerState{
		constraintCopy,
	}
}

func (state *ConstraintCheckerState) addSumGroups(newRules []*boundschecking.SumGroup) (bool, error) {
	var nextRules []*boundschecking.KnownConstraints = nil

	for _, sumGroup := range newRules {
		for _, existingConstraint := range state.knownConstraints {
			isValid, err := existingConstraint.InsertSumGroup(sumGroup)

			if err != nil {
				return false, err
			} else if isValid {
				nextRules = append(nextRules, existingConstraint)
			}
		}
	}

	state.knownConstraints = nextRules

	return true, nil
}

func (state *ConstraintCheckerState) addRules(newRules []*boundschecking.AndGroup) (bool, error) {
	if len(newRules) == 0 {
		return true, nil
	}

	var nextRules []*boundschecking.KnownConstraints = nil

	for newRuleIndex, andGroup := range newRules {
		for _, sumGroup := range andGroup.SumGroups {
			for _, existingConstraint := range state.knownConstraints {
				var nextConstraint = existingConstraint

				if newRuleIndex != len(newRules)-1 {
					nextConstraint = nextConstraint.Copy()
				}

				isValid, err := nextConstraint.InsertSumGroup(sumGroup)

				if err != nil {
					return false, err
				} else if isValid {
					nextRules = append(nextRules, nextConstraint)
				}
			}
		}
	}

	state.knownConstraints = nextRules

	return len(state.knownConstraints) != 0, nil
}

func checkAndGroup(knownConstraint *boundschecking.KnownConstraints, sumGroupCache map[uint32]bool, rulesCheck *boundschecking.AndGroup) ([]*boundschecking.SumGroup, error) {
	var falseSumGroups []*boundschecking.SumGroup = nil
	for _, sumGroup := range rulesCheck.SumGroups {
		var isSumGroupTrue, ok = sumGroupCache[sumGroup.GetUniqueId()]

		if !ok {
			checkResult, err := knownConstraint.CheckSumGroup(sumGroup)

			if err != nil {
				return nil, err
			}

			isSumGroupTrue = checkResult.IsTrue
			sumGroupCache[sumGroup.GetUniqueId()] = checkResult.IsTrue
		}

		if !isSumGroupTrue {
			falseSumGroups = append(falseSumGroups, sumGroup)
		}
	}

	return falseSumGroups, nil
}

func checkOrGroup(knownConstraint *boundschecking.KnownConstraints, sumGroupCache map[uint32]bool, rulesCheck *boundschecking.OrGroup) ([]*boundschecking.SumGroup, error) {
	var result []*boundschecking.SumGroup = nil

	for _, andGroup := range rulesCheck.AndGroups {
		andCheck, err := checkAndGroup(knownConstraint, sumGroupCache, andGroup)

		if err != nil {
			return nil, nil
		}

		if len(andCheck) == 0 {
			return nil, nil
		} else {
			result = append(result, andCheck...)
		}
	}

	return result, nil
}

func (state *ConstraintCheckerState) checkOrGroup(rulesCheck *boundschecking.OrGroup) ([]*boundschecking.SumGroup, error) {
	var result []*boundschecking.SumGroup = nil
	for _, knownConstraints := range state.knownConstraints {
		var sumGroupCache = make(map[uint32]bool)
		checkResult, err := checkOrGroup(knownConstraints, sumGroupCache, rulesCheck)

		if err != nil {
			return nil, err
		}

		if len(checkResult) == 0 {
			return nil, nil
		} else {
			result = append(result, checkResult...)
		}
	}

	return result, nil
}

func (state *ConstraintCheckerState) checkAndGroup(rulesCheck *boundschecking.AndGroup) ([]*boundschecking.SumGroup, error) {
	var result []*boundschecking.SumGroup = nil
	for _, knownConstraints := range state.knownConstraints {
		var sumGroupCache = make(map[uint32]bool)
		checkResult, err := checkAndGroup(knownConstraints, sumGroupCache, rulesCheck)

		if err != nil {
			return nil, nil
		}

		if len(checkResult) > 0 {
			result = append(result, checkResult...)
		}
	}

	return result, nil
}
