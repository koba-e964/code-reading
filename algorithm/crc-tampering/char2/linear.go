package char2

import (
	"fmt"
)

// SolveLinear solves a linear equation xA = b, where x and b are row vectors.
func SolveLinear(A [][]byte, b []byte) ([]bool, bool) {
	n := len(A)
	m := len(b)
	if n == 0 {
		return nil, true
	}
	acp := make([][]byte, len(A))
	for i := range acp {
		if len(A[i]) != m {
			panic(fmt.Sprintf("len(A[%d]) != len(b)", i))
		}
		acp[i] = make([]byte, len(A[i]))
		copy(acp[i], A[i])
	}
	bcp := make([]byte, len(b))
	copy(bcp, b)

	col := 0
	rows := make([]int, 0)
	for i := 0; i < n; i++ {
		found := false
		for j := col; j < 8*m; j++ {
			if acp[i][j/8]&(1<<(j%8)) != 0 {
				found = true
				for k := i; k < n; k++ {
					tmp := acp[k][j/8]>>(j%8) ^ acp[k][col/8]>>(col%8)
					acp[k][col/8] ^= tmp << (col % 8)
					acp[k][j/8] ^= tmp << (j % 8)
				}
				tmp := bcp[j/8]>>(j%8) ^ bcp[col/8]>>(col%8)
				bcp[col/8] ^= tmp << (col % 8)
				bcp[j/8] ^= tmp << (j % 8)
				break
			}
		}
		if !found {
			continue
		}
		if acp[i][col/8]&(1<<(col%8)) == 0 {
			panic("Error!")
		}
		for j := col + 1; j < 8*m; j++ {
			if acp[i][j/8]&(1<<(j%8)) != 0 {
				for k := i; k < n; k++ {
					tmp := acp[k][col/8] >> (col % 8)
					acp[k][j/8] ^= tmp << (j % 8)
				}
				bcp[j/8] ^= (bcp[col/8] >> (col % 8)) << (j % 8)
			}
		}
		rows = append(rows, i)
		col++
	}

	x := make([]bool, n)
	for i := len(rows) - 1; i >= 0; i-- {
		row := rows[i]
		if acp[row][i/8]&(1<<(i%8)) == 0 {
			panic("Error!")
		}
		x[row] = bcp[i/8]&(1<<(i%8)) != 0
		for j := 0; j <= i/8; j++ {
			bcp[j] ^= acp[row][j]
		}
	}
	// is b all zero?
	for i := 0; i < len(bcp); i++ {
		if bcp[i] != 0 {
			return nil, false
		}
	}
	return x, true
}
