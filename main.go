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

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("DublClear")
	w.Resize(fyne.NewSize(500, 250))

	label := widget.NewLabel("Обработано файлов: ")
	startText := widget.NewLabel("Начать поиск дубликатов?")
	progress := widget.NewProgressBar()
	progress.Min = 0
	progress.Max = 100

	var startBtn *widget.Button

	startBtn = widget.NewButton("Да", func() {
		startBtn.Disable()
		progress.SetValue(0)
		go search_file(label, progress, startBtn, w)
	})

	quitBtn := widget.NewButton("Нет", func() { a.Quit() })

	w.SetContent(container.NewVBox(
		startText,
		progress,
		label,
		container.NewHBox(startBtn, quitBtn),
	))

	w.ShowAndRun()
}

func search_file(label *widget.Label, progress *widget.ProgressBar, startBtn *widget.Button, w fyne.Window) {

	defer startBtn.Enable()

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

	filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
		file_cout++
		if file_cout%1000 == 0 {
			label.SetText(fmt.Sprintf("Обработано файлов: %d", file_cout))
		}

		if err != nil {
			return nil
		}

		if !d.IsDir() {
			parts := strings.Split(p, string(os.PathSeparator))
			is_system := false
			for _, part := range parts {
				for _, folder := range system_files {
					if strings.EqualFold(part, folder) {
						is_system = true
						break
					}
				}
				if is_system {
					break
				}
			}
			if is_system {
				return nil
			}

			fileInfo, err := os.Stat(p)
			if err != nil {
				return nil
			}
			size := fileInfo.Size()
			size_map[size] = append(size_map[size], p)
		}

		return nil
	})

	processed := 0
	label.SetText("Проверка дубликатов...")

	total_groups := 0
	for _, paths := range size_map {
		if len(paths) > 1 {
			total_groups++
		}
	}

	processed_groups := 0

	for _, paths := range size_map {
		if len(paths) > 1 {
			for _, p := range paths {
				file, err := os.Open(p)
				if err != nil {
					continue
				}
				m := md5.New()
				if _, err := io.Copy(m, file); err != nil {
					file.Close()
					continue
				}
				file.Close()
				hashBytes := m.Sum(nil)
				hash := hex.EncodeToString(hashBytes)
				hash_map[hash] = append(hash_map[hash], p)
				processed++
				if processed%10 == 0 {
					label.SetText(fmt.Sprintf("Проверено файлов: %d", processed))
				}
			}
		}

		if len(paths) > 1 {
			processed_groups++
			if total_groups > 0 {
				progress.SetValue(float64(processed_groups) / float64(total_groups) * 100)
			}
		}
	}

	label.SetText(fmt.Sprintf("Проверка завершена! Обработано файлов: %d", processed))
	progress.SetValue(100)

	find_dubl(hash_map, label, progress, startBtn, w)
}

func find_dubl(hash_map map[string][]string, label *widget.Label, progress *widget.ProgressBar, startBtn *widget.Button, w fyne.Window) {

	total_dubl_file := 0
	error_file := 0
	success_file := 0

	for _, paths := range hash_map {
		if len(paths) > 1 {
			total_dubl_file += len(paths) - 1
		}
	}

	if total_dubl_file == 0 {
		label.SetText("Дубликаты не найдены")
		dialog.ShowInformation("Результат", "Дубликаты не найдены!", w)
		return
	}

	label.SetText(fmt.Sprintf("Найдено дубликатов: %d", total_dubl_file))

	dialog.ShowConfirm("Удаление дубликатов",
		fmt.Sprintf("Найдено %d дубликатов. Удалить?", total_dubl_file),
		func(ok bool) {
			if !ok {
				label.SetText("Удаление отменено")
				return
			}

			label.SetText("Удаление дубликатов...")

			for _, paths := range hash_map {
				if len(paths) > 1 {
					for i := 1; i < len(paths); i++ {
						del_path := paths[i]
						err := os.Remove(del_path)
						if err != nil {
							error_file++
						} else {
							success_file++
						}
					}
				}
			}

			label.SetText(fmt.Sprintf("Готово! Удалено: %d, ошибок: %d", success_file, error_file))
			dialog.ShowInformation("Результат удаления",
				fmt.Sprintf("Удалено файлов: %d\nОшибок: %d", success_file, error_file), w)
		}, w)
}
