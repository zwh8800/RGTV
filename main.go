package main

import app "github.com/zwh8800/rgbili/app"

func main() {
	app, err := app.NewApp()
	if err != nil {
		panic(err)
	}
	defer app.Quit()

	app.Run()
}
