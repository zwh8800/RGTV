package conf_qrcode

import (
	"fmt"
	"github.com/skip2/go-qrcode"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/zwh8800/RGTV/component"
	"github.com/zwh8800/RGTV/conf"
	"github.com/zwh8800/RGTV/consts"
	"github.com/zwh8800/RGTV/embeddata"
	"github.com/zwh8800/RGTV/text"
	"github.com/zwh8800/RGTV/util"
	"image"
)

const (
	size = 300
)

type ConfQrcode struct {
	img   *image.RGBA
	shown bool
}

func New() *ConfQrcode {
	ip := util.GetIP()

	q, err := qrcode.New(fmt.Sprintf("http://%s:8080/", ip), qrcode.Medium)
	if err != nil {
		fmt.Println("qrcode.New err:", err.Error())
		return &ConfQrcode{}
	}
	img := q.Image(size)

	return &ConfQrcode{
		img: util.ImageToRGBA(img),
	}
}

func (c *ConfQrcode) HandleEvent(e sdl.Event) {
	if !c.shown {
		return
	}

	switch event := e.(type) {
	case *sdl.JoyButtonEvent:
		if event.State == sdl.RELEASED {
			if event.Button == consts.ButtonY {
				c.buttonY()
			}
		}

	case *sdl.KeyboardEvent:
		if event.Type == sdl.KEYUP {
			if event.Keysym.Sym == sdl.K_y {
				c.buttonY()
			}
		}
	}
}

func (c *ConfQrcode) buttonY() {
	c.Hide()
}

func (c *ConfQrcode) Draw(renderer *sdl.Renderer) {
	if !c.shown {
		return
	}

	if c.img == nil {
		return
	}
	x := (640 - size) / 2
	y := (480 - size) / 2
	util.DrawGoImage(renderer, c.img,
		image.Rect(
			x,
			y,
			x+size,
			y+size,
		))

	textDrawer, err := text.GetDrawerFromData(embeddata.FontName, embeddata.FontData)
	if err != nil {
		panic(err)
	}

	img, err := textDrawer.Draw("请使用手机扫码进行配置", 32, image.White)
	if err != nil {
		panic(err)
	}

	x = (640 - img.Bounds().Dx()) / 2
	y = 400

	util.DrawGoImage(renderer, img,
		image.Rect(
			x,
			y,
			x+img.Bounds().Dx(),
			y+img.Bounds().Dy(),
		))
}

func (c *ConfQrcode) Dispose() {}

func (c *ConfQrcode) Show() {
	if c.img == nil {
		return
	}
	c.shown = true
	conf.StartNetConfServer()
}

func (c *ConfQrcode) Hide() {
	c.shown = false
	conf.StopNetConfServer()
}

func (c *ConfQrcode) IsShown() bool {
	return c.shown
}

var _ component.Component = (*ConfQrcode)(nil)
