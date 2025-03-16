package database

import (
	"encoding/json"

	"github.com/izquiratops/tango/common/utils"
)

func (be *WordSearchable) UnmarshalJSON(data []byte) error {
	// Define a temporary struct to unmarshal the JSON data
	type Alias WordSearchable
	temp := &struct {
		KanjiExact any `json:"kanji_exact"`
		KanjiChar  any `json:"kanji_char"`
		KanaExact  any `json:"kana_exact"`
		KanaChar   any `json:"kana_char"`
		Meanings   any `json:"meanings"`
		Romaji     any `json:"romaji"`
		*Alias
	}{
		Alias: (*Alias)(be),
	}

	// Unmarshal the JSON data into the temporary struct
	if err := json.Unmarshal(data, temp); err != nil {
		return err
	}

	// Convert Kanji, Kana, and Meanings to []string
	be.KanjiExact = utils.EnsureSlice(temp.KanjiExact)
	be.KanjiChar = utils.EnsureSlice(temp.KanjiChar)
	be.KanaExact = utils.EnsureSlice(temp.KanaExact)
	be.KanaChar = utils.EnsureSlice(temp.KanaChar)
	be.Meanings = utils.EnsureSlice(temp.Meanings)
	be.Romaji = utils.EnsureSlice(temp.Romaji)

	return nil
}
