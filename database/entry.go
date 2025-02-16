package database

import (
	"encoding/json"
	"fmt"
	TangoUtil "tango/util"
)

type BleveEntry struct { // Simplified version of JMdictWord
	ID       string   `json:"id"`
	Kanji    []string `json:"kanji"`
	Kana     []string `json:"kana"`
	Meanings []string `json:"meanings"`
}

func (d *JMdictWord) ToBleveEntry() (BleveEntry, error) {
	entry := BleveEntry{
		ID:       d.ID, // Ids not indexed
		Kanji:    make([]string, 0),
		Kana:     make([]string, 0),
		Meanings: make([]string, 0),
	}

	for _, k := range d.Kanji {
		if k.Text == "" {
			return entry, fmt.Errorf("emtpy field at %v", d.ID)
		}

		entry.Kanji = append(entry.Kanji, k.Text)
	}

	for _, k := range d.Kana {
		if k.Text == "" {
			return entry, fmt.Errorf("emtpy field at %v", d.ID)
		}

		entry.Kana = append(entry.Kana, k.Text)
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
		Kanji    interface{} `json:"kanji"`
		Kana     interface{} `json:"kana"`
		Meanings interface{} `json:"meanings"`
		*Alias
	}{
		Alias: (*Alias)(be),
	}

	// Unmarshal the JSON data into the temporary struct
	if err := json.Unmarshal(data, temp); err != nil {
		return err
	}

	// Convert Kanji, Kana, and Meanings to []string
	be.Kanji = TangoUtil.EnsureSlice(temp.Kanji)
	be.Kana = TangoUtil.EnsureSlice(temp.Kana)
	be.Meanings = TangoUtil.EnsureSlice(temp.Meanings)

	return nil
}
