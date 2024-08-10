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

const (
	posX      = 0
	posY      = 0
	width     = 200
	height    = 480
	splitPosX = 50
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

func New(channelData *model.ChannelData) *ChannelList {
	return &ChannelList{
		channelData: channelData,
		eventBus:    evbus.New(),
	}
}

func (c *ChannelList) HandleEvent(e sdl.Event) {
	if !c.shown {
		return
	}

	switch event := e.(type) {
	case *sdl.JoyHatEvent:
		if event.Value&sdl.HAT_LEFT != 0 {
			c.hatLeft()
		} else if event.Value&sdl.HAT_RIGHT != 0 {
			c.hatRight()
		} else if event.Value&sdl.HAT_UP != 0 {
			c.hatUp()
		} else if event.Value&sdl.HAT_DOWN != 0 {
			c.hatDown()
		}

	case *sdl.JoyButtonEvent:
		if event.State == sdl.RELEASED {
			if event.Button == consts.ButtonA {
				c.buttonA()
			} else if event.Button == consts.ButtonB {
				c.buttonB()
			}
		}
	case *sdl.KeyboardEvent:
		if event.Type == sdl.KEYUP {
			if event.Keysym.Sym == sdl.K_LEFT {
				c.hatLeft()
			} else if event.Keysym.Sym == sdl.K_RIGHT {
				c.hatRight()
			} else if event.Keysym.Sym == sdl.K_UP {
				c.hatUp()
			} else if event.Keysym.Sym == sdl.K_DOWN {
				c.hatDown()
			} else if event.Keysym.Sym == sdl.K_a {
				c.buttonA()
			} else if event.Keysym.Sym == sdl.K_b {
				c.buttonB()
			}
		}
	}
}

func (c *ChannelList) hatLeft() {
	c.closeTimer.Reset(channelListCloseTimeout)
}

func (c *ChannelList) hatRight() {
	c.closeTimer.Reset(channelListCloseTimeout)
}

func (c *ChannelList) hatUp() {
	c.closeTimer.Reset(channelListCloseTimeout)
}

func (c *ChannelList) hatDown() {
	c.closeTimer.Reset(channelListCloseTimeout)

}

func (c *ChannelList) buttonA() {
	c.closeTimer.Reset(channelListCloseTimeout)

}

func (c *ChannelList) buttonB() {
	c.Hide()
}

func (c *ChannelList) Draw(renderer *sdl.Renderer) {
	if !c.shown {
		return
	}
	c.drawBorder(renderer)
	c.drawSplitLine(renderer)
}

func (c *ChannelList) drawBorder(renderer *sdl.Renderer) {
	renderer.SetDrawColor(0, 0, 0, 128)
	renderer.FillRect(&sdl.Rect{
		X: posX,
		Y: posY,
		W: width,
		H: height,
	})

	renderer.SetDrawColor(255, 255, 255, 255)
	renderer.DrawRect(&sdl.Rect{
		X: posX + 1,
		Y: posY + 1,
		W: width - 2,
		H: height - 2,
	})
}

func (c *ChannelList) drawSplitLine(renderer *sdl.Renderer) {
	renderer.SetDrawColor(255, 255, 255, 255)
	renderer.DrawLine(splitPosX, posY, splitPosX, posY+height)
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
