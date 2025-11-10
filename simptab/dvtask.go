package simptab

import (
	"fmt"
	"math"
)

// значения вставяляются те же что и в прямой задаче
func DualNewTable(c_vector []float64, b_vector []float64, a_matrix [][]float64, sign_vector []bool, view bool) *SimplexTable {

	for i := range sign_vector {
		sign_vector[i] = !sign_vector[i]
	}

	table := NewTable(b_vector, c_vector, transpon(a_matrix), sign_vector, !view)

	for i := range table.Basis {
		table.Basis[i] = replaceChar(table.Basis[i], rune('X'), rune('Y'))
		table.Basis[i] = replaceChar(table.Basis[i], rune('F'), rune('G'))
	}
	for i := range table.Headers {
		table.Headers[i] = replaceChar(table.Headers[i], rune('X'), rune('Y'))
		table.Headers[i] = replaceChar(table.Headers[i], rune('F'), rune('G'))
	}

	return table
}

func transpon(a_matrix [][]float64) [][]float64 {
	a_matrix_transp := make([][]float64, len(a_matrix[0]))

	for j := 0; j < len(a_matrix[0]); j++ {

		a_matrix_transp[j] = make([]float64, len(a_matrix))
		for i := 0; i < len(a_matrix); i++ {
			a_matrix_transp[j][i] = a_matrix[i][j]
		}
	}
	return a_matrix_transp
}

func replaceChar(s string, old_char rune, new_char rune) string {
	result := ""
	for _, val := range s {
		if val == old_char {
			result += string(new_char)
		} else {
			result += string(val)
		}
	}
	return result
}

func (simptab *SimplexTable) LiteFindSupportSolution() bool {

	for !simptab.CheckSupportSolution() {

		preRow := -1
		min := 0.0
		razrech_string, razrech_stolb := -1, -1
		var sizeTable int = len(simptab.Table)

		for i := 0; i < len(simptab.Table)-1; i++ {
			if simptab.Table[i][0] < min {
				min = simptab.Table[i][0]
				preRow = i
				
			}
		}
		//fmt.Println(preRow)
		min1 := 0.0
		for i, val := range simptab.Table[preRow] {
			//fmt.Println(val,i)
			if val < min1 && i != 0 {
				//fmt.Println(val,i)
				razrech_stolb = i
				min1 = val
			}
		}
		if razrech_stolb == -1 {
			//fmt.Println("не найден столбец с отрицательным элементом для постановки опорного решения")
			return false
		}

		max := math.MaxFloat64

		for i := range simptab.Table {
			target := simptab.Table[i][0] / simptab.Table[i][razrech_stolb]
			if simptab.Table[i][0] != 0 && target < max && target > 0 {
				max = target
				razrech_string = i
			}
		}

		pivot := simptab.Table[razrech_string][razrech_stolb] //pivot - разрешающее значение
		var targetLine []float64 = simptab.Table[sizeTable-1]

		newTable := make([][]float64, sizeTable)
		for i := range newTable {
			newTable[i] = make([]float64, len(targetLine))
		}

		for i := 0; i < sizeTable; i++ {
			for j := 0; j < len(targetLine); j++ {
				switch {
				case i == razrech_string && j == razrech_stolb:
					newTable[i][j] = 1.0 / pivot
				case i == razrech_string && j != razrech_stolb:
					newTable[i][j] = simptab.Table[i][j] / pivot
				case i != razrech_string && j == razrech_stolb:
					newTable[i][j] = -simptab.Table[i][j] / pivot
				default:
					newTable[i][j] = simptab.Table[i][j] - (simptab.Table[i][razrech_stolb] * simptab.Table[razrech_string][j] / pivot)
				}
			}
		}

		//заменяем таблицу шапки
		simptab.Table = newTable
		simptab.Basis[razrech_string], simptab.Headers[razrech_stolb] = simptab.Headers[razrech_stolb], simptab.Basis[razrech_string]

	}
	return true
}

func (simptab *SimplexTable) DoLiteSimplexMethod() {

	fmt.Println()

	var count int = 1
	var sizeTable int = len(simptab.Table)

	for !simptab.CheckOptimized() {

		if count > 100 {
			fmt.Println("Превышено число итераций")
			return
		}

		var targetLine []float64 = simptab.Table[sizeTable-1]
		var max float64 = 0
		razrech_stolb := -1
		var razrech_string int

		for j := 1; j < len(targetLine); j++ {
			if targetLine[j] > max {
				max = targetLine[j]
				razrech_stolb = j
			}
		}

		if razrech_stolb == -1 {
			fmt.Println("Не нашел разрешающий столбец") // нет входящего столбца
			break
		}

		max = math.MaxFloat64
		// первый столбец это столбец начальных ограничений
		for i, val := range simptab.Table {
			if i != sizeTable-1 && val[razrech_stolb] > 0 && val[0] >= 0 && val[0]/val[razrech_stolb] < max { ////////////////////
				max = val[0] / val[razrech_stolb]
				razrech_string = i
			}
		}

		if max == math.MaxFloat64 { // если разрешающей строки нет
			fmt.Println("Решения нет")
			return
		}

		// новая таблица
		pivot := simptab.Table[razrech_string][razrech_stolb] //pivot - разрешающее значение
		newTable := make([][]float64, sizeTable)
		for i := range newTable {
			newTable[i] = make([]float64, len(targetLine))
		}

		for i := 0; i < sizeTable; i++ {
			for j := 0; j < len(targetLine); j++ {
				switch {
				case i == razrech_string && j == razrech_stolb:
					newTable[i][j] = 1.0 / pivot
				case i == razrech_string && j != razrech_stolb:
					newTable[i][j] = simptab.Table[i][j] / pivot
				case i != razrech_string && j == razrech_stolb:
					newTable[i][j] = -simptab.Table[i][j] / pivot
				default:
					newTable[i][j] = simptab.Table[i][j] - (simptab.Table[i][razrech_stolb] * simptab.Table[razrech_string][j] / pivot)
				}
			}
		}

		//заменяем таблицу шапки
		simptab.Table = newTable
		simptab.Basis[razrech_string], simptab.Headers[razrech_stolb] = simptab.Headers[razrech_stolb], simptab.Basis[razrech_string]
	}
	//simptab.Print(-1, -1)
}

func (simptab SimplexTable) Printindent(rowPrint int, colPrint int, indent string) { // -1 без подсветки
	red := "\033[31m"
	reset := "\033[0m"

	
	fmt.Printf("%s%-4s", indent, "B\\F")
	for i, val := range simptab.Headers {
		if i == colPrint {
			fmt.Printf("%s%8s%s", red, val, reset)
		} else {
			fmt.Printf("%8s", val)
		}
	}
	fmt.Println()

	
	for i, row := range simptab.Table {
		fmt.Printf("%s", indent) // каждый ряд начинается с отступа

		if i == rowPrint {
			fmt.Printf("%s%-4s%s", red, simptab.Basis[i], reset)
			for _, val := range row {
				fmt.Printf("%s%8.2f%s", red, val, reset)
			}
			fmt.Println()
		} else {
			fmt.Printf("%-4s", simptab.Basis[i])
			for j, val := range row {
				if j == colPrint {
					fmt.Printf("%s%8.2f%s", red, val, reset)
				} else {
					fmt.Printf("%8.2f", val)
				}
			}
			fmt.Println()
		}
	}
}
