package conf

import (
	"encoding/json"
	"os"
)

type Type struct {
	ResX                int32  `json:"res_x"`
	ResY                int32  `json:"res_y"`
	LiveSourceUrl       string `json:"live_source_url"`
	EPGUrl              string `json:"epg_url"`
	FFMPEGPath          string `json:"ffmpeg_path"`
	RevertSwitchChannel bool   `json:"revert_switch_channel"`
}

var config = &Type{}

func init() {
	data, err := os.ReadFile("tv.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, config)
	if err != nil {
		panic(err)
	}
	if config.ResX == 0 {
		config.ResX = 640
	}
	if config.ResY == 0 {
		config.ResY = 480
	}
	if config.FFMPEGPath == "" {
		config.FFMPEGPath = "ffmpeg"
	}
}

func GetConfig() Type {
	return *config
}
