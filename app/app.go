package app

import "github.com/veandco/go-sdl2/sdl"

type App struct {
	window  *sdl.Window
	surface *sdl.Surface
	running bool

	testRect sdl.Rect
}

func NewApp() (*App, error) {
	if err := sdl.Init(sdl.INIT_AUDIO | sdl.INIT_VIDEO); err != nil {
		return nil, err
	}

	window, err := sdl.CreateWindow("rgbili", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		640, 480, sdl.WINDOW_SHOWN)
	if err != nil {
		return nil, err
	}

	surface, err := window.GetSurface()
	if err != nil {
		return nil, err
	}

	return &App{
		window:  window,
		surface: surface,
	}, nil
}

func (app *App) Run() {
	app.running = true
	app.init()
	for app.running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			app.handleEvent(event)
			app.draw()
		}
	}
}

func (app *App) init() {
	app.testRect = sdl.Rect{0, 0, 200, 200}
}

func (app *App) handleEvent(event sdl.Event) {
	switch event.(type) {
	case *sdl.KeyboardEvent:
		kEvent := event.(*sdl.KeyboardEvent)
		if kEvent.Type == sdl.KEYDOWN {
			switch kEvent.Keysym.Sym {
			case RawUp:
				app.testRect.Y -= 10
			case RawDown:
				app.testRect.Y += 10
			case RawLeft:
				app.testRect.X -= 10
			case RawRight:
				app.testRect.X += 10
			}
		}

	case *sdl.QuitEvent:
		println("Quit")
		app.running = false
		break
	}
}

func (app *App) draw() {
	app.surface.FillRect(nil, 0)

	colour := sdl.Color{R: 255, G: 0, B: 255, A: 255} // purple
	pixel := sdl.MapRGBA(app.surface.Format, colour.R, colour.G, colour.B, colour.A)
	app.surface.FillRect(&app.testRect, pixel)
	app.window.UpdateSurface()
}

func (app *App) Quit() {
	app.window.Destroy()
	sdl.Quit()
}
