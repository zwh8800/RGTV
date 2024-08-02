package model

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/jamesnetherton/m3u"
	"github.com/wk8/go-ordered-map/v2"
)

type ChannelData struct {
	Groups []*ChannelGroup
}

type ChannelGroup struct {
	Name     string
	Channels []*Channel
}

type Channel struct {
	Name   string
	Url    string
	Source string
}

func ParseChannelFromDIYP(path string) (*ChannelData, error) {
	var f io.ReadCloser

	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		data, err := http.Get(path)
		if err != nil {
			return nil, fmt.Errorf("unable to open playlist URL: %v", err)
		}
		f = data.Body
	} else {
		file, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("unable to open playlist file: %v", err)
		}
		f = file
	}
	defer f.Close()

	groups := orderedmap.New[string, *ChannelGroup]()

	curGroup := "直播"
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, ",#genre#") {
			curGroup = strings.Split(line, ",")[0]
		} else {
			part := strings.Split(line, ",")
			if len(part) < 2 {
				continue
			}
			name := part[0]
			url := part[1]
			source := ""
			part2 := strings.Split(url, "$")
			if len(part2) == 2 {
				url = part2[0]
				source = part2[1]
			}

			group, ok := groups.Get(curGroup)
			if !ok {
				groups.Set(curGroup, &ChannelGroup{
					Name:     curGroup,
					Channels: make([]*Channel, 0),
				})
				group, _ = groups.Get(curGroup)
			}
			group.Channels = append(group.Channels, &Channel{
				Name:   name,
				Url:    url,
				Source: source,
			})
		}
	}
	channelData := &ChannelData{
		Groups: make([]*ChannelGroup, 0, groups.Len()),
	}
	for pair := groups.Oldest(); pair != nil; pair = pair.Next() {
		channelData.Groups = append(channelData.Groups, pair.Value)
	}
	return channelData, nil
}

func ParseChannelFromM3U8(path string) (*ChannelData, error) {
	m, err := m3u.Parse(path)
	if err != nil {
		return nil, err
	}
	groups := orderedmap.New[string, *ChannelGroup]()
	for _, track := range m.Tracks {
		groupTitle := getGroupTitle(track)
		if groupTitle == "" {
			groupTitle = "直播"
		}
		group, ok := groups.Get(groupTitle)
		if !ok {
			groups.Set(groupTitle, &ChannelGroup{
				Name:     groupTitle,
				Channels: make([]*Channel, 0),
			})
			group, _ = groups.Get(groupTitle)
		}
		name := getTvgName(track)
		if name == "" {
			name = track.Name
		}
		group.Channels = append(group.Channels, &Channel{
			Name:   name,
			Url:    track.URI,
			Source: "",
		})
	}
	channelData := &ChannelData{
		Groups: make([]*ChannelGroup, 0, groups.Len()),
	}
	for pair := groups.Oldest(); pair != nil; pair = pair.Next() {
		channelData.Groups = append(channelData.Groups, pair.Value)
	}
	return channelData, nil
}

func getGroupTitle(track m3u.Track) string {
	for _, tag := range track.Tags {
		if tag.Name == "group-title" {
			return tag.Value
		}
	}
	return ""
}

func getTvgName(track m3u.Track) string {
	for _, tag := range track.Tags {
		if tag.Name == "tvg-name" {
			return tag.Value
		}
	}
	return ""
}
