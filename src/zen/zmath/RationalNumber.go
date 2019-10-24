package zmath

import (
	"strconv"
)

type RationalNumberi64 struct {
	Numerator   int64
	Denominator int64
}

func Ri64Fromi64(value int64) RationalNumberi64 {
	return RationalNumberi64{
		value,
		1,
	}
}

func Ri64_0() RationalNumberi64 {
	return RationalNumberi64{0, 1}
}

func Ri64_1() RationalNumberi64 {
	return RationalNumberi64{1, 1}
}

func Absi64(a int64) int64 {
	if a > 0 {
		return a
	} else {
		return -a
	}
}

func Signi64(a int64) int64 {
	if a > 0 {
		return 1
	} else if a < 0 {
		return -1
	} else {
		return 0
	}
}

func GcdRi64(a int64, b int64) int64 {
	biggerValue := Absi64(a)
	smallerValue := Absi64(b)

	if biggerValue == smallerValue {
		return biggerValue
	} else if biggerValue == 0 || smallerValue == 0 {
		return 0
	}

	var tmp int64

	if biggerValue < smallerValue {
		tmp = biggerValue
		biggerValue = smallerValue
		biggerValue = tmp
	}

	for smallerValue > 0 {
		tmp = biggerValue % smallerValue
		biggerValue = smallerValue
		smallerValue = tmp
	}

	return biggerValue
}

func (number RationalNumberi64) SimplifyRi64() RationalNumberi64 {
	gdc := GcdRi64(number.Numerator, number.Denominator)

	if number.Denominator < 0 {
		gdc = -gdc
	}

	if gdc != 0 {
		number.Numerator /= gdc
		number.Denominator /= gdc
	} else {
		number.Numerator = Signi64(number.Numerator)
		number.Denominator = Signi64(number.Denominator)
	}

	return number
}

func (number RationalNumberi64) IsZero() bool {
	return number.Numerator == 0
}

func (number RationalNumberi64) IsOne() bool {
	return number.Numerator == number.Denominator
}

func (number RationalNumberi64) ToString() string {
	if number.Denominator == 0 {
		return "NaN{" + strconv.Itoa(int(number.Numerator)) + "}"
	} else if number.Numerator == 0 {
		return "0"
	} else if number.Denominator == 1 {
		return strconv.Itoa(int(number.Numerator))
	} else {
		return strconv.Itoa(int(number.Numerator)) + "/" + strconv.Itoa(int(number.Denominator))
	}
}

func AbsRi64(a RationalNumberi64) RationalNumberi64 {
	if a.Numerator < 0 {
		return RationalNumberi64{
			-a.Numerator,
			a.Denominator,
		}
	} else {
		return a
	}
}

func AddRi64(a RationalNumberi64, b RationalNumberi64) RationalNumberi64 {
	return RationalNumberi64{
		a.Numerator*b.Denominator + b.Numerator*a.Denominator,
		a.Denominator * b.Denominator,
	}
}

func SubRi64(a RationalNumberi64, b RationalNumberi64) RationalNumberi64 {
	return RationalNumberi64{
		a.Numerator*b.Denominator - b.Numerator*a.Denominator,
		a.Denominator * b.Denominator,
	}
}

func MulRi64(a RationalNumberi64, b RationalNumberi64) RationalNumberi64 {
	return RationalNumberi64{
		a.Numerator * b.Numerator,
		a.Denominator * b.Denominator,
	}
}

func Muli64(a RationalNumberi64, scalar int64) RationalNumberi64 {
	return RationalNumberi64{
		a.Numerator * scalar,
		a.Denominator,
	}
}

func DivRi64(a RationalNumberi64, b RationalNumberi64) RationalNumberi64 {
	return RationalNumberi64{
		a.Numerator * b.Denominator,
		a.Denominator * b.Numerator,
	}
}

func InvRi64(a RationalNumberi64) RationalNumberi64 {
	return RationalNumberi64{
		a.Denominator,
		a.Numerator,
	}
}

func NegateRi64(a RationalNumberi64) RationalNumberi64 {
	return RationalNumberi64{
		-a.Numerator,
		a.Denominator,
	}
}

func (a RationalNumberi64) Compare(b RationalNumberi64) int {
	return int(a.Numerator*b.Denominator - b.Numerator*a.Denominator)
}
