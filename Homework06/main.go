package main

import (
	"bufio"
	"fmt"
	"os"

	"mod.go/internal/model"
)

func main() {
	var строки, столбцы int

	fmt.Print("Введите кол-во строк: ")
	fmt.Scan(&строки)
	fmt.Print("Введите кол-во столбцов: ")
	fmt.Scan(&столбцы)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	fmt.Print("Введите имя первого шахматиста: ")
	scanner.Scan()
	имя1 := scanner.Text()

	fmt.Print("Введите имя второго шахматиста: ")
	scanner.Scan()
	имя2 := scanner.Text()

	игра := model.НоваяИгра(имя1, имя2, строки, столбцы)

	fmt.Print(игра.Отобразить())
}
