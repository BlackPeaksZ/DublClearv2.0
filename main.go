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
			m := md5.New()
			io.Copy(m, file)
			file.Close()
			hashBytes := m.Sum(nil)
			md5 := hex.EncodeToString(hashBytes)
			hash_map[md5] = append(hash_map[md5], path)
		}
		return nil
	})
	find_dubl(hash_map)

}
func find_dubl(hash_map map[string][]string) {
	fmt.Println(hash_map)
	dubl_file := 0
	for md5 := range hash_map {
		fmt.Println(hash_map[md5])
		if len(hash_map[md5]) > 1 {
			dubl_file++
			safe_file := len(hash_map[md5]) - 1
			fmt.Println("Файлы для удаления: ", dubl_file+safe_file)
		}
		if len(hash_map[md5]) == 1 {
			continue
		}
	}

}
