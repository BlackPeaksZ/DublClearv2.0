package main

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("DublClear")
	myWindow.Resize(fyne.NewSize(500, 400))

	statusLabel := widget.NewLabel("Готов к работе")
	progressLabel := widget.NewLabel("")

	var startButton *widget.Button
	startButton = widget.NewButton("Начать поиск дубликатов", func() {
		startButton.Disable()
		statusLabel.SetText("Идёт сканирование...")
		go search_file(statusLabel, progressLabel, startButton, myWindow)
	})

	content := container.NewVBox(
		widget.NewLabel("DublClear - очистка от дубликатов"),
		widget.NewLabel(""),
		startButton,
		statusLabel,
		progressLabel,
	)

	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}

func search_file(statusLabel *widget.Label, progressLabel *widget.Label, startButton *widget.Button, myWindow fyne.Window) {
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

	statusLabel.SetText("Сканирование файлов...")

	filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		file_cout++
		if file_cout%1000 == 0 {
			progressLabel.SetText("Обработано " + strconv.Itoa(file_cout) + " файлов")
		}

		if err != nil {
			return nil
		}

		if !d.IsDir() {
			is_system := false
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

	statusLabel.SetText("Проверка дубликатов...")
	processed := 0

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
					progressLabel.SetText("Проверено файлов: " + strconv.Itoa(processed))
				}
			}
		}
	}

	progressLabel.SetText("Проверка завершена! Обработано файлов: " + strconv.Itoa(processed))
	statusLabel.SetText("Поиск дубликатов завершён")
	find_dubl(hash_map, statusLabel, progressLabel, startButton, myWindow)
}

func find_dubl(hash_map map[string][]string, statusLabel *widget.Label, progressLabel *widget.Label, startButton *widget.Button, myWindow fyne.Window) {
	total_dubl_file := 0
	error_file := 0
	success_file := 0

	for md5 := range hash_map {
		if len(hash_map[md5]) > 1 {
			safe_file := len(hash_map[md5]) - 1
			total_dubl_file += safe_file
		}
	}

	if total_dubl_file == 0 {
		statusLabel.SetText("Дубликаты не найдены!")
		progressLabel.SetText("")
		startButton.Enable()
		dialog.ShowInformation("Результат", "Дубликаты не найдены!", myWindow)
		return
	}

	statusLabel.SetText("Найдено дубликатов: " + strconv.Itoa(total_dubl_file))
	progressLabel.SetText("")

	dialog.ShowConfirm("Удаление дубликатов",
		"Найдено "+strconv.Itoa(total_dubl_file)+" дубликатов. Удалить?",
		func(confirmed bool) {
			if confirmed {
				statusLabel.SetText("Удаление дубликатов...")
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
				statusLabel.SetText("Готово! Удалено " + strconv.Itoa(success_file) + " файлов, " + strconv.Itoa(error_file) + " ошибок")
				progressLabel.SetText("")
			} else {
				statusLabel.SetText("Удаление отменено")
			}
			startButton.Enable()
		}, myWindow)
}
