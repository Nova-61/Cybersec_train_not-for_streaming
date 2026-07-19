package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func listFilesHandler(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir("uploads")
	if err != nil {
		http.Error(w, "Ошибка чтения папки", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<h2>Список файлов в uploads/</h2><ul>")
	for _, file := range files {
		if !file.IsDir() {
			fmt.Fprintf(w, `<li><a href="/download/%s">%s</a> | <a href="/delete/%s">🗑️</a></li>`,
				file.Name(), file.Name(), file.Name())
		}
	}
	fmt.Fprintf(w, "</ul>")
	fmt.Fprintf(w, `<br><a href="/upload">Загрузить файл</a>`)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, `<!DOCTYPE html>
<html>
<head><title>Загрузка файла</title></head>
<body>
	<h2>Загрузите файл</h2>
	<form method="POST" enctype="multipart/form-data">
		<input type="file" name="file">
		<button type="submit">Загрузить</button>
	</form>
	<br><a href="/">Назад</a>
</body>
</html>`)

	case "POST":
		err := r.ParseMultipartForm(10 << 20) // 10 MB
		if err != nil {
			http.Error(w, "Файл слишком большой (макс. 10 МБ)", http.StatusBadRequest)
			return
		}

		file, handler, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Файл не выбран", http.StatusBadRequest)
			return
		}
		defer file.Close()

		filename := handler.Filename

		// Проверка на обход папок
		if strings.Contains(filename, "..") {
			http.Error(w, "Недопустимое имя файла", http.StatusBadRequest)
			return
		}

		// Разрешённые расширения
		ext := strings.ToLower(filepath.Ext(filename))
		allowedExts := map[string]bool{
			".txt": true,
			".jpg": true,
			".png": true,
			".pdf": true,
		}
		if !allowedExts[ext] {
			http.Error(w, "Разрешены только: .txt, .jpg, .png, .pdf", http.StatusBadRequest)
			return
		}

		err = os.MkdirAll("uploads", 0755)
		if err != nil {
			http.Error(w, "Ошибка создания папки", http.StatusInternalServerError)
			return
		}

		filePath := filepath.Join("uploads", filename)
		dst, err := os.Create(filePath)
		if err != nil {
			http.Error(w, "Ошибка сохранения файла", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		_, err = io.Copy(dst, file)
		if err != nil {
			http.Error(w, "Ошибка записи файла", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Файл %s загружен!<br><a href=\"/\">Назад</a>", filename)
	}
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	filename := strings.TrimPrefix(r.URL.Path, "/download/")
	if filename == "" {
		http.Error(w, "Не указан файл", http.StatusBadRequest)
		return
	}

	if strings.Contains(filename, "..") {
		http.Error(w, "Недопустимое имя файла", http.StatusBadRequest)
		return
	}

	filename = filepath.Join("uploads", filename)

	file, err := os.Open(filename)
	if err != nil {
		http.Error(w, "Файл не найден", http.StatusNotFound)
		return
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		http.Error(w, "Ошибка чтения файла", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(filename))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", info.Size()))

	io.Copy(w, file)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	filename := strings.TrimPrefix(r.URL.Path, "/delete/")
	if filename == "" {
		http.Error(w, "Не указан файл", http.StatusBadRequest)
		return
	}

	if strings.Contains(filename, "..") {
		http.Error(w, "Недопустимое имя файла", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join("uploads", filename)
	err := os.Remove(filePath)
	if err != nil {
		http.Error(w, "Файл не найден", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "Файл %s удалён!<br><a href=\"/\">Назад</a>", filename)
}

func main() {
	os.MkdirAll("uploads", 0755)

	http.HandleFunc("/", listFilesHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/download/", downloadHandler)
	http.HandleFunc("/delete/", deleteHandler)

	fmt.Println("Server running on http://localhost:8080")
	fmt.Println("  /                 - список файлов")
	fmt.Println("  /upload           - загрузка файла")
	fmt.Println("  /download/{file}  - скачать файл")
	fmt.Println("  /delete/{file}    - удалить файл")

	http.ListenAndServe(":8080", nil)
}
