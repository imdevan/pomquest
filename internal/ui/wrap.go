package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// breakLongWord breaks a single long word (with no spaces) into sub-strings
// that each have visual cell width <= maxWidth.
func breakLongWord(word string, maxWidth int) []string {
	if maxWidth <= 0 {
		return []string{word}
	}
	var lines []string
	runes := []rune(word)
	var current []rune

	for _, r := range runes {
		testStr := string(append(current, r))
		if lipgloss.Width(testStr) > maxWidth && len(current) > 0 {
			lines = append(lines, string(current))
			current = []rune{r}
		} else {
			current = append(current, r)
		}
	}
	if len(current) > 0 {
		lines = append(lines, string(current))
	}
	return lines
}

// WrapLine wraps a single line of text so that no returned line exceeds maxWidth visual cells.
func WrapLine(line string, maxWidth int) []string {
	if maxWidth <= 0 || lipgloss.Width(line) <= maxWidth {
		return []string{line}
	}

	words := strings.Split(line, " ")
	var wrapped []string
	var currentLine string

	for i, word := range words {
		if i == 0 {
			if lipgloss.Width(word) <= maxWidth {
				currentLine = word
			} else {
				subLines := breakLongWord(word, maxWidth)
				if len(subLines) > 0 {
					wrapped = append(wrapped, subLines[:len(subLines)-1]...)
					currentLine = subLines[len(subLines)-1]
				}
			}
		} else {
			candidate := currentLine + " " + word
			if lipgloss.Width(candidate) <= maxWidth {
				currentLine = candidate
			} else {
				wrapped = append(wrapped, currentLine)
				if lipgloss.Width(word) <= maxWidth {
					currentLine = word
				} else {
					subLines := breakLongWord(word, maxWidth)
					if len(subLines) > 0 {
						wrapped = append(wrapped, subLines[:len(subLines)-1]...)
						currentLine = subLines[len(subLines)-1]
					} else {
						currentLine = ""
					}
				}
			}
		}
	}

	if currentLine != "" {
		wrapped = append(wrapped, currentLine)
	}

	if len(wrapped) == 0 {
		return []string{line}
	}

	return wrapped
}

// WrapText splits text by newlines and wraps each line to fit within maxWidth.
func WrapText(text string, maxWidth int) []string {
	if maxWidth <= 0 {
		return strings.Split(text, "\n")
	}
	rawLines := strings.Split(text, "\n")
	var result []string
	for _, line := range rawLines {
		result = append(result, WrapLine(line, maxWidth)...)
	}
	return result
}
