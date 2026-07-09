package main

import "fmt"

func main() {
	var chessGrid string

	var a, b uint // Размер доски не может быть отрицательным

	fmt.Print("Введите кол-во строк: ")
	fmt.Scan(&a)

	fmt.Print("Введите кол-во стобцов: ")
	fmt.Scan(&b)

	for строка := uint(0); строка < a; строка++ {
		for столбец := uint(0); столбец < b; столбец++ {
			if (строка+столбец)%2 == 0 {
				chessGrid += " "
			} else {
				chessGrid += "#"
			}
		}
		chessGrid += "\n"

	}
	fmt.Print(chessGrid)
}
