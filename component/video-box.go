package component

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"unsafe"

	ffmpeg "github.com/u2takey/ffmpeg-go"
	"github.com/veandco/go-sdl2/sdl"
)

type VideoBox struct {
	rawVideoStream *io.PipeReader
}

func NewVideoBox() *VideoBox {
	v := &VideoBox{}

	pr, pw := io.Pipe()
	v.rawVideoStream = pr

	go func() {
		defer pw.Close()
		err := ffmpeg.Input("/Users/wastecat/Downloads/IMG_0981.MOV").
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
	texture, err := renderer.CreateTexture(sdl.PIXELFORMAT_RGB24, sdl.TEXTUREACCESS_STREAMING, 640, 480)
	if err != nil {
		panic(err)
	}

	frameSize := 640 * 480 * 3
	buf := make([]byte, frameSize)
	n, err := io.ReadFull(v.rawVideoStream, buf)
	if n == 0 || err == io.EOF {
		return
	} else if n != frameSize || err != nil {
		panic(fmt.Sprintf("read error: %d, %s", n, err))
	}

	err = texture.Update(nil, unsafe.Pointer(&buf[0]), 640*3)
	if err != nil {
		panic(err)
	}
	renderer.Clear()
	renderer.Copy(texture, nil, nil)
	renderer.Present()
}

func getVideoSize(fileName string) (int, int) {
	log.Println("Getting video size for", fileName)
	data, err := ffmpeg.Probe(fileName)
	if err != nil {
		panic(err)
	}
	log.Println("got video info", data)
	type VideoInfo struct {
		Streams []struct {
			CodecType string `json:"codec_type"`
			Width     int
			Height    int
		} `json:"streams"`
	}
	vInfo := &VideoInfo{}
	err = json.Unmarshal([]byte(data), vInfo)
	if err != nil {
		panic(err)
	}
	for _, s := range vInfo.Streams {
		if s.CodecType == "video" {
			return s.Width, s.Height
		}
	}
	return 0, 0
}
