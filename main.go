package main

import (
	"github.com/veandco/go-sdl2/sdl"

	"github.com/zwh8800/rgbili/app"
)

func main() {
	sdl.Main(func() {
		sdl.Do(func() {
			app, err := app.NewApp()
			if err != nil {
				panic(err)
			}
			defer app.Quit()

			app.Run()
		})
	})
}
