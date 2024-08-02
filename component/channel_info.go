package component

import (
	"github.com/veandco/go-sdl2/sdl"
)

type ChannelInfo struct {
}

func NewChannelInfo() *ChannelInfo {
	return &ChannelInfo{}
}

func (c *ChannelInfo) HandleEvent(e sdl.Event) {

}

func (c *ChannelInfo) Draw(renderer *sdl.Renderer) {

}

func (c *ChannelInfo) Dispose() {

}
