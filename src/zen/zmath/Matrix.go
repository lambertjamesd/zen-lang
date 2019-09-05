package zmath

import "errors"

type Matrixi64 struct {
	data        []RationalNumberi64
	Rows        uint32
	Cols        uint32
	rowCapacity uint32
	colCapacity uint32
}

func PowOfTwoi32(a uint32) uint32 {
	var result uint32 = 1

	if a == 0 {
		return 0
	}

	a = a - 1

	for a > 0 {
		result = result << 1
		a = a >> 1
	}

	return result
}

func NewMatrixi64(Rows uint32, Cols uint32) *Matrixi64 {
	rowCapacity := PowOfTwoi32(Rows)
	colCapacity := PowOfTwoi32(Cols)

	return &Matrixi64{
		make([]RationalNumberi64, rowCapacity*colCapacity, rowCapacity*colCapacity),
		Rows,
		Cols,
		rowCapacity,
		colCapacity,
	}
}

func NewMatrixi64WithData(Rows uint32, Cols uint32, data []RationalNumberi64) *Matrixi64 {
	rowCapacity := PowOfTwoi32(Rows)
	colCapacity := PowOfTwoi32(Cols)

	var result = &Matrixi64{
		make([]RationalNumberi64, rowCapacity*colCapacity, rowCapacity*colCapacity),
		Rows,
		Cols,
		rowCapacity,
		colCapacity,
	}

	for row := uint32(0); row < result.Rows; row = row + 1 {
		for col := uint32(0); col < result.Cols; col = col + 1 {
			result.SetEntryi64(row, col, data[row*Cols+col])
		}
	}

	return result
}

func (matrix *Matrixi64) InitialzeIdentityi64() {
	for row := uint32(0); row < matrix.Rows; row = row + 1 {
		for col := uint32(0); col < matrix.Cols; col = col + 1 {
			if row == col {
				matrix.data[row*matrix.colCapacity+col] = Ri64_1()
			} else {
				matrix.data[row*matrix.colCapacity+col] = Ri64_0()
			}
		}
	}
}

func (matrix *Matrixi64) Resize(Rows uint32, Cols uint32) {
	rowCapacity := PowOfTwoi32(Rows)
	colCapacity := PowOfTwoi32(Cols)

	if rowCapacity > matrix.rowCapacity || colCapacity > matrix.colCapacity {
		newData := make([]RationalNumberi64, rowCapacity, colCapacity)

		for row := uint32(0); row < matrix.Rows; row = row + 1 {
			for col := uint32(0); col < matrix.Cols; col = col + 1 {
				newData[row*colCapacity+col] = matrix.data[row*matrix.colCapacity+col]
			}
		}

		matrix.data = newData
		matrix.rowCapacity = rowCapacity
		matrix.colCapacity = colCapacity
	}

	if Rows > matrix.Rows {
		for row := matrix.Rows; row < Rows; row = row + 1 {
			for col := uint32(0); col < Cols; col = col + 1 {
				if row == col {
					matrix.SetEntryi64(row, col, Ri64_1())
				} else {
					matrix.SetEntryi64(row, col, Ri64_0())
				}
			}
		}
	}

	if Cols > matrix.Cols {
		for row := uint32(0); row < matrix.Rows; row = row + 1 {
			for col := matrix.Cols; col < Cols; col = col + 1 {
				if row == col {
					matrix.SetEntryi64(row, col, Ri64_1())
				} else {
					matrix.SetEntryi64(row, col, Ri64_0())
				}
			}
		}
	}

	matrix.Rows = Rows
	matrix.Cols = Cols
}

func (matrix *Matrixi64) GetEntryi64(row uint32, col uint32) (result RationalNumberi64) {
	return matrix.data[row*matrix.colCapacity+col]
}

func (matrix *Matrixi64) SetEntryi64(row uint32, col uint32, value RationalNumberi64) {
	matrix.data[row*matrix.colCapacity+col] = value
}

func (a *Matrixi64) Muli64(b *Matrixi64) (result *Matrixi64, err error) {
	if a.Cols != b.Rows {
		return nil, errors.New("Matrix sizes are not compatible for multiplication")
	}

	result = NewMatrixi64(a.Rows, b.Cols)

	for row := uint32(0); row < result.Rows; row = row + 1 {
		for col := uint32(0); col < result.Cols; col = col + 1 {
			var rowValue RationalNumberi64 = Ri64_0()

			for span := uint32(0); span < a.Cols; span = span + 1 {
				rowValue = AddRi64(rowValue, MulRi64(a.GetEntryi64(row, span), b.GetEntryi64(span, col)))
			}

			rowValue.SimplifyRi64()

			result.SetEntryi64(row, col, rowValue)
		}
	}

	return result, nil
}

func (matrix *Matrixi64) Scalei64(value RationalNumberi64) *Matrixi64 {
	var result = NewMatrixi64(matrix.Rows, matrix.Cols)

	for row := uint32(0); row < result.Rows; row = row + 1 {
		for col := uint32(0); col < result.Cols; col = col + 1 {
			result.SetEntryi64(row, col, MulRi64(matrix.GetEntryi64(row, col), value).SimplifyRi64())
		}
	}

	return result
}

func (matrix *Matrixi64) GetRow(row uint32) *Matrixi64 {
	var result = NewMatrixi64(1, matrix.Cols)

	for col := uint32(0); col < matrix.Cols; col = col + 1 {
		result.SetEntryi64(0, col, matrix.GetEntryi64(row, col))
	}

	return result
}

func (matrix *Matrixi64) ScaleRowi64(row uint32, value RationalNumberi64) {
	for col := uint32(0); col < matrix.Cols; col = col + 1 {
		matrix.SetEntryi64(row, col, MulRi64(matrix.GetEntryi64(row, col), value).SimplifyRi64())
	}
}

func (matrix *Matrixi64) AddRowToRow(fromRow uint32, toRow uint32, scalar RationalNumberi64) {
	for col := uint32(0); col < matrix.Cols; col = col + 1 {
		matrix.SetEntryi64(
			toRow,
			col,
			AddRi64(
				matrix.GetEntryi64(toRow, col),
				MulRi64(matrix.GetEntryi64(fromRow, col), scalar),
			).SimplifyRi64(),
		)
	}
}

func (matrix *Matrixi64) Copy() *Matrixi64 {
	result := NewMatrixi64(matrix.Rows, matrix.Cols)

	for row := uint32(0); row < result.Rows; row = row + 1 {
		for col := uint32(0); col < result.Cols; col = col + 1 {
			result.SetEntryi64(row, col, matrix.GetEntryi64(row, col))
		}
	}

	return result
}
