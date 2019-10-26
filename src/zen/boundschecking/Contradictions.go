package boundschecking

func FindContraditions(orGroup *OrGroup) ([]*SumGroup, error) {
	var result []*SumGroup = nil

	for _, andGroup := range orGroup.AndGroups {
		var checker = NewKnownConstraints()

		for _, sumGroup := range andGroup.SumGroups {
			addResult, err := checker.InsertSumGroup(sumGroup)

			if err != nil {
				return nil, err
			}

			if !addResult {
				result = append(result, sumGroup)
			}
		}
	}

	return result, nil
}
