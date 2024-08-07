package exit_mask

import (
	"image"
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
)

type ExitMask struct {
	shown      bool
	closeTimer *time.Timer
}

func New() *ExitMask {
	return &ExitMask{}
}

func (em *ExitMask) HandleEvent(e sdl.Event) {
	if !em.shown {
		return
	}

	switch event := e.(type) {
	case *sdl.JoyButtonEvent:
		if event.State == sdl.RELEASED {
			if event.Button == consts.ButtonA {
				em.buttonA()
			} else if event.Button == consts.ButtonB {
				em.buttonB()
			}
		}

	case *sdl.KeyboardEvent:
		if event.Type == sdl.KEYUP {
			if event.Keysym.Sym == sdl.K_a {
				em.buttonA()
			} else if event.Keysym.Sym == sdl.K_b {
				em.buttonB()
			}
		}
	}
}

func (em *ExitMask) Draw(renderer *sdl.Renderer) {
	if !em.shown {
		return
	}
	renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
	renderer.SetDrawColor(0, 0, 0, 128)
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

	em.drawExitMsg(renderer, textDrawer)
	em.drawExitButton(renderer, textDrawer)
}

func (em *ExitMask) drawExitMsg(renderer *sdl.Renderer, textDrawer *text.Drawer) {
	img, err := textDrawer.Draw("是否要退出RGTV", 64, image.White)
	if err != nil {
		panic(err)
	}

	x := (640 - img.Bounds().Dx()) / 2
	y := (480-img.Bounds().Dy())/2 - 40

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

func (em *ExitMask) drawExitButton(renderer *sdl.Renderer, textDrawer *text.Drawer) {
	img, err := textDrawer.Draw("A：确认  B：取消", 48, image.White)
	if err != nil {
		panic(err)
	}

	x := (640 - img.Bounds().Dx()) / 2
	y := (480-img.Bounds().Dy())/2 + 60

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

func (em *ExitMask) Dispose() {}

func (em *ExitMask) Show() {
	em.shown = true
	if em.closeTimer != nil {
		em.closeTimer.Stop()
	}
	em.closeTimer = time.AfterFunc(closeTimeout, func() {
		em.shown = false
	})
}

func (em *ExitMask) Hide() {
	em.shown = false
}

func (em *ExitMask) IsShown() bool {
	return em.shown
}

func (em *ExitMask) buttonA() {
	sdl.PushEvent(&sdl.QuitEvent{
		Type:      sdl.QUIT,
		Timestamp: sdl.GetTicks(),
	})
}

func (em *ExitMask) buttonB() {
	em.Hide()
}

var _ component.Component = (*ExitMask)(nil)
