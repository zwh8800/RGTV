package channel_list

import (
	"image"
	"image/color"
	"time"

	evbus "github.com/asaskevich/EventBus"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/zwh8800/RGTV/component"
	"github.com/zwh8800/RGTV/consts"
	"github.com/zwh8800/RGTV/embeddata"
	"github.com/zwh8800/RGTV/model"
	"github.com/zwh8800/RGTV/text"
	"github.com/zwh8800/RGTV/util"
)

const (
	eventHide          = "ChannelList:Hide"
	eventChannelChange = "ChannelList:ChannelChange"
)

const (
	closeTimeout = 10 * time.Second
)

const (
	posX         = 0
	posY         = 0
	width        = 240
	height       = 480
	genreWidth   = 70
	channelWidth = width - genreWidth
	genreCount   = 9
	channelCount = 7
)

var (
	colorActive = color.RGBA{66, 144, 245, 204}
)

type ChannelList struct {
	channelData *model.ChannelData

	selectedGroup   int
	selectedChannel int
	playingGroup    int
	playingChannel  int
	focusOnGenre    bool

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
	c.closeTimer.Reset(closeTimeout)
	c.focusOnGenre = true
}

func (c *ChannelList) hatRight() {
	c.closeTimer.Reset(closeTimeout)
	c.focusOnGenre = false
}

func (c *ChannelList) hatUp() {
	c.closeTimer.Reset(closeTimeout)
	if c.focusOnGenre {
		c.selectedGroup--
		c.selectedGroup %= len(c.channelData.Groups)
	} else {
		c.selectedChannel--
		if c.selectedChannel < 0 {
			c.selectedGroup--
			if c.selectedGroup < 0 {
				c.selectedGroup = len(c.channelData.Groups) - 1
			}
			c.selectedChannel = len(c.channelData.Groups[c.selectedGroup].Channels) - 1
		}
	}
}

func (c *ChannelList) hatDown() {
	c.closeTimer.Reset(closeTimeout)
	if c.focusOnGenre {
		c.selectedGroup++
		c.selectedGroup %= len(c.channelData.Groups)
	} else {
		c.selectedChannel++
		if c.selectedChannel >= len(c.channelData.Groups[c.selectedGroup].Channels) {
			c.selectedGroup++
			if c.selectedGroup >= len(c.channelData.Groups) {
				c.selectedGroup = 0
			}
			c.selectedChannel = 0
		}
	}
}

func (c *ChannelList) buttonA() {
	c.closeTimer.Reset(closeTimeout)
	if c.focusOnGenre {
		c.focusOnGenre = false
	} else {
		c.playingChannel = c.selectedChannel
		c.playingGroup = c.selectedGroup
		c.eventBus.Publish(eventChannelChange, c)
		c.Hide()
	}
}

func (c *ChannelList) buttonB() {
	c.Hide()
}

func (c *ChannelList) Draw(renderer *sdl.Renderer) {
	if !c.shown {
		return
	}
	textDrawer, err := text.GetDrawerFromData(embeddata.FontName, embeddata.FontData)
	if err != nil {
		panic(err)
	}
	c.drawBorder(renderer)
	c.drawSplitLine(renderer)
	c.drawGenre(renderer, textDrawer)
	c.drawChannel(renderer, textDrawer)
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
	renderer.DrawLine(genreWidth, posY+1, genreWidth, posY+height-2)
}

func (c *ChannelList) drawGenre(renderer *sdl.Renderer, textDrawer *text.Drawer) {
	genreHeight := height / genreCount
	drawCount := min(len(c.channelData.Groups), genreCount)
	startGroupIdx := min(max(0, c.selectedGroup-genreCount/2), max(0, len(c.channelData.Groups)-genreCount))
	for i := 0; i < drawCount; i++ {
		idx := startGroupIdx + i

		iPosX := posX
		iPosY := posY + i*genreHeight
		renderer.SetDrawColor(255, 255, 255, 255)
		renderer.DrawRect(&sdl.Rect{
			X: int32(iPosX + 4),
			Y: int32(iPosY + 4),
			W: int32(genreWidth - 8),
			H: int32(genreHeight - 8),
		})

		if idx == c.selectedGroup && c.focusOnGenre {
			renderer.SetDrawColor(colorActive.R, colorActive.G, colorActive.B, colorActive.A)
			renderer.FillRect(&sdl.Rect{
				X: int32(iPosX + 6),
				Y: int32(iPosY + 6),
				W: int32(genreWidth - 12),
				H: int32(genreHeight - 12),
			})
		}

		group := c.channelData.Groups[idx]

		img, err := textDrawer.Draw(group.Name, 12, image.White)
		if err != nil {
			panic(err)
		}
		x := iPosX + (genreWidth-img.Bounds().Dx())/2
		y := iPosY + (genreHeight-img.Bounds().Dy())/2
		util.DrawGoImage(renderer, img,
			image.Rect(
				x,
				y,
				x+img.Bounds().Dx(),
				y+img.Bounds().Dy(),
			))
	}
}

func (c *ChannelList) drawChannel(renderer *sdl.Renderer, textDrawer *text.Drawer) {
	channelHeight := height / channelCount
	drawCount := min(len(c.channelData.Groups[c.selectedGroup].Channels), channelCount)
	startGroupIdx := min(max(0, c.selectedChannel-channelCount/2), max(0, len(c.channelData.Groups[c.selectedGroup].Channels)-channelCount))
	for i := 0; i < drawCount; i++ {
		idx := startGroupIdx + i

		iPosX := genreWidth
		iPosY := posY + i*channelHeight

		renderer.SetDrawColor(255, 255, 255, 255)
		renderer.DrawRect(&sdl.Rect{
			X: int32(iPosX + 4),
			Y: int32(iPosY + 4),
			W: int32(channelWidth - 8),
			H: int32(channelHeight - 8),
		})

		if idx == c.selectedChannel && !c.focusOnGenre {
			renderer.SetDrawColor(colorActive.R, colorActive.G, colorActive.B, colorActive.A)
			renderer.FillRect(&sdl.Rect{
				X: int32(iPosX + 6),
				Y: int32(iPosY + 6),
				W: int32(channelWidth - 12),
				H: int32(channelHeight - 12),
			})
		}

		channel := c.channelData.Groups[c.selectedGroup].Channels[idx]

		img, err := textDrawer.Draw(channel.Name, 20, image.White)
		if err != nil {
			panic(err)
		}
		x := iPosX + (channelWidth-img.Bounds().Dx())/2
		y := iPosY + (channelHeight-img.Bounds().Dy())/2
		util.DrawGoImage(renderer, img,
			image.Rect(
				x,
				y,
				x+img.Bounds().Dx(),
				y+img.Bounds().Dy(),
			))
	}
}

func (c *ChannelList) Dispose() {

}

func (c *ChannelList) Show() {
	c.shown = true
	if c.closeTimer != nil {
		c.closeTimer.Stop()
	}
	c.closeTimer = time.AfterFunc(closeTimeout, func() {
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
