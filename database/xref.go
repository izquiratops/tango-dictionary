package main

import (
	"encoding/json"
	"fmt"
	JishoUtil "jisho-clone-database/util"
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
		x.Kanji = JishoUtil.ToStringPtr(fullRef[0])
		x.Kana = JishoUtil.ToStringPtr(fullRef[1])
		x.SenseIndex = JishoUtil.ToIntPtr(fullRef[2])
		return nil
	}

	if fullRef[1] != nil {
		kana := JishoUtil.ToStringPtr(fullRef[1])
		senseIndex := JishoUtil.ToIntPtr(fullRef[1])

		if kana != nil {
			// 2. Try to unmarshal as [kanji, kana]
			x.Kanji = JishoUtil.ToStringPtr(fullRef[0])
			x.Kana = kana
			return nil
		} else if senseIndex != nil {
			// 3. Try to unmarshal as [kanjiOrKana, senseIndex]
			x.KanjiOrKana = JishoUtil.ToStringPtr(fullRef[0])
			x.SenseIndex = senseIndex
			return nil
		}
	}

	if fullRef[0] != nil {
		// 4. Try to unmarshal as [kanjiOrKana]
		x.KanjiOrKana = JishoUtil.ToStringPtr(fullRef[0])
		return nil
	}

	return fmt.Errorf("invalid JSON format for Xref")
}
