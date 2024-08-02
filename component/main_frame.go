package component

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/zwh8800/rgbili/conf"
	"github.com/zwh8800/rgbili/consts"
	"github.com/zwh8800/rgbili/model"
)

type MainFrame struct {
	videoBox    *VideoBox
	channelList *ChannelList
	channelInfo *ChannelInfo
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

	videoBox, err := NewVideoBox(channelData.Groups[0].Channels[0].Url)
	if err != nil {
		panic(err)
	}
	channelList := NewChannelList(channelData)
	channelInfo := NewChannelInfo()

	m := &MainFrame{
		videoBox:    videoBox,
		channelList: channelList,
		channelInfo: channelInfo,
	}

	m.channelList.OnChannelChange(m.OnChannelChange)

	return m
}

func (m *MainFrame) HandleEvent(e sdl.Event) {
	switch e.(type) {
	case *sdl.JoyHatEvent:
		event := e.(*sdl.JoyHatEvent)
		if event.Value&sdl.HAT_UP != 0 {
			m.HatUp()
		} else if event.Value&sdl.HAT_DOWN != 0 {
			m.HatDown()
		}
	case *sdl.JoyButtonEvent:
		event := e.(*sdl.JoyButtonEvent)
		if event.State == sdl.RELEASED {
			if event.Button == consts.ButtonA {
				m.buttonA()
			}
		}

	case *sdl.KeyboardEvent:
		event := e.(*sdl.KeyboardEvent)
		if event.Type == sdl.KEYUP {
			if event.Keysym.Sym == sdl.K_UP {
				m.HatUp()
			} else if event.Keysym.Sym == sdl.K_DOWN {
				m.HatDown()
			} else if event.Keysym.Sym == sdl.K_a {
				m.buttonA()
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

func (m *MainFrame) HatUp() {
	if !m.channelList.IsShown() {
		m.channelList.MoveUp()
	}
}

func (m *MainFrame) HatDown() {
	if !m.channelList.IsShown() {
		m.channelList.MoveDown()
	}
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
	m.videoBox, err = NewVideoBox(channel.Url)
	if err != nil {
		panic(err)
	}
}