package database

import (
	"encoding/json"

	"github.com/izquiratops/tango/common/utils"
)

func (be *WordSearchable) UnmarshalJSON(data []byte) error {
	// Define a temporary struct to unmarshal the JSON data
	type Alias WordSearchable
	temp := &struct {
		KanjiExact interface{} `json:"kanji_exact"`
		KanjiChar  interface{} `json:"kanji_char"`
		KanaExact  interface{} `json:"kana_exact"`
		KanaChar   interface{} `json:"kana_char"`
		Meanings   interface{} `json:"meanings"`
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

	return nil
}
