package main

import (
	"math/rand"
	"time"

	"github.com/gdamore/tcell/v2"
)

// the minimum length of a vertical string of characters
const MIN_STRING_LENGTH = 8

func initMatrix(xmax int, ymax int, config *Config) ([][]rune, []tcell.Style, tcell.Style, []uint64) {
	matrix := make([][]rune, xmax)
	colorGradient := make([]tcell.Style, ymax)
	columnSpeed := make([]uint64, xmax)
	colors := config.colors

	styleAttribute := tcell.AttrNone
	if config.bold {
		styleAttribute = tcell.AttrBold
	}

	for i := 0; i < ymax; i++ {
		if colors.start == colors.end && i > 0 {
			colorGradient[i] = colorGradient[0]
		} else {
			colorGradient[i] = tcell.StyleDefault.Foreground(
				pickBetweenGradient(
					colors.start,
					colors.end,
					float32(i)/float32(ymax),
				),
			).Attributes(styleAttribute)
		}
	}

	for i := range matrix {
		matrix[i] = make([]rune, ymax)
		if config.async {
			columnSpeed[i] = uint64(rand.Intn(6) + 5)
		} else {
			columnSpeed[i] = 7
		}
	}

	whiteStyle := tcell.StyleDefault.Foreground(
		tcell.ColorWhite,
	).Attributes(styleAttribute)

	return matrix, colorGradient, whiteStyle, columnSpeed
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Matrix(xmax *int, ymax *int, waitTimeMs *uint64, _config *Config, _s *tcell.Screen) {
	xmaxOld := *xmax
	ymaxOld := *ymax
	config := *_config
	minStringLength := min(MIN_STRING_LENGTH, *ymax)

	s := *_s

	matrix, colorGradient, whiteStyle, columnSpeed := initMatrix(*xmax, *ymax, &config)

	createHead := func(column int, row int) {
		matrix[column][row] = rune(rand.Intn(94) + 33)
		s.SetContent(column, row, matrix[column][row], nil, whiteStyle)
	}

	for loopCounter := uint64(0); ; loopCounter++ {
		s.Show()
		afterLastDraw := time.Now()

		if (xmaxOld != *xmax) || (ymaxOld != *ymax) {
			s.Clear()
			xmaxOld = *xmax
			ymaxOld = *ymax
			minStringLength = min(MIN_STRING_LENGTH, *ymax)
			matrix, colorGradient, whiteStyle, columnSpeed = initMatrix(*xmax, *ymax, &config)
		}

		if minStringLength < 4 {
			continue
		}

		for column := range matrix {
			columnShouldMove := loopCounter%columnSpeed[column] == 0
			if !columnShouldMove {
				continue
			}
			last := len(matrix[column]) - 1
			for row := last; row >= 0; row-- {
				if row != 0 {
					if matrix[column][row-1] == 0 {
						// if the character above is empty, move it down (chop the tail)
						if matrix[column][row] != 0 {
							matrix[column][row] = 0
							s.SetContent(column, row, matrix[column][row], nil, whiteStyle)
						}
					} else if matrix[column][row] == 0 {
						// if the character above is not empty and the current character is empty, create more head
						createHead(column, row)
						// change previous head to the appropriate color
						s.SetContent(column, row-1, matrix[column][row-1], nil, colorGradient[row-1])

						// fuck the police 😎
						// here we are at the head so we know that we have at least minStringLength chars that we can skip because they are drawn
						row -= (minStringLength - 1)
					}

					if (row == last) && (matrix[column][row] != 0) {
						s.SetContent(column, row, matrix[column][row], nil, colorGradient[row])
					}
				} else {
					// row == 0
					if matrix[column][row] == 0 {
						// empty cell
						if rand.Intn(*xmax/2) == 1 {
							// begin new head
							createHead(column, row)
						}
					} else {
						// cell with content
						if rand.Intn(*ymax/2) == 1 {
							// this vertical-string has been chosen to be ended if it has the minimum length

							hasMinLength := true
							for i := 0; i < minStringLength; i++ {
								if matrix[column][row+i] == 0 {
									hasMinLength = false
									break
								}
							}

							if hasMinLength {
								// finish this column-string
								matrix[column][row] = 0
								s.SetContent(column, row, matrix[column][row], nil, whiteStyle)
							}
						}
					}
				}
			}
		}

		duration := time.Since(afterLastDraw)
		time.Sleep(time.Duration((*waitTimeMs)-uint64(duration.Milliseconds())) * time.Millisecond)
	}
}
