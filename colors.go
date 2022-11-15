// the only package of this app
package main

import (
	"strings"

	"github.com/gdamore/tcell/v2"
)

// Colors contains the start and end colors of the gradient
type Colors struct {
	start tcell.Color
	end   tcell.Color
}

func parseColors(color string) Colors {
	var startColor tcell.Color
	var endColor tcell.Color

	if strings.Contains(color, ",") {
		colorsArr := strings.Split(color, ",")

		startColor = tcell.GetColor("#" + colorsArr[0])
		endColor = tcell.GetColor("#" + colorsArr[1])
	} else {
		startColor = tcell.GetColor("#" + color)
		endColor = startColor
	}

	return Colors{
		start: startColor,
		end:   endColor,
	}
}

func pickBetweenGradient(color1 tcell.Color, color2 tcell.Color, weight float32) tcell.Color {
	w2 := weight
	w1 := 1.0 - w2
	c1r, c1g, c1b := color1.RGB()
	c2r, c2g, c2b := color2.RGB()
	rgb := tcell.NewRGBColor(
		int32(float32(c1r)*w1+float32(c2r)*w2),
		int32(float32(c1g)*w1+float32(c2g)*w2),
		int32(float32(c1b)*w1+float32(c2b)*w2),
	)

	return rgb
}
