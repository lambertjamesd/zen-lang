package zmath

type RationalNumberi64 struct {
	numerator   int64
	denominator int64
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
	gdc := GcdRi64(number.numerator, number.denominator)

	if number.denominator < 0 {
		gdc = -gdc
	}

	if gdc != 0 {
		number.numerator /= gdc
		number.denominator /= gdc
	} else {
		number.numerator = Signi64(number.numerator)
		number.denominator = Signi64(number.denominator)
	}

	return number
}

func AddRi64(a RationalNumberi64, b RationalNumberi64) RationalNumberi64 {
	return RationalNumberi64{
		a.numerator*b.denominator + b.numerator*a.denominator,
		a.denominator * b.denominator,
	}
}

func SubRi64(a RationalNumberi64, b RationalNumberi64) RationalNumberi64 {
	return RationalNumberi64{
		a.numerator*b.denominator - b.numerator*a.denominator,
		a.denominator * b.denominator,
	}
}

func MulRi64(a RationalNumberi64, b RationalNumberi64) RationalNumberi64 {
	return RationalNumberi64{
		a.numerator * b.numerator,
		a.denominator * b.denominator,
	}
}

func DivRi64(a RationalNumberi64, b RationalNumberi64) RationalNumberi64 {
	return RationalNumberi64{
		a.numerator * b.denominator,
		a.denominator * b.numerator,
	}
}
