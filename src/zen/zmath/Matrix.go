package zmath

import (
	"errors"
	"strconv"
	"strings"
)

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

func (matrix *Matrixi64) InitialzeZeroi64() {
	for row := uint32(0); row < matrix.Rows; row = row + 1 {
		for col := uint32(0); col < matrix.Cols; col = col + 1 {
			matrix.data[row*matrix.colCapacity+col] = Ri64_0()
		}
	}
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
		newData := make([]RationalNumberi64, rowCapacity*colCapacity)

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

func MatrixDot(a *Matrixi64, b *Matrixi64) RationalNumberi64 {
	result := Ri64_0()

	for row := uint32(0); row < a.Rows; row = row + 1 {
		for col := uint32(0); col < b.Cols; col = col + 1 {
			result = AddRi64(result, MulRi64(a.GetEntryi64(row, col), b.GetEntryi64(row, col)))
		}
	}

	return result.SimplifyRi64()
}

func MatrixSum(a *Matrixi64) RationalNumberi64 {
	result := Ri64_0()

	for row := uint32(0); row < a.Rows; row = row + 1 {
		for col := uint32(0); col < a.Cols; col = col + 1 {
			result = AddRi64(result, a.GetEntryi64(row, col))
		}
	}

	return result.SimplifyRi64()
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

			rowValue = rowValue.SimplifyRi64()

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

func (matrix *Matrixi64) subDeterminant(col uint32, size uint32, rowIndices []uint32) RationalNumberi64 {
	if size == uint32(0) {
		return Ri64_0()
	} else if size == uint32(1) {
		return matrix.GetEntryi64(rowIndices[0], col)
	} else {
		var firstIndex = rowIndices[0]

		for i := uint32(0); i < size-1; i = i + 1 {
			rowIndices[i] = rowIndices[i+1]
		}

		rowIndices[size-1] = firstIndex

		var result = Ri64_0()

		for i := uint32(0); i < size; i = i + 1 {
			var scalar RationalNumberi64 = matrix.GetEntryi64(rowIndices[size-1], col)

			if i%2 == 1 {
				scalar = NegateRi64(scalar)
			}

			result = AddRi64(result, MulRi64(scalar, matrix.subDeterminant(col+1, size-1, rowIndices)))

			rowIndices[i], rowIndices[size-1] = rowIndices[size-1], rowIndices[i]
		}

		return result
	}
}

func OrthogonalVector(vectors []*Matrixi64) (result *Matrixi64, err error) {
	var vectorCount = uint32(len(vectors))
	var dimensionCount = vectorCount + 1

	var matrixData = NewMatrixi64(dimensionCount, vectorCount)

	for vectorIndex, vector := range vectors {
		if vector.Cols != 1 {
			return nil, errors.New("Input vectors should be column vectors")
		} else if vector.Rows != dimensionCount {
			return nil, errors.New("Input vector heights should match the number of dimensions")
		}

		for row := uint32(0); row < dimensionCount; row = row + 1 {
			matrixData.SetEntryi64(row, uint32(vectorIndex), vector.GetEntryi64(row, 0))
		}
	}

	result = NewMatrixi64(dimensionCount, 1)

	var rowIndices = make([]uint32, dimensionCount)

	for i := uint32(0); i < dimensionCount; i = i + 1 {
		if i == 0 {
			rowIndices[dimensionCount-1] = i
		} else {
			rowIndices[i-1] = i
		}
	}

	for i := uint32(0); i < dimensionCount; i = i + 1 {
		var subResult = matrixData.subDeterminant(0, vectorCount, rowIndices).SimplifyRi64()

		if i%2 == 1 {
			result.SetEntryi64(i, 0, NegateRi64(subResult))
		} else {
			result.SetEntryi64(i, 0, subResult)
		}

		rowIndices[i], rowIndices[dimensionCount-1] = rowIndices[dimensionCount-1], rowIndices[i]
	}

	return result, nil
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

func (matrix *Matrixi64) String() string {
	var result strings.Builder
	matrix.BuildString(&result, "")
	return result.String()
}

func (matrix *Matrixi64) BuildString(stringBuilder *strings.Builder, indent string) {
	stringBuilder.WriteString(indent + "Matrix " + strconv.Itoa(int(matrix.Rows)) + "x" + strconv.Itoa(int(matrix.Cols)) + "\n")

	for row := uint32(0); row < matrix.Rows; row = row + 1 {
		stringBuilder.WriteString(indent)
		stringBuilder.WriteString("|")
		for col := uint32(0); col < matrix.Cols; col = col + 1 {
			var asString = matrix.GetEntryi64(row, col).ToString()

			for index := len(asString); index < 8; index = index + 1 {
				stringBuilder.WriteString(" ")
			}
			stringBuilder.WriteString(asString + " ")
		}
		stringBuilder.WriteString("|\n")
	}
}
