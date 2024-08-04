package component

import (
	"image"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/zwh8800/RGTV/embeddata"
	"github.com/zwh8800/RGTV/text"
	"github.com/zwh8800/RGTV/util"
)

const (
	width  = 560
	height = 120

	channelInfoCloseList = 3 * time.Second
)

type ChannelInfo struct {
	ChannelNumber  int
	ChannelName    string
	CurrentProgram string
	NextProgram    string

	shown      bool
	closeTimer *time.Timer
}

func NewChannelInfo() *ChannelInfo {
	return &ChannelInfo{}
}

func (c *ChannelInfo) HandleEvent(e sdl.Event) {

}

func (c *ChannelInfo) Draw(renderer *sdl.Renderer) {
	if !c.shown {
		return
	}
	textDrawer, err := text.GetDrawerFromData(embeddata.FontName, embeddata.FontData)
	if err != nil {
		panic(err)
	}
	c.drawChannelName(renderer, textDrawer)

}

func (c *ChannelInfo) drawChannelName(renderer *sdl.Renderer, textDrawer *text.Drawer) {
	channelNameImg, err := textDrawer.Draw(c.ChannelName, 64, image.White)
	if err != nil {
		panic(err)
	}

	err = util.DrawGoImage(renderer, channelNameImg,
		image.Rect(200, 200, 200+channelNameImg.Bounds().Dx(), 200+channelNameImg.Bounds().Dy()))
	if err != nil {
		panic(err)
	}
}

func (c *ChannelInfo) Dispose() {

}

func (c *ChannelInfo) Show() {
	c.shown = true
	if c.closeTimer != nil {
		c.closeTimer.Stop()
	}
	c.closeTimer = time.AfterFunc(channelInfoCloseList, func() {
		c.shown = false
	})
}

func (c *ChannelInfo) Hide() {
	c.shown = false
}

func (c *ChannelInfo) IsShown() bool {
	return c.shown
}

var _ Component = (*ChannelInfo)(nil)
