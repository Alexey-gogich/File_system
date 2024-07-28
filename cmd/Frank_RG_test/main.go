package main

import (
	"net/http"

	"Frank_rg/internal/router"
	"Frank_rg/internal/storage/database"
)

// Загружать / удалять файлы
// Создавать папки
// Скачивать файлы
// Переименовывать файлы / папки
// Просматривать содержимое папок

func main() {
	folders := []*database.Folder{}
	folders = append(folders, database.NewFolder("root", nil))

	http.Handle("/", router.Routes(folders))
	http.ListenAndServe(":80", nil)
}
