package model

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/jamesnetherton/m3u"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

type ChannelData struct {
	Groups []*ChannelGroup
}

type ChannelGroup struct {
	Index    int
	Name     string
	Channels []*Channel
}

type Channel struct {
	Index   int
	Name    string
	Sources []*Source
}

type Source struct {
	Name string
	Url  string
}

type sortByName []*Source

func (s sortByName) Len() int {
	return len(s)
}
func (s sortByName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s sortByName) Less(i, j int) bool {
	return s[i].Name < s[j].Name
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
	groupIndex := 1
	channelIndex := 1
	channelMap := make(map[string]*Channel)
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
			source := "默认线路"
			part2 := strings.Split(url, "$")
			if len(part2) == 2 {
				url = part2[0]
				source = part2[1]
			}

			group, ok := groups.Get(curGroup)
			if !ok {
				groups.Set(curGroup, &ChannelGroup{
					Index:    groupIndex,
					Name:     curGroup,
					Channels: make([]*Channel, 0),
				})
				group, _ = groups.Get(curGroup)
				groupIndex++
			}
			channel, ok := channelMap[name]
			if !ok {
				channel = &Channel{
					Index: channelIndex,
					Name:  name,
				}
				channelMap[name] = channel
				channelIndex++
			}
			channel.Sources = append(channel.Sources, &Source{
				Name: source,
				Url:  url,
			})

			group.Channels = append(group.Channels, channel)
		}
	}
	channelData := &ChannelData{
		Groups: make([]*ChannelGroup, 0, groups.Len()),
	}
	for pair := groups.Oldest(); pair != nil; pair = pair.Next() {
		channelData.Groups = append(channelData.Groups, pair.Value)
		for _, channel := range pair.Value.Channels {
			sort.Sort(sortByName(channel.Sources))
		}
	}
	return channelData, nil
}

func ParseChannelFromM3U8(path string) (*ChannelData, error) {
	m, err := m3u.Parse(path)
	if err != nil {
		return nil, err
	}
	groups := orderedmap.New[string, *ChannelGroup]()
	channelMap := make(map[string]*Channel)
	channelIndex := 1
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
		channel, ok := channelMap[name]
		if !ok {
			channel = &Channel{
				Index: channelIndex,
				Name:  name,
			}
			channelMap[name] = channel
			channelIndex++
		}
		channel.Sources = append(channel.Sources, &Source{
			Name: "默认线路",
			Url:  track.URI,
		})
		group.Channels = append(group.Channels, channel)
	}
	channelData := &ChannelData{
		Groups: make([]*ChannelGroup, 0, groups.Len()),
	}
	for pair := groups.Oldest(); pair != nil; pair = pair.Next() {
		channelData.Groups = append(channelData.Groups, pair.Value)
		for _, channel := range pair.Value.Channels {
			sort.Sort(sortByName(channel.Sources))
		}
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
