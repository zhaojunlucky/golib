package collection

func CreateTwoDimArray(row, column int) [][]int {
	var twoDimArr = make([][]int, row)

	for i := range twoDimArr {
		twoDimArr[i] = make([]int, column)
	}

	return twoDimArr
}
