package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var answer string
	fmt.Print("Начать поиск дубликатов? (д/н): ")
	fmt.Scan(&answer)

	if answer == "д" {
		search_file()
	} else if answer == "н" {
		os.Exit(0)
	} else {
		fmt.Println("Ошибка комманды, введите 'д' или 'н'")
		os.Exit(0)
	}
}
func search_file() {
	system_files := []string{
		"Windows",
		"Program Files",
		"Program Files (x86)",
		"System Volume Information",
		"$Recycle.Bin",
		"Boot",
		"ProgramData",
		"System32",
		"SysWOW64",
		"WinSxS",
		"Microsoft.NET",
		"Common Files",
	}

	hash_map := make(map[string][]string)
	file_cout := 0
	path := "C:\\"
	size_map := make(map[int64][]string)
	filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		file_cout++
		if file_cout%1000 == 0 {
			fmt.Println("Обработано ", file_cout, " файлов")
		}

		if err != nil {
			return nil
		}

		if !d.IsDir() {
			is_system := false
			//!-считает кол-во файлов, без !-считает
			for _, folder := range system_files {
				if strings.Contains(path, folder) {
					is_system = true
					break
				}
			}
			if is_system == true {
				return nil
			}

			fileInfo, err := os.Stat(path)
			if err != nil {
				return nil
			}
			size := fileInfo.Size()
			size_map[size] = append(size_map[size], path)

		}

		return nil
	})
	processed := 0
	fmt.Println("Проверка дубликатов")

	for _, paths := range size_map {
		if len(paths) > 1 {
			for _, path := range paths {
				file, err := os.Open(path)
				if err != nil {
					continue
				}
				m := md5.New()
				io.Copy(m, file)
				file.Close()
				hashBytes := m.Sum(nil)
				md5 := hex.EncodeToString(hashBytes)
				hash_map[md5] = append(hash_map[md5], path)
				processed++
				if processed%10 == 0 {
					fmt.Println("Проверено файлов:", processed)
				}
			}

		}
	}
	fmt.Println("Проверка завершена! Обработано файлов:", processed)
	find_dubl(hash_map)

}
func find_dubl(hash_map map[string][]string) {
	total_dubl_file := 0
	error_file := 0
	success_file := 0
	for md5 := range hash_map {

		if len(hash_map[md5]) > 1 {
			safe_file := len(hash_map[md5]) - 1
			total_dubl_file += safe_file
		}
		if len(hash_map[md5]) == 1 {
			continue
		}

	}
	if total_dubl_file == 0 {
		fmt.Println("Дубликаты не найдены")
		os.Exit(0)
	}
	fmt.Println("Колличество дубликатов: ", total_dubl_file)
	var answer2 string
	fmt.Println("Удалить дубликаты?(д/н): ")
	fmt.Scan(&answer2)
	if answer2 == "н" {
		fmt.Println("Дубликаты не будут удалены")
		os.Exit(0)
	}
	if answer2 == "д" {
		fmt.Println("Дубликаты будут удалены")
		for md5 := range hash_map {
			if len(hash_map[md5]) > 1 {
				paths := hash_map[md5]
				for i := 1; i < len(paths); i++ {
					del_path := paths[i]
					err := os.Remove(del_path)
					if err != nil {
						error_file++
					}
					if err == nil {
						success_file++
					}
				}

			}
		}
		fmt.Println(error_file, "ошибок")
		fmt.Println("Удалено", success_file, "файлов")

	}
}
