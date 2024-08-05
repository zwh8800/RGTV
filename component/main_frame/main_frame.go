package main_frame

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/zwh8800/RGTV/component"
	"github.com/zwh8800/RGTV/component/channel_info"
	"github.com/zwh8800/RGTV/component/channel_list"
	"github.com/zwh8800/RGTV/component/video_box"
	"github.com/zwh8800/RGTV/conf"
	"github.com/zwh8800/RGTV/consts"
	"github.com/zwh8800/RGTV/model"
)

type MainFrame struct {
	videoBox    *video_box.VideoBox
	channelList *channel_list.ChannelList
	channelInfo *channel_info.ChannelInfo
}

func NewMainFrame() *MainFrame {
	url := conf.GetConfig().LiveSourceUrl
	channelData, err := model.ParseChannelFromM3U8(url)
	if err != nil {
		channelData, err = model.ParseChannelFromDIYP(url)
		if err != nil {
			panic(err)
		}
	}

	videoBox, err := video_box.NewVideoBox(channelData.Groups[0].Channels[0].Url)
	if err != nil {
		panic(err)
	}
	channelList := channel_list.NewChannelList(channelData)
	channelInfo := channel_info.NewChannelInfo()
	channelInfo.ChannelName = channelData.Groups[0].Channels[0].Name

	m := &MainFrame{
		videoBox:    videoBox,
		channelList: channelList,
		channelInfo: channelInfo,
	}

	m.channelList.OnChannelChange(m.OnChannelChange)

	return m
}

func (m *MainFrame) HandleEvent(e sdl.Event) {
	switch event := e.(type) {
	case *sdl.JoyHatEvent:
		if event.Value&sdl.HAT_UP != 0 {
			m.hatUp()
		} else if event.Value&sdl.HAT_DOWN != 0 {
			m.hatDown()
		}
	case *sdl.JoyButtonEvent:
		if event.State == sdl.RELEASED {
			if event.Button == consts.ButtonA {
				m.buttonA()
			} else if event.Button == consts.ButtonX {
				m.buttonX()
			} else if event.Button == consts.ButtonVolumeUp {
				m.buttonVolumeUp()
			} else if event.Button == consts.ButtonVolumeDown {
				m.buttonVolumeDown()
			}
		}

	case *sdl.KeyboardEvent:
		if event.Type == sdl.KEYUP {
			if event.Keysym.Sym == sdl.K_UP {
				m.hatUp()
			} else if event.Keysym.Sym == sdl.K_DOWN {
				m.hatDown()
			} else if event.Keysym.Sym == sdl.K_a {
				m.buttonA()
			} else if event.Keysym.Sym == sdl.K_x {
				m.buttonX()
			} else if event.Keysym.Sym == sdl.K_u {
				m.buttonVolumeUp()
			} else if event.Keysym.Sym == sdl.K_d {
				m.buttonVolumeDown()
			}
		}
	}

	m.videoBox.HandleEvent(e)
	m.channelList.HandleEvent(e)
	m.channelInfo.HandleEvent(e)
}

func (m *MainFrame) buttonA() {
	if !m.channelList.IsShown() {
		m.channelList.Show()
	}
}

func (m *MainFrame) buttonX() {
	m.channelInfo.Show()
}

func (m *MainFrame) hatUp() {
	if !m.channelList.IsShown() {
		if conf.GetConfig().RevertSwitchChannel {
			m.channelList.MoveUp()
		} else {
			m.channelList.MoveDown()
		}
	}
}

func (m *MainFrame) hatDown() {
	if !m.channelList.IsShown() {
		if conf.GetConfig().RevertSwitchChannel {
			m.channelList.MoveDown()
		} else {
			m.channelList.MoveUp()
		}
	}
}

func (m *MainFrame) buttonVolumeUp() {
	m.videoBox.VolumeUp()
}
func (m *MainFrame) buttonVolumeDown() {
	m.videoBox.VolumeDown()
}

func (m *MainFrame) Draw(renderer *sdl.Renderer) {
	err := renderer.Clear()
	if err != nil {
		panic(err)
	}
	m.videoBox.Draw(renderer)
	m.channelList.Draw(renderer)
	m.channelInfo.Draw(renderer)
	renderer.Present()
}

func (m *MainFrame) Dispose() {
	m.videoBox.Dispose()
	m.channelList.Dispose()
	m.channelInfo.Dispose()
}

func (m *MainFrame) OnChannelChange(_ any) {
	_, channel := m.channelList.GetCurChannel()
	m.videoBox.Dispose()
	var err error
	m.videoBox, err = video_box.NewVideoBox(channel.Url)
	if err != nil {
		panic(err)
	}
	m.channelInfo.ChannelName = channel.Name
	m.channelInfo.Show()
}

var _ component.Component = (*MainFrame)(nil)
