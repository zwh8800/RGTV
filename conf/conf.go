package conf

import (
	"encoding/json"
	"os"
)

type Type struct {
	LiveSourceUrl string `json:"live_source_url"`
	EPGUrl        string `json:"epg_url"`
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
}

func GetConfig() Type {
	return *config
}
