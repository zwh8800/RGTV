package model

import (
	"encoding/json"
	"math"
	"testing"
)

func TestParseChannelFromM3U8(t *testing.T) {
	data, err := ParseChannelFromM3U8("https://raw.githubusercontent.com/imDazui/Tvlist-awesome-m3u-m3u8/master/m3u/%E7%99%BE%E8%A7%86TV.m3u")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tojson(data))

	data, err = ParseChannelFromM3U8("https://raw.githubusercontent.com/imDazui/Tvlist-awesome-m3u-m3u8/master/m3u/%E5%9B%BD%E5%86%85%E7%94%B5%E8%A7%86%E5%8F%B02023.m3u8")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tojson(data))
}

func TestParseChannelFromDIYP(t *testing.T) {
	data, err := ParseChannelFromDIYP("https://x-x-xxx.github.io/diyp/tv.txt")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tojson(data))
}

func TestParseChannelFromM3U8Err(t *testing.T) {
	data, err := ParseChannelFromM3U8("https://x-x-xxx.github.io/diyp/tv.txt")
	if err != nil {
		t.Log("ok", data)
		return
	}
	t.Fatal("err is nil")
}

func tojson(data interface{}) string {
	j, _ := json.Marshal(data)
	return string(j)
}

func TestSSss(t *testing.T) {
	y := -5
	t.Log(math.Pow(10, float64(y)/10+1))
}
