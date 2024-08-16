package exit_mask

import (
	"bytes"
	"image"
	"image/png"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/zwh8800/RGTV/component"
	"github.com/zwh8800/RGTV/consts"
	"github.com/zwh8800/RGTV/embeddata"
	"github.com/zwh8800/RGTV/text"
	"github.com/zwh8800/RGTV/util"
)

const (
	closeTimeout = 10 * time.Second

	logoSize = 144
)

type ExitMask struct {
	shown      bool
	closeTimer *time.Timer

	logoImg *image.RGBA
}

func New() *ExitMask {
	img, err := png.Decode(bytes.NewReader(embeddata.LogoData))
	if err != nil {
		panic(err)
	}
	logoImg := util.ImageToRGBA(img)

	return &ExitMask{
		logoImg: logoImg,
	}
}

func (c *ExitMask) HandleEvent(e sdl.Event) {
	if !c.shown {
		return
	}

	switch event := e.(type) {
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
			if event.Keysym.Sym == sdl.K_a {
				c.buttonA()
			} else if event.Keysym.Sym == sdl.K_b {
				c.buttonB()
			}
		}
	}
}

func (c *ExitMask) Draw(renderer *sdl.Renderer) {
	if !c.shown {
		return
	}
	renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
	renderer.SetDrawColor(0, 0, 0, 180)
	renderer.FillRect(&sdl.Rect{
		X: 0,
		Y: 0,
		W: 640,
		H: 480,
	})
	renderer.SetDrawColor(255, 255, 255, 255)
	renderer.DrawRect(&sdl.Rect{
		X: 0 + 1,
		Y: 0 + 1,
		W: 640 - 2,
		H: 480 - 2,
	})

	textDrawer, err := text.GetDrawerFromData(embeddata.FontName, embeddata.FontData)
	if err != nil {
		panic(err)
	}

	c.drawLogo(renderer)
	c.drawExitMsg(renderer, textDrawer)
	c.drawExitButton(renderer, textDrawer)
}

func (c *ExitMask) drawLogo(render *sdl.Renderer) {
	util.DrawGoImage(render, c.logoImg,
		image.Rect(
			(640-logoSize)/2,
			(480-logoSize)/2,
			(640-logoSize)/2+logoSize,
			(480-logoSize)/2+logoSize,
		))
}

func (c *ExitMask) drawExitMsg(renderer *sdl.Renderer, textDrawer *text.Drawer) {
	img, err := textDrawer.Draw("是否要退出RGTV", 32, image.White)
	if err != nil {
		panic(err)
	}

	x := (640 - img.Bounds().Dx()) / 2
	y := (480-img.Bounds().Dy())/2 + 120

	util.DrawGoImage(renderer, img,
		image.Rect(
			x,
			y,
			x+img.Bounds().Dx(),
			y+img.Bounds().Dy(),
		))
}

func (c *ExitMask) drawExitButton(renderer *sdl.Renderer, textDrawer *text.Drawer) {
	img, err := textDrawer.Draw("A：确认  B：取消", 24, image.White)
	if err != nil {
		panic(err)
	}

	x := (640 - img.Bounds().Dx()) / 2
	y := (480-img.Bounds().Dy())/2 + 160

	util.DrawGoImage(renderer, img,
		image.Rect(
			x,
			y,
			x+img.Bounds().Dx(),
			y+img.Bounds().Dy(),
		))
}

func (c *ExitMask) Dispose() {}

func (c *ExitMask) Show() {
	c.shown = true
	if c.closeTimer != nil {
		c.closeTimer.Stop()
	}
	c.closeTimer = time.AfterFunc(closeTimeout, func() {
		c.shown = false
	})
}

func (c *ExitMask) Hide() {
	c.shown = false
}

func (c *ExitMask) IsShown() bool {
	return c.shown
}

func (c *ExitMask) buttonA() {
	sdl.PushEvent(&sdl.QuitEvent{
		Type:      sdl.QUIT,
		Timestamp: sdl.GetTicks(),
	})
}

func (c *ExitMask) buttonB() {
	c.Hide()
}

var _ component.Component = (*ExitMask)(nil)
