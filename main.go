package main

import (
	"fmt"
	"os"
)

func main() {

	var answer string

	fmt.Println("Начать поиск дубликатов? (д/н)")
	fmt.Scan(&answer)

	if answer == "д" {
		file_read()
	} else if answer == "н" {
		os.Exit(0)
	} else {
		fmt.Println("Ошибка комманды, введите 'д' или 'н'")
	}
}
func file_read() {

	path := "C:\\"


}
