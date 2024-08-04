package component

import (
	"fmt"
	"os/exec"

	"github.com/veandco/go-sdl2/sdl"
)

type TestComp struct {
	testRect sdl.Rect
	cmd      *exec.Cmd
	videoBox *VideoBox
}

func NewTestComp() *TestComp {
	return &TestComp{
		testRect: sdl.Rect{0, 0, 200, 200},
	}
}

func (c *TestComp) HandleEvent(e sdl.Event) {
	fmt.Printf("%#v\n", e)
	switch e.(type) {
	case *sdl.KeyboardEvent:
		event := e.(*sdl.KeyboardEvent)
		if event.Type == sdl.KEYUP {
			if event.Keysym.Sym == sdl.K_a {
				fmt.Println("run videobox")
				c.videoBox, _ = NewVideoBox("")
			}
		}

	case *sdl.JoyHatEvent:
		event := e.(*sdl.JoyHatEvent)
		if event.Value&sdl.HAT_LEFT != 0 {
			c.testRect.X -= 10
		} else if event.Value&sdl.HAT_RIGHT != 0 {
			c.testRect.X += 10
		} else if event.Value&sdl.HAT_UP != 0 {
			c.testRect.Y -= 10
		} else if event.Value&sdl.HAT_DOWN != 0 {
			c.testRect.Y += 10
		}

	case *sdl.JoyButtonEvent:
		event := e.(*sdl.JoyButtonEvent)
		if event.State == sdl.PRESSED {
			if event.Button == 0x0 { // A
				c.testRect.X = 0
				c.testRect.Y = 0
			}

			if event.Button == 0x2 { // Y
				fmt.Println("run cmd")
				if c.cmd != nil {
					err := c.cmd.Process.Kill()
					if err != nil {
						fmt.Println("cmd kill err:", err.Error())
					}
				}

				c.cmd = exec.Command(
					"/mnt/vendor/bin/video/ffplay",
					"-fs", "-autoexit", "-vf", "scale=640:-2", "-i", "/mnt/mmc/Video/猫和老鼠/002.mp4",
				)
				err := c.cmd.Run()
				if err != nil {
					fmt.Println("cmd run err:", err.Error())
				}
			}

			if event.Button == 0x3 { // X
				fmt.Println("run videobox")
				c.videoBox, _ = NewVideoBox("")
			}
		}

	case *sdl.JoyAxisEvent:
		event := e.(*sdl.JoyAxisEvent)
		if event.Axis == 0 { // 左摇杆左右
			c.testRect.X = int32(event.Value / 1000)
		} else if event.Axis == 1 { // 左摇杆上下
			c.testRect.Y = int32(event.Value / 1000)
		}
	}
}

func (c *TestComp) Draw(renderer *sdl.Renderer) {
	if c.videoBox != nil {
		c.videoBox.Draw(renderer)
	} else {
		renderer.Clear()
		renderer.SetDrawColor(0, 0, 0, 255)
		renderer.FillRect(nil)
		color := sdl.Color{R: 255, G: 0, B: 255, A: 255}
		renderer.SetDrawColor(color.R, color.G, color.B, color.A)
		renderer.FillRect(&c.testRect)
		renderer.Present()
	}
}

func (c *TestComp) Dispose() {}
