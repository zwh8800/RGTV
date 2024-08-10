package util

import (
	"fmt"
	"github.com/zwh8800/RGTV/conf"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	gocache "github.com/patrickmn/go-cache"
	"github.com/zwh8800/RGTV/model"
)

var epgCache = gocache.New(10*time.Minute, 10*time.Minute)

func GetEpg(ch string, date time.Time) (*model.EpgResp, error) {
	cli := resty.New()
	var epgResp model.EpgResp
	resp, err := cli.R().
		SetQueryParam("ch", ch).
		SetQueryParam("date", date.Format("2006-01-02")).
		SetResult(&epgResp).
		Get(conf.GetConfig().EPGUrl)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("status code: %d", resp.StatusCode())
	}
	return &epgResp, nil
}

func GetEpgFromCache(ch string, date time.Time) (*model.EpgResp, error) {
	key := ch + date.Format("2006-01-02")
	if v, ok := epgCache.Get(key); ok {
		return v.(*model.EpgResp), nil
	}
	epgResp, err := GetEpg(ch, date)
	if err != nil {
		return nil, err
	}
	epgCache.Set(key, epgResp, gocache.DefaultExpiration)
	return epgResp, nil
}

func GetCurrentProgram(ch string) string {
	now := time.Now()
	epgResp, err := GetEpgFromCache(ch, now)
	if err != nil {
		return ""
	}
	for _, program := range epgResp.EpgData {
		start, err := time.ParseInLocation("15:04", program.Start, now.Location())
		if err != nil {
			return ""
		}
		end, err := time.ParseInLocation("15:04", program.End, now.Location())
		if err != nil {
			return ""
		}
		start = start.AddDate(now.Year(), int(now.Month())-1, now.Day()-1)
		end = end.AddDate(now.Year(), int(now.Month())-1, now.Day()-1)
		if start.Before(now) && end.After(now) {
			return program.Title
		}
	}
	return ""
}

func GetNextProgram(ch string) string {
	now := time.Now()
	epgResp, err := GetEpgFromCache(ch, now)
	if err != nil {
		return ""
	}
	for _, program := range epgResp.EpgData {
		start, err := time.ParseInLocation("15:04", program.Start, now.Location())
		if err != nil {
			return ""
		}
		end, err := time.ParseInLocation("15:04", program.End, now.Location())
		if err != nil {
			return ""
		}
		start = start.AddDate(now.Year(), int(now.Month())-1, now.Day()-1)
		end = end.AddDate(now.Year(), int(now.Month())-1, now.Day()-1)
		if start.After(now) {
			return program.Title
		}
	}
	return ""
}
