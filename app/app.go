package app

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/zwh8800/rgbili/component"
)

type App struct {
	window     *sdl.Window
	surface    *sdl.Surface
	renderer   *sdl.Renderer
	running    bool
	windowHide bool

	cur component.Component // 当前组件
}

func NewApp() (*App, error) {
	if err := sdl.Init(sdl.INIT_AUDIO | sdl.INIT_VIDEO | sdl.INIT_GAMECONTROLLER); err != nil {
		return nil, err
	}

	joystick := sdl.JoystickOpen(0)
	fmt.Printf("joystick: %#v\n", joystick)

	window, err := sdl.CreateWindow("rgbili", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		640, 480, sdl.WINDOW_SHOWN|sdl.WINDOW_OPENGL)
	if err != nil {
		return nil, err
	}

	fmt.Printf("window: %#v\n", window)

	surface, err := window.GetSurface()
	if err != nil {
		return nil, err
	}

	renderer, err := window.GetRenderer()
	if err != nil {
		return nil, err
	}

	fmt.Printf("surface: %#v\n", surface)

	return &App{
		window:   window,
		surface:  surface,
		renderer: renderer,

		cur: component.NewTestComp(),
	}, nil
}

func (app *App) Run() {
	app.running = true
	for app.running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			app.handleEvent(event)
		}
		app.draw()
	}
}

func (app *App) handleEvent(e sdl.Event) {
	switch e.(type) {
	case *sdl.JoyButtonEvent:
		event := e.(*sdl.JoyButtonEvent)
		if event.State == sdl.PRESSED {
			if event.Button == 0x8 { // M
				app.running = false
			}

			if event.Button == 0x1 { // B
				if app.windowHide {
					app.window.Show()
					app.windowHide = false
				} else {
					app.window.Hide()
					app.windowHide = true
				}
			}
		}
	case *sdl.QuitEvent:
		println("Quit")
		app.running = false
		break
	}

	app.cur.HandleEvent(e)
}

func (app *App) draw() {
	app.cur.Draw(app.renderer)
}

func (app *App) Quit() {
	app.window.Destroy()
	sdl.Quit()
}
