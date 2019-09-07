package zmath

import (
	"testing"
)

func checkValue(t *testing.T, actual RationalNumberi64, expected RationalNumberi64) {
	if actual.Numerator != expected.Numerator || actual.Denominator != expected.Denominator {
		t.Errorf("Expected %d/%d to equal %d/%d", expected.Numerator, expected.Denominator, actual.Numerator, actual.Denominator)
	}
}

func TestAddition(t *testing.T) {
	var one = Ri64Fromi64(1)
	var two = Ri64Fromi64(2)

	checkValue(t, AddRi64(one, two), Ri64Fromi64(3))
}
