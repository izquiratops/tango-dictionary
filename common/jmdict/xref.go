package jmdict

import (
	"encoding/json"
	"fmt"

	"github.com/izquiratops/tango/common/utils"
)

// Xref is a struct that can be 4 different kinds:
//
//	XrefWordReadingIndex | XrefWordReading | XrefWordIndex | XrefWord
//
// It's a valid thing in TypeScript, but Go doesn't let you have it.
// This custom unmarshaller parse Xref to a single struct as a workaround.
func (x *Xref) UnmarshalJSON(data []byte) error {
	var fullRef [3]interface{}
	if err := json.Unmarshal(data, &fullRef); err != nil {
		return err
	}

	if fullRef[2] != nil {
		// 1. Try to unmarshal as [kanji, kana, senseIndex]
		x.Kanji = utils.ToStringPtr(fullRef[0])
		x.Kana = utils.ToStringPtr(fullRef[1])
		x.SenseIndex = utils.ToIntPtr(fullRef[2])
		return nil
	}

	if fullRef[1] != nil {
		kana := utils.ToStringPtr(fullRef[1])
		senseIndex := utils.ToIntPtr(fullRef[1])

		if kana != nil {
			// 2. Try to unmarshal as [kanji, kana]
			x.Kanji = utils.ToStringPtr(fullRef[0])
			x.Kana = kana
			return nil
		} else if senseIndex != nil {
			// 3. Try to unmarshal as [kanjiOrKana, senseIndex]
			x.KanjiOrKana = utils.ToStringPtr(fullRef[0])
			x.SenseIndex = senseIndex
			return nil
		}
	}

	if fullRef[0] != nil {
		// 4. Try to unmarshal as [kanjiOrKana]
		x.KanjiOrKana = utils.ToStringPtr(fullRef[0])
		return nil
	}

	return fmt.Errorf("invalid JSON format for Xref")
}
