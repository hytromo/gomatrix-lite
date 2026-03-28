package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	figure "github.com/common-nighthawk/go-figure"
	"github.com/mattn/go-runewidth"
	"golang.org/x/term"
)

const (
	defaultBannerColumns = 80
	defaultBannerRows    = 24
)

func printBannerFitted(phrase string, font string, replaceSimilar bool) {
	columns, rows := getBannerTerminalSize()
	lines := buildFittedBannerLines(phrase, font, columns, rows)
	if len(lines) == 0 {
		return
	}

	// vertical centering
	if len(lines) < rows {
		slack := rows - len(lines)
		top := slack / 2
		bottom := slack - top
		out := make([]string, 0, top+len(lines)+bottom)
		for range top {
			out = append(out, "")
		}
		out = append(out, lines...)
		for range bottom {
			out = append(out, "")
		}
		lines = out
	}

	output := strings.Join(lines, "\n")
	if replaceSimilar {
		output = replaceAsciiSimilar(output)
	}

	fmt.Print(output)
}

func getBannerTerminalSize() (int, int) {
	columns, rows, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || columns <= 0 || rows <= 0 {
		return defaultBannerColumns, defaultBannerRows
	}

	return columns, rows
}

func buildFittedBannerLines(phrase string, font string, columns int, rows int) []string {
	if columns <= 0 {
		columns = defaultBannerColumns
	}
	if rows <= 0 {
		rows = defaultBannerRows
	}

	chunks := splitBannerChunks(phrase, font, columns)
	if len(chunks) == 0 {
		return nil
	}

	lines := make([]string, 0, rows)

	for _, chunk := range chunks {
		block := centerBannerBlock(renderBannerBlock(chunk, font), columns)
		if len(block) == 0 {
			continue
		}

		remainingRows := rows - len(lines)
		if remainingRows <= 0 {
			break
		}

		if len(lines) > 0 {
			if remainingRows == 1 {
				break
			}
			lines = append(lines, "")
			remainingRows--
		}

		if len(block) > remainingRows {
			block = block[:remainingRows]
		}

		lines = append(lines, block...)
	}

	return lines
}

func splitBannerChunks(phrase string, font string, columns int) []string {
	words := strings.Fields(phrase)
	if len(words) == 0 {
		return nil
	}

	chunks := make([]string, 0, len(words))
	current := ""

	for _, word := range words {
		if current == "" {
			if measureRenderedWidth(word, font) <= columns {
				current = word
			} else {
				chunks = append(chunks, splitLongToken(word, font, columns)...)
			}
			continue
		}

		candidate := current + " " + word
		if measureRenderedWidth(candidate, font) <= columns {
			current = candidate
			continue
		}

		chunks = append(chunks, current)
		if measureRenderedWidth(word, font) <= columns {
			current = word
		} else {
			chunks = append(chunks, splitLongToken(word, font, columns)...)
			current = ""
		}
	}

	if current != "" {
		chunks = append(chunks, current)
	}

	return chunks
}

func splitLongToken(token string, font string, columns int) []string {
	runes := []rune(token)
	chunks := make([]string, 0, len(runes))

	for len(runes) > 0 {
		prefixLength := longestFittingPrefix(runes, font, columns)
		if prefixLength <= 0 {
			prefixLength = 1
		}

		chunks = append(chunks, string(runes[:prefixLength]))
		runes = runes[prefixLength:]
	}

	return chunks
}

func longestFittingPrefix(runes []rune, font string, columns int) int {
	low := 1
	high := len(runes)
	best := 0

	for low <= high {
		middle := low + (high-low)/2
		if measureRenderedWidth(string(runes[:middle]), font) <= columns {
			best = middle
			low = middle + 1
		} else {
			high = middle - 1
		}
	}

	return best
}

func measureRenderedWidth(phrase string, font string) int {
	return maxRenderedLineWidth(renderBannerBlock(phrase, font))
}

func renderBannerBlock(phrase string, font string) []string {
	if strings.TrimSpace(phrase) == "" {
		return nil
	}

	var buffer bytes.Buffer
	figure.Write(&buffer, figure.NewFigure(phrase, font, false))

	rendered := strings.TrimRight(buffer.String(), "\n")
	if rendered == "" {
		return nil
	}

	return strings.Split(rendered, "\n")
}

func centerBannerBlock(lines []string, columns int) []string {
	if len(lines) == 0 {
		return nil
	}

	clipped := make([]string, 0, len(lines))
	for _, line := range lines {
		clipped = append(clipped, truncateToDisplayWidth(line, columns))
	}

	blockWidth := maxRenderedLineWidth(clipped)
	padding := 0
	if blockWidth < columns {
		padding = (columns - blockWidth) / 2
	}

	centered := make([]string, 0, len(clipped))
	for _, line := range clipped {
		centered = append(centered, strings.Repeat(" ", padding)+line)
	}

	return centered
}

func maxRenderedLineWidth(lines []string) int {
	maxWidth := 0
	for _, line := range lines {
		lineWidth := runewidth.StringWidth(line)
		if lineWidth > maxWidth {
			maxWidth = lineWidth
		}
	}

	return maxWidth
}

func truncateToDisplayWidth(line string, width int) string {
	if width <= 0 {
		return ""
	}

	currentWidth := 0
	var builder strings.Builder

	for _, r := range line {
		runeWidth := runewidth.RuneWidth(r)
		if currentWidth+runeWidth > width {
			break
		}

		builder.WriteRune(r)
		currentWidth += runeWidth
	}

	return builder.String()
}
