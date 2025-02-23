package model

import (
	"encoding/json"
	"fmt"
	"tango/util"
)

type Xref struct {
	Kanji       *string `json:"kanji,omitempty" bson:"kanji,omitempty"`
	Kana        *string `json:"kana,omitempty" bson:"kana,omitempty"`
	KanjiOrKana *string `json:"kanji_or_kana,omitempty" bson:"kanji_or_kana,omitempty"`
	SenseIndex  *int    `json:"senseIndex,omitempty" bson:"senseIndex,omitempty"`
}

func (x *Xref) UnmarshalJSON(data []byte) error {
	var fullRef [3]interface{}
	if err := json.Unmarshal(data, &fullRef); err != nil {
		return err
	}

	if fullRef[2] != nil {
		// 1. Try to unmarshal as [kanji, kana, senseIndex]
		x.Kanji = util.ToStringPtr(fullRef[0])
		x.Kana = util.ToStringPtr(fullRef[1])
		x.SenseIndex = util.ToIntPtr(fullRef[2])
		return nil
	}

	if fullRef[1] != nil {
		kana := util.ToStringPtr(fullRef[1])
		senseIndex := util.ToIntPtr(fullRef[1])

		if kana != nil {
			// 2. Try to unmarshal as [kanji, kana]
			x.Kanji = util.ToStringPtr(fullRef[0])
			x.Kana = kana
			return nil
		} else if senseIndex != nil {
			// 3. Try to unmarshal as [kanjiOrKana, senseIndex]
			x.KanjiOrKana = util.ToStringPtr(fullRef[0])
			x.SenseIndex = senseIndex
			return nil
		}
	}

	if fullRef[0] != nil {
		// 4. Try to unmarshal as [kanjiOrKana]
		x.KanjiOrKana = util.ToStringPtr(fullRef[0])
		return nil
	}

	return fmt.Errorf("invalid JSON format for Xref")
}
