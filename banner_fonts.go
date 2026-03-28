package main

import (
	"fmt"
	"os"
	"strings"
)

var allowedBannerFonts = []string{
	"univers", "thick", "tanja", "stellar", "starwars", "trek", "tinker-toy",
	"thin", "standard", "sblood", "rozzo", "rowancap", "roman", "poison",
	"pebbles", "o8", "ntgreek", "nancyj", "nancyj-underlined", "nancyj-fancy",
	"lean", "larry3d", "kban", "jazmine", "ivrit", "isometric4", "isometric3",
	"isometric2", "isometric1", "hollywood", "graffiti", "gothic", "fender",
	"epic", "dotmatrix", "doom", "doh", "computer", "colossal", "caligraphy",
	"calgphy2", "big", "bell", "basic", "banner4", "banner3", "banner3-D",
	"banner", "alligator2", "alligator", "acrobatic",
}

var allowedBannerFontSet = func() map[string]struct{} {
	fonts := make(map[string]struct{}, len(allowedBannerFonts))
	for _, font := range allowedBannerFonts {
		fonts[font] = struct{}{}
	}

	return fonts
}()

func bannerFontsHelpText() string {
	return strings.Join(allowedBannerFonts, ", ")
}

func isAllowedBannerFont(font string) bool {
	if font == "" {
		return true
	}

	_, ok := allowedBannerFontSet[font]
	return ok
}

func printBannerFontHelp() {
	fmt.Printf("\nAllowed --banner-font values:\n%s\n", bannerFontsHelpText())
}

func exitInvalidBannerFont(font string) {
	fmt.Printf("Error: invalid --banner-font %q\nAllowed values: %s\n", font, bannerFontsHelpText())
	os.Exit(1)
}
