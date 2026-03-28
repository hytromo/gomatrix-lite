package main

import (
	_ "embed"
	"encoding/json"
	"math/rand"
	"strings"
	"sync"
)

//go:embed ascii_similar.json
var asciiSimilarJSON []byte

var (
	asciiSimilarOnce sync.Once
	asciiSimilarMap  map[string][]string
	asciiSimilarErr  error
)

func loadAsciiSimilar() (map[string][]string, error) {
	asciiSimilarOnce.Do(func() {
		asciiSimilarErr = json.Unmarshal(asciiSimilarJSON, &asciiSimilarMap)
	})

	return asciiSimilarMap, asciiSimilarErr
}

func replaceAsciiSimilar(s string) string {
	asciiSimilar, err := loadAsciiSimilar()
	if err != nil {
		return s
	}

	var builder strings.Builder
	builder.Grow(len(s))

	for _, r := range s {
		if r == '\n' {
			builder.WriteRune(r)
			continue
		}

		replacements := asciiSimilar[string(r)]
		if len(replacements) == 0 {
			builder.WriteRune(r)
			continue
		}

		builder.WriteString(replacements[rand.Intn(len(replacements))])
	}

	return builder.String()
}
