package utils

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var nonAlphanumericSpace = regexp.MustCompile(`[^a-zA-Z0-9]+`)

// NormalizeLocationName prepara un nombre geográfico para comparación flexible:
// quita tildes, pasa a mayúsculas, elimina puntuación y colapsa espacios.
// Ej: "Bogotá D.C." → "BOGOTA DC", "MEDELLÍN" → "MEDELLIN".
func NormalizeLocationName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return ""
	}

	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	withoutAccents, _, _ := transform.String(t, name)
	upper := strings.ToUpper(withoutAccents)
	cleaned := nonAlphanumericSpace.ReplaceAllString(upper, " ")
	return strings.Join(strings.Fields(cleaned), " ")
}

// ScoreCityNameMatch puntúa qué tan bien normalizedCity representa normalizedInput.
// Valores más altos indican mejor coincidencia.
func ScoreCityNameMatch(normalizedInput, normalizedCity string) int {
	if normalizedInput == "" || normalizedCity == "" {
		return 0
	}

	switch {
	case normalizedInput == normalizedCity:
		return 1000
	case strings.HasPrefix(normalizedCity, normalizedInput):
		return 900 - (len(normalizedCity) - len(normalizedInput))
	case strings.HasPrefix(normalizedInput, normalizedCity):
		return 850 - (len(normalizedInput) - len(normalizedCity))
	case strings.Contains(normalizedCity, normalizedInput):
		return 700 + (len(normalizedInput)*100)/len(normalizedCity)
	case strings.Contains(normalizedInput, normalizedCity):
		return 650 + (len(normalizedCity)*100)/len(normalizedInput)
	}

	inputTokens := strings.Fields(normalizedInput)
	cityTokens := strings.Fields(normalizedCity)
	if len(inputTokens) == 0 {
		return 0
	}

	matched := 0
	for _, inputToken := range inputTokens {
		for _, cityToken := range cityTokens {
			if inputToken == cityToken ||
				strings.HasPrefix(cityToken, inputToken) ||
				strings.HasPrefix(inputToken, cityToken) {
				matched++
				break
			}
		}
	}

	if matched == 0 {
		return 0
	}
	return 400 + (matched*100)/len(inputTokens)
}
