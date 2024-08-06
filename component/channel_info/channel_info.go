package channel_info

import (
	"image"
	"strconv"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/zwh8800/RGTV/component"
	"github.com/zwh8800/RGTV/embeddata"
	"github.com/zwh8800/RGTV/text"
	"github.com/zwh8800/RGTV/util"
)

const (
	posX   = 40
	posY   = 320
	width  = 560
	height = 120

	midY = posY + height/2

	channelNameWidth = 200
	epgWidth         = width - channelNameWidth

	closeTimeout = 3 * time.Second
)

type ChannelInfo struct {
	ChannelNumber  int
	ChannelName    string
	CurrentProgram string
	NextProgram    string

	shown      bool
	closeTimer *time.Timer
}

func New() *ChannelInfo {
	return &ChannelInfo{}
}

func (c *ChannelInfo) HandleEvent(e sdl.Event) {

}

func (c *ChannelInfo) Draw(renderer *sdl.Renderer) {
	if !c.shown {
		return
	}
	renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)

	c.drawBorder(renderer)
	c.drawSplitLine(renderer)

	textDrawer, err := text.GetDrawerFromData(embeddata.FontName, embeddata.FontData)
	if err != nil {
		panic(err)
	}
	c.drawChannelName(renderer, textDrawer)
	c.drawChannelNumber(renderer, textDrawer)
	c.drawCurProgram(renderer, textDrawer)
	c.drawNextProgram(renderer, textDrawer)
	c.drawHelp(renderer, textDrawer)
}

func (c *ChannelInfo) drawChannelNumber(renderer *sdl.Renderer, textDrawer *text.Drawer) {
	img, err := textDrawer.Draw(strconv.Itoa(c.ChannelNumber), 46, image.White)
	if err != nil {
		panic(err)
	}

	x := posX + (channelNameWidth-img.Bounds().Dx())/2
	y := posY + (height/2-img.Bounds().Dy())/2 + 10

	err = util.DrawGoImage(renderer, img,
		image.Rect(
			x,
			y,
			x+img.Bounds().Dx(),
			y+img.Bounds().Dy(),
		))
	if err != nil {
		panic(err)
	}
}

func (c *ChannelInfo) drawChannelName(renderer *sdl.Renderer, textDrawer *text.Drawer) {
	img, err := textDrawer.Draw(c.ChannelName, 32, image.White)
	if err != nil {
		panic(err)
	}

	x := posX + (channelNameWidth-img.Bounds().Dx())/2
	y := midY + (height/2-img.Bounds().Dy())/2 - 10

	err = util.DrawGoImage(renderer, img,
		image.Rect(
			x,
			y,
			x+img.Bounds().Dx(),
			y+img.Bounds().Dy(),
		))
	if err != nil {
		panic(err)
	}
}

func (c *ChannelInfo) drawCurProgram(renderer *sdl.Renderer, textDrawer *text.Drawer) {
	txt := "当前播放：" + cutProgram(c.CurrentProgram)

	img, err := textDrawer.Draw(txt, 20, image.White)
	if err != nil {
		panic(err)
	}

	x := posX + channelNameWidth + 20
	y := posY + (height/2-img.Bounds().Dy())/2 + 10

	err = util.DrawGoImage(renderer, img,
		image.Rect(
			x,
			y,
			x+img.Bounds().Dx(),
			y+img.Bounds().Dy(),
		))
	if err != nil {
		panic(err)
	}
}

func (c *ChannelInfo) drawNextProgram(renderer *sdl.Renderer, textDrawer *text.Drawer) {
	txt := "当前播放：" + cutProgram(c.NextProgram)

	img, err := textDrawer.Draw(txt, 20, image.White)
	if err != nil {
		panic(err)
	}

	x := posX + channelNameWidth + 20
	y := midY + (height/2-img.Bounds().Dy())/2 - 10

	err = util.DrawGoImage(renderer, img,
		image.Rect(
			x,
			y,
			x+img.Bounds().Dx(),
			y+img.Bounds().Dy(),
		))
	if err != nil {
		panic(err)
	}
}

func (c *ChannelInfo) drawHelp(renderer *sdl.Renderer, textDrawer *text.Drawer) {
	txt := "▲▼：换台  ◁▷：换源  A：列表  B：退出  X：屏显"

	img, err := textDrawer.Draw(txt, 12, image.White)
	if err != nil {
		panic(err)
	}

	x := posX + channelNameWidth + 20
	y := posY + height - 20

	err = util.DrawGoImage(renderer, img,
		image.Rect(
			x,
			y,
			x+img.Bounds().Dx(),
			y+img.Bounds().Dy(),
		))
	if err != nil {
		panic(err)
	}
}

func (c *ChannelInfo) drawSplitLine(renderer *sdl.Renderer) {
	renderer.SetDrawColor(255, 255, 255, 255)

	// 竖线
	renderer.DrawLine(
		posX+channelNameWidth,
		posY+10,
		posX+channelNameWidth,
		posY+height-10,
	)

	// 横线
	renderer.DrawLine(
		posX+channelNameWidth+20,
		midY,
		posX+width-20,
		midY,
	)
}

func (c *ChannelInfo) drawBorder(renderer *sdl.Renderer) {
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

func (c *ChannelInfo) Dispose() {

}

func (c *ChannelInfo) Show() {
	c.shown = true
	if c.closeTimer != nil {
		c.closeTimer.Stop()
	}
	c.closeTimer = time.AfterFunc(closeTimeout, func() {
		c.shown = false
	})
}

func (c *ChannelInfo) Hide() {
	c.shown = false
}

func (c *ChannelInfo) IsShown() bool {
	return c.shown
}

func cutProgram(program string) string {
	if program == "" {
		program = "精彩节目"
	} else if len([]rune(program)) > 10 {
		r := []rune(program)
		program = string(r[:10]) + "……"
	}
	return program
}

var _ component.Component = (*ChannelInfo)(nil)
