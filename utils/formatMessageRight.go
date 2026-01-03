package utils

import (
	"strings"
)

func FormatMessageRight(msg string, width int) string {
	lines := strings.Split(msg, "\n")
	for i, line := range lines {
		lineLen := len(line)
		if lineLen > width {
			line = line[:width] // обрезаем, если слишком длинное
			lineLen = width
		}
		space := width - lineLen
		lines[i] = strings.Repeat(" ", space) + line
	}
	return strings.Join(lines, "\n")
}

// форматируем сообщение слева
func FormatMessageLeft(msg string, width int) string {
	lines := strings.Split(msg, "\n")
	for i, line := range lines {
		if len(line) > width {
			lines[i] = line[:width]
		}
	}
	return strings.Join(lines, "\n")
}

