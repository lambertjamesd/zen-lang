package boundschecking

import (
	"testing"
	"zen/test"
	"zen/zmath"
)

func TestEqualityChecks(t *testing.T) {
	var a = &VariableReference{
		"a",
		0,
	}

	test.Assert(t, a.Compare(&VariableReference{
		"a",
		0,
	}) == 0, "Variable references should match")

	test.Assert(t, a.Compare(&VariableReference{
		"a",
		10,
	}) < 0, "Variable references compare on reference id")

	test.Assert(t, a.Compare(&VariableReference{
		"a",
		-10,
	}) > 0, "Variable references compare on reference id")

	var sumGroupA = &SumGroup{
		[]*ProductGroup{&ProductGroup{
			&NormalizedNodeArray{
				[]NormalizedNode{
					a,
				},
				10,
			},
			zmath.Ri64Fromi64(-5),
		}},
		int64(10),
	}

	var sumGroupB = &SumGroup{
		[]*ProductGroup{&ProductGroup{
			&NormalizedNodeArray{
				[]NormalizedNode{
					a,
				},
				20,
			},
			zmath.Ri64Fromi64(-5),
		}},
		int64(10),
	}

	var sumGroupC = &SumGroup{
		[]*ProductGroup{&ProductGroup{
			&NormalizedNodeArray{
				[]NormalizedNode{
					a,
				},
				20,
			},
			zmath.Ri64Fromi64(-5),
		}},
		int64(11),
	}

	var sumGroupD = &SumGroup{
		[]*ProductGroup{&ProductGroup{
			&NormalizedNodeArray{
				[]NormalizedNode{
					a,
				},
				20,
			},
			zmath.Ri64Fromi64(-10),
		}},
		int64(10),
	}

	test.Assert(t, sumGroupA.Compare(sumGroupB) == 0, "Sum group equality")
	test.Assert(t, sumGroupB.GetHashCode() == sumGroupB.GetHashCode(), "Sum group hash codes should match")
	test.Assert(t, sumGroupB.GetHashCode() != sumGroupC.GetHashCode(), "Different sum group hash codes should not match")
	test.Assert(t, sumGroupB.GetHashCode() != sumGroupD.GetHashCode(), "Different sum group hash codes should not match")
}
