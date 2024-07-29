package component

import "github.com/veandco/go-sdl2/sdl"

type Component interface {
	HandleEvent(e sdl.Event)
	Draw(renderer *sdl.Renderer)
}
