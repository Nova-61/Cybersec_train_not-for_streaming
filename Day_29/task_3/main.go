// ============================================================
// ЗАДАНИЕ 3: ФИЛЬТРАЦИЯ И СТАТИСТИКА
// ============================================================
// Цель: научиться фильтровать логи по дате и считать статистику.
//
// Что делаем:
// 1. Добавляем фильтрацию по дате (флаги -from и -to)
// 2. Считаем статистику по уровням (флаг -stats)
// 3. Выводим отфильтрованные записи и статистику
// ============================================================

package main

import (
	"bufio"   // Для построчного чтения
	"flag"    // Для флагов
	"fmt"     // Для вывода
	"os"      // Для работы с файлами
	"strings" // Для работы со строками
	"time"    // Для работы с датой/временем
)

// ============================================================
// СТРУКТУРА LogEntry (описание см. в Задании 2)
// ============================================================
type LogEntry struct {
	Timestamp time.Time
	Level     string
	Message   string
}

// ============================================================
// ФУНКЦИЯ parseLogLine() (описание см. в Задании 2)
// ============================================================
func parseLogLine(line string) (LogEntry, error) {
	parts := strings.SplitN(line, " ", 3)
	if len(parts) < 3 {
		return LogEntry{}, fmt.Errorf("invalid log format: %s", line)
	}

	timestamp, err := time.Parse("2006-01-02 15:04:05", parts[0]+" "+parts[1])
	if err != nil {
		return LogEntry{}, err
	}

	return LogEntry{
		Timestamp: timestamp,
		Level:     parts[2],
		Message:   parts[3],
	}, nil
}

// ============================================================
// ФУНКЦИЯ filterByDate()
// ============================================================
// filterByDate - фильтрует записи по диапазону дат.
// Принимает:
//   - entries []LogEntry - список всех записей
//   - from time.Time - начальная дата (если пустая - не фильтруем)
//   - to time.Time   - конечная дата (если пустая - не фильтруем)
//
// Возвращает: []LogEntry - отфильтрованный список.
//
// Алгоритм:
//
//	Проходим по всем записям и оставляем только те,
//	у которых дата в диапазоне [from, to].
//
// ============================================================
func filterByDate(entries []LogEntry, from, to time.Time) []LogEntry {
	var result []LogEntry // Слайс для результата

	for _, entry := range entries {
		// Проверяем, что запись в диапазоне [from, to]
		// from.IsZero() - true, если from не задана (пустая дата)
		// entry.Timestamp.After(from) - true, если запись ПОСЛЕ from
		// entry.Timestamp.Equal(from) - true, если запись РАВНА from
		// Аналогично для to (Before / Equal)
		if (from.IsZero() || entry.Timestamp.After(from) || entry.Timestamp.Equal(from)) &&
			(to.IsZero() || entry.Timestamp.Before(to) || entry.Timestamp.Equal(to)) {
			result = append(result, entry)
		}
	}
	return result
}

// ============================================================
// ФУНКЦИЯ getStats()
// ============================================================
// getStats - считает статистику по уровням логов.
// Принимает: entries []LogEntry - список записей.
// Возвращает: map[string]int - количество записей по каждому уровню.
//
// Пример результата:
//
//	{
//	  "ERROR": 3,
//	  "INFO": 2,
//	  "WARNING": 1
//	}
//
// ============================================================
func getStats(entries []LogEntry) map[string]int {
	stats := make(map[string]int) // Создаём пустую мапу

	for _, entry := range entries {
		// Увеличиваем счётчик для данного уровня
		// Если ключа нет - он создаётся автоматически со значением 0
		stats[entry.Level]++
	}
	return stats
}

// ============================================================
// ФУНКЦИЯ printStats()
// ============================================================
// printStats - выводит статистику в консоль в красивом виде.
// ============================================================
func printStats(stats map[string]int) {
	fmt.Println("\n=== СТАТИСТИКА ЛОГОВ ===")
	total := 0

	// Перебираем все уровни и выводим их количество
	for level, count := range stats {
		fmt.Printf("  %s: %d\n", level, count)
		total += count
	}

	// Выводим итоговое количество
	fmt.Printf("  Всего: %d\n", total)
	fmt.Println("=========================")
}

// ============================================================
// ОСНОВНАЯ ФУНКЦИЯ
// ============================================================
func main() {
	// -----------------------------------------------------------------
	// ШАГ 1: ОПРЕДЕЛЯЕМ ФЛАГИ
	// -----------------------------------------------------------------
	// Добавляем новые флаги: -from, -to, -stats
	// -----------------------------------------------------------------
	filePath := flag.String("file", "logs.txt", "путь к файлу логов")
	fromDate := flag.String("from", "", "начальная дата (YYYY-MM-DD)")
	toDate := flag.String("to", "", "конечная дата (YYYY-MM-DD)")
	showStats := flag.Bool("stats", false, "показать статистику")

	flag.Parse()

	// -----------------------------------------------------------------
	// ШАГ 2: ЧИТАЕМ ФАЙЛ
	// -----------------------------------------------------------------
	file, err := os.Open(*filePath)
	if err != nil {
		fmt.Printf("Ошибка открытия файла: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var entries []LogEntry
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		if line == "" {
			continue
		}

		entry, err := parseLogLine(line)
		if err != nil {
			fmt.Printf("Ошибка в строке %d: %v\n", lineNum, err)
			continue
		}
		entries = append(entries, entry)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Ошибка чтения файла: %v\n", err)
		os.Exit(1)
	}

	// -----------------------------------------------------------------
	// ШАГ 3: ПАРСИМ ДАТЫ ДЛЯ ФИЛЬТРАЦИИ
	// -----------------------------------------------------------------
	// Превращаем строки из флагов в объекты time.Time
	// -----------------------------------------------------------------
	var from, to time.Time

	if *fromDate != "" {
		from, err = time.Parse("2006-01-02", *fromDate)
		if err != nil {
			fmt.Printf("Ошибка парсинга даты 'from': %v\n", err)
			os.Exit(1)
		}
	}

	if *toDate != "" {
		to, err = time.Parse("2006-01-02", *toDate)
		if err != nil {
			fmt.Printf("Ошибка парсинга даты 'to': %v\n", err)
			os.Exit(1)
		}
	}

	// -----------------------------------------------------------------
	// ШАГ 4: ФИЛЬТРУЕМ ЗАПИСИ ПО ДАТЕ
	// -----------------------------------------------------------------
	filtered := filterByDate(entries, from, to)

	// Выводим статистику по количеству записей
	fmt.Printf("Всего записей: %d\n", len(entries))
	fmt.Printf("Отфильтровано: %d\n", len(filtered))

	// -----------------------------------------------------------------
	// ШАГ 5: ВЫВОДИМ СТАТИСТИКУ ПО УРОВНЯМ
	// -----------------------------------------------------------------
	if *showStats {
		stats := getStats(filtered)
		printStats(stats)
	}

	// -----------------------------------------------------------------
	// ШАГ 6: ВЫВОДИМ ПЕРВЫЕ 5 ОТФИЛЬТРОВАННЫХ ЗАПИСЕЙ
	// -----------------------------------------------------------------
	if len(filtered) > 0 {
		fmt.Println("\n=== ПЕРВЫЕ 5 ЗАПИСЕЙ ===")

		// Ограничиваем вывод 5 записями
		limit := 5
		if len(filtered) < limit {
			limit = len(filtered)
		}

		for i := 0; i < limit; i++ {
			entry := filtered[i]
			fmt.Printf("[%s] %s: %s\n",
				entry.Timestamp.Format("2006-01-02 15:04:05"),
				entry.Level,
				entry.Message)
		}
	}
}
