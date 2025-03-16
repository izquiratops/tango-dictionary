package server

import "unicode"

type SearchTermType string

const (
	Romaji SearchTermType = "romaji"
	Kana   SearchTermType = "kana"
	Kanji  SearchTermType = "kanji"
)

func DetectSearchTermType(text string) SearchTermType {
	hasKana := false

	for _, r := range text {
		// Check for Kanji (CJK Unified Ideographs)
		if unicode.Is(unicode.Han, r) {
			return Kanji
		}

		if unicode.Is(unicode.Hiragana, r) {
			hasKana = true
			continue
		}

		if unicode.Is(unicode.Katakana, r) {
			hasKana = true
			continue
		}

		// For non-Japanese characters like spaces, punctuation, etc.
		// Just continue without marking as any specific script
		// This allows for things like "こんにちは!" to still be classified as kana
		if unicode.IsPunct(r) || unicode.IsSpace(r) || unicode.IsNumber(r) {
			continue
		}
	}

	// If we found any kana characters and no kanji, it's kana
	if hasKana {
		return Kana
	}

	// Default to romaji
	return Romaji
}
