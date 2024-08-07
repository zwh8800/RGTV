package main_frame

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/zwh8800/RGTV/component"
	channelinfo "github.com/zwh8800/RGTV/component/channel_info"
	channellist "github.com/zwh8800/RGTV/component/channel_list"
	exitmask "github.com/zwh8800/RGTV/component/exit_mask"
	videobox "github.com/zwh8800/RGTV/component/video_box"
	volumebar "github.com/zwh8800/RGTV/component/volume_bar"
	"github.com/zwh8800/RGTV/conf"
	"github.com/zwh8800/RGTV/consts"
	"github.com/zwh8800/RGTV/model"
)

type MainFrame struct {
	videoBox    *videobox.VideoBox
	channelList *channellist.ChannelList
	channelInfo *channelinfo.ChannelInfo
	volumeBar   *volumebar.VolumeBar
	exitMask    *exitmask.ExitMask
}

func New() *MainFrame {
	url := conf.GetConfig().LiveSourceUrl
	channelData, err := model.ParseChannelFromM3U8(url)
	if err != nil {
		channelData, err = model.ParseChannelFromDIYP(url)
		if err != nil {
			panic(err)
		}
	}

	videoBox, err := videobox.New(channelData.Groups[0].Channels[0].Url)
	if err != nil {
		panic(err)
	}
	channelList := channellist.New(channelData)

	channelInfo := channelinfo.New()
	channelInfo.ChannelNumber = 1
	channelInfo.ChannelName = channelData.Groups[0].Channels[0].Name

	volumeBar := volumebar.New()

	exitMask := exitmask.New()

	m := &MainFrame{
		videoBox:    videoBox,
		channelList: channelList,
		channelInfo: channelInfo,
		volumeBar:   volumeBar,
		exitMask:    exitMask,
	}

	m.channelList.OnChannelChange(m.OnChannelChange)
	m.channelInfo.Show()

	return m
}

func (m *MainFrame) HandleEvent(e sdl.Event) {
	captured := false
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
				captured = m.buttonA()
			} else if event.Button == consts.ButtonB {
				captured = m.buttonB()
			} else if event.Button == consts.ButtonX {
				captured = m.buttonX()
			} else if event.Button == consts.ButtonVolumeUp {
				captured = m.buttonVolumeUp()
			} else if event.Button == consts.ButtonVolumeDown {
				captured = m.buttonVolumeDown()
			}
		}

	case *sdl.KeyboardEvent:
		if event.Type == sdl.KEYUP {
			if event.Keysym.Sym == sdl.K_UP {
				captured = m.hatUp()
			} else if event.Keysym.Sym == sdl.K_DOWN {
				captured = m.hatDown()
			} else if event.Keysym.Sym == sdl.K_a {
				captured = m.buttonA()
			} else if event.Keysym.Sym == sdl.K_b {
				captured = m.buttonB()
			} else if event.Keysym.Sym == sdl.K_x {
				captured = m.buttonX()
			} else if event.Keysym.Sym == sdl.K_u {
				captured = m.buttonVolumeUp()
			} else if event.Keysym.Sym == sdl.K_d {
				captured = m.buttonVolumeDown()
			}
		}
	}

	if captured {
		return
	}

	m.videoBox.HandleEvent(e)
	m.channelList.HandleEvent(e)
	m.channelInfo.HandleEvent(e)
	m.exitMask.HandleEvent(e)
}

func (m *MainFrame) buttonA() bool {
	if m.channelList.IsShown() {
		return false
	}
	if m.exitMask.IsShown() {
		return false
	}
	m.channelList.Show()
	return true
}

func (m *MainFrame) buttonB() bool {
	if m.channelList.IsShown() {
		return false
	}
	if m.exitMask.IsShown() {
		return false
	}
	m.exitMask.Show()
	return true
}

func (m *MainFrame) buttonX() bool {
	if m.channelList.IsShown() {
		return false
	}
	m.channelInfo.Show()
	return true
}

func (m *MainFrame) hatUp() bool {
	if m.channelList.IsShown() {
		return false
	}
	if m.exitMask.IsShown() {
		return false
	}
	if conf.GetConfig().RevertSwitchChannel {
		m.channelList.MoveUp()
	} else {
		m.channelList.MoveDown()
	}
	return true
}

func (m *MainFrame) hatDown() bool {
	if m.channelList.IsShown() {
		return false
	}
	if m.exitMask.IsShown() {
		return false
	}
	if conf.GetConfig().RevertSwitchChannel {
		m.channelList.MoveDown()
	} else {
		m.channelList.MoveUp()
	}
	return true
}

func (m *MainFrame) buttonVolumeUp() bool {
	m.volumeBar.VolumeUp()
	m.videoBox.SetVolume(m.volumeBar.GetVolume())
	return true
}

func (m *MainFrame) buttonVolumeDown() bool {
	m.volumeBar.VolumeDown()
	m.videoBox.SetVolume(m.volumeBar.GetVolume())
	return true
}

func (m *MainFrame) Draw(renderer *sdl.Renderer) {
	err := renderer.Clear()
	if err != nil {
		panic(err)
	}
	m.videoBox.Draw(renderer)
	m.channelList.Draw(renderer)
	m.channelInfo.Draw(renderer)
	m.volumeBar.Draw(renderer)
	m.exitMask.Draw(renderer)
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
	m.videoBox, err = videobox.New(channel.Url)
	if err != nil {
		panic(err)
	}
	m.videoBox.SetVolume(m.volumeBar.GetVolume())
	m.channelInfo.ChannelNumber = channel.Index
	m.channelInfo.ChannelName = channel.Name
	m.channelInfo.Show()
}

var _ component.Component = (*MainFrame)(nil)
