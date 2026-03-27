// the only package of this app
package main

import (
	"math/rand"
	"time"

	"github.com/gdamore/tcell/v2"
)

// the minimum length of a vertical string of characters
const minStringOfCharsLength = 8

func initMatrix(xmax int, ymax int, config *Config) ([][]rune, []tcell.Style, tcell.Style, []uint64) {
	// the characters
	matrix := make([][]rune, xmax)
	// the color of each row
	colorGradient := make([]tcell.Style, ymax)
	// controls how often each column should move
	columnMovesEveryNLoops := make([]uint64, xmax)
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

	for column := range matrix {
		matrix[column] = make([]rune, ymax)
		if config.async {
			columnMovesEveryNLoops[column] = uint64(rand.Intn(6) + 5)
		} else {
			columnMovesEveryNLoops[column] = 7
		}
	}

	whiteStyle := tcell.StyleDefault.Foreground(
		tcell.ColorWhite,
	).Attributes(characterStyleAttr)

	return matrix, colorGradient, whiteStyle, columnMovesEveryNLoops
}

func matrix(xmax *int, ymax *int, waitTimeMs *uint64, _config *Config, _s *tcell.Screen) {
	config := *_config

	var matrix [][]rune
	var colorGradient []tcell.Style
	var whiteStyle tcell.Style
	var columnMovesEveryNLoops []uint64
	var xmaxOld, ymaxOld, currentMinStringOfCharsLength int

	resetTermVariables := func() {
		xmaxOld = *xmax
		ymaxOld = *ymax
		currentMinStringOfCharsLength = min(minStringOfCharsLength, *ymax)
		matrix, colorGradient, whiteStyle, columnMovesEveryNLoops = initMatrix(*xmax, *ymax, &config)
	}

	resetTermVariables()

	s := *_s

	createHead := func(column int, row int) {
		matrix[column][row] = rune(rand.Intn(94) + 33)
		s.SetContent(column, row, matrix[column][row], nil, whiteStyle)
	}

	for loopCounter := uint64(0); ; loopCounter++ {
		s.Show()
		timeAfterLastDraw := time.Now()

		if (xmaxOld != *xmax) || (ymaxOld != *ymax) {
			// terminal has been resized
			s.Clear()
			resetTermVariables()
		}

		if *ymax < 4 {
			// terminal is too small to run the matrix
			continue
		}

		for columnIndex := range matrix {
			if loopCounter%columnMovesEveryNLoops[columnIndex] != 0 {
				// this column should not move this loop
				continue
			}
			lastRowIndex := len(matrix[columnIndex]) - 1
			for rowIndex := lastRowIndex; rowIndex >= 0; rowIndex-- {
				if rowIndex != 0 {
					// we are not at the first line

					if matrix[columnIndex][rowIndex-1] == 0 { // 0 means nil / empty rune = empty cell
						if matrix[columnIndex][rowIndex] != 0 {
							// if the character above is empty and the current is not, move the current down (chop the tail)
							matrix[columnIndex][rowIndex] = 0
							s.SetContent(columnIndex, rowIndex, matrix[columnIndex][rowIndex], nil, whiteStyle)
						}
					} else if matrix[columnIndex][rowIndex] == 0 {
						// if the character above is not empty and the current character is empty, create more head
						// this just means that the column moves downwards (new character at the end of the column)
						createHead(columnIndex, rowIndex)
						// previous head was white (it always is), so we need to change it to the appropriate color now
						// that it became body
						s.SetContent(columnIndex, rowIndex-1, matrix[columnIndex][rowIndex-1], nil, colorGradient[rowIndex-1])

						// here we are at the head so we know that we have at least minStringLength chars above us that we can skip because they are drawn already correctly
						rowIndex -= (currentMinStringOfCharsLength - 1)
					} else if (rowIndex == lastRowIndex) && (matrix[columnIndex][rowIndex] != 0) {
						// we are at the last row and we already have content from the previous loop, so let's remove the white
						// (head moves off screen)
						_, _, style, _ := s.GetContent(columnIndex, rowIndex)
						if style == whiteStyle {
							s.SetContent(columnIndex, rowIndex, matrix[columnIndex][rowIndex], nil, colorGradient[rowIndex])
						}
					}
				} else {
					// rowIndex == 0
					if matrix[columnIndex][rowIndex] == 0 {
						// empty cell, add a chance to create a new head
						// 60 is just a nice number that produces not too many or too few heads
						if rand.Intn(60) == 42 {
							// begin new head
							createHead(columnIndex, rowIndex)
						}
					} else {
						// cell with content
						if rand.Intn(10) == 1 {
							// this vertical-string has been chosen to be ended if it has at least the minimum length

							hasMinLength := true
							for i := 0; i < currentMinStringOfCharsLength; i++ {
								if matrix[columnIndex][rowIndex+i] == 0 {
									hasMinLength = false
									break
								}
							}

							if hasMinLength {
								// finish this column-string
								matrix[columnIndex][rowIndex] = 0
								s.SetContent(columnIndex, rowIndex, matrix[columnIndex][rowIndex], nil, whiteStyle)
							}
						}
					}
				}
			}
		}

		duration := time.Since(timeAfterLastDraw)
		time.Sleep(time.Duration((*waitTimeMs)-uint64(duration.Milliseconds())) * time.Millisecond)
	}
}
