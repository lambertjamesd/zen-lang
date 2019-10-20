package zmath

import (
	"testing"
)

func checkMatrix(t *testing.T, actual *Matrixi64, expected []RationalNumberi64) {
	for row := uint32(0); row < actual.Rows; row = row + 1 {
		for col := uint32(0); col < actual.Cols; col = col + 1 {
			var actualVal = actual.GetEntryi64(row, col)
			var expectedVal = expected[row*actual.Cols+col]
			if actualVal.Numerator != expectedVal.Numerator || actualVal.Denominator != expectedVal.Denominator {
				t.Errorf(
					"Expected %d/%d to equal %d/%d at %d %d",
					expectedVal.Numerator,
					expectedVal.Denominator,
					actualVal.Numerator,
					actualVal.Denominator,
					row,
					col,
				)
			}
		}
	}
}

func TestPowOfTwoi32(t *testing.T) {
	if PowOfTwoi32(0) != 0 {
		t.Errorf("Bad A")
	}

	if PowOfTwoi32(1) != 1 {
		t.Errorf("Bad B")
	}

	if PowOfTwoi32(2) != 2 {
		t.Errorf("Bad C")
	}

	if PowOfTwoi32(3) != 4 {
		t.Errorf("Bad D")
	}
}

func TestMatrixCreation(t *testing.T) {
	var mat = NewMatrixi64(4, 3)

	if mat.Rows != 4 {
		t.Errorf("Wrong number of rows")
	}

	if mat.Cols != 3 {
		t.Errorf("Wrong number of rows")
	}

	if mat.rowCapacity != 4 {
		t.Errorf("Wrong row capacity")
	}

	if mat.colCapacity != 4 {
		t.Errorf("Wrong col capacity got %d", mat.colCapacity)
	}

	mat.InitialzeIdentityi64()

	checkMatrix(t, mat, []RationalNumberi64{
		Ri64Fromi64(1), Ri64Fromi64(0), Ri64Fromi64(0),
		Ri64Fromi64(0), Ri64Fromi64(1), Ri64Fromi64(0),
		Ri64Fromi64(0), Ri64Fromi64(0), Ri64Fromi64(1),
		Ri64Fromi64(0), Ri64Fromi64(0), Ri64Fromi64(0),
	})

	var withData = NewMatrixi64WithData(2, 3, []RationalNumberi64{
		Ri64Fromi64(-1), Ri64Fromi64(4), Ri64Fromi64(0),
		Ri64Fromi64(3), Ri64Fromi64(-1), Ri64Fromi64(1),
	})

	checkMatrix(t, withData, []RationalNumberi64{
		Ri64Fromi64(-1), Ri64Fromi64(4), Ri64Fromi64(0),
		Ri64Fromi64(3), Ri64Fromi64(-1), Ri64Fromi64(1),
	})
}

func TestRowAddition(t *testing.T) {
	var mat = NewMatrixi64(3, 3)

	mat.InitialzeIdentityi64()

	mat.AddRowToRow(0, 1, Ri64Fromi64(-1))
	mat.AddRowToRow(0, 2, Ri64Fromi64(2))
	mat.ScaleRowi64(0, Ri64Fromi64(-2))

	checkMatrix(t, mat, []RationalNumberi64{
		Ri64Fromi64(-2), Ri64Fromi64(0), Ri64Fromi64(0),
		Ri64Fromi64(-1), Ri64Fromi64(1), Ri64Fromi64(0),
		Ri64Fromi64(2), Ri64Fromi64(0), Ri64Fromi64(1),
	})

	mat.AddRowToRow(1, 0, Ri64Fromi64(-1))
	mat.ScaleRowi64(1, Ri64Fromi64(-1))

	checkMatrix(t, mat, []RationalNumberi64{
		Ri64Fromi64(-1), Ri64Fromi64(-1), Ri64Fromi64(0),
		Ri64Fromi64(1), Ri64Fromi64(-1), Ri64Fromi64(0),
		Ri64Fromi64(2), Ri64Fromi64(0), Ri64Fromi64(1),
	})
}

func TestMatrixCopy(t *testing.T) {
	var withData = NewMatrixi64WithData(2, 3, []RationalNumberi64{
		Ri64Fromi64(-1), Ri64Fromi64(4), Ri64Fromi64(0),
		Ri64Fromi64(3), Ri64Fromi64(-1), Ri64Fromi64(1),
	})

	var copy = withData.Copy()

	checkMatrix(t, copy, []RationalNumberi64{
		Ri64Fromi64(-1), Ri64Fromi64(4), Ri64Fromi64(0),
		Ri64Fromi64(3), Ri64Fromi64(-1), Ri64Fromi64(1),
	})
}

func TestMatrixResize(t *testing.T) {
	var withData = NewMatrixi64WithData(3, 3, []RationalNumberi64{
		Ri64Fromi64(-1), Ri64Fromi64(4), Ri64Fromi64(0),
		Ri64Fromi64(3), Ri64Fromi64(-1), Ri64Fromi64(1),
		Ri64Fromi64(-2), Ri64Fromi64(1), Ri64Fromi64(1),
	})

	withData.Resize(4, 4)

	checkMatrix(t, withData, []RationalNumberi64{
		Ri64Fromi64(-1), Ri64Fromi64(4), Ri64Fromi64(0), Ri64Fromi64(0),
		Ri64Fromi64(3), Ri64Fromi64(-1), Ri64Fromi64(1), Ri64Fromi64(0),
		Ri64Fromi64(-2), Ri64Fromi64(1), Ri64Fromi64(1), Ri64Fromi64(0),
		Ri64Fromi64(0), Ri64Fromi64(0), Ri64Fromi64(0), Ri64Fromi64(1),
	})

	withData.Resize(2, 1)

	checkMatrix(t, withData, []RationalNumberi64{
		Ri64Fromi64(-1),
		Ri64Fromi64(3),
	})

	withData.Resize(4, 4)

	checkMatrix(t, withData, []RationalNumberi64{
		Ri64Fromi64(-1), Ri64Fromi64(0), Ri64Fromi64(0), Ri64Fromi64(0),
		Ri64Fromi64(3), Ri64Fromi64(1), Ri64Fromi64(0), Ri64Fromi64(0),
		Ri64Fromi64(0), Ri64Fromi64(0), Ri64Fromi64(1), Ri64Fromi64(0),
		Ri64Fromi64(0), Ri64Fromi64(0), Ri64Fromi64(0), Ri64Fromi64(1),
	})
}

func TestMatrixOrthogonal(t *testing.T) {
	orthoResult, err := OrthogonalVector([]*Matrixi64{NewMatrixi64WithData(2, 1, []RationalNumberi64{
		Ri64Fromi64(1), Ri64Fromi64(1),
	})})

	if err != nil {
		t.Error(err.Error())
	}

	checkMatrix(t, orthoResult, []RationalNumberi64{
		Ri64Fromi64(1), Ri64Fromi64(-1),
	})

	orthoResult, err = OrthogonalVector([]*Matrixi64{NewMatrixi64WithData(2, 1, []RationalNumberi64{
		Ri64Fromi64(-2), Ri64Fromi64(3),
	})})

	if err != nil {
		t.Error(err.Error())
	}

	checkMatrix(t, orthoResult, []RationalNumberi64{
		Ri64Fromi64(3), Ri64Fromi64(2),
	})

	orthoResult, err = OrthogonalVector([]*Matrixi64{
		NewMatrixi64WithData(3, 1, []RationalNumberi64{
			Ri64Fromi64(1), Ri64Fromi64(0), Ri64Fromi64(0),
		}),
		NewMatrixi64WithData(3, 1, []RationalNumberi64{
			Ri64Fromi64(0), Ri64Fromi64(1), Ri64Fromi64(0),
		}),
	})

	if err != nil {
		t.Error(err.Error())
	}

	checkMatrix(t, orthoResult, []RationalNumberi64{
		Ri64Fromi64(0), Ri64Fromi64(0), Ri64Fromi64(1),
	})
}
