package filter

import (
	"strings"

	"siem-agent/internal/models"
)

var levelPriority = map[string]int{
	"DEBUG":    0,
	"INFO":     1,
	"WARNING":  2,
	"ERROR":    3,
	"CRITICAL": 4,
}

type Filter struct {
	minLevel string
}

func NewFilter(minLevel string) *Filter {
	return &Filter{minLevel: strings.ToUpper(minLevel)}
}

// ShouldSend true, если уровень записи не ниже минимального порога
func (f *Filter) ShouldSend(entry models.LogEntry) bool {
	if entry.Level == "" {
		return false
	}
	return levelPriority[strings.ToUpper(entry.Level)] >= levelPriority[f.minLevel]
}
