package main_frame

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/zwh8800/RGTV/component"
	channelinfo "github.com/zwh8800/RGTV/component/channel_info"
	channellist "github.com/zwh8800/RGTV/component/channel_list"
	channelsource "github.com/zwh8800/RGTV/component/channel_source"
	exitmask "github.com/zwh8800/RGTV/component/exit_mask"
	loadingbar "github.com/zwh8800/RGTV/component/loading_bar"
	videobox "github.com/zwh8800/RGTV/component/video_box"
	volumebar "github.com/zwh8800/RGTV/component/volume_bar"
	"github.com/zwh8800/RGTV/conf"
	"github.com/zwh8800/RGTV/consts"
	"github.com/zwh8800/RGTV/model"
	"github.com/zwh8800/RGTV/util"
)

type MainFrame struct {
	videoBox      *videobox.VideoBox
	channelList   *channellist.ChannelList
	channelInfo   *channelinfo.ChannelInfo
	volumeBar     *volumebar.VolumeBar
	exitMask      *exitmask.ExitMask
	loadingBar    *loadingbar.LoadingBar
	channelSource *channelsource.ChannelSource
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

	videoBox, err := videobox.New(channelData.Groups[0].Channels[0].Sources[0].Url)
	if err != nil {
		panic(err)
	}
	channelList := channellist.New(channelData)

	channelInfo := channelinfo.New()
	channelInfo.ChannelNumber = 1
	channelInfo.ChannelName = channelData.Groups[0].Channels[0].Name
	go func() {
		channelInfo.CurrentProgram = util.GetCurrentProgram(channelInfo.ChannelName)
		channelInfo.NextProgram = util.GetNextProgram(channelInfo.ChannelName)
	}()

	volumeBar := volumebar.New()

	exitMask := exitmask.New()

	loadBar := loadingbar.New()

	channelSource := channelsource.New(channelData.Groups[0].Channels[0])

	m := &MainFrame{
		videoBox:      videoBox,
		channelList:   channelList,
		channelInfo:   channelInfo,
		volumeBar:     volumeBar,
		exitMask:      exitMask,
		loadingBar:    loadBar,
		channelSource: channelSource,
	}

	m.channelList.OnChannelChange(m.onChannelChange)
	m.channelInfo.Show()

	m.channelSource.OnSourceChange(m.onSourceChange)

	return m
}

func (m *MainFrame) HandleEvent(e sdl.Event) {
	captured := false
	switch event := e.(type) {
	case *sdl.JoyHatEvent:
		if event.Value&sdl.HAT_UP != 0 {
			captured = m.hatUp()
		} else if event.Value&sdl.HAT_DOWN != 0 {
			captured = m.hatDown()
		} else if event.Value&sdl.HAT_LEFT != 0 {
			captured = m.hatLeft()
		} else if event.Value&sdl.HAT_RIGHT != 0 {
			captured = m.hatRight()
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
			} else if event.Keysym.Sym == sdl.K_LEFT {
				captured = m.hatLeft()
			} else if event.Keysym.Sym == sdl.K_RIGHT {
				captured = m.hatRight()
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
	m.volumeBar.HandleEvent(e)
	m.exitMask.HandleEvent(e)
	m.loadingBar.HandleEvent(e)
	m.channelSource.HandleEvent(e)
}

func (m *MainFrame) buttonA() bool {
	if m.channelList.IsShown() {
		return false
	}
	if m.exitMask.IsShown() {
		return false
	}
	m.channelList.Show()
	m.channelInfo.Hide()
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
	m.channelInfo.Hide()
	return true
}

func (m *MainFrame) buttonX() bool {
	if m.channelList.IsShown() {
		return false
	}
	if m.exitMask.IsShown() {
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

func (m *MainFrame) hatLeft() bool {
	if m.channelList.IsShown() {
		return false
	}
	if m.exitMask.IsShown() {
		return false
	}
	m.channelSource.PrevSource()

	return true
}

func (m *MainFrame) hatRight() bool {
	if m.channelList.IsShown() {
		return false
	}
	if m.exitMask.IsShown() {
		return false
	}
	m.channelSource.NextSource()

	return true
}

func (m *MainFrame) buttonVolumeUp() bool {
	if m.exitMask.IsShown() {
		return false
	}
	m.volumeBar.VolumeUp()
	m.videoBox.SetVolume(m.volumeBar.GetVolume())
	return true
}

func (m *MainFrame) buttonVolumeDown() bool {
	if m.exitMask.IsShown() {
		return false
	}
	m.volumeBar.VolumeDown()
	m.videoBox.SetVolume(m.volumeBar.GetVolume())
	return true
}

func (m *MainFrame) Draw(renderer *sdl.Renderer) {
	renderer.Clear()
	m.videoBox.Draw(renderer)
	m.channelList.Draw(renderer)
	m.channelInfo.Draw(renderer)
	m.volumeBar.Draw(renderer)
	m.exitMask.Draw(renderer)
	m.loadingBar.Draw(renderer)
	m.channelSource.Draw(renderer)
	renderer.Present()
}

func (m *MainFrame) Dispose() {
	m.videoBox.Dispose()
	m.channelList.Dispose()
	m.channelInfo.Dispose()
}

func (m *MainFrame) onChannelChange(_ any) {
	_, channel := m.channelList.GetCurChannel()
	m.channelSource.SetChannel(channel)

	m.videoBox.Dispose()
	var err error
	m.videoBox, err = videobox.New(m.channelSource.GetSource().Url)
	if err != nil {
		panic(err)
	}
	m.videoBox.SetVolume(m.volumeBar.GetVolume())
	m.videoBox.OnVideoLag(m.OnVideoLag)
	m.channelInfo.ChannelNumber = channel.Index
	m.channelInfo.ChannelName = channel.Name
	go func() {
		m.channelInfo.CurrentProgram = util.GetCurrentProgram(channel.Name)
		m.channelInfo.NextProgram = util.GetNextProgram(channel.Name)
	}()
	m.channelInfo.Show()
}

func (m *MainFrame) onSourceChange(_ any) {
	m.videoBox.Dispose()
	var err error
	m.videoBox, err = videobox.New(m.channelSource.GetSource().Url)
	if err != nil {
		panic(err)
	}
	m.videoBox.SetVolume(m.volumeBar.GetVolume())
	m.videoBox.OnVideoLag(m.OnVideoLag)
}

func (m *MainFrame) OnVideoLag(_ any) {
	m.loadingBar.Show()
}

var _ component.Component = (*MainFrame)(nil)
