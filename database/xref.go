package main

import (
	"encoding/json"
	"fmt"
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
		x.Kanji = ToStringPtr(fullRef[0])
		x.Kana = ToStringPtr(fullRef[1])
		x.SenseIndex = ToIntPtr(fullRef[2])
		return nil
	}

	if fullRef[1] != nil {
		kana := ToStringPtr(fullRef[1])
		senseIndex := ToIntPtr(fullRef[1])

		if kana != nil {
			// 2. Try to unmarshal as [kanji, kana]
			x.Kanji = ToStringPtr(fullRef[0])
			x.Kana = kana
			return nil
		} else if senseIndex != nil {
			// 3. Try to unmarshal as [kanjiOrKana, senseIndex]
			x.KanjiOrKana = ToStringPtr(fullRef[0])
			x.SenseIndex = senseIndex
			return nil
		}
	}

	if fullRef[0] != nil {
		// 4. Try to unmarshal as [kanjiOrKana]
		x.KanjiOrKana = ToStringPtr(fullRef[0])
		return nil
	}

	return fmt.Errorf("invalid JSON format for Xref")
}

func ToStringPtr(v interface{}) *string {
	if str, ok := v.(string); ok {
		return &str
	}
	return nil
}

func ToIntPtr(v interface{}) *int {
	if floatNum, ok := v.(float64); ok {
		num := int(floatNum)
		return &num
	}
	return nil
}
