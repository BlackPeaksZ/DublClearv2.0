package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func main() {
	io.ByteReader()
	hex.AppendDecode()

	var answer string

	fmt.Println("Начать поиск дубликатов? (д/н)")
	fmt.Scan(&answer)

	if answer == "д" {
		search_file()
	} else if answer == "н" {
		os.Exit(0)
	} else {
		fmt.Println("Ошибка комманды, введите 'д' или 'н'")
	}
}
func search_file() {

	hash_map := make(map[string][]string)

	path := "C:\\"
	file_cout, err_cout := 0, 0

	filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			err_cout++
			return nil
		} else if !d.IsDir() == true {
			//!-считает кол-во файлов, без !-считает папки
			file_cout++
			file, err := os.Open(path)
			if err != nil {
				err_cout++
				return nil
			}
			defer file.Close()
			m := md5.New()
			io.Copy(m, file)
			hashBytes := m.Sum(nil)

		}
		return nil

	})

	fmt.Println("Колличество файлов", file_cout)
	fmt.Println("Колличество ошибок", err_cout)
}
