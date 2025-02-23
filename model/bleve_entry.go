package model

import (
	"encoding/json"
	"fmt"
	"tango/util"
)

type BleveEntry struct { // Simplified version of JMdictWord
	ID         string   `json:"id"`
	KanjiExact []string `json:"kanji_exact"`
	KanjiChar  []string `json:"kanji_char"`
	KanaExact  []string `json:"kana_exact"`
	KanaChar   []string `json:"kana_char"`
	Meanings   []string `json:"meanings"`
}

func (d *JMdictWord) ToBleveEntry() (BleveEntry, error) {
	entry := BleveEntry{
		ID:         d.ID, // ID not indexed
		KanjiExact: make([]string, 0),
		KanjiChar:  make([]string, 0),
		KanaExact:  make([]string, 0),
		KanaChar:   make([]string, 0),
		Meanings:   make([]string, 0),
	}

	for _, k := range d.Kanji {
		if k.Text == "" {
			return entry, fmt.Errorf("emtpy field at %v", d.ID)
		}

		entry.KanjiExact = append(entry.KanjiExact, k.Text)
		entry.KanjiChar = append(entry.KanjiChar, k.Text)
	}

	for _, k := range d.Kana {
		if k.Text == "" {
			return entry, fmt.Errorf("emtpy field at %v", d.ID)
		}

		entry.KanaExact = append(entry.KanaExact, k.Text)
		entry.KanaChar = append(entry.KanaChar, k.Text)
	}

	for _, s := range d.Sense {
		for _, g := range s.Gloss {
			if g.Lang == "eng" {
				if g.Text == "" {
					return entry, fmt.Errorf("emtpy field at %v", d.ID)
				}

				entry.Meanings = append(entry.Meanings, g.Text)
			}
		}
	}

	return entry, nil
}

func (be *BleveEntry) UnmarshalJSON(data []byte) error {
	// Define a temporary struct to unmarshal the JSON data
	type Alias BleveEntry
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
	be.KanjiExact = util.EnsureSlice(temp.KanjiExact)
	be.KanjiChar = util.EnsureSlice(temp.KanjiChar)
	be.KanaExact = util.EnsureSlice(temp.KanaExact)
	be.KanaChar = util.EnsureSlice(temp.KanaChar)
	be.Meanings = util.EnsureSlice(temp.Meanings)

	return nil
}
