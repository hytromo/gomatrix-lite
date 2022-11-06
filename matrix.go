package main

import (
	"math/rand"
	"time"

	"github.com/gdamore/tcell/v2"
)

// the minimum length of a vertical string of characters
const MIN_STRING_LENGTH = 8

func initMatrix(xmax int, ymax int, colors *Colors) ([][]rune, []tcell.Style) {
	matrix := make([][]rune, xmax)
	colorGradient := make([]tcell.Style, ymax)

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
			)
		}
	}

	for i := range matrix {
		matrix[i] = make([]rune, ymax)
	}

	return matrix, colorGradient
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Matrix(xmax *int, ymax *int, waitTimeMs *int64, colors *Colors, _s *tcell.Screen) {
	xmaxOld := *xmax
	ymaxOld := *ymax
	minStringLength := min(MIN_STRING_LENGTH, *ymax)

	s := *_s

	matrix, colorGradient := initMatrix(*xmax, *ymax, colors)

	whiteStyle := tcell.StyleDefault.Foreground(
		tcell.ColorWhite,
	)

	createHead := func(column int, row int) {
		matrix[column][row] = rune(rand.Intn(94) + 33)
		s.SetContent(column, row, matrix[column][row], nil, whiteStyle)
	}

	for {
		s.Show()
		afterLastDraw := time.Now()

		if (xmaxOld != *xmax) || (ymaxOld != *ymax) {
			s.Clear()
			xmaxOld = *xmax
			ymaxOld = *ymax
			minStringLength = min(MIN_STRING_LENGTH, *ymax)
			matrix, colorGradient = initMatrix(*xmax, *ymax, colors)
		}

		if minStringLength < 4 {
			continue
		}

		for column := range matrix {
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
		time.Sleep(time.Duration((*waitTimeMs)-duration.Milliseconds()) * time.Millisecond)
	}
}
