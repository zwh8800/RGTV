package volume_bar

import (
	"fmt"
	"image"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/zwh8800/RGTV/component"
	"github.com/zwh8800/RGTV/embeddata"
	"github.com/zwh8800/RGTV/text"
	"github.com/zwh8800/RGTV/util"
)

const (
	closeTimeout = 1 * time.Second
)

type VolumeBar struct {
	volume     int
	shown      bool
	closeTimer *time.Timer
}

func New() *VolumeBar {
	return &VolumeBar{
		volume: 5,
	}
}

func (v *VolumeBar) HandleEvent(e sdl.Event) {}
func (v *VolumeBar) Dispose()                {}

func (v *VolumeBar) Draw(renderer *sdl.Renderer) {
	if !v.shown {
		return
	}

	textDrawer, err := text.GetDrawerFromData(embeddata.FontName, embeddata.FontData)
	if err != nil {
		panic(err)
	}

	txt := fmt.Sprintf("♫%02d", v.volume*10)
	img, err := textDrawer.Draw(txt, 46, image.White)
	if err != nil {
		panic(err)
	}

	x := (640 - img.Bounds().Dx()) / 2
	y := (480 - img.Bounds().Dy()) / 2

	renderer.SetDrawColor(0, 0, 0, 128)
	renderer.FillRect(&sdl.Rect{
		X: int32(x - 20),
		Y: int32(y - 20),
		W: int32(img.Bounds().Dx() + 40),
		H: int32(img.Bounds().Dy() + 40),
	})

	renderer.SetDrawColor(255, 255, 255, 255)
	renderer.DrawRect(&sdl.Rect{
		X: int32(x-20) + 1,
		Y: int32(y-20) + 1,
		W: int32(img.Bounds().Dx()+40) - 2,
		H: int32(img.Bounds().Dy()+40) - 2,
	})

	util.DrawGoImage(renderer, img,
		image.Rect(
			x,
			y,
			x+img.Bounds().Dx(),
			y+img.Bounds().Dy(),
		))
}

func (v *VolumeBar) GetVolume() int {
	return v.volume
}

func (v *VolumeBar) VolumeUp() {
	if v.volume >= 10 {
		v.volume = 10
		return
	}
	v.volume++
	v.Show()
}

func (v *VolumeBar) VolumeDown() {
	if v.volume <= 0 {
		v.volume = 0
		return
	}
	v.volume--
	v.Show()
}

func (v *VolumeBar) Show() {
	v.shown = true
	if v.closeTimer != nil {
		v.closeTimer.Stop()
	}
	v.closeTimer = time.AfterFunc(closeTimeout, func() {
		v.shown = false
	})
}

var _ component.Component = (*VolumeBar)(nil)
