package main

import (
	"fmt"

	"github.com/Ramil-caramel/LaboratornayaMethodOptimization-3/simptab"
)

//var matrica [][]float64

func main() {

	c_vector := []float64{7, 8, 3}
	b_vector := []float64{4, 7, 8}

	a_matrix := [][]float64{{3, 1, 1}, {1, 4, 0}, {0, 0.5, 2}}
	answer := make([]simptab.SimplexTable, 0, 10)
	simptab.Simplex(c_vector, b_vector, a_matrix, []bool{true, true, true}, true, 0, &answer)

	table := simptab.NewTable(c_vector, b_vector, a_matrix, []bool{true, true, true}, true)
	for _,val:= range answer {
		val.GetAnswerAndCheck(table)
	}

	maxVals := make([]int, len(c_vector))
	for j := 0; j < len(c_vector); j++ {
		maxVal := 1_000_000.0
		for i := 0; i < len(b_vector); i++ {
			if a_matrix[i][j] > 0 {
				val := b_vector[i] / a_matrix[i][j]
				if val < maxVal {
					maxVal = val
				}
			}
		}
		maxVals[j] = int(maxVal)
	}

	fmt.Println("Верхние границы переменных:", maxVals)

	bestF := -1e9
	bestX := [3]int{}
	count := 0
	fmt.Println("\nВсе допустимые целочисленные решения:")
	fmt.Println("X1\tX2\tX3\t  F")

	for x1 := 0; x1 <= maxVals[0]; x1++ {
		for x2 := 0; x2 <= maxVals[1]; x2++ {
			for x3 := 0; x3 <= maxVals[2]; x3++ {

				ok := true
				for i := 0; i < len(b_vector); i++ {
					left := a_matrix[i][0]*float64(x1) + a_matrix[i][1]*float64(x2) + a_matrix[i][2]*float64(x3)
					if left > b_vector[i]+1e-9 {
						ok = false
						break
					}
				}

				if ok {
					f := 7*float64(x1) + 8*float64(x2) + 3*float64(x3)
					fmt.Printf("%d\t%d\t%d\t%6.2f\n", x1, x2, x3, f)
					count++
					if f > bestF {
						bestF = f
						bestX = [3]int{x1, x2, x3}
					}
				}
			}
		}
	}
	fmt.Printf("\nНайдено %d допустимых целочисленных решений.\n", count)
	fmt.Printf("\nЛучшее решение:\n")
	fmt.Printf("X1 = %d, X2 = %d, X3 = %d\n", bestX[0], bestX[1], bestX[2])
	fmt.Printf("Максимум целевой функции F = %.2f\n", bestF)

	/*
		table := simptab.NewTable(c_vector, b_vector, a_matrix, []bool{true, true, true}, true)
		copy1 := table.DeepCopy()
		copy1.FindSupportSolution()
		copy1.DoSimplexMethod()
		copy1.GetAnswerAndCheck(table)
	*/
	//table := simptab.NewTable(c_vector, b_vector, a_matrix, []bool{true, true, true}, true)
	/*
		table.Print(-1, -1)
		fmt.Println()
		table.MakeKanonView(true)

		dualtable := simptab.DualNewTable(c_vector, b_vector, a_matrix, false)
		dualtable.Print(-1,-1)

		table.Print(-1, -1)
		table.MakeKanonView()
		fmt.Println()

		copy := table.DeepCopy()
		copy.DoSimplexMethod()

		copy.GetAnswerAndCheck(table)

		fmt.Println()


		dualtable := simptab.DualNewTable(c_vector, b_vector, a_matrix, []bool{true, true, true}, true)
		dualtable.MakeKanonView()
		fmt.Println()

		copy1 := dualtable.DeepCopy()
		copy1.FindSupportSolution()
		copy1.DoSimplexMethod()
		copy1.GetAnswerAndCheck(dualtable)
	*/

}
