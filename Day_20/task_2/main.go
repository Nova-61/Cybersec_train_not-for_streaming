package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, `<!DOCTYPE html>
						<html>
						<head>
							<title>Загрузка файла</title>
						</head>
						<body>
							<h2>Загрузите файл</h2>
							<form method="POST" enctype="multipart/form-data">
								<input type="file" name="file">
								<button type="submit">Загрузить</button>
							</form>
						</body>
						</html>`)

	case "POST":
		// Проверяем, что файл есть
		file, handler, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Файл не выбран", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Создаём папку uploads, если её нет
		err = os.MkdirAll("uploads", 0755)
		if err != nil {
			http.Error(w, "Ошибка создания папки", http.StatusInternalServerError)
			return
		}

		// Создаём файл на сервере
		filePath := filepath.Join("uploads", handler.Filename)
		dst, err := os.Create(filePath)
		if err != nil {
			http.Error(w, "Ошибка сохранения файла", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		// Копируем содержимое
		_, err = io.Copy(dst, file)
		if err != nil {
			http.Error(w, "Ошибка записи файла", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Файл %s загружен!", handler.Filename)

	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/upload", uploadHandler)

	fmt.Println("Server running on http://localhost:8080")
	fmt.Println("  /upload - загрузка файлов")

	http.ListenAndServe(":8080", nil)
}
