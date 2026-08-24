package models

import "time"

// LogEntry — Является переносчиком состояния, который растет по пути,
// тоесть вместо того что бы заполнять одной проверкой,
// мы делим эту проверку на части что бы в дальнейшем понять
// в чем была проблема
type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"` // INFO, WARNING, ERROR, DEBUG, CRITICAL
	Message   string    `json:"message"`
	Source    string    `json:"source"` // из какого файла пришла запись
}

// LogBatch — Сохраняет информацию о том какая менно это запись "Count" и о том
// откуда она пришла "Source".
// Если бы мы просто отправляли []LogEntry без обёртки — это работало бы,
// но мы бы потеряли место, куда естественно добавляются метаданные о самой
// отправке, а не о конкретной записи.
type LogBatch struct {
	Entries []LogEntry `json:"entries"`
	Count   int        `json:"count"`
	Source  string     `json:"source"`
}
