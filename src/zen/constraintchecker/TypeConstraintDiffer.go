package constraintchecker

import (
	"sort"
	"strings"
	"zen/boundschecking"
	"zen/parser"
)

type typeVariables struct {
	name string
	at   int
}

type TypeConstraints struct {
	constraints       *boundschecking.OrGroup
	variables         []typeVariables
	contradictions    []*boundschecking.SumGroup
	expressionMapping map[uint32]parser.Expression
}

type TypeConstraintDiff struct {
	additionalConstrants *boundschecking.OrGroup
	expressionMapping    map[uint32]parser.Expression
}

type TypeConstraintDifferCache struct {
	state     *boundschecking.NormalizerState
	mapping   map[int]TypeConstraints
	diffCache map[int64]TypeConstraintDiff
}

func NewTypeConstraintDifferCache(normalizerState *boundschecking.NormalizerState) *TypeConstraintDifferCache {
	return &TypeConstraintDifferCache{
		normalizerState,
		make(map[int]TypeConstraints),
		make(map[int64]TypeConstraintDiff),
	}
}

func variablesForStructure(structure *parser.StructureTypeType) []typeVariables {
	var result []typeVariables = make([]typeVariables, 0, len(structure.Entries))

	for _, entry := range structure.Entries {
		if entry.Name != "" {
			result = append(result, typeVariables{
				entry.Name,
				entry.UniqueId,
			})
		}
	}

	sort.Slice(result, func(a, b int) bool {
		return strings.Compare(result[a].name, result[b].name) < 0
	})

	return result
}

func typeConstraintsFromType(state *boundschecking.NormalizerState, typeNode parser.TypeNode) (TypeConstraints, error) {
	state.UseIdentifierMapping("self", typeNode.UniqueId())

	asStruct, isStruct := typeNode.(*parser.StructureTypeType)

	if isStruct {
		for _, entry := range asStruct.Entries {
			state.UseIdentifierMapping(entry.Name, entry.UniqueId)
		}
	}

	var expressionMapping = state.StartTrackingExpressionMapping()
	var constraints = state.NormalizeToOrGroup(typeNode.GetWhereExpression())
	state.StopTrackingExpressionMapping()

	contradictions, err := boundschecking.FindContraditions(constraints)

	if err != nil {
		return TypeConstraints{}, err
	}

	if isStruct {
		return TypeConstraints{
			constraints,
			variablesForStructure(asStruct),
			contradictions,
			expressionMapping,
		}, nil
	} else {
		return TypeConstraints{
			constraints,
			[]typeVariables{typeVariables{"self", typeNode.UniqueId()}},
			contradictions,
			expressionMapping,
		}, nil
	}
}

func getVariableMapping(from []typeVariables, to []typeVariables) map[int]int {
	var result = make(map[int]int)

	var aIndex = 0
	var bIndex = 0

	for aIndex < len(from) && bIndex < len(to) {
		var nameCompare = strings.Compare(from[aIndex].name, to[bIndex].name)

		if nameCompare == 0 {
			result[from[aIndex].at] = to[bIndex].at
			aIndex = aIndex + 1
			bIndex = bIndex + 1
		} else if nameCompare < 0 {
			aIndex = aIndex + 1
		} else {
			bIndex = bIndex + 1
		}
	}

	return result
}

func (typeConstraintCache *TypeConstraintDifferCache) GetConstraintsForType(typeNode parser.TypeNode) (TypeConstraints, error) {
	result, ok := typeConstraintCache.mapping[typeNode.UniqueId()]

	if ok {
		return result, nil
	} else {
		result, err := typeConstraintsFromType(typeConstraintCache.state, typeNode)

		if err != nil {
			return TypeConstraints{}, nil
		}

		typeConstraintCache.mapping[typeNode.UniqueId()] = result
		return result, nil
	}
}

func (typeConstraintCache *TypeConstraintDifferCache) GetTypeDiff(from parser.TypeNode, to parser.TypeNode) (TypeConstraintDiff, error) {
	var cacheKey = int64(from.UniqueId()) | (int64(to.UniqueId()) << 32)

	result, ok := typeConstraintCache.diffCache[cacheKey]

	if ok {
		return result, nil
	} else {
		result, err := typeConstraintCache.createTypeDiff(from, to)

		if err != nil {
			return TypeConstraintDiff{}, err
		}

		typeConstraintCache.diffCache[cacheKey] = result
		return result, nil
	}
}

func (typeConstraintCache *TypeConstraintDifferCache) createTypeDiff(from parser.TypeNode, to parser.TypeNode) (TypeConstraintDiff, error) {
	if from == to {
		return TypeConstraintDiff{}, nil
	}

	fromConstraints, err := typeConstraintCache.GetConstraintsForType(from)

	if err != nil {
		return TypeConstraintDiff{}, err
	}

	toConstraints, err := typeConstraintCache.GetConstraintsForType(to)

	if err != nil {
		return TypeConstraintDiff{}, err
	}

	var variableMapping = getVariableMapping(toConstraints.variables, fromConstraints.variables)

	var exprMapping = make(map[uint32]parser.Expression)
	var toConstraintsInFromSpace = toConstraints.constraints.MapNodes(
		typeConstraintCache.state,
		variableMapping,
		toConstraints.expressionMapping,
		exprMapping,
	)
	var constraintsDiff []*boundschecking.AndGroup = nil

	if toConstraintsInFromSpace != nil {
		for _, andGroup := range toConstraintsInFromSpace.AndGroups {
			resultProductGroup, err := typeConstraintCache.findConstraintsForAndGroup(andGroup, fromConstraints)

			if err != nil {
				return TypeConstraintDiff{}, err
			}

			if resultProductGroup != nil {
				constraintsDiff = append(constraintsDiff, resultProductGroup)
			}
		}
	}

	return TypeConstraintDiff{
		typeConstraintCache.state.CreateOrGroup(constraintsDiff),
		exprMapping,
	}, nil
}

func (typeConstraintCache *TypeConstraintDifferCache) findConstraintsForAndGroup(andGroup *boundschecking.AndGroup, fromConstraints TypeConstraints) (*boundschecking.AndGroup, error) {
	var resultProductGroup []*boundschecking.SumGroup = nil

	if fromConstraints.constraints == nil || len(fromConstraints.contradictions) == 0 {
		return andGroup, nil
	}

	for _, fromProductGroup := range fromConstraints.constraints.AndGroups {
		var subConstraints = make([]*boundschecking.SumGroup, 0, len(fromProductGroup.SumGroups))

		var constraints = boundschecking.NewKnownConstraints()

		for _, fromSumGroup := range fromProductGroup.SumGroups {
			_, err := constraints.InsertSumGroup(fromSumGroup)

			if err != nil {
				return nil, err
			}
		}

		for _, sumGroup := range andGroup.SumGroups {
			checkResult, err := constraints.CheckSumGroup(sumGroup)

			if err != nil {
				return nil, err
			}

			if !checkResult.IsTrue {
				subConstraints = append(subConstraints, sumGroup)
			}
		}

		resultProductGroup = boundschecking.ZipSumGroups(resultProductGroup, subConstraints)
	}

	if len(resultProductGroup) != 0 {
		return typeConstraintCache.state.CreateAndGroup(resultProductGroup), nil
	} else {
		return nil, nil
	}
}
