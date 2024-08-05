package video_box

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"sync/atomic"
	"unsafe"

	ffmpeg "github.com/u2takey/ffmpeg-go"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/zwh8800/RGTV/component"
	"github.com/zwh8800/RGTV/conf"
	"github.com/zwh8800/RGTV/util"
)

const (
	videoBufSize = 640 * 480 * 3
	audioBufSize = 4096 * 2 * 2
)

type VideoBox struct {
	url string

	cmd *exec.Cmd

	rawVideoStream *os.File
	rawAudioStream *os.File

	once    sync.Once
	texture *sdl.Texture

	videoBuf    *[2][videoBufSize]byte
	videoBufIdx atomic.Int32

	audioBuf *util.RingBuffer

	audioVolume int

	pinner runtime.Pinner
}

func NewVideoBox(url string) (*VideoBox, error) {
	initAudio()

	v := &VideoBox{
		url:         url,
		audioVolume: 5,
	}

	v.pinner.Pin(v)
	v.videoBuf = new([2][videoBufSize]byte)
	v.pinner.Pin(v.videoBuf)
	v.audioBuf = util.NewLimitedBuffer(2 * audioBufSize)

	pr1, pw1, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	v.rawVideoStream = pr1
	pr2, pw2, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	v.rawAudioStream = pr2

	go v.runFFMPEG(pw1, pw2)
	go v.asyncReadVideo()
	go v.asyncReadAudio()

	vf.Set(v)
	return v, nil
}

func (v *VideoBox) runFFMPEG(pw1, pw2 *os.File) {
	defer pw1.Close()
	defer pw2.Close()

	i := ffmpeg.Input(v.url, ffmpeg.KwArgs{
		"re":      "",
		"headers": "User-Agent: Lavf",
	})

	out1 := i.Get("v").
		Filter("scale", ffmpeg.Args{"640:480:force_original_aspect_ratio=decrease"}).
		Filter("pad", ffmpeg.Args{"640:480:(ow-iw)/2:(oh-ih)/2"}).
		Filter("fps", ffmpeg.Args{"30"}).
		Output("pipe:3",
			ffmpeg.KwArgs{
				"format": "rawvideo", "pix_fmt": "rgb24",
			})
	out2 := i.Get("a").
		Output("pipe:4",
			ffmpeg.KwArgs{
				"format": "s16le",
				"ar":     44100,
				"ac":     2,
			})

	v.cmd = ffmpeg.MergeOutputs(out1, out2).
		ErrorToStdOut().
		SetFfmpegPath(conf.GetConfig().FFMPEGPath).
		Compile()

	v.cmd.ExtraFiles = []*os.File{
		pw1, pw2,
	}

	err := v.cmd.Run()
	if err != nil && !errors.As(err, new(*exec.ExitError)) {
		log.Printf("ffmpeg error: %s", err)
	}
}

func (v *VideoBox) asyncReadVideo() {
	for {
		bufSize := videoBufSize
		i := v.videoBufIdx.Load()
		nextIdx := (i + 1) % int32(len(v.videoBuf))

		n, err := io.ReadFull(v.rawVideoStream, v.videoBuf[nextIdx][0:bufSize])
		if n == 0 || err == io.EOF {
			return
		} else if n != bufSize || err != nil {
			fmt.Printf("read rawVideoStream error: %d, %s\n", n, err)
			return
		}
		v.videoBufIdx.Store(nextIdx)
	}
}

func (v *VideoBox) asyncReadAudio() {
	_, err := io.Copy(v.audioBuf, v.rawAudioStream)
	if err != nil {
		fmt.Printf("read rawAudioStream error: %s\n", err)
		return
	}
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

	err := v.texture.Update(nil, unsafe.Pointer(&v.videoBuf[v.videoBufIdx.Load()][0]), 640*3)
	if err != nil {
		panic(err)
	}
	err = renderer.Copy(v.texture, nil, nil)
	if err != nil {
		panic(err)
	}
}

func (v *VideoBox) Dispose() {
	vf.Set(nil)
	v.pinner.Unpin()
	v.rawVideoStream.Close()
	v.rawAudioStream.Close()
	v.texture.Destroy()
	if v.cmd != nil && v.cmd.Process != nil {
		v.cmd.Process.Kill()
	}
}

func (v *VideoBox) VolumeUp() {
	if v.audioVolume >= 10 {
		v.audioVolume = 10
		return
	}
	v.audioVolume++
}
func (v *VideoBox) VolumeDown() {
	if v.audioVolume <= 0 {
		v.audioVolume = 0
		return
	}
	v.audioVolume--
}

var _ component.Component = (*VideoBox)(nil)
