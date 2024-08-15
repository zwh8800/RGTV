package conf

import (
	"encoding/json"
	"fmt"
	"os"
)

type Type struct {
	LiveSourceUrl       string `json:"live_source_url,omitempty"`
	EPGUrl              string `json:"epg_url,omitempty"`
	FFMPEGPath          string `json:"ffmpeg_path,omitempty"`
	RevertSwitchChannel bool   `json:"revert_switch_channel,omitempty"`
}

var config = &Type{}

func init() {
	data, err := os.ReadFile("tv.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, config)
	if err != nil {
		fmt.Println(err)
	}
	if config.FFMPEGPath == "" {
		config.FFMPEGPath = "ffmpeg"
	}
	if config.EPGUrl == "" {
		config.EPGUrl = "http://epg.112114.xyz/"
	}
}

func GetConfig() Type {
	return *config
}
