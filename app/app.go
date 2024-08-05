package app

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/zwh8800/RGTV/component"
	"github.com/zwh8800/RGTV/component/main_frame"
	"github.com/zwh8800/RGTV/conf"
	"github.com/zwh8800/RGTV/consts"
	"github.com/zwh8800/RGTV/util"
)

type App struct {
	window   *sdl.Window
	surface  *sdl.Surface
	renderer *sdl.Renderer
	running  bool
	fpsm     *util.FPSManager

	cur component.Component // 当前组件
}

func NewApp() (*App, error) {
	if err := sdl.Init(sdl.INIT_AUDIO | sdl.INIT_VIDEO | sdl.INIT_GAMECONTROLLER); err != nil {
		return nil, err
	}

	printDebugInfo()

	joystick := sdl.JoystickOpen(0)
	fmt.Printf("joystick: %#v\n", joystick)

	window, err := sdl.CreateWindow("RGTV", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		conf.GetConfig().ResX, conf.GetConfig().ResY, sdl.WINDOW_SHOWN)
	if err != nil {
		return nil, err
	}

	fmt.Printf("window: %#v\n", window)

	sdl.ShowSimpleMessageBox(sdl.MESSAGEBOX_INFORMATION, "测试", "gggggg", window)

	surface, err := window.GetSurface()
	if err != nil {
		return nil, err
	}

	renderer, err := window.GetRenderer()
	if err != nil {
		return nil, err
	}

	fmt.Printf("surface: %#v\n", surface)

	fpsm := util.NewFPSManager(30)

	return &App{
		window:   window,
		surface:  surface,
		renderer: renderer,
		fpsm:     fpsm,

		cur: main_frame.NewMainFrame(),
	}, nil
}

func (app *App) Run() {
	app.running = true
	for app.running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			app.handleEvent(event)
		}
		app.draw()
		app.fpsm.Wait()
	}
}

func (app *App) handleEvent(e sdl.Event) {
	switch event := e.(type) {
	case *sdl.JoyButtonEvent:
		if event.State == sdl.PRESSED {
			if event.Button == consts.ButtonMenu {
				app.running = false
			}
		}
	case *sdl.QuitEvent:
		fmt.Println("Quit")
		app.running = false
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

func printDebugInfo() {
	n, _ := sdl.GetNumVideoDrivers()
	for i := 0; i < n; i++ {
		fmt.Println("VideoDriver:", sdl.GetVideoDriver(i))
	}
	n = sdl.GetNumAudioDrivers()
	for i := 0; i < n; i++ {
		fmt.Println("AudioDriver:", sdl.GetAudioDriver(i))
	}

	n, _ = sdl.GetNumRenderDrivers()
	for i := 0; i < n; i++ {
		var renderDriverInfo sdl.RendererInfo
		sdl.GetRenderDriverInfo(i, &renderDriverInfo)
		fmt.Printf("RenderDriver:%#v\n", renderDriverInfo)
	}

	n, _ = sdl.GetNumDisplayModes(0)
	for i := 0; i < n; i++ {
		mode, _ := sdl.GetDisplayMode(0, i)
		fmt.Printf("DisplayMode:%#v\n", mode)
	}

	n, _ = sdl.GetNumVideoDisplays()
	fmt.Println("NumVideoDisplays:", n)

	n = sdl.GetNumAudioDevices(false)
	for i := 0; i < n; i++ {
		name := sdl.GetAudioDeviceName(i, false)
		spec, _ := sdl.GetAudioDeviceSpec(i, false)
		fmt.Printf("AudioDevice:%s, %#v\n", name, spec)
	}

	vd, _ := sdl.GetCurrentVideoDriver()
	fmt.Println("VideoDriver:", vd)
	ad := sdl.GetCurrentAudioDriver()
	fmt.Println("AudioDriver:", ad)
}
