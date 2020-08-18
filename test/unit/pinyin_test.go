package main

import (
	"github.com/Chain-Zhang/pinyin"
	"testing"
)

func TestPinyin(t *testing.T) {
	sent := "我是中国人"
	str, err := pinyin.New(sent).Split(" ").Mode(pinyin.WithoutTone).Convert()
	if err == nil {
		t.Log(str)
	} else {
		t.Error("fail to parse " + sent)
	}
}