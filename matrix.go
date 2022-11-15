package main

import (
	"math/rand"
	"time"

	"github.com/gdamore/tcell/v2"
)

// the minimum length of a vertical string of characters
const MIN_STRING_LENGTH = 8

func initMatrix(xmax int, ymax int, config *Config) ([][]rune, []tcell.Style, tcell.Style, []uint64) {
	// the characters
	matrix := make([][]rune, xmax)
	// the color of each row
	colorGradient := make([]tcell.Style, ymax)
	// controls the speed
	columnDrag := make([]uint64, xmax)
	// the start and end colors
	colors := config.colors

	characterStyleAttr := tcell.AttrNone
	if config.bold {
		characterStyleAttr = tcell.AttrBold
	}

	if config.pride {
		prideColors := []string{"#cc3516", "#eb9528", "#faee3b", "#3d7e2a", "#2d4efa", "#691a85"}
		prideColorsCount := len(prideColors)
		prideStyles := make([]tcell.Style, len(prideColors))

		for i, color := range prideColors {
			prideStyles[i] = tcell.StyleDefault.Foreground(tcell.GetColor(color)).Attributes(characterStyleAttr)
		}

		for i := 0; i < ymax; i++ {
			colorGradient[i] = prideStyles[int((float32(i)/float32(ymax))*float32(prideColorsCount))]
		}
	} else {
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
				).Attributes(characterStyleAttr)
			}
		}
	}

	for i := range matrix {
		matrix[i] = make([]rune, ymax)
		if config.async {
			columnDrag[i] = uint64(rand.Intn(6) + 5)
		} else {
			columnDrag[i] = 7
		}
	}

	whiteStyle := tcell.StyleDefault.Foreground(
		tcell.ColorWhite,
	).Attributes(characterStyleAttr)

	return matrix, colorGradient, whiteStyle, columnDrag
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func matrix(xmax *int, ymax *int, waitTimeMs *uint64, _config *Config, _s *tcell.Screen) {
	xmaxOld := *xmax
	ymaxOld := *ymax
	config := *_config
	currentMinStringLength := min(MIN_STRING_LENGTH, *ymax)

	s := *_s

	matrix, colorGradient, whiteStyle, columnDrag := initMatrix(*xmax, *ymax, &config)

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
			currentMinStringLength = min(MIN_STRING_LENGTH, *ymax)
			matrix, colorGradient, whiteStyle, columnDrag = initMatrix(*xmax, *ymax, &config)
		}

		if *ymax < 4 {
			// too small of a terminal to do something here
			continue
		}

		for column := range matrix {
			columnShouldMove := loopCounter%columnDrag[column] == 0
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
						row -= (currentMinStringLength - 1)
					}

					if (row == last) && (matrix[column][row] != 0) {
						s.SetContent(column, row, matrix[column][row], nil, colorGradient[row])
					}
				} else {
					// row == 0
					if matrix[column][row] == 0 {
						// empty cell
						if rand.Intn(350/int(columnDrag[column])) == 1 {
							// begin new head
							createHead(column, row)
						}
					} else {
						// cell with content
						if rand.Intn(100/int(columnDrag[column])) == 1 {
							// this vertical-string has been chosen to be ended if it has the minimum length

							hasMinLength := true
							for i := 0; i < currentMinStringLength; i++ {
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
