package conf_qrcode

import (
	"fmt"
	"image"
	"image/draw"

	"github.com/skip2/go-qrcode"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/zwh8800/RGTV/component"
	"github.com/zwh8800/RGTV/conf"
	"github.com/zwh8800/RGTV/consts"
	"github.com/zwh8800/RGTV/util"
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
		img: imageToRGBA(img),
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

// imageToRGBA 将image.Image转换为image.RGBA
func imageToRGBA(img image.Image) *image.RGBA {
	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, img, bounds.Min, draw.Src)
	return rgba
}

var _ component.Component = (*ConfQrcode)(nil)
