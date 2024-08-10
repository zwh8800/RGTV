package model

type EpgResp struct {
	Date        string `json:"date"`
	ChannelName string `json:"channel_name"`
	URL         string `json:"url"`
	EpgData     []struct {
		Start string `json:"start"`
		Desc  string `json:"desc"`
		End   string `json:"end"`
		Title string `json:"title"`
	} `json:"epg_data"`
}
