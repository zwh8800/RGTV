package component

import (
	"fmt"
	"io"
	"sync"
	"unsafe"

	ffmpeg "github.com/u2takey/ffmpeg-go"
	"github.com/veandco/go-sdl2/sdl"
)

type VideoBox struct {
	rawVideoStream *io.PipeReader

	once    sync.Once
	texture *sdl.Texture
}

func NewVideoBox() *VideoBox {
	v := &VideoBox{}

	pr, pw := io.Pipe()
	v.rawVideoStream = pr

	go func() {
		defer pw.Close()
		err := ffmpeg.Input("/mnt/mmc/Video/猫和老鼠/005.mp4", ffmpeg.KwArgs{"re": ""}).
			SetFfmpegPath("/root/code/go/rgbili/ffmpeg").
			Filter("scale", ffmpeg.Args{"640:480"}).
			Output("pipe:",
				ffmpeg.KwArgs{
					"format": "rawvideo", "pix_fmt": "rgb24",
				}).
			WithOutput(pw).
			ErrorToStdOut().
			Run()
		if err != nil {
			panic(err)
		}
	}()

	return v
}

func (v *VideoBox) HandleEvent(e sdl.Event) {

}

func (v *VideoBox) Draw(renderer *sdl.Renderer) {
	v.once.Do(func() {
		texture, err := renderer.CreateTexture(sdl.PIXELFORMAT_RGB24, sdl.TEXTUREACCESS_STREAMING, 640, 480)
		if err != nil {
			panic(err)
		}
		v.texture = texture
	})

	frameSize := 640 * 480 * 3
	buf := make([]byte, frameSize)
	n, err := io.ReadFull(v.rawVideoStream, buf)
	if n == 0 || err == io.EOF {
		return
	} else if n != frameSize || err != nil {
		panic(fmt.Sprintf("read error: %d, %s", n, err))
	}

	err = v.texture.Update(nil, unsafe.Pointer(&buf[0]), 640*3)
	if err != nil {
		panic(err)
	}
	renderer.Clear()
	renderer.Copy(v.texture, nil, nil)
	renderer.Present()
}
