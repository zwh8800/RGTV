package app

import (
	"fmt"
	"os/exec"

	"github.com/veandco/go-sdl2/sdl"
)

type App struct {
	window  *sdl.Window
	surface *sdl.Surface
	running bool

	testRect   sdl.Rect
	windowHide bool
	cmd        *exec.Cmd
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

	fmt.Printf("surface: %#v\n", surface)

	return &App{
		window:  window,
		surface: surface,

		testRect: sdl.Rect{0, 0, 200, 200},
	}, nil
}

func (app *App) Run() {
	app.running = true
	for app.running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			app.handleEvent(event)
			app.draw()
		}
	}
}

func (app *App) handleEvent(e sdl.Event) {
	fmt.Printf("%#v\n", e)
	switch e.(type) {
	case *sdl.JoyHatEvent:
		event := e.(*sdl.JoyHatEvent)
		if event.Value&sdl.HAT_LEFT != 0 {
			app.testRect.X -= 10
		} else if event.Value&sdl.HAT_RIGHT != 0 {
			app.testRect.X += 10
		} else if event.Value&sdl.HAT_UP != 0 {
			app.testRect.Y -= 10
		} else if event.Value&sdl.HAT_DOWN != 0 {
			app.testRect.Y += 10
		}

	case *sdl.JoyButtonEvent:
		event := e.(*sdl.JoyButtonEvent)
		if event.State == sdl.PRESSED {
			if event.Button == 0x0 { // A
				app.testRect.X = 0
				app.testRect.Y = 0
			}

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

			if event.Button == 0x2 { // Y
				fmt.Println("run cmd")
				if app.cmd != nil {
					err := app.cmd.Process.Kill()
					if err != nil {
						fmt.Println("cmd kill err:", err.Error())
					}
				}

				app.cmd = exec.Command(
					"/mnt/vendor/bin/video/ffplay",
					"-fs", "-autoexit", "-vf", "scale=640:-2", "-i", "/mnt/mmc/Video/猫和老鼠/002.mp4",
				)
				err := app.cmd.Run()
				if err != nil {
					fmt.Println("cmd run err:", err.Error())
				}
			}
		}

	case *sdl.JoyAxisEvent:
		event := e.(*sdl.JoyAxisEvent)
		if event.Axis == 0 { // 左摇杆左右
			app.testRect.X = int32(event.Value / 1000)
		} else if event.Axis == 1 { // 左摇杆上下
			app.testRect.Y = int32(event.Value / 1000)
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
