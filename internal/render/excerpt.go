package render

import (
	"regexp"
	"strings"
)

// AutoExcerpt extracts the first maxRunes runes of plain text from Markdown.
// It strips formatting for use as a summary.
func AutoExcerpt(markdown string, maxRunes int) string {
	if maxRunes <= 0 {
		maxRunes = 200
	}

	text := markdown
	text = stripCodeBlocks(text)
	text = stripImages(text)
	text = stripLinks(text)
	text = stripHeadings(text)
	text = stripEmphasis(text)
	text = stripBlockquotes(text)
	text = collapseWhitespace(text)

	runes := []rune(text)
	if len(runes) > maxRunes {
		return string(runes[:maxRunes]) + "…"
	}
	return text
}

// --- strip helpers ---

func stripCodeBlocks(s string) string {
	var lines []string
	inBlock := false
	for _, line := range strings.Split(s, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "```") {
			inBlock = !inBlock
			continue
		}
		if !inBlock {
			lines = append(lines, line)
		}
	}
	return strings.Join(lines, "\n")
}

func stripImages(s string) string {
	re := regexp.MustCompile(`!\[.*?\]\(.*?\)`)
	return re.ReplaceAllString(s, "")
}

func stripLinks(s string) string {
	re := regexp.MustCompile(`\[(.*?)\]\(.*?\)`)
	return re.ReplaceAllString(s, "$1")
}

func stripHeadings(s string) string {
	re := regexp.MustCompile(`(?m)^#{1,6}\s+`)
	return re.ReplaceAllString(s, "")
}

func stripEmphasis(s string) string {
	s = regexp.MustCompile(`\*{1,2}(.*?)\*{1,2}`).ReplaceAllString(s, "$1")
	s = regexp.MustCompile(`_{1,2}(.*?)_{1,2}`).ReplaceAllString(s, "$1")
	s = regexp.MustCompile(`~~(.*?)~~`).ReplaceAllString(s, "$1")
	return s
}

func stripBlockquotes(s string) string {
	var lines []string
	for _, line := range strings.Split(s, "\n") {
		lines = append(lines, strings.TrimPrefix(line, "> "))
	}
	return strings.Join(lines, "\n")
}

func collapseWhitespace(s string) string {
	s = strings.TrimSpace(s)
	for strings.Contains(s, "  ") {
		s = strings.ReplaceAll(s, "  ", " ")
	}
	for strings.Contains(s, "\n\n\n") {
		s = strings.ReplaceAll(s, "\n\n\n", "\n\n")
	}
	return s
}
