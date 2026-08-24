package reader

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"

	"siem-agent/internal/models"
)

type LogReader struct {
	paths   []string
	watcher *fsnotify.Watcher
	output  chan models.LogEntry
	errors  chan error

	// mu защищает offsets от гонки данных: watchEvents может вызвать readFile
	// для одного и того же файла, пока предыдущий вызов ещё не закончился
	mu      sync.Mutex
	offsets map[string]int64 // путь -> позиция, до которой уже прочитано
}

func NewLogReader(paths []string) (*LogReader, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create watcher: %w", err)
	}

	return &LogReader{
		paths:   paths,
		watcher: watcher,
		output:  make(chan models.LogEntry, 1000),
		errors:  make(chan error, 10),
		offsets: make(map[string]int64),
	}, nil
}

// Start начинает чтение логов: разово читает существующее содержимое
// каждого файла и подписывается на дальнейшие изменения
func (r *LogReader) Start() {
	for _, path := range r.paths {
		if err := r.watcher.Add(path); err != nil {
			r.errors <- fmt.Errorf("failed to watch %s: %w", path, err)
			continue
		}
		fmt.Printf("Watching: %s\n", path)
		go r.readFile(path)
	}

	go r.watchEvents()
}

// readFile читает файл начиная с сохранённой позиции (offset) и отправляет
// каждую новую строку в канал. Защищено мьютексом: если watchEvents вызовет
// readFile для того же файла, пока предыдущий вызов ещё работает, второй
// вызов подождёт своей очереди, а не будет читать параллельно с первым.
func (r *LogReader) readFile(path string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	file, err := os.Open(path)
	if err != nil {
		r.errors <- fmt.Errorf("failed to open %s: %w", path, err)
		return
	}
	defer file.Close()

	// Переходим на позицию, до которой уже читали в прошлый раз.
	// При первом вызове offsets[path] == 0 (нулевое значение map по умолчанию),
	// поэтому первый вызов честно читает файл с самого начала.
	offset := r.offsets[path]
	if _, err := file.Seek(offset, io.SeekStart); err != nil {
		r.errors <- fmt.Errorf("failed to seek %s: %w", path, err)
		return
	}

	scanner := bufio.NewScanner(file)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	var bytesRead int64
	for scanner.Scan() {
		line := scanner.Text()
		bytesRead += int64(len(scanner.Bytes())) + 1 // +1 за символ перевода строки
		if line == "" {
			continue
		}
		r.output <- models.LogEntry{
			Timestamp: time.Now(),
			Message:   line,
			Source:    path,
		}
	}

	if err := scanner.Err(); err != nil {
		r.errors <- fmt.Errorf("error reading %s: %w", path, err)
		return
	}

	// Запоминаем новую позицию для следующего вызова
	r.offsets[path] = offset + bytesRead
}

// watchEvents слушает события файловой системы и реагирует на изменения
func (r *LogReader) watchEvents() {
	for event := range r.watcher.Events {
		if event.Op&fsnotify.Write == fsnotify.Write {
			for _, path := range r.paths {
				if event.Name == path {
					go r.readFile(path)
					break
				}
			}
		}
	}
}

func (r *LogReader) Output() <-chan models.LogEntry {
	return r.output
}

func (r *LogReader) Errors() <-chan error {
	return r.errors
}

func (r *LogReader) Close() error {
	return r.watcher.Close()
}