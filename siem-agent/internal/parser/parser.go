package parser

import (
	"regexp"
	"strings"
	"time"

	"siem-agent/internal/models"
)

var (
	// Apache/Nginx Common Log Format
	apacheRegex = regexp.MustCompile(`^(\S+) (\S+) (\S+) \[([^\]]+)\] "([^"]*)" (\d{3}) (\S+)(?: "([^"]*)")?$`)
	// Syslog формат
	syslogRegex = regexp.MustCompile(`^(\w{3}\s+\d{1,2}\s+\d{2}:\d{2}:\d{2})\s+(\S+)\s+([^:]+):\s+(.*)$`)
	// Простой формат: "УРОВЕНЬ: сообщение"
	simpleRegex = regexp.MustCompile(`^(ERROR|WARNING|INFO|DEBUG|CRITICAL)\s*:\s*(.*)$`)
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

// Parse пытается по очереди применить известные форматы; если ни один
// не подошёл — определяет уровень по ключевым словам в тексте
func (p *Parser) Parse(line string, source string) models.LogEntry {
	entry := models.LogEntry{
		Source:  source,
		Message: line,
		Level:   "INFO",
	}

	if matches := syslogRegex.FindStringSubmatch(line); len(matches) == 5 {
		if t, err := time.Parse("Jan _2 15:04:05", matches[1]); err == nil {
			entry.Timestamp = t
		}
		entry.Message = matches[4]
		entry.Level = extractLevel(matches[4])
		return entry
	}

	if matches := simpleRegex.FindStringSubmatch(line); len(matches) == 3 {
		entry.Level = matches[1]
		entry.Message = matches[2]
		return entry
	}

	if matches := apacheRegex.FindStringSubmatch(line); len(matches) == 8 {
		if t, err := time.Parse("02/Jan/2006:15:04:05 -0700", matches[4]); err == nil {
			entry.Timestamp = t
		}
		entry.Message = matches[5]
		if len(matches[6]) == 3 && (matches[6][0] == '4' || matches[6][0] == '5') {
			entry.Level = "ERROR"
		}
		return entry
	}

	entry.Level = extractLevel(line)
	return entry
}

func extractLevel(msg string) string {
	msg = strings.ToUpper(msg)
	switch {
	case strings.Contains(msg, "CRITICAL"):
		return "CRITICAL"
	case strings.Contains(msg, "ERROR"), strings.Contains(msg, "FAIL"):
		return "ERROR"
	case strings.Contains(msg, "WARNING"), strings.Contains(msg, "WARN"):
		return "WARNING"
	case strings.Contains(msg, "DEBUG"):
		return "DEBUG"
	default:
		return "INFO"
	}
}
