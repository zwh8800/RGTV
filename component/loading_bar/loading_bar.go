package loading_bar

import (
	"image"
	"image/color"
	"math"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/zwh8800/RGTV/component"
	"github.com/zwh8800/RGTV/embeddata"
	"github.com/zwh8800/RGTV/text"
	"github.com/zwh8800/RGTV/util"
)

const (
	closeTimeout = 200 * time.Millisecond
	interval     = 60
	half         = interval / 2
	lineLength   = 200
	linePosY     = 250
	textPosY     = 200
	textMoveY    = 40
)

var (
	charList = []string{"L", "O", "A", "D", "I", "N", "G"}
)

type LoadingBar struct {
	shown      bool
	closeTimer *time.Timer

	tick int
}

func New() *LoadingBar {
	return &LoadingBar{}
}

func (c *LoadingBar) HandleEvent(e sdl.Event) {}
func (c *LoadingBar) Dispose()                {}

func (c *LoadingBar) Draw(renderer *sdl.Renderer) {
	if !c.shown {
		return
	}

	c.drawText(renderer)
	c.drawLine(renderer)

	c.tick++
	c.tick %= interval
}

func (c *LoadingBar) drawText(renderer *sdl.Renderer) {
	textDrawer, err := text.GetDrawerFromData(embeddata.FontName, embeddata.FontData)
	if err != nil {
		panic(err)
	}

	txtLen := len(charList)
	tickPerChar := half / float64(txtLen)
	for i, t := range charList {
		y := 0.0
		alpha := uint8(0)

		if c.tick < half {
			if float64(c.tick) <= float64(i)*tickPerChar { // 起始位置
				y = textPosY + textMoveY
				alpha = 0
			} else if float64(c.tick) >= float64(i+1)*tickPerChar { // 终点位置
				y = textPosY
				alpha = 255
			} else { // 活动
				factor := util.EaseOut((float64(c.tick) - float64(i)*tickPerChar) / tickPerChar)
				y = math.Round(textPosY + textMoveY -
					float64(textMoveY)*factor)
				alpha = uint8(math.Round(255 * factor))
			}
		} else {
			if float64(c.tick) <= half+float64(i)*tickPerChar { // 起始位置
				y = textPosY
				alpha = 255
			} else if float64(c.tick) >= half+float64(i+1)*tickPerChar { // 终点位置
				y = textPosY - textMoveY
				alpha = 0
			} else { // 活动
				factor := util.EaseOut((float64(c.tick-half) - float64(i)*tickPerChar) / tickPerChar)
				y = math.Round(textPosY -
					float64(textMoveY)*factor)
				alpha = uint8(math.Round(255 - 255*factor))
			}
		}

		img, err := textDrawer.Draw(t, 64, image.NewUniform(color.RGBA{
			R: 255,
			G: 255,
			B: 255,
			A: alpha,
		}))
		if err != nil {
			panic(err)
		}

		x := (640-txtLen*img.Bounds().Dx())/2 + i*img.Bounds().Dx()

		util.DrawGoImage(renderer, img,
			image.Rect(
				x,
				int(y),
				x+img.Bounds().Dx(),
				int(y)+img.Bounds().Dy(),
			))
	}

}

func (c *LoadingBar) drawLine(renderer *sdl.Renderer) {
	var x1, x2 int
	if c.tick < half {
		width := math.Round(lineLength * util.EaseOut(float64(c.tick)/half))
		x1 = (640 - lineLength) / 2
		x2 = (640-lineLength)/2 + int(width)
	} else {
		width := math.Round(lineLength * util.EaseOut(float64(c.tick-half)/half))
		x1 = (640-lineLength)/2 + int(width)
		x2 = (640-lineLength)/2 + lineLength
	}

	renderer.DrawLine(int32(x1), linePosY, int32(x2), linePosY)
}

func (c *LoadingBar) Show() {
	c.shown = true
	if c.closeTimer != nil {
		c.closeTimer.Stop()
	}
	c.closeTimer = time.AfterFunc(closeTimeout, func() {
		c.shown = false
		c.tick = 0
	})
}

var _ component.Component = (*LoadingBar)(nil)
