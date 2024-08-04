package channel_list

import (
	"time"

	evbus "github.com/asaskevich/EventBus"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/zwh8800/RGTV/component"
	"github.com/zwh8800/RGTV/consts"
	"github.com/zwh8800/RGTV/model"
)

const (
	eventHide          = "ChannelList:Hide"
	eventChannelChange = "ChannelList:ChannelChange"
)

const (
	channelListCloseTimeout = 10 * time.Second
)

type ChannelList struct {
	channelData *model.ChannelData

	selectedGroup   int
	selectedChannel int
	playingGroup    int
	playingChannel  int

	shown      bool
	closeTimer *time.Timer
	eventBus   evbus.Bus
}

func NewChannelList(channelData *model.ChannelData) *ChannelList {
	return &ChannelList{
		channelData: channelData,
		eventBus:    evbus.New(),
	}
}

func (c *ChannelList) HandleEvent(e sdl.Event) {
	if !c.shown {
		return
	}

	switch e.(type) {
	case *sdl.JoyHatEvent:
		event := e.(*sdl.JoyHatEvent)
		if event.Value&sdl.HAT_LEFT != 0 {

		} else if event.Value&sdl.HAT_RIGHT != 0 {

		} else if event.Value&sdl.HAT_UP != 0 {
			c.closeTimer.Reset(channelListCloseTimeout)

		} else if event.Value&sdl.HAT_DOWN != 0 {
			c.closeTimer.Reset(channelListCloseTimeout)

		}
	case *sdl.JoyButtonEvent:
		event := e.(*sdl.JoyButtonEvent)
		if event.State == sdl.PRESSED {
			if event.Button == consts.ButtonA {

			} else if event.Button == consts.ButtonB {
				c.Hide()
			}
		}
	}
}

func (c *ChannelList) Draw(renderer *sdl.Renderer) {
	if !c.shown {
		return
	}
}

func (c *ChannelList) Dispose() {

}

func (c *ChannelList) Show() {
	c.shown = true
	if c.closeTimer != nil {
		c.closeTimer.Stop()
	}
	c.closeTimer = time.AfterFunc(channelListCloseTimeout, func() {
		c.shown = false
	})
}

func (c *ChannelList) Hide() {
	c.shown = false
	c.eventBus.Publish(eventHide, c)
}

func (c *ChannelList) IsShown() bool {
	return c.shown
}

func (c *ChannelList) OnHide(f model.EventHandler) {
	c.eventBus.Subscribe(eventHide, f)
}

func (c *ChannelList) OnChannelChange(f model.EventHandler) {
	c.eventBus.Subscribe(eventChannelChange, f)
}

func (c *ChannelList) GetCurChannel() (*model.ChannelGroup, *model.Channel) {
	if c.selectedGroup >= len(c.channelData.Groups) {
		return nil, nil
	}
	if c.selectedChannel >= len(c.channelData.Groups[c.selectedGroup].Channels) {
		return nil, nil
	}
	return c.channelData.Groups[c.selectedGroup], c.channelData.Groups[c.selectedGroup].Channels[c.selectedChannel]
}

func (c *ChannelList) MoveUp() {
	c.playingChannel++
	if c.playingChannel >= len(c.channelData.Groups[c.playingGroup].Channels) {
		c.playingGroup++
		if c.playingGroup >= len(c.channelData.Groups) {
			c.playingGroup = 0
		}
		c.playingChannel = 0
	}
	c.selectedChannel = c.playingChannel
	c.selectedGroup = c.playingGroup
	c.eventBus.Publish(eventChannelChange, c)
}

func (c *ChannelList) MoveDown() {
	c.playingChannel--
	if c.playingChannel < 0 {
		c.playingGroup--
		if c.playingGroup < 0 {
			c.playingGroup = len(c.channelData.Groups) - 1
		}
		c.playingChannel = len(c.channelData.Groups[c.playingGroup].Channels) - 1
	}
	c.selectedChannel = c.playingChannel
	c.selectedGroup = c.playingGroup
	c.eventBus.Publish(eventChannelChange, c)
}

var _ component.Component = (*ChannelList)(nil)
